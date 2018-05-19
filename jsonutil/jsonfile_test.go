package jsonutil

import (
	"testing"

	"io/ioutil"
	"os"
	"reflect"
)

type testSturct struct {
	Number  int
	Boolean bool
	Text    string
}

var testWriteData = &testSturct{
	Text: "be stronger",
}

const testWriteReadable = `{
    "Number": 0,
    "Boolean": false,
    "Text": "be stronger"
}`

func TestReadNotExist(t *testing.T) {
	const filename = "testdata/rumpelstilzchen"
	obj := &struct{}{}
	if err := ReadFile(filename, obj); err == nil {
		t.Errorf(
			"Read %s: want not-exist err, got nil",
			filename,
		)
	} else if !os.IsNotExist(err) {
		t.Errorf(
			"Read %s: want not-exist err, got %s",
			filename, err,
		)
	}
}

func TestReadNotJson(t *testing.T) {
	const filename = "testdata/invalid.json"
	obj := &struct{}{}
	if err := ReadFile(filename, obj); err == nil {
		t.Errorf(
			"Read %s: want unmarshal error, got %s",
			filename, err,
		)
	}
}

func TestRead(t *testing.T) {
	data := new(testSturct)
	const filename = "testdata/stronger.json"
	if err := ReadFile(filename, data); err != nil {
		t.Fatalf(
			"Read %q: got error: %s", filename, err,
		)
	}

	if !reflect.DeepEqual(data, testWriteData) {
		t.Errorf(
			"Read %s: want %v, got %v",
			filename, testWriteData, data,
		)
	}
}

func TestWrite(t *testing.T) {
	f, err := ioutil.TempFile("", "jsonfile-test")
	if err != nil {
		t.Fatal(err)
	}
	filename := f.Name()
	f.Close()
	defer os.Remove(filename)

	if err := WriteFile(filename, testWriteData); err != nil {
		t.Fatalf("Failed to Write %s: %s", filename, err)
	}
	dat := new(testSturct)
	if err := ReadFile(filename, dat); err != nil {
		t.Fatalf("Failed to Read %s: %s", filename, err)
	}

	if !reflect.DeepEqual(dat, testWriteData) {
		t.Errorf("expect %v, got %v", testWriteData, dat)
	}
}

func TestWriteReadable(t *testing.T) {
	f, err := ioutil.TempFile("", "jsonfile-test")
	if err != nil {
		t.Fatal(err)
	}
	filename := f.Name()
	f.Close()
	defer os.Remove(filename)

	if err := WriteFileReadable(filename, testWriteData); err != nil {
		t.Fatalf("Failed to WriteReadable %s: %s", filename, err)
	}

	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	got := string(bs)
	if got != testWriteReadable {
		t.Errorf("WriteReadable want %q, got %q", testWriteReadable, got)
	}
}
