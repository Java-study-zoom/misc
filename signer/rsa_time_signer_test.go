package signer

import (
	"testing"
	"time"

	"crypto/rand"
	"crypto/rsa"
)

func TestRsaTimeSigner(t *testing.T) {
	size := 1024
	key, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		t.Fatal(err)
	}
	wrongKey, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		t.Fatal(err)
	}
	s := NewRSATimeSigner(&key.PublicKey, time.Second)
	b, err := RSASignTime(key)
	signedTime := time.Now()
	if err != nil {
		t.Fatal(err)
	}
	if s.Check(b) != nil && time.Since(signedTime) < time.Second {
		t.Errorf("signer should be valid")
	}
	time.Sleep(2 * time.Second)
	if s.Check(b) == nil {
		t.Errorf("signer should time out")
	}
	b, err = RSASignTime(wrongKey)
	if err != nil {
		t.Fatal(err)
	}
	if s.Check(b) == nil {
		t.Errorf("signer should not valid")
	}

}
