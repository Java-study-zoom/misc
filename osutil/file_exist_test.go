package osutil

import (
	"testing"

	"io/ioutil"
	"os"
	"path"
)

func TestExist(t *testing.T) {
	ne := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}

	d, err := ioutil.TempDir("", "osutil")
	ne(err)
	defer os.RemoveAll(d)

	f := path.Join(d, "post")
	ne(ioutil.WriteFile(f, []byte("post"), 0600))

	ok, err := Exist(f)
	ne(err)
	if !ok {
		t.Errorf("file %q should exist", f)
	}

	ghost := path.Join(d, "ghost")
	ok, err = Exist(ghost)
	ne(err)
	if ok {
		t.Errorf("file %q should not exist", f)
	}
}
