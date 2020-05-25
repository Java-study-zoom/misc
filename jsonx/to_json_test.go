package jsonx

import (
	"testing"
)

func TestToJSON(t *testing.T) {
	for _, test := range []struct {
		in, out string
	}{
		{`1234`, `1234`},
		{`true`, `true`},
		{`false`, `false`},
		{`{value:42}`, `{"value":42}`},
		{`{bool:false}`, `{"bool":false}`},
		{`{a:42,b:true}`, `{"a":42,"b":true}`},
		{`{a:42,}`, `{"a":42}`},
		// {"{a:42,\n}", `{"a":42}`},
	} {
		bs, errs := ToJSON([]byte(test.in))
		if errs != nil {
			t.Errorf("convert %q, got error %q", test.in, errs[0])
			continue
		}
		if got := string(bs); got != test.out {
			t.Errorf("convert %q, got %q, want %q", test.in, got, test.out)
		}
	}
}
