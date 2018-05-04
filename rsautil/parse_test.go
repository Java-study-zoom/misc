package rsautil

import (
	"testing"

	"reflect"
)

func TestParseKey(t *testing.T) {
	privateKey, err := ParsePrivateKeyFile("testdata/test.pem")
	if err != nil {
		t.Fatal(err)
	}

	publicKey, err := ParsePublicKeyFile("testdata/test.pub")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(&privateKey.PublicKey, publicKey) {
		t.Error("public/private key pair not matching")
	}
}
