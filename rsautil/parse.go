package rsautil

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

// ParsePrivateKey parses the PEM encoded RSA private key.
func ParsePrivateKey(bs []byte) (*rsa.PrivateKey, error) {
	b, _ := pem.Decode(bs)
	if b == nil {
		return nil, fmt.Errorf("pem decode failed")
	}

	if !x509.IsEncryptedPEMBlock(b) {
		return nil, fmt.Errorf("key is encrypted")
	}

	return x509.ParsePKCS1PrivateKey(b.Bytes)
}
