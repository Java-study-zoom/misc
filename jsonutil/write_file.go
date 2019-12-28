package jsonutil

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
)

// WriteFile marshals a JSON object and writes it into a file.
func WriteFile(file string, obj interface{}) error {
	bs, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, bs, 0644)
}

// WriteFileReadable marshals a JSON object with indents and writes it into a
// file.
func WriteFileReadable(f string, v interface{}) error {
	buf := new(bytes.Buffer)
	bs, err := json.MarshalIndent(v, "", formatIndent)
	if err != nil {
		return err
	}
	buf.Write(bs)
	buf.Write([]byte("\n"))

	return ioutil.WriteFile(f, bs, 0644)
}
