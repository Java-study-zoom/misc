package jsonx

import (
	"testing"

	"reflect"
	"strings"
)

func TestDecoder(t *testing.T) {
	input := strings.NewReader(`"a""b";"c"`)

	dec := NewDecoder(input)
	var got []string
	for dec.More() {
		var s string
		if err := dec.Decode(&s); err != nil {
			t.Fatal(err)
		}
		got = append(got, s)
	}

	want := []string{"a", "b", "c"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestUnmarshal(t *testing.T) {
	var v int
	if err := Unmarshal([]byte("1234"), &v); err != nil {
		t.Fatal(err)
	}
	if v != 1234 {
		t.Errorf("got %d, want 1234", v)
	}
}
