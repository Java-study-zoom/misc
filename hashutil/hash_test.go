package hashutil

import (
	"testing"
)

func TestHash(t *testing.T) {
	m := make(map[string]bool)
	addHash := func(h string) {
		if m[h] {
			t.Fatalf("hash conflict: %s", h)
		}
		m[h] = true
	}

	addHash(Hash(nil))
	addHash(HashStr("a"))
	addHash(HashStr("A"))
	addHash(HashStr("A "))
	addHash(HashStr("Hello"))
}
