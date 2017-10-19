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
