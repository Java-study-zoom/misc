package jwt

import (
	"crypto/hmac"
	"crypto/sha256"

	"shanhu.io/misc/errcode"
)

// HS256 implements the HS256 signing algorithm. It uses SHA256 hash and HMAC
// signing.
type HS256 struct {
	key    []byte
	header *Header
}

// NewHS256 creates a new HS256 signer using the given key and key ID.
func NewHS256(key []byte, kid string) *HS256 {
	return &HS256{
		key: key,
		header: &Header{
			Alg:   algHS256,
			Typ:   jwtType,
			KeyID: kid,
		},
	}
}

// Header returns the JWT header for this signer.
func (h *HS256) Header() *Header {
	cp := *h.header
	return &cp
}

func (h *HS256) mac(data []byte) []byte {
	hash := hmac.New(sha256.New, h.key)
	hash.Write(data)
	return hash.Sum(nil)
}

// Sign signs the HS256 signature.
func (h *HS256) Sign(_ *Header, data []byte) ([]byte, error) {
	return h.mac(data), nil
}

// Verify verifies the HS256 signature.
func (h *HS256) Verify(header *Header, data, sig []byte) error {
	if err := checkHeader(header, h.header); err != nil {
		return err
	}
	want := h.mac(data)
	if !hmac.Equal(want, sig) {
		return errcode.InvalidArgf("wrong signature")
	}
	return nil
}
