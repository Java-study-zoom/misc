package jsonutil

import (
	"encoding/json"
	"io/ioutil"
)

// ReadFile reads and unmarshals a JSON file.
func ReadFile(file string, obj interface{}) error {
	bs, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(bs, obj)
}
