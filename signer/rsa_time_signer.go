package signer

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"time"
)

// SignedRSABlock is a signed RSA block.
type SignedRSABlock struct {
	Data []byte
	Hash []byte
	Sig  []byte
}

// RSATimeSigner signes the current time, or checks if a signed time
// is within a time window of the current time reading.
type RSATimeSigner struct {
	k      *rsa.PublicKey
	window time.Duration

	TimeFunc func() time.Time
}

// NewRSATimeSigner creates a new time signer that uses an RSA key.
func NewRSATimeSigner(k *rsa.PublicKey, w time.Duration) *RSATimeSigner {
	if w < 0 {
		w = -w
	}

	return &RSATimeSigner{
		k:      k,
		window: w,
	}
}

func rsaSignTime(k *rsa.PrivateKey, t time.Time) (*SignedRSABlock, error) {
	buf := make([]byte, timestampLen)
	binary.LittleEndian.PutUint64(buf, uint64(t.UnixNano()))
	hash := sha256.Sum256(buf)
	sig, err := rsa.SignPKCS1v15(rand.Reader, k, crypto.SHA256, hash[:])
	if err != nil {
		return nil, err
	}

	return &SignedRSABlock{
		Data: buf,
		Hash: hash[:],
		Sig:  sig,
	}, nil
}

// RSASignTime signes the current time with the given RSA key.
func RSASignTime(k *rsa.PrivateKey) (*SignedRSABlock, error) {
	return rsaSignTime(k, time.Now())
}

// Check checks if the timestamp is with in the time window.
func (s *RSATimeSigner) Check(b *SignedRSABlock) error {
	t := time.Unix(0, int64(binary.LittleEndian.Uint64(b.Data)))
	timeNow := now(s.TimeFunc)
	if !inWindow(t, timeNow, s.window) {
		return fmt.Errorf("time out of window")
	}

	hash := sha256.Sum256(b.Data)
	if !bytes.Equal(hash[:], b.Hash) {
		return fmt.Errorf("hash incorrect")
	}
	return rsa.VerifyPKCS1v15(s.k, crypto.SHA256, b.Hash, b.Sig)
}
