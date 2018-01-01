package osutil

import (
	"testing"

	"io/ioutil"
	"os"
	"path"
)

func TestIsDir(t *testing.T) {
	ne := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}

	d, err := ioutil.TempDir("", "osutil")
	ne(err)
	defer os.RemoveAll(d)

	ok, err := IsDir(d)
	ne(err)
	if !ok {
		t.Errorf("IsDir(%q) should return true", d)
	}

	f := path.Join(d, "post")
	ne(ioutil.WriteFile(f, []byte("post"), 0600))

	ok, err = IsDir(f)
	ne(err)
	if ok {
		t.Errorf("IsDir(%q) should return false", f)
	}

	ok, err = IsDir(path.Join(d, "ghost"))
	ne(err)
	if ok {
		t.Errorf("IsDir(%q) should return false", f)
	}
}
