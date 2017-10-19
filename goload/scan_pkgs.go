package goload

import (
	"go/build"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"shanhu.io/misc/strutil"
)

func isNoGoError(e error) bool {
	if e == nil {
		return false
	}
	_, hit := e.(*build.NoGoError)
	return hit
}

// MakeContext makes a build context for the given GOPATH.
func MakeContext(gopath string) *build.Context {
	ctx := build.Default
	if gopath != "" {
		ctx.GOPATH = gopath
	}
	return &ctx
}

// ScanOptions provides the options for scanning a
// Go language repository.
type ScanOptions struct {
	Context *build.Context
}

type scanner struct {
	path string
	ctx  *build.Context
	opts *ScanOptions
}

func newScanner(p string, opts *ScanOptions) *scanner {
	if opts == nil {
		opts = new(ScanOptions)
	}

	ret := &scanner{
		path: p,
		opts: opts,
	}

	if opts.Context != nil {
		ret.ctx = opts.Context
	} else {
		ret.ctx = &build.Default
	}

	return ret
}

// some (bad) repos use testdata folder to save code and import them.
var testdataWhiteList = func() map[string]bool {
	return strutil.MakeSet([]string{
		"github.com/golang/protobuf/proto/testdata",
		"google.golang.org/grpc/testdata",
	})
}()

// ScanPkgs scans all packages under a package path.
func ScanPkgs(p string, opts *ScanOptions) (*ScanResult, error) {
	s := newScanner(p, opts)

	// First check if the folder can be found.
	pkg, e := s.ctx.Import(p, "", build.FindOnly)
	if e != nil && !isNoGoError(e) {
		return nil, e
	}

	ret := newScanResult(p)
	walk := func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}

		path, e := filepath.Rel(pkg.SrcRoot, p)
		if e != nil {
			return e
		}
		base := filepath.Base(path)

		if strings.HasPrefix(base, "_") || strings.HasPrefix(base, ".") {
			return filepath.SkipDir
		}

		switch base {
		case "testdata":
			if testdataWhiteList[path] {
				break
			}
			return filepath.SkipDir
		case "vendor":
			ret.HasVendor = true
			return filepath.SkipDir
		case "internal":
			ret.HasInternal = true
		}

		pkg, err := s.ctx.Import(path, "", 0) // check if it is a package
		if err != nil {
			if isNoGoError(err) { // not a go language package
				return nil
			}
			return err
		}

		ret.Pkgs[path] = newPkg(path, pkg.Imports)
		return nil
	}

	if err := filepath.Walk(pkg.Dir, walk); err != nil {
		return nil, err
	}

	return ret, nil
}

// ListPkgs list all packages under a package path.
func ListPkgs(p string) ([]string, error) {
	res, err := ScanPkgs(p, nil)
	if err != nil {
		return nil, err
	}

	var lst []string
	for pkg := range res.Pkgs {
		lst = append(lst, pkg)
	}
	sort.Strings(lst)
	return lst, nil
}
