package ziputil

import (
	"testing"

	"archive/zip"
	"io"
	"io/ioutil"
	"os"
	"path"
	"reflect"

	"shanhu.io/misc/tempfile"
)

func testDiffFile(t *testing.T, f1, f2 string) bool {
	ne := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}

	bs1, err := ioutil.ReadFile(f1)
	ne(err)

	bs2, err := ioutil.ReadFile(f2)
	ne(err)

	if !reflect.DeepEqual(bs1, bs2) {
		return false
	}

	s1, err := os.Stat(f1)
	ne(err)

	s2, err := os.Stat(f2)
	ne(err)

	if s1.Mode() != s2.Mode() {
		return false
	}
	return true
}

func TestZipFile(t *testing.T) {
	ne := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}

	temp, err := tempfile.NewFile("", "ziputil")
	ne(err)
	defer temp.CleanUp()

	const p = "testdata/testfile"
	ne(ZipFile(p, temp))

	size, err := temp.Seek(0, io.SeekCurrent)
	ne(err)

	ne(temp.Reset())

	output, err := ioutil.TempDir("", "ziputil")
	ne(err)

	defer os.RemoveAll(output)

	z, err := zip.NewReader(temp, size)
	ne(err)

	ne(UnzipDir(output, z, true))

	outPath := path.Join(output, "testfile")
	if !testDiffFile(t, outPath, p) {
		t.Error("zip loop back failed")
	}
}

func TestZipDir(t *testing.T) {
	ne := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}

	temp, err := tempfile.NewFile("", "ziputil")
	ne(err)
	defer temp.CleanUp()

	const p = "testdata/testdir"
	ne(ZipDir(p, temp))

	size, err := temp.Seek(0, io.SeekCurrent)
	ne(err)

	ne(temp.Reset())

	output, err := ioutil.TempDir("", "ziputil")
	ne(err)

	defer os.RemoveAll(output)

	z, err := zip.NewReader(temp, size)
	ne(err)

	ne(UnzipDir(output, z, true))

	for _, name := range []string{
		"bin-file", "private-file", "text-file",
	} {
		outPath := path.Join(output, name)
		target := path.Join(p, name)
		if !testDiffFile(t, outPath, target) {
			t.Errorf("zip loop back failed for file %q", name)
		}
	}
}
