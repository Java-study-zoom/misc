package goload

import (
	"go/build"
)

// Package is a package in the scan result.
type Package struct {
	Build     *build.Package
	ImportMap map[string]string

	ModRoot string
	ModPath string // alias import path when module is enabled
}

// ScanResult has the scanning result
type ScanResult struct {
	Repo        string
	Pkgs        map[string]*Package
	HasVendor   bool
	HasInternal bool
}

func newScanResult(repo string) *ScanResult {
	return &ScanResult{
		Repo: repo,
		Pkgs: make(map[string]*Package),
	}
}
