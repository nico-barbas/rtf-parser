package main

const (
	LayoutNodeInvalid LayoutNodeKind = iota
	LayoutNodeParagraph
	LayoutNodeText
)

type (
	LayoutNodeKind int

	LayoutNode interface {
		kind() LayoutNodeKind
		getFormat() layoutFormat
		getParent() LayoutNode
	}

	LayoutParagraph struct {
		format   layoutFormat
		parent   LayoutNode
		children []LayoutNode
	}

	LayoutText struct {
		format layoutFormat
		parent LayoutNode
		value  string
	}
)

func (p *LayoutParagraph) kind() LayoutNodeKind {
	return LayoutNodeParagraph
}

func (p *LayoutParagraph) getFormat() layoutFormat {
	return p.format
}

func (p *LayoutParagraph) getParent() LayoutNode {
	return p.parent
}

func (t *LayoutText) kind() LayoutNodeKind {
	return LayoutNodeText
}

func (t *LayoutText) getFormat() layoutFormat {
	return t.format
}

func (t *LayoutText) getParent() LayoutNode {
	return t.parent
}

const (
	layoutFormatColor layoutFormatKind = iota
	layoutFormatTextStyle
	layoutFormatFont
	layoutFormatFontSize
	layoutFormatFontWeight
	layoutFormatTextAlign
	layoutFormatTextIndent
	layoutFormatMAX
)

const (
	layoutTextStyleItalic layoutTextStyleKind = iota
	layoutTextStyleStrike
	layoutTextStyleMAX
)

const (
	layoutFontWeightBold layoutFontWeight = iota
)

var (
	layoutFontWeightStr = map[layoutFontWeight]string{
		layoutFontWeightBold: "bold",
	}
)

const (
	layoutTextAlignCenter layoutTextAlign = iota
	layoutTextAlignJustify
	layoutTextAlignRight
)

var (
	layoutTextAlignStr = map[layoutTextAlign]string{
		layoutTextAlignCenter:  "center",
		layoutTextAlignJustify: "justify",
		layoutTextAlignRight:   "right",
	}
)

type (
	layoutFormatKind int

	layoutFormat [layoutFormatMAX]layoutFormatOp

	layoutFormatOp interface {
		kind() layoutFormatKind
		concat(layoutFormatOp) layoutFormatOp
	}

	layoutFont struct {
		name string
	}

	layoutColor struct {
		r, g, b, a uint8
	}

	layoutTextStyleKind byte

	layoutTextStyle byte

	layoutFontSize int

	layoutFontWeight int

	layoutTextAlign int

	layoutTextIndent struct {
		dir             int
		unit            MeasuringUnit
		value           int
		firstLineOffset int
	}
)

func checkLayoutFormatOpConcat(op layoutFormatOp) (ok bool) {
	switch op.(type) {
	// case layoutFont:
	// 	ok = false
	// case layoutColor:
	// 	ok = false
	// case layoutFontSize:
	// 	ok = false
	// case layoutFontWeight:
	// 	ok = false
	// case layoutTextAlign:
	// 	ok = false
	case layoutTextStyle:
		ok = true
	case layoutTextIndent:
		ok = true
	default:
		ok = false
	}
	return ok
}

func (f layoutFont) kind() layoutFormatKind {
	return layoutFormatFont
}

func (f layoutFont) concat(other layoutFormatOp) layoutFormatOp {
	return f
}

func (c layoutColor) kind() layoutFormatKind {
	return layoutFormatColor
}

func (c layoutColor) concat(other layoutFormatOp) layoutFormatOp {
	return c
}

func (s layoutTextStyle) kind() layoutFormatKind {
	return layoutFormatTextStyle
}

func (s layoutTextStyle) concat(other layoutFormatOp) layoutFormatOp {
	return s | other.(layoutTextStyle)
}

func (f layoutFontSize) kind() layoutFormatKind {
	return layoutFormatFontSize
}

func (f layoutFontSize) concat(other layoutFormatOp) layoutFormatOp {
	return f
}

func (f layoutFontWeight) kind() layoutFormatKind {
	return layoutFormatFontWeight
}

func (f layoutFontWeight) concat(other layoutFormatOp) layoutFormatOp {
	return f
}

func (a layoutTextAlign) kind() layoutFormatKind {
	return layoutFormatTextAlign
}

func (a layoutTextAlign) concat(other layoutFormatOp) layoutFormatOp {
	return a
}

func (i layoutTextIndent) kind() layoutFormatKind {
	return layoutFormatTextIndent
}

func (i layoutTextIndent) concat(other layoutFormatOp) layoutFormatOp {
	o := other.(layoutTextIndent)
	result := layoutTextIndent{
		dir:             i.dir,
		unit:            i.unit,
		value:           i.value + o.value,
		firstLineOffset: i.firstLineOffset + o.firstLineOffset,
	}
	return result
}
