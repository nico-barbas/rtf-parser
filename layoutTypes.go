package main

const (
	LayoutNodeInvalid LayoutNodeKind = iota
	LayoutNodeParagraph
	LayoutNodeText
	LayoutNodeBody
)

type (
	LayoutNodeKind int

	LayoutNode interface {
		kind() LayoutNodeKind
		getFormat() layoutFormat
	}

	LayoutParagraph struct {
		format   layoutFormat
		children []LayoutNode
	}

	LayoutText struct {
		format layoutFormat
		value  string
	}
)

func (p *LayoutParagraph) kind() LayoutNodeKind {
	return LayoutNodeParagraph
}

func (p *LayoutParagraph) getFormat() layoutFormat {
	return p.format
}

func (t *LayoutText) kind() LayoutNodeKind {
	return LayoutNodeText
}

func (t *LayoutText) getFormat() layoutFormat {
	return t.format
}

const (
	layoutFormatColor layoutFormatKind = iota
	layoutFormatFont
	layoutFormatFontSize
	layoutFormatMAX
)

type (
	layoutFormatKind int

	layoutFormat [layoutFormatMAX]layoutFormatOp

	layoutFormatOp interface {
		kind() layoutFormatKind
	}

	layoutColor struct {
		r, g, b, a uint8
	}

	layoutFontIndex int

	layoutFontSize int
)

func (c layoutColor) kind() layoutFormatKind {
	return layoutFormatColor
}

func (f layoutFontSize) kind() layoutFormatKind {
	return layoutFormatFontSize
}
