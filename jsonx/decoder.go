package jsonx

import (
	"encoding/json"
	"io"

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

// Decode decodes a JSON value from the parser.
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
