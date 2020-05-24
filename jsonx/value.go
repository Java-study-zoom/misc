package jsonx

import (
	"shanhu.io/smlvm/lexing"
)

type value interface{}

type basic struct {
	lead  *lexing.Token
	token *lexing.Token
}

type object struct {
	left    *lexing.Token
	entries []*objectEntry
	right   *lexing.Token
}

type objectEntry struct {
	key   *lexing.Token
	colon *lexing.Token
	value value
	comma *lexing.Token
}

type list struct {
	left    *lexing.Token
	entries []*listEntry
	right   *lexing.Token
}

type listEntry struct {
	value value
	comma *lexing.Token
}
