package jsonx

import (
	"shanhu.io/smlvm/lexing"
)

type value interface{}

type basic struct {
	lead  *lexing.Token
	token *lexing.Token
	value interface{}
}

type boolean struct {
	keyword *lexing.Token
}

type object struct {
	left    *lexing.Token
	entries []*objectEntry
	right   *lexing.Token
}

type objectKey struct {
	token *lexing.Token
	value interface{}
}

type objectEntry struct {
	key   *objectKey
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

type identList struct {
	entries []*lexing.Token
	dots []*lexing.Token
}

type trunk struct {
	values []value
	semi   *lexing.Token
}
