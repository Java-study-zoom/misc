package strutil

import (
	"reflect"
	"testing"
)

func TestDefault(t *testing.T) {
	want := "some string"
	if got := Default("", want); got != want {
		t.Errorf("Default test failed, expect %q, got %q", want, got)
	}
	if got := Default(want, "def"); got != want {
		t.Errorf("Default test failed, expect %q, got %q", want, got)
	}
}

func TestMakeSet(t *testing.T) {
	for _, test := range []struct {
		list []string
		want map[string]bool
	}{
		{nil, make(map[string]bool)},
		{make([]string, 0), make(map[string]bool)},
		{[]string{"a", "B"}, map[string]bool{"a": true, "B": true}},
	} {
		got := MakeSet(test.list)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf(
				"moment string for %s: got %v, want %v",
				test.list, got, test.want,
			)
		}
	}
}

func TestCountLines(t *testing.T) {
	for _, test := range []struct {
		bs   []byte
		want int
	}{
		{nil, 0},
		{make([]byte, 0), 0},
		{[]byte("abcd"), 1},
		{[]byte("\n"), 1},
		{[]byte(" \n"), 1},
		{[]byte("\n\n"), 2},
		{[]byte("\n\nabc"), 3},
		{[]byte(""), 0},
	} {
		got := CountLines(test.bs)
		if test.want != got {
			t.Errorf(
				"const lines for %s: got %d, want %d",
				string(test.bs), got, test.want,
			)
		}
	}
}
