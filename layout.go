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
		colorTable  []LayoutColor
		builder     strings.Builder
	}

	LayoutColor struct {
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
