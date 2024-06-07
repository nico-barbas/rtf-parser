package main

import (
	"slices"
	"strings"
)

type (
	Layout struct {
		ops               []Entity
		previous          Entity
		current           Entity
		formatStack       []layoutFormat
		formatStackFrames []int
		formatMask        layoutFormatMask
		colorTable        []layoutColor
		builder           strings.Builder
	}

	layoutFormat interface{}

	layoutFormatMask uint64

	layoutColor struct {
		r, g, b, a uint8
	}

	layoutFontIndex int

	layoutFontSize int
)

func BuildLayoutHTML(ops []Entity) string {
	layout := Layout{ops: slices.Clone(ops)}

	for len(layout.ops) > 0 {
		layout.previous = layout.current
		layout.current = layout.popOperation()

		switch e := layout.current.(type) {
		case ControlGroup:
			switch e.groupKind {
			case ControlGroupKindBegin:
				layout.pushFormatStackFrame()
			case ControlGroupKindEnd:
				layout.popFormatStackFrame()
			}
		case ColorTableEntry:
			layout.storeColor(e)
		case TextFormat:
			layout.applyFormat(e)
		case Text:
			layout.buildTextHTML(e)
		default:
		}
	}

	return layout.builder.String()
}

func (layout *Layout) popOperation() Entity {
	op := layout.ops[len(layout.ops)-1]
	layout.ops = layout.ops[:len(layout.ops)-1]
	return op
}

func (layout *Layout) pushFormatStackFrame() {
	layout.formatStackFrames = append(layout.formatStackFrames, len(layout.formatStack))
}

func (layout *Layout) popFormatStackFrame() {
	last := len(layout.formatStackFrames) - 1
	stackIdx := layout.formatStackFrames[last]
	layout.formatStackFrames = layout.formatStackFrames[:last]
	layout.formatStack = layout.formatStack[:stackIdx]
}

func (layout *Layout) clearFormatStack() {
	layout.formatStack = layout.formatStack[:0]
	layout.formatStackFrames = layout.formatStackFrames[:0]
}

func (layout *Layout) storeColor(c ColorTableEntry) {
	clr := layoutColor{}

	clr.r = c.args[0].(ColorComponent).value
	clr.g = c.args[1].(ColorComponent).value
	clr.b = c.args[2].(ColorComponent).value

	if len(c.args) == 4 {
		clr.b = c.args[3].(ColorComponent).value
	}

	layout.colorTable = append(layout.colorTable, clr)
}

func (layout *Layout) applyFormat(t TextFormat) {
	switch t.formatKind {
	case TextFormatColor:
	case TextFormatFontIndex:
	case TextFormatFontSize:
	}
}

// TODO(nb): Absolutely disgusting way to stich it
func (layout *Layout) buildTextHTML(t Text) {
	layout.builder.WriteString("<span>")
	for _, tok := range t.tokens {
		layout.builder.WriteString(tok.text)
	}
	layout.builder.WriteString("</span>")
}
