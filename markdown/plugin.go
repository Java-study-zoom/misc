package markdown

import (
	"bytes"
	"fmt"

	"golang.org/x/net/html"
)

// Compiler compiles a source string into a series of bytes.
type Compiler interface {
	Compile(src string) ([]byte, error)
}

// Compile goes through the given HTML and compiles the smallrepo code plugins
// using the given compiler.
func Compile(src []byte, c Compiler) ([]byte, error) {
	r := bytes.NewReader(src)
	doc, err := html.Parse(r)
	if err != nil {
		return nil, fmt.Errorf("html parse: %s", err)
	}

	w := new(bytes.Buffer)
	if err := html.Render(w, doc); err != nil {
		return nil, fmt.Errorf("html render: %s", err)
	}

	return w.Bytes(), nil
}
