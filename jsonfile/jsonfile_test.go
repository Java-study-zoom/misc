package jsonfile

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestJsonRead(t *testing.T) {
	filename := "rumpelstilzchen"
	obj := &struct{}{}
	err := Read(filename, obj)
	if err == nil {
		t.Fatalf(
			"Read %s: file not exist err expected, none found",
			filename,
		)
	}

	filename = "jsonfile_test.go"
	if err = Read(filename, obj); err == nil {
		t.Fatalf(
			"Read %s: failed unmarshal err expected, none found",
			filename,
		)
	}
	f, err := ioutil.TempFile("", "jsonfile-test")
	if err != nil {
		t.Fatal(err)
	}
	filename = f.Name()
	data := &struct {
		Number  int
		Boolean bool
		Test    string
	}{
		Test: "be stronger",
	}

	if err := Write(filename, data); err != nil {
		t.Fatalf("Failed to Write %s: %v", filename, err)
	}

	if err := Read(filename, obj); err != nil {
		t.Fatalf("Failed to Read %s: %v", filename, err)
	}

	if err := WriteReadable(filename, data); err != nil {
		t.Fatalf("Failed to WriteReadable %s: %v", filename, err)
	}

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	expected := "{\n    \"Number\": 0,\n" +
		"    \"Boolean\": false,\n" +
		"    \"Test\": \"be stronger\"\n" +
		"}"
	if string(content) != expected {
		t.Fatalf("expect WriteReadable as \n%s, get \n%s", expected, content)
	}

	f.Close()
	os.Remove(filename) // ignore error
}
