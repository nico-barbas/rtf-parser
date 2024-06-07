package main

import (
	"fmt"
	"strings"
)

var (
	tokenKindStr = map[TokenKind]string{
		TokenInvalid:      "TokenInvalid",
		TokenNewline:      "TokenNewline",
		TokenEOF:          "TokenEOF",
		TokenOpenBracket:  "TokenOpenBracket",
		TokenCloseBracket: "TokenCloseBracket",
		TokenBackslash:    "TokenBackslash",
		TokenString:       "TokenString",
		TokenNumber:       "TokenNumber",
	}
)

type (
	Debugger struct {
		ops           []Entity
		indent        int
		builder       strings.Builder
		entityBuilder strings.Builder
	}
)

func PrintToken(t Token) {
	fmt.Printf("[%s] %s (from %d to %d)\n", tokenKindStr[t.kind], t.text, t.start, t.end)
}

func PrettyPrintEntities(ops []Entity) {
	d := Debugger{
		ops:     ops,
		builder: strings.Builder{},
	}

	for _, op := range d.ops {
		d.printOperation(op)
	}

	fmt.Println(d.builder.String())
}

func (d *Debugger) printOperation(op Entity) {
	nextIndent := d.indent

	k := op.kind()

	fmt.Fprintf(&d.entityBuilder, "%s", entityKindStr[k])

	switch e := op.(type) {
	case ControlGroup:
		fmt.Fprintf(&d.entityBuilder, " %s", controlGroupKindStr[e.groupKind])
		if e.groupKind == ControlGroupKindBegin {
			nextIndent += 2
		} else {
			nextIndent -= 2
			d.indent = nextIndent
		}

	case ControlWord:
		fmt.Fprintf(&d.entityBuilder, " %s", e.wordToken.text)

	case CharacterSet:
		fmt.Fprintf(&d.entityBuilder, " %s", characterSetKindStr[e.setKind])
		if e.setKind == CharacterSetANSICPG {
			fmt.Fprintf(&d.entityBuilder, " (code page: %d)", e.codePage)
		}

	case ColorTableEntry:
		for _, arg := range e.args {
			d.printOperation(arg)
		}

	case ColorComponent:
		fmt.Fprintf(&d.entityBuilder, "(channel: %s, value: %d)", e.wordToken.text, e.value)

	case TextFormat:
		fmt.Fprintf(&d.entityBuilder, " %s (arg: %d)", textFormatKindStr[e.formatKind], e.arg)

	case Text:
		fmt.Fprintf(&d.entityBuilder, ` (value: "`)
		for _, t := range e.tokens {
			d.entityBuilder.WriteString(t.text)
		}
		d.entityBuilder.WriteString(`")`)
	default:
	}

	for j := 0; j < d.indent; j += 1 {
		d.builder.WriteByte(' ')
	}
	d.builder.WriteString(d.entityBuilder.String())
	d.builder.WriteByte('\n')
	d.entityBuilder.Reset()
	d.indent = nextIndent
}
