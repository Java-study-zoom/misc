package jsonx

func parseObjectEntries(p *parser) []*objectEntry {
	var entries []*objectEntry
	for !p.seeOp("}") {
		if !(p.See(tokIdent) || p.See(tokString)) {
			p.CodeErrorfHere("jsonx.expectObjectEntry", "expect object entry")
			break
		}

		k := p.Shift()
		colon := p.expectOp(":")
		v := parseValue(p)
		entry := &objectEntry{
			key:   k,
			colon: colon,
			value: v,
		}

		if p.seeOp(",") {
			entry.comma = p.Shift()
		} else if !p.seeOp("}") {
			p.expectOp(",")
		}
		entries = append(entries, entry)

		if p.InError() {
			break
		}
	}

	return entries
}

func parseListEntries(p *parser) []*listEntry {
	var entries []*listEntry
	for !p.seeOp("]") {
		v := parseValue(p)
		entry := &listEntry{value: v}
		if p.seeOp(",") {
			entry.comma = p.Shift()
		} else if !p.seeOp("]") {
			p.expectOp(",")
		}
		entries = append(entries, entry)
		if p.InError() {
			break
		}
	}

	return entries
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
