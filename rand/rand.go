package rand

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	mrand "math/rand"
	"sync"
	"time"
)

var randMutex sync.Mutex
var fallbackRand = mrand.New(
	mrand.NewSource(time.Now().UnixNano()),
)

// Bytes returns a byte slice of random bytes.
func Bytes(n int) []byte {
	ret := make([]byte, n)
	if _, err := rand.Read(ret); err == nil {
		return ret
	}

	randMutex.Lock()
	defer randMutex.Unlock()
	if _, err := fallbackRand.Read(ret); err != nil {
		panic(err)
	}

	return ret
}

// HexBytes returns the hex encoding of a random hex bytes
func HexBytes(n int) string {
	return hex.EncodeToString(Bytes(n))
}

// Letters returns a random ID of n random letters.
func Letters(n int) string {
	seed := int64(binary.LittleEndian.Uint64(Bytes(8)))
	src := mrand.NewSource(seed)
	r := mrand.New(src)
	var ret bytes.Buffer

	for i := 0; i < n; i++ {
		x := r.Int31n(52)
		if x < 26 {
			ret.WriteRune('a' + x)
		} else {
			ret.WriteRune('A' + x - 26)
		}
	}
	return ret.String()
}
