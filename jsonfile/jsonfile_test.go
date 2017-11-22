package jsonfile

import (
	"io/ioutil"
	"os"
<<<<<<< HEAD
	"reflect"
=======
>>>>>>> a25c218d5ef15982ce5d4a9d432319df4fe1374d
	"testing"
)

func TestJsonRead(t *testing.T) {
<<<<<<< HEAD
	filename := "testdata/rumpelstilzchen"
=======
	filename := "rumpelstilzchen"
>>>>>>> a25c218d5ef15982ce5d4a9d432319df4fe1374d
	obj := &struct{}{}
	err := Read(filename, obj)
	if err == nil {
		t.Fatalf(
			"Read %s: file not exist err expected, none found",
			filename,
		)
	}
<<<<<<< HEAD
	data := &struct {
		Number  int
		Boolean bool
		Text    string
	}{}
=======
>>>>>>> a25c218d5ef15982ce5d4a9d432319df4fe1374d

	filename = "jsonfile_test.go"
	if err = Read(filename, obj); err == nil {
		t.Fatalf(
<<<<<<< HEAD
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

=======
			"Read %s: failed unmarshal err expected, none found",
			filename,
		)
	}
>>>>>>> a25c218d5ef15982ce5d4a9d432319df4fe1374d
	f, err := ioutil.TempFile("", "jsonfile-test")
	if err != nil {
		t.Fatal(err)
	}
<<<<<<< HEAD

	expectedData := &struct {
		Number  int
		Boolean bool
		Text    string
	}{
		Text: "be stronger",
	}

	filename = f.Name()

=======
	filename = f.Name()
	data := &struct {
		Number  int
		Boolean bool
		Test    string
	}{
		Test: "be stronger",
	}

>>>>>>> a25c218d5ef15982ce5d4a9d432319df4fe1374d
	if err := Write(filename, data); err != nil {
		t.Fatalf("Failed to Write %s: %v", filename, err)
	}

	if err := Read(filename, obj); err != nil {
		t.Fatalf("Failed to Read %s: %v", filename, err)
	}

<<<<<<< HEAD
	if !reflect.DeepEqual(data, expectedData) {
		t.Errorf("expect %v, got %v", expectedData, data)
	}

=======
>>>>>>> a25c218d5ef15982ce5d4a9d432319df4fe1374d
	if err := WriteReadable(filename, data); err != nil {
		t.Fatalf("Failed to WriteReadable %s: %v", filename, err)
	}

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	expected := "{\n    \"Number\": 0,\n" +
		"    \"Boolean\": false,\n" +
<<<<<<< HEAD
		"    \"Text\": \"be stronger\"\n" +
		"}"
	if string(content) != expected {
		t.Errorf("expect WriteReadable as \n%s, get \n%s", expected, content)
=======
		"    \"Test\": \"be stronger\"\n" +
		"}"
	if string(content) != expected {
		t.Fatalf("expect WriteReadable as \n%s, get \n%s", expected, content)
>>>>>>> a25c218d5ef15982ce5d4a9d432319df4fe1374d
	}

	f.Close()
	os.Remove(filename) // ignore error
}
