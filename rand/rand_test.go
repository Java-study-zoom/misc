package rand

import (
	"testing"

	"bytes"
	"strings"
)

func TestBytes(t *testing.T) {
	bs1 := Bytes(8)
	bs2 := Bytes(8)

	if bytes.Equal(bs1, bs2) {
		t.Errorf("not so random: %v == %v", bs1, bs2)
	}
}

func TestLowerLetters(t *testing.T) {
	s1 := LowerLetters(16)
	s2 := LowerLetters(16)
	if s1 == s2 {
		t.Errorf("not so random: %q == %q", s1, s2)
	}
	if strings.ToLower(s1) != s1 {
		t.Errorf("contains non-lower case: %q", s1)
	}
}

func TestLetters(t *testing.T) {
	s1 := Letters(16)
	s2 := Letters(16)
	if s1 == s2 {
		t.Errorf("not so random: %q == %q", s1, s2)
	}
}
