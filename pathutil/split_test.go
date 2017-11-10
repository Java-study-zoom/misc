package pathutil

import (
	"testing"
)

func TestSplit(t *testing.T) {
	for _, p := range []string{
		"",
		"/x",
		"x//y",
		"/",
		"a/b/c/",
	} {
		parts, err := Split(p)
		if err == nil {
			t.Errorf(
				"split path %q got parts %v, want error",
				p, parts,
			)
		}
	}
}
