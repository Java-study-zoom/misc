package jsonfile

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
)

// Write marshals a JSON object and writes it into a file.
func Write(file string, obj interface{}) error {
	bs, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, bs, 0644)
}

// WriteReadable marshals a JSON object with indents and writes it into
// a file.
func WriteReadable(f string, v interface{}) error {
	buf := new(bytes.Buffer)
	bs, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		return err
	}
	buf.Write(bs)
	buf.Write([]byte("\n"))

	return ioutil.WriteFile(f, bs, 0644)
}
