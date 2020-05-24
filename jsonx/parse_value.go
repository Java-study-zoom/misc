package jsonx

func parseObjectEntries(p *parser) []*objectEntry {
	panic("todo")
}

func parseListEntries(p *parser) []value {
	panic("todo")
}

func parseValue(p *parser) value {
	if p.See(tokString) || p.See(tokInt) || p.See(tokFloat) {
		return &basic{token: p.Shift()}
	}
	if p.seeOp("{") {
		left := p.Shift()
		entries := parseObjectEntries(p)
		right := p.expectOp("}")
		return &object{
			left:    left,
			entries: entries,
			right:   right,
		}
	}
	if p.seeOp("[") {
		left := p.Shift()
		entries := parseListEntries(p)
		right := p.expectOp("]")
		return &list{
			left:    left,
			entries: entries,
			right:   right,
		}
	}

	t := p.Token()
	p.CodeErrorf(
		t.Pos, "jsonx.expectOperand",
		"expect an operand, got %s", typeStr(t.Type),
	)
	return nil
}
