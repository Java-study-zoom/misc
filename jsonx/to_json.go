package jsonx

import (
	"bytes"

	"shanhu.io/smlvm/lexing"
)

// ToJSON converts a JSONX stream into a JSON stream.
func ToJSON(input []byte) ([]byte, []*lexing.Error) {
	r := bytes.NewReader(input)
	p, _ := newParser("", r)
	t := parseTrunk(p)
	if errs := p.Errs(); errs != nil {
		return nil, errs
	}

	buf := new(bytes.Buffer)
	if err := encodeTrunk(buf, t); err != nil {
		return nil, lexing.SingleErr(err)
	}
	return buf.Bytes(), nil
}
