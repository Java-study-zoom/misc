package jsonfile

import (
	"encoding/json"
	"io/ioutil"
)

// Read reads and unmarshals a JSON file.
func Read(file string, obj interface{}) error {
	bs, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(bs, obj)
}
