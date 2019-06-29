package optline

import (
	"testing"
)

func TestParse(t *testing.T) {
	for _, test := range []struct {
		line string
		k, v string
	}{
		{`key: "value"`, "key", "value"},
		{"key: `value`\n\n", "key", "value"},
		{`key: "value 0"`, "key", "value 0"},
		{`k1: "value"`, "k1", "value"},
		{`_k: "value"`, "_k", "value"},
		{"k: `value`", "k", "value"},
	} {
		opt, err := Parse(test.line)
		if err != nil {
			t.Errorf("%q: unexpected error: %s", test.line, err)
			continue
		}
		if opt.Key != test.k {
			t.Errorf("%q: want key %q, got %q", test.line, test.k, opt.Key)
		}
		if opt.Value != test.v {
			t.Errorf("%q: want value %q, got %q", test.line, test.v, opt.Value)
		}
	}
}
