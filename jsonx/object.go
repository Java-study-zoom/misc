package jsonx

import (
	"shanhu.io/smlvm/lexing"
)

type object struct {
	left  *lexing.Token
	right *lexing.Token

	entries []*objectEntry
}

type value struct {
	lead  *lexing.Token
	token *lexing.Token
}

type objectEntry struct {
	key   *lexing.Token
	value interface{}
}

type list struct {
	left    *lexing.Token
	right   *lexing.Token
	entires []interface{}
}
