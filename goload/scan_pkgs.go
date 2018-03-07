package goload

import (
	"go/build"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
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

// ScanOptions provides the options for scanning a Go language repository.
type ScanOptions struct {
	Context *build.Context

	// TestdataWhiteList provides a whitelist of "testdata" packages that are
	// valid ones and being imported.
	TestdataWhiteList map[string]bool

	// PkgBlackList is a list of packages that will be skipped. It will also
	// skip its sub packages.
	PkgBlackList map[string]bool
}

func inSet(s map[string]bool, k string) bool {
	if s == nil {
		return false
	}
	return s[k]
}

type scanner struct {
	path    string
	ctx     *build.Context
	srcRoot string
	opts    *ScanOptions
	res     *ScanResult
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

	ret.srcRoot = filepath.Join(ret.ctx.GOPATH, "src")

	return ret
}

func (s *scanner) handleDir(dir, path string) error {
	if inSet(s.opts.PkgBlackList, path) {
		return filepath.SkipDir
	}

	base := filepath.Base(path)
	if strings.HasPrefix(base, "_") || strings.HasPrefix(base, ".") {
		return filepath.SkipDir
	}

	switch base {
	case "testdata":
		if !inSet(s.opts.TestdataWhiteList, path) {
			return filepath.SkipDir
		}
	case "vendor":
		s.res.HasVendor = true
	case "internal":
		s.res.HasInternal = true
	}

	pkg, err := s.ctx.Import(path, "", 0) // check if it is a package
	if err != nil {
		if isNoGoError(err) {
			return nil
		}
		return err
	}

	if len(pkg.GoFiles) == 0 && len(pkg.CgoFiles) == 0 {
		return nil
	}

	s.res.Pkgs[path] = pkg
	return nil
}

func (s *scanner) walk(dir, p string) error {
	info, err := os.Lstat(dir)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return nil
	}

	if err := s.handleDir(dir, p); err != nil {
		if err == filepath.SkipDir {
			return nil
		}
		return err
	}

	f, err := os.Open(dir)
	if err != nil {
		return err
	}

	defer f.Close()
	names, err := f.Readdirnames(-1)
	if err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}

	sort.Strings(names)
	for _, name := range names {
		subDir := filepath.Join(dir, name)
		subPath := path.Join(p, name)
		if err := s.walk(subDir, subPath); err != nil {
			return err
		}
	}

	return nil
}

// ScanPkgs scans all packages under a package path.
func ScanPkgs(p string, opts *ScanOptions) (*ScanResult, error) {
	s := newScanner(p, opts)

	// First check if the folder can be found.
	s.res = newScanResult(p)
	dir := filepath.Join(s.srcRoot, filepath.ToSlash(p))
	if err := s.walk(dir, p); err != nil && err != filepath.SkipDir {
		return nil, err
	}
	return s.res, nil
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
