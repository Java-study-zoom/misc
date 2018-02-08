package trie

import (
	"testing"

	"strings"
)

func pathSplit(s string) []string {
	return strings.Split(s, "/")
}

func trieAddPath(t *Trie, p string) bool {
	return t.Add(pathSplit(p), p)
}

func trieFindPath(t *Trie, p string) (string, string) {
	route, v := t.Find(pathSplit(p))
	if v == "" {
		return "", ""
	}
	return strings.Join(route, "/"), v
}

func TestTrie(t *testing.T) {
	tr := New()
	as := func(cond bool) {
		if !cond {
			t.Error("assertion failed")
		}
	}

	as(trieAddPath(tr, "a/b/c"))
	as(trieAddPath(tr, "a/b"))
	as(trieAddPath(tr, "abc"))
	as(trieAddPath(tr, "a/c"))

	if trieAddPath(tr, "a/c") {
		t.Error("should fail to add duplicate path")
	}

	for _, test := range []struct {
		input string
		output string
	} {
		{ "a/c", "a/c"},
		{ "a/b/c", "a/b/c"},
		{ "a/b", "a/b"},
		{"abc", "abc"},
		{ "a/c/d", "a/c"},
		{"a/b/c/d", "a/b/c"},
	} {
		r, v := trieFindPath(tr, test.input)
		if r != v {
			t.Errorf("find %q, got %q != %q", test.input, r, v)
		}
		if r != test.output {
			t.Errorf("find %q, want %q, got %q", test.input, test.output, r)
		}
	}

	for _, p := range []string {
		"def",
		"a",
	} {
		r, v := trieFindPath(tr, p)
		if r != "" || v != "" {
			t.Errorf("find %q, want not found, got %q, %q", p, r, v)
		}
	}
}
