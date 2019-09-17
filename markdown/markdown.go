package markdown

import (
	"bufio"
	"bytes"
	"regexp"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

var langSmlrepo = regexp.MustCompile(`^language\-smlrepo$`)

func newPolicy() *bluemonday.Policy {
	p := bluemonday.UGCPolicy()
	p.AllowDataURIImages()
	p.AllowAttrs("class").Matching(langSmlrepo).OnElements("code")
	return p
}

// ToHTMLWithTitle parses the text that uses the first H1 as a title.
func ToHTMLWithTitle(text []byte) (string, []byte) {
	return toHTMLWithTitle(text, true)
}

// ToWildHTMLWithTitle parses the text taht uses the first H1 as a title.
// It does not sanitize the HTML.
func ToWildHTMLWithTitle(text []byte) (string, []byte) {
	return toHTMLWithTitle(text, false)
}

func toHTMLWithTitle(text []byte, sanitize bool) (string, []byte) {
	if len(text) == 0 {
		return "", nil
	}

	if text[0] != '#' {
		return "", ToHTML(text)
	}

	pos := bytes.IndexRune(text, '\n')
	if pos < 0 {
		pos = len(text)
	}

	title := parseTitle(string(text[:pos]))
	return title, toHTML(text[pos:], sanitize)
}

func toHTML(text []byte, sanitize bool) []byte {
	unsanitized := blackfriday.MarkdownCommon(text)
	if !sanitize {
		return unsanitized
	}
	sanitized := newPolicy().SanitizeBytes(unsanitized)
	if len(sanitized) == 0 {
		return nil
	}
	return sanitized
}

// ToHTML converts a markdown file to an HTML.
func ToHTML(text []byte) []byte {
	return toHTML(text, true)
}

func firstLine(text []byte) string {
	r := bytes.NewReader(text)
	s := bufio.NewScanner(r)
	if s.Scan() {
		return s.Text()
	}
	return ""
}

func parseTitle(line string) string {
	line = strings.TrimLeft(line, "#")
	line = strings.TrimSpace(line)
	return line
}

// ParseTitle parses the title of a markdown file.
func ParseTitle(text []byte) string {
	return parseTitle(firstLine(text))
}
