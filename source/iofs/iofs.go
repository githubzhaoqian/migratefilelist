//go:build go1.16
// +build go1.16

package iofs

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"strconv"

	"github.com/golang-migrate/migrate/v4/source"

	flsource "github.com/githubzhaoqian/migrate-filelist/source"
)

const (
	errSchemaRecordIsEmptyFormat   = "schema record %d is empty"
	errSchemaItemIsDuplicateFormat = "schema record %s is duplicate"
)

type driver struct {
	PartialDriver
}

// New returns a new Driver from io/fs#FS and a relative path.
func New(fileName, dir string) (source.Driver, error) {
	var i driver
	if err := i.Init(fileName, dir); err != nil {
		return nil, fmt.Errorf("failed to init driver with file %s: %w", fileName, err)
	}
	return &i, nil
}

// Open is part of source.Driver interface implementation.
// Open cannot be called on the iofs passthrough driver.
func (d *driver) Open(url string) (source.Driver, error) {
	return nil, errors.New("Open() cannot be called on the iofs passthrough driver")
}

// PartialDriver is a helper service for creating new source drivers working with
// io/fs.FS instances. It implements all source.Driver interface methods
// except for Open(). New driver could embed this struct and add missing Open()
// method.
//
// To prepare PartialDriver for use Init() function.
type PartialDriver struct {
	migrations *source.Migrations
	schemaFile string
	path       string
}

// Init prepares not initialized IoFS instance to read migrations from a
// io/fs#FS instance and a relative path.
func (d *PartialDriver) Init(fileName, dir string) error {
	schemaFile, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer schemaFile.Close()
	scanner := bufio.NewScanner(schemaFile)
	itemMap := make(map[string]struct{})

	ms := source.NewMigrations()
	var (
		index uint
	)
	for scanner.Scan() {
		index++
		line := scanner.Text()
		if line == "" {
			if !scanner.Scan() {
				break
			}
			return fmt.Errorf(errSchemaRecordIsEmptyFormat, index)
		}
		if _, ok := itemMap[line]; ok {
			return fmt.Errorf(errSchemaItemIsDuplicateFormat, line)
		}
		file, err := os.Open(line)
		if err != nil {
			return err
		}
		fileInfo, err := file.Stat()
		if err != nil {
			return err
		}
		err = file.Close()
		if err != nil {
			return err
		}
		m, err := flsource.DefaultParse(index, line)
		if err != nil {
			return err
		}
		if !ms.Append(m) {
			return source.ErrDuplicateMigration{
				Migration: *m,
				FileInfo:  fileInfo,
			}
		}
	}
	d.path = dir
	d.schemaFile = fileName
	d.migrations = ms
	return nil
}

// Close is part of source.Driver interface implementation.
// Closes the file system if possible.
func (d *PartialDriver) Close() error {
	return nil
}

// First is part of source.Driver interface implementation.
func (d *PartialDriver) First() (version uint, err error) {
	if version, ok := d.migrations.First(); ok {
		return version, nil
	}
	return 0, &fs.PathError{
		Op:   "first",
		Path: d.schemaFile,
		Err:  fs.ErrNotExist,
	}
}

// Prev is part of source.Driver interface implementation.
func (d *PartialDriver) Prev(version uint) (prevVersion uint, err error) {
	if version, ok := d.migrations.Prev(version); ok {
		return version, nil
	}
	return 0, &fs.PathError{
		Op:  "prev for version " + strconv.FormatUint(uint64(version), 10),
		Err: fs.ErrNotExist,
	}
}

// Next is part of source.Driver interface implementation.
func (d *PartialDriver) Next(version uint) (nextVersion uint, err error) {
	if version, ok := d.migrations.Next(version); ok {
		return version, nil
	}
	return 0, &fs.PathError{
		Op:  "next for version " + strconv.FormatUint(uint64(version), 10),
		Err: fs.ErrNotExist,
	}
}

// ReadUp is part of source.Driver interface implementation.
func (d *PartialDriver) ReadUp(version uint) (r io.ReadCloser, identifier string, err error) {
	if m, ok := d.migrations.Up(version); ok {
		body, err := d.open(path.Join(d.path, m.Raw))
		if err != nil {
			return nil, "", err
		}
		return body, m.Identifier, nil
	}
	return nil, "", &fs.PathError{
		Op:  "read up for version " + strconv.FormatUint(uint64(version), 10),
		Err: fs.ErrNotExist,
	}
}

// ReadDown is part of source.Driver interface implementation.
func (d *PartialDriver) ReadDown(version uint) (r io.ReadCloser, identifier string, err error) {
	if m, ok := d.migrations.Down(version); ok {
		body, err := d.open(path.Join(d.path, m.Raw))
		if err != nil {
			return nil, "", err
		}
		return body, m.Identifier, nil
	}
	return nil, "", &fs.PathError{
		Op:   "read down for version " + strconv.FormatUint(uint64(version), 10),
		Path: d.schemaFile,
		Err:  fs.ErrNotExist,
	}
}

func (d *PartialDriver) open(schemaFile string) (fs.File, error) {
	f, err := os.Open(schemaFile)
	if err == nil {
		return f, nil
	}
	// Some non-standard file systems may return errors that don't include the path, that
	// makes debugging harder.
	if !errors.As(err, new(*fs.PathError)) {
		err = &fs.PathError{
			Op:   "open",
			Path: schemaFile,
			Err:  err,
		}
	}
	return nil, err
}
