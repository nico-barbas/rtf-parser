package main

import (
	"slices"
	"strings"
)

type (
	Layout struct {
		ops         []Entity
		previous    Entity
		current     Entity
		formatStack []int
		colorTable  []layoutColor
		builder     strings.Builder
	}

	layoutColor struct {
		r, g, b, a uint8
	}
)

func (layout *Layout) popOperation() Entity {
	op := layout.ops[len(layout.ops)-1]
	layout.ops = layout.ops[:len(layout.ops)-1]
	return op
}

func BuildLayoutHTML(ops []Entity) string {
	layout := Layout{ops: slices.Clone(ops)}

	for len(layout.ops) > 0 {
		layout.previous = layout.current
		layout.current = layout.popOperation()

		switch e := layout.current.(type) {
		case ColorTableEntry:
			layout.storeColor(e)
		case Text:
			layout.buildTextHTML(e)
		default:
		}
	}

	return layout.builder.String()
}

// TODO(nb): Absolutely disgusting way to stich it
func (layout *Layout) buildTextHTML(t Text) {
	layout.builder.WriteString("<span>")
	for _, tok := range t.tokens {
		layout.builder.WriteString(tok.text)
	}
	layout.builder.WriteString("</span>")
}

func (layout *Layout) storeColor(c ColorTableEntry) {
	clr := layoutColor{}

	clr.r = c.args[0].(ColorComponent).value
	clr.g = c.args[1].(ColorComponent).value
	clr.b = c.args[2].(ColorComponent).value

	if c.args[3] != nil {
		clr.b = c.args[3].(ColorComponent).value
	}

	layout.colorTable = append(layout.colorTable, clr)
}
