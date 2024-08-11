package filelist

import (
	nurl "net/url"
	"os"
	"path/filepath"

	"github.com/githubzhaoqian/migratefilelist/source/iofs"

	"github.com/golang-migrate/migrate/v4/source"
)

func init() {
	source.Register("filelist", &FileList{})
}

type FileList struct {
	iofs.PartialDriver
	url  string
	path string
}

func (f *FileList) Open(url string) (source.Driver, error) {
	fileName, err := parseURL(url)
	if err != nil {
		return nil, err
	}
	nf := &FileList{
		url:  url,
		path: fileName,
	}
	dir := filepath.Dir(fileName)
	if err := nf.Init(fileName, dir); err != nil {
		return nil, err
	}
	return nf, nil
}

func parseURL(url string) (string, error) {
	u, err := nurl.Parse(url)
	if err != nil {
		return "", err
	}
	// concat host and path to restore full path
	// host might be `.`
	p := u.Opaque
	if len(p) == 0 {
		p = u.Host + u.Path
	}

	if len(p) == 0 {
		// default to current directory if no path
		wd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		p = wd

	} else if p[0:1] == "." || p[0:1] != "/" {
		// make path absolute if relative
		abs, err := filepath.Abs(p)
		if err != nil {
			return "", err
		}
		p = abs
	}
	return p, nil
}
