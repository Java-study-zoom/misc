package jsonx

import (
	"shanhu.io/smlvm/lexing"
)

type value interface{}

type object struct {
	left    *lexing.Token
	entries []*objectEntry
	right   *lexing.Token
}

type basic struct {
	lead  *lexing.Token
	token *lexing.Token
}

type objectEntry struct {
	key   *lexing.Token
	value value
}

type list struct {
	left    *lexing.Token
	entries []value
	right   *lexing.Token
}
