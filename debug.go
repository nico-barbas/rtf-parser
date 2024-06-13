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
	OpDebugger struct {
		ops     []Entity
		output  []debugInfo
		builder strings.Builder
	}

	debugInfo struct {
		indent    int
		op        Entity
		userIndex int
	}
)

func PrintToken(t Token) {
	fmt.Printf("[%s] %s (from %d to %d)\n", tokenKindStr[t.kind], t.text, t.start, t.end)
}

func DebugOpStream(ops []Entity) {
	d := OpDebugger{
		ops:     ops,
		builder: strings.Builder{},
	}

	d.buildInfoIndentation()

	for _, info := range d.output {
		for i := 0; i < info.indent; i += 1 {
			d.builder.WriteByte('\t')
		}

		d.buildDebugMessage(info)
		d.builder.WriteByte('\n')
	}

	fmt.Println(d.builder.String())
}

func (d *OpDebugger) buildInfoIndentation() {
	indent := 0

	for _, op := range d.ops {
		idx := len(d.output)
		d.output = append(d.output, debugInfo{
			op:     op,
			indent: indent,
		})

		switch e := op.(type) {
		case ControlGroup:
			if e.groupKind == ControlGroupKindBegin {
				indent += 1
			} else {
				d.output[idx].indent -= 1
				indent -= 1
			}

		case ColorTable:
			for i, clr := range e.colors {
				d.output = append(d.output, debugInfo{
					op:        clr,
					indent:    d.output[idx].indent + 1,
					userIndex: i,
				})
			}

		case FontTableEntry:

		default:
		}
	}
}

func (d *OpDebugger) buildDebugMessage(info debugInfo) {
	k := info.op.kind()
	fmt.Fprintf(&d.builder, "%s", entityKindStr[k])

	switch e := info.op.(type) {
	case ControlGroup:
		fmt.Fprintf(&d.builder, " %s", controlGroupKindStr[e.groupKind])

	case ControlWord:
		fmt.Fprintf(&d.builder, " %s", e.wordToken.text)

	case CharacterSet:
		fmt.Fprintf(&d.builder, " %s", characterSetKindStr[e.setKind])
		if e.setKind == CharacterSetANSICPG {
			fmt.Fprintf(&d.builder, " (code page: %d)", e.codePage)
		}

	case ColorTableEntry:
		fmt.Fprintf(&d.builder, "%d = (", info.userIndex)
		for _, channel := range e.channels {
			if channel == nil {
				continue
			}
			c := channel.(ColorComponent)
			fmt.Fprintf(&d.builder, "%s: %d; ", c.wordToken.text, c.value)
		}
		d.builder.WriteByte(')')

	case FontTableEntry:
		fmt.Fprintf(
			&d.builder,
			"(name: %s, index: %d, charset: %d, default fallback: %t)",
			e.fontNameToken.text,
			e.index,
			e.charset,
			e.defaultFallback,
		)

	case ColorComponent:
		fmt.Fprintf(&d.builder, "(channel: %s, value: %d)", e.wordToken.text, e.value)

	case TextFormat:
		if e.arg != -1 {
			fmt.Fprintf(&d.builder, " %s (arg: %d)", textFormatKindStr[e.formatKind], e.arg)
		} else {
			fmt.Fprintf(&d.builder, " %s", textFormatKindStr[e.formatKind])
		}

	case Text:
		fmt.Fprintf(&d.builder, ` (value: "`)
		for _, t := range e.tokens {
			d.builder.WriteString(t.text)
		}
		d.builder.WriteString(`")`)
	default:
	}
}
