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
		formatStack       []layoutFormatOp
		formatStackFrames []int
		fontTable         map[int]layoutFont
		colorTable        []layoutColor

		// Output
		roots       []LayoutNode
		currentNode *LayoutParagraph
		builder     strings.Builder
	}
)

func BuildLayout(ops []Entity) []LayoutNode {
	layout := Layout{ops: slices.Clone(ops), fontTable: map[int]layoutFont{}}

	for _, op := range layout.ops {
		layout.previous = layout.current
		layout.current = op
		// layout.current = layout.popOperation()

		switch e := layout.current.(type) {
		case ControlGroup:
			switch e.groupKind {
			case ControlGroupKindBegin:
				layout.pushFormatStackFrame()
			case ControlGroupKindEnd:
				layout.popFormatStackFrame()
			}
		case FontTableEntry:
			layout.storeFont(e)
		case ColorTable:
			for _, clr := range e.colors {
				layout.storeColor(clr.(ColorTableEntry))
			}
		case TextFormat:
			layout.processFormat(e)
		case Text:
			text := layout.buildText(e)
			if layout.currentNode != nil {
				layout.currentNode.children = append(layout.currentNode.children, text)
			}
		default:
		}
	}

	return layout.roots
}

func (layout *Layout) pushFormat(format layoutFormatOp) {
	layout.formatStack = append(layout.formatStack, format)
}

func (layout *Layout) pushFormatStackFrame() {
	layout.formatStackFrames = append(layout.formatStackFrames, len(layout.formatStack))
}

func (layout *Layout) popFormatStackFrame() {
	if len(layout.formatStackFrames) == 0 {
		return
	}

	last := len(layout.formatStackFrames) - 1
	stackIdx := layout.formatStackFrames[last]
	layout.formatStackFrames = layout.formatStackFrames[:last]
	layout.formatStack = layout.formatStack[:stackIdx]
}

func (layout *Layout) clearFormatStack() {
	layout.formatStack = layout.formatStack[:0]
	layout.formatStackFrames = layout.formatStackFrames[:0]
}

func (layout *Layout) storeFont(f FontTableEntry) {
	layout.fontTable[f.index] = layoutFont{
		name: f.fontNameToken.text,
	}
}

func (layout *Layout) storeColor(c ColorTableEntry) {
	clr := layoutColor{}

	clr.r = c.channels[0].(ColorComponent).value
	clr.g = c.channels[1].(ColorComponent).value
	clr.b = c.channels[2].(ColorComponent).value

	if c.channels[3] != nil {
		clr.a = c.channels[3].(ColorComponent).value
	} else {
		clr.a = 255
	}

	layout.colorTable = append(layout.colorTable, clr)
}

func (layout *Layout) processFormat(t TextFormat) {
	switch t.formatKind {
	case TextFormatColor:
		layout.pushFormat(layout.colorTable[t.arg-1])
	case TextFormatFontIndex:
		layout.pushFormat(layout.fontTable[t.arg])
	case TextFormatFontSize:
		layout.pushFormat(layoutFontSize(t.arg))

	case TextFormatParagraphClear:
		layout.clearFormatStack()
		layout.pushFormatStackFrame()
		p := &LayoutParagraph{}

		if layout.currentNode != nil {
			layout.currentNode.children = append(layout.currentNode.children, p)
		} else {
			layout.roots = append(layout.roots, p)
		}
		layout.currentNode = p

	case TextFormatParagraphEnd:
		layout.currentNode.format = layout.buildFormat()
		layout.currentNode = nil
	}
}

func (layout *Layout) buildFormat() layoutFormat {
	format := layoutFormat{}

	// Walk the format stack backward and skip any format that is already set in the bitmask
	for i := len(layout.formatStack) - 1; i >= 0; i -= 1 {
		f := layout.formatStack[i]
		if format[f.kind()] != nil {
			continue
		}

		format[f.kind()] = f
	}

	return format
}

func (layout *Layout) buildText(t Text) *LayoutText {
	layout.builder.Reset()

	for _, token := range t.tokens {
		layout.builder.WriteString(token.text)
	}

	return &LayoutText{
		value: layout.builder.String(),
	}
}
