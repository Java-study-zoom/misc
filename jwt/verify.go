package jwt

import (
	"strings"

	"shanhu.io/misc/errcode"
)

// Verifier verifies the token.
type Verifier interface {
	Verify(h *Header, data, sig []byte) error
}

// Token is a parsed JWT token.
type Token struct {
	Header    *Header
	ClaimSet  *ClaimSet
	Signature []byte
}

// DecodeAndVerify decodes and verifies a token.
func DecodeAndVerify(token string, v Verifier) (*Token, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errcode.InvalidArgf(
			"invalid token: %d parts", len(parts),
		)
	}

	h, c, sig := parts[0], parts[1], parts[2]
	header := new(Header)
	if err := decodeSegment(h, header); err != nil {
		return nil, errcode.InvalidArgf("decode header: %s", err)
	}

	payload := []byte(token[:len(h)+1+len(c)])
	sigBytes, err := decodeSegmentBytes(sig)
	if err != nil {
		return nil, errcode.InvalidArgf("decode signature: %s", err)
	}
	if err := v.Verify(header, payload, sigBytes); err != nil {
		return nil, errcode.Annotate(err, "verify signature")
	}

	// TODO(h8liu): Decode claim set.
	return &Token{
		Header:    header,
		Signature: sigBytes,
	}, nil
}
