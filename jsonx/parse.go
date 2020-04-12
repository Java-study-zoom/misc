package jsonx

import (
	"io"

	"shanhu.io/smlvm/lexing"
)

type parser struct {
	f string
	x lexing.Tokener
	*lexing.Parser
}

func newParser(f string, r io.Reader) (*parser, *lexing.Recorder) {
	panic("todo")
}
