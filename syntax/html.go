package syntax

import (
	"bytes"
	"fmt"
	"html"
)

func runeHTML(r rune) string {
	switch r {
	case '\t':
		return "&nbsp;&nbsp;&nbsp;&nbsp;"
	case ' ':
		return "&nbsp;"
	case '\n':
		return "<br>\n"
	}
	return html.EscapeString(string(r))
}

func writeToken(buf *bytes.Buffer, tok *Token) {
	fmt.Fprintf(buf, `<span class="%s">`, tok.Type)
	for _, r := range tok.Lit {
		fmt.Fprint(buf, runeHTML(r))
	}
	fmt.Fprint(buf, "</span>")
}

// RenderHTML renders a token series into a HTML file.
func RenderHTML(toks []*Token) string {
	buf := new(bytes.Buffer)
	for _, t := range toks {
		writeToken(buf, t)
	}
	return buf.String()
}
