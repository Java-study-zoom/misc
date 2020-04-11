package jsonx

import (
	"shanhu.io/smlvm/lexing"
)

func lexOperator(x *lexing.Lexer, r rune) *lexing.Token {
	switch r {
	case '{', '}', '[', ']', ',', ':', '+', '-':
		/* do nothing */
	default:
		return nil
	}
	return x.MakeToken(tokOperator)
}

func lexJSONX(x *lexing.Lexer) *lexing.Token {
	r := x.Rune()
	if x.IsWhite(r) {
		panic("incorrect token start")
	}

	switch r {
	case '\n':
		x.Next()
		return x.MakeToken(tokEndl)
	case '"':
		return lexing.LexString(x, tokString, '"')
	case '`':
		return lexing.LexRawString(x, tokString)
	case '/':
		r2 := x.Rune()
		if r2 == '/' || r2 == '*' {
			return lexing.LexComment(x)
		}
	}

	if lexing.IsDigit(r) {
		return lexing.LexNumber(x, tokInt, tokFloat)
	}
	if lexing.IsIdentLetter(r) {
		return lexing.LexIdent(x, tokIdent)
	}

	x.Next()
	t := lexOperator(x, r)
	if t != nil {
		return t
	}

	x.CodeErrorf("jsonx.illegalChar", "illegal char %q", r)
	return x.MakeToken(lexing.Illegal)
}
