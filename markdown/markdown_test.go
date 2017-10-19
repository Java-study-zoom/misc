package markdown

import (
	"testing"
)

func TestToHTML(t *testing.T) {
	for _, d := range []struct {
		in, out string
	}{
		{"", ""},
	} {
		out := ToHTML([]byte(d.in))
		if string(out) != d.out {
			t.Errorf("with title for %q", d.in)
			t.Logf("got output %q", string(out))
			t.Logf("want output %q", string(d.out))
		}
	}

	for _, d := range []struct {
		in, out, title string
	}{
		{"", "", ""},
	} {
		title, out := ToHTMLWithTitle([]byte(d.in))
		if string(out) != d.out || title != d.title {
			t.Errorf("with title for %q", d.in)
			t.Logf("got title %q and output %q", title, string(out))
			t.Logf("want title %q and output %q", d.title, string(d.out))
		}
	}
}
