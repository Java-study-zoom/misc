package hashutil

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
)

// Hash hashes a blob into a hex hash that is assumed to be unique in the
// entire universe.
func Hash(bs []byte) string {
	ret := sha256.Sum256(bs)
	return hex.EncodeToString(ret[:])
}

// HashStr hashes a string into a hash in hex.
func HashStr(s string) string {
	h := sha256.New()
	io.WriteString(h, s)
	return hex.EncodeToString(h.Sum(nil))
}
