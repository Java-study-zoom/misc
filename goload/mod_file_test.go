package goload

import (
	"testing"
)

func TestModulePath(t *testing.T) {
	for _, test := range []struct {
		content, mod string
	}{
		{`module shanhu.io/misc`, "shanhu.io/misc"},
		{"  module    shanhu.io/misc\t\t\t\n\nextra", "shanhu.io/misc"},
		{`module "shanhu.io/misc/v1"`, "shanhu.io/misc/v1"},
		{`module "shanhu.io/misc"`, "shanhu.io/misc"},
		{"// comment\nmodule x // tail\nnext line", "x"},
		{"module `x` // tail", "x"},
	} {
		got, err := modulePath([]byte(test.content))
		if err != nil {
			t.Errorf("modulePath(%q) got error: %s", test.content, err)
		} else if got != test.mod {
			t.Errorf(
				"modulePath(%q), want %q, got %q",
				test.content, test.mod, got,
			)
		}
	}
}
