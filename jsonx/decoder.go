package jsonx

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"

	"shanhu.io/misc/errcode"
	"shanhu.io/smlvm/lexing"
)

// Decoder is a decoder that is capable of parsing a stream.
type Decoder struct {
	p *parser
}

// NewDecoder creates a new decoder that can parse a stream
// of jsonx objects.
func NewDecoder(r io.Reader) *Decoder {
	p, _ := newParser("", r)
	return &Decoder{p: p}
}

// More returns true if there is more stuff.
func (d *Decoder) More() bool {
	return !(d.p.See(lexing.EOF) || d.p.Token() == nil)
}

// Decode decodes a JSON value from the parser. When there is
// error on parsing JSONx, v is always unchanged.
func (d *Decoder) Decode(v interface{}) []*lexing.Error {
	value := parseValue(d.p)
	if errs := d.p.Errs(); errs != nil {
		return errs
	}
	if d.p.See(tokSemi) {
		d.p.Shift()
	}

	bs, errs := marshalValueLexing(value)
	if errs != nil {
		return errs
	}
	if err := json.Unmarshal(bs, v); err != nil {
		return lexing.SingleErr(err)
	}
	return nil
}

// Unmarshal unmarshals a value into a JSON object. When there is an error on
// parsing JSONx, v is always unchagned.
func Unmarshal(bs []byte, v interface{}) error {
	dec := NewDecoder(bytes.NewReader(bs))
	if errs := dec.Decode(v); errs != nil {
		return errs[0]
	}
	if dec.More() {
		return errcode.InvalidArgf("expect EOF, got more")
	}
	return nil
}

// ReadFile reads a file and unmarshals the content into the JSON object.
func ReadFile(file string, v interface{}) error {
	bs, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return Unmarshal(bs, v)
}

// ReadFileMaybeJSON reads a file that might be JSONx or JSON.
func ReadFileMaybeJSON(file string, v interface{}) error {
	bs, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	if err := Unmarshal(bs, v); err != nil {
		// JSONx fails to decode, maybe it is plain JSON?
		if json.Unmarshal(bs, v) == nil {
			return nil
		}
		return err
	}
	return nil
}
