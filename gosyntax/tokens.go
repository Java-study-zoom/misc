package gosyntax

import (
	"go/scanner"
	"go/token"

	"shanhu.io/misc/syntax"
)

func tokType(t token.Token, lit string) string {
	switch t {
	case token.COMMENT:
		return "cm"
	case token.IDENT:
		if builtInFuncMap[lit] {
			return "bfunc"
		}
		if builtInTypeMap[lit] {
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

// Tokens breaks a Go language program in a token stream.
func Tokens(bs []byte) ([]*syntax.Token, error) {
	fset := token.NewFileSet()
	f := fset.AddFile("a.go", fset.Base(), len(bs))
	s := new(scanner.Scanner)
	var errs scanner.ErrorList
	s.Init(f, bs, errs.Add, scanner.ScanComments)

	var ret []*syntax.Token
	endPos := 0
	for {
		pos, t, lit := s.Scan()
		if t == token.SEMICOLON && lit != ";" {
			continue // this is actually white space
		}

		offset := f.Offset(pos)
		if offset > endPos {
			ret = append(ret, &syntax.Token{
				Type: "ws",
				Lit:  string(bs[endPos:offset]),
			})
		}

		ret = append(ret, &syntax.Token{
			Type: tokType(t, lit),
			Lit:  lit,
		})

		if t == token.EOF {
			break
		}
		endPos = offset + len(lit)
	}

	return ret, nil
}

// HTML renders a G language file into HTML file.
func HTML(bs []byte) (string, error) {
	toks, err := Tokens(bs)
	if err != nil {
		return "", err
	}
	return syntax.RenderHTML(toks), nil
}
