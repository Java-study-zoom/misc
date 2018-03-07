package goload

import (
	"go/build"
)

// ScanResult has the scanning result
type ScanResult struct {
	Repo        string
	Pkgs        map[string]*build.Package
	ImportMap   map[string]map[string]string
	HasVendor   bool
	HasInternal bool
}

func newScanResult(repo string) *ScanResult {
	return &ScanResult{
		Repo:      repo,
		Pkgs:      make(map[string]*build.Package),
		ImportMap: make(map[string]map[string]string),
	}
}
