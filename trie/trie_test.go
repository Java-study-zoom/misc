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

func trieFindPath(t *Trie, p string) string {
	return t.Find(pathSplit(p))
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

	as(trieFindPath(tr, "a/c") == "a/c")
	as(trieFindPath(tr, "a/b/c") == "a/b/c")
	as(trieFindPath(tr, "a/b") == "a/b")
	as(trieFindPath(tr, "abc") == "abc")
	
	as(trieFindPath(tr, "def") == "")
	as(trieFindPath(tr, "a/c/d") == "")
	as(trieFindPath(tr, "a") == "")
	as(trieFindPath(tr, "a/b/c/d") == "")
}
