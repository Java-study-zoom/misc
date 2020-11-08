package jwt

import (
	"bytes"
	"io"

	"shanhu.io/misc/errcode"
)

// Signer signs the token, returns the signature and the header.
type Signer interface {
	Header() *Header
	Sign(h *Header, data []byte) ([]byte, error)
}

// SignAndEncode signs and encodes a claim set and signs it.
func SignAndEncode(c *ClaimSet, s Signer) (string, error) {
	h := s.Header()
	hb, err := h.encode()
	if err != nil {
		return "", errcode.Annotate(err, "encode header")
	}

	cb, err := c.encode()
	if err != nil {
		return "", errcode.Annotate(err, "encode claims")
	}
	buf := new(bytes.Buffer)
	io.WriteString(buf, hb)
	io.WriteString(buf, ".")
	io.WriteString(buf, cb)
	sig, err := s.Sign(h, buf.Bytes())
	if err != nil {
		return "", errcode.Annotate(err, "signing token")
	}
	io.WriteString(buf, ".")
	io.WriteString(buf, encodeSegmentBytes(sig))
	return buf.String(), nil
}
