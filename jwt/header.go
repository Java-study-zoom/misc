package jwt

import (
	"encoding/json"
)

// Header is the JWT header.
type Header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
	Kid string `json:"kid,omitempty"` // Key ID.
}

func (h *Header) encode() (string, error) {
	return encodeSegment(h)
}

func decodeHeader(s string) (*Header, error) {
	bs, err := decodeSegmentBytes(s)
	if err != nil {
		return nil, err
	}
	h := new(Header)
	if err := json.Unmarshal(bs, h); err != nil {
		return nil, err
	}
	return h, nil
}
