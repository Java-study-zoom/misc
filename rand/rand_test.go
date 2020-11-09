package rand

import (
	"testing"

	"reflect"
)

func TestBytes(t *testing.T) {
	bs1 := Bytes(8)
	bs2 := Bytes(8)

	if reflect.DeepEqual(bs1, bs2) {
		t.Errorf("not so random: %v == %v", bs1, bs2)
	}
}

func TestLowerLetters(t *testing.T) {
	s1 := LowerLetters(16)
	s2 := LowerLetters(16)
	if s1 == s2 {
		t.Errorf("not so random: %q == %q", s1, s2)
	}
}
