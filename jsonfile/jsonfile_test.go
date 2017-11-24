package jsonfile

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestJsonRead(t *testing.T) {
	filename := "testdata/rumpelstilzchen"
	obj := &struct{}{}
	err := Read(filename, obj)
	if err == nil {
		t.Fatalf(
			"Read %s: file not exist err expected, none found",
			filename,
		)
	}
	data := &struct {
		Number  int
		Boolean bool
		Text    string
	}{}

	filename = "jsonfile_test.go"
	if err = Read(filename, obj); err == nil {
		t.Fatalf(
			"Read %s: failed unmarshal err expected, %s found",
			filename,
			err,
		)
	}
	filename = "testdata/test_json"
	if err = Read(filename, data); err != nil {
		t.Fatalf(
			"Failed to unmarshal from ./testdata/%s: %s",
			filename,
			err,
		)
	}

	f, err := ioutil.TempFile("", "jsonfile-test")
	if err != nil {
		t.Fatal(err)
	}

	expectedData := &struct {
		Number  int
		Boolean bool
		Text    string
	}{
		Text: "be stronger",
	}

	filename = f.Name()

	if err := Write(filename, data); err != nil {
		t.Fatalf("Failed to Write %s: %v", filename, err)
	}

	if err := Read(filename, obj); err != nil {
		t.Fatalf("Failed to Read %s: %v", filename, err)
	}

	if !reflect.DeepEqual(data, expectedData) {
		t.Errorf("expect %v, got %v", expectedData, data)
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
		"    \"Text\": \"be stronger\"\n" +
		"}"
	if string(content) != expected {
		t.Errorf("expect WriteReadable as \n%s, get \n%s", expected, content)
	}

	f.Close()
	os.Remove(filename) // ignore error
}
