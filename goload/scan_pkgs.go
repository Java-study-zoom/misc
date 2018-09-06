package goload

import (
	"fmt"
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

type scanDir struct {
	dir    string
	path   string // import path
	base   string
	vendor *vendorLayer

	modPath string // import path with module enabled
	modRoot string // root of module
}

func (d *scanDir) sub(sub string) *scanDir {
	ret := &scanDir{
		dir:     filepath.Join(d.dir, sub),
		path:    path.Join(d.path, sub),
		base:    sub,
		vendor:  d.vendor,
		modRoot: d.modRoot,
	}

	if d.modPath != "" {
		ret.modPath = path.Join(d.modPath, sub)
	}
	return ret
}

type scanner struct {
	path    string
	ctx     *build.Context
	srcRoot string
	opts    *ScanOptions
	res     *ScanResult

	vendorStack    *vendorStack
	vendorLayers   map[string]*vendorLayer
	vendorScanning bool
}

func newScanner(p string, opts *ScanOptions) *scanner {
	if opts == nil {
		opts = new(ScanOptions)
	}

	ret := &scanner{
		path:         p,
		opts:         opts,
		vendorStack:  new(vendorStack),
		vendorLayers: make(map[string]*vendorLayer),
	}

	if opts.Context != nil {
		ret.ctx = opts.Context
	} else {
		ret.ctx = &build.Default
	}

	ret.srcRoot = filepath.Join(ret.ctx.GOPATH, "src")

	return ret
}

func (s *scanner) skipDir(dir *scanDir) bool {
	if inSet(s.opts.PkgBlackList, dir.path) {
		return true
	}
	base := dir.base
	if strings.HasPrefix(base, "_") || strings.HasPrefix(base, ".") {
		return true
	}
	if base == "testdata" && !inSet(s.opts.TestdataWhiteList, dir.path) {
		return true
	}
	return false
}

func (s *scanner) handleDir(dir *scanDir) error {
	switch dir.base {
	case "vendor":
		s.res.HasVendor = true
	case "internal":
		s.res.HasInternal = true
	}

	mode := build.ImportComment

	if s.vendorScanning {
		if inSet(s.opts.PkgBlackList, "!"+dir.path) {
			return nil
		}

		// check if it is a package
		pkg, err := s.ctx.Import(dir.path, "", mode)
		if err != nil {
			if isNoGoError(err) {
				return nil
			}
			return err
		}

		if len(pkg.GoFiles) == 0 && len(pkg.CgoFiles) == 0 {
			return nil
		}

		if dir.vendor != nil {
			dir.vendor.addPkg(dir.path)
		}

		s.res.Pkgs[dir.path] = &Package{
			Build: pkg,
		}
		fmt.Printf("%s - %s\n", dir.path, pkg.ImportPath)
	} else {
		pkg, found := s.res.Pkgs[dir.path]
		if !found {
			return nil
		}

		importMap := make(map[string]string)
		for _, imp := range pkg.Build.Imports {
			mapped, hit := s.vendorStack.mapImport(imp)
			if !hit {
				continue
			}
			importMap[imp] = mapped
		}

		if len(importMap) > 0 {
			pkg.ImportMap = importMap
		}
	}
	return nil
}

func (s *scanner) walk(dir *scanDir) error {
	info, err := os.Lstat(dir.dir)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return nil
	}

	if s.skipDir(dir) {
		return nil
	}

	f, err := os.Open(dir.dir)
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

	if !s.vendorScanning {
		index := sort.SearchStrings(names, "vendor")
		hasVendor := index < len(names) && names[index] == "vendor"
		if hasVendor {
			ly := s.vendorLayers[dir.path]
			if ly == nil {
				panic(fmt.Sprintf("vendor layer missing: %s", dir.path))
			}
			if len(ly.pkgs) > 0 {
				s.vendorStack.push(ly)
				defer s.vendorStack.pop()
			}
		}
	}

	if err := s.handleDir(dir); err != nil {
		if err == filepath.SkipDir {
			return nil
		}
		return err
	}

	for _, name := range names {
		sub := dir.sub(name)

		if s.vendorScanning && name == "vendor" {
			ly := newVendorLayer(dir.path)
			s.vendorLayers[dir.path] = ly
			sub.vendor = ly
		}

		if err := s.walk(sub); err != nil {
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
	dir := &scanDir{
		dir:  filepath.Join(s.srcRoot, filepath.ToSlash(p)),
		path: p,
		base: path.Base(p),
	}

	for _, scanning := range []bool{true, false} {
		s.vendorScanning = scanning
		if err := s.walk(dir); err != nil && err != filepath.SkipDir {
			return nil, err
		}
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
