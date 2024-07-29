package source

import (
	"fmt"
	"regexp"

	"github.com/golang-migrate/migrate/v4/source"
)

var (
	ErrParse = fmt.Errorf("no match")
)

var (
	DefaultParse = Parse
	DefaultRegex = Regex
)

// Regex matches the following pattern:
//
//	123_name.up.ext
//	123_name.down.ext
var Regex = regexp.MustCompile(`^([0-9]+)_(.*)\.(` + string(source.Down) + `|` + string(source.Up) + `)\.(.*)$`)

// Parse returns Migration for matching Regex pattern.
func Parse(version uint, raw string) (*source.Migration, error) {
	m := Regex.FindStringSubmatch(raw)
	if len(m) == 5 {
		return &source.Migration{
			Version:    version,
			Identifier: m[2],
			Direction:  source.Direction(m[3]),
			Raw:        raw,
		}, nil
	}
	return nil, ErrParse
}
