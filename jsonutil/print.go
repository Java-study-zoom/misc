package jsonutil

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

// Fprint pretty prints a JSON data blob into a writer.
func Fprint(w io.Writer, v interface{}) error {
	bs, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	if _, err := w.Write(bs); err != nil {
		return err
	}
	_, err = fmt.Fprintln(w)
	return err
}

// Print pretty prints a JSON data blob into stdout.
func Print(v interface{}) {
	if err := Fprint(os.Stdout, v); err != nil {
		log.Println(err)
	}
}
