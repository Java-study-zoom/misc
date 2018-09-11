package goload

import (
	"path"
	"path/filepath"
)

type scanDir struct {
	dir    string
	path   string // import path
	base   string
	vendor *vendorLayer
}

func (d *scanDir) sub(sub string) *scanDir {
	ret := &scanDir{
		dir:    filepath.Join(d.dir, sub),
		path:   path.Join(d.path, sub),
		base:   sub,
		vendor: d.vendor,
	}

	return ret
}
