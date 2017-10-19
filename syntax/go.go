package syntax

import (
	"go/scanner"
	"go/token"
)

func goTokType(t token.Token, lit string) string {
	switch t {
	case token.COMMENT:
		return "cm"
	case token.IDENT:
		if _, found := goBuiltInFuncMap[lit]; found {
			return "bfunc"
		}
		if _, found := goBuiltInTypeMap[lit]; found {
			return "btype"
		}
		return "ident"
	case token.INT, token.FLOAT, token.IMAG:
		return "num"
	case token.CHAR, token.STRING:
		return "str"
	}
	if t.IsOperator() {
		return "op"
	}
	if t.IsKeyword() {
		return "kw"
	}
	return "na"
}

// RenderGo renders a Go language program in a token stream.
func RenderGo(bs []byte) ([]*Token, error) {
	fset := token.NewFileSet()
	f := fset.AddFile("a.go", fset.Base(), len(bs))
	s := new(scanner.Scanner)
	var errs scanner.ErrorList
	s.Init(f, bs, errs.Add, scanner.ScanComments)

	var ret []*Token
	endPos := 0
	for {
		pos, t, lit := s.Scan()
		if t == token.SEMICOLON && lit != ";" {
			continue // this is actually white space
		}

		offset := f.Offset(pos)
		if offset > endPos {
			ret = append(ret, &Token{
				Type: "ws",
				Lit:  string(bs[endPos:offset]),
			})
		}

		ret = append(ret, &Token{
			Type: goTokType(t, lit),
			Lit:  lit,
		})

		if t == token.EOF {
			break
		}
		endPos = offset + len(lit)
	}

	return ret, nil
}
