package goload

// Pkg is a Go language package.
type Pkg struct {
	Path    string
	Imports []string
}

func newPkg(path string, imports []string) *Pkg {
	return &Pkg{
		Path:    path,
		Imports: imports,
	}
}

// ScanResult has the scanning result
type ScanResult struct {
	Repo        string
	Pkgs        map[string]*Pkg
	HasVendor   bool
	HasInternal bool
}

func newScanResult(repo string) *ScanResult {
	return &ScanResult{
		Repo: repo,
		Pkgs: make(map[string]*Pkg),
	}
}
