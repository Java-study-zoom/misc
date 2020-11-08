package jwt

import (
	"testing"

	"time"

	"shanhu.io/misc/rand"
)

func TestHS256(t *testing.T) {
	key := rand.Bytes(32) // 256 bits
	h := NewHS256(key, "")
	now := time.Now()
	c := &ClaimSet{
		Iss: "shanhu.io",
		Aud: "nextcloud",
		Exp: now.Unix(),
		Iat: now.Add(time.Hour).Unix(),
		Sub: "h8liu",
	}

	tokStr, err := EncodeAndSign(c, h)
	if err != nil {
		t.Fatal("encode: ", err)
	}
	t.Log(tokStr)

	tok, err := DecodeAndVerify(tokStr, h)
	if err != nil {
		t.Fatal("decode: ", err)
	}

	if got, want := tok.ClaimSet.Iss, c.Iss; got != want {
		t.Errorf("got issuer %q, want %q", got, want)
	}
}
