package main

import (
	"fmt"
	"strings"
)

type (
	Builder struct {
		opt      BuilderOptions
		buf      strings.Builder
		styleBuf strings.Builder
	}

	BuilderOptions struct {
		prettyOutput bool
	}
)

func OutputHTML(nodes []LayoutNode, options BuilderOptions) string {
	builder := Builder{}

	for _, root := range nodes {
		builder.outputNodeHTML(root)
	}

	return builder.buf.String()
}

// TODO(nico): Indent the html correctly
func (builder *Builder) outputNodeHTML(node LayoutNode) {
	switch r := node.(type) {
	case *LayoutParagraph:
		builder.openHTMLTag("p", builder.outputStyleCSS(r.format))
		defer builder.closeHTMLTag("p")
		for _, child := range r.children {
			builder.outputNodeHTML(child)
		}

	case *LayoutText:
		builder.openHTMLTag("span", builder.outputStyleCSS(r.format))
		builder.buf.WriteString(r.value)
		builder.closeHTMLTag("span")
	}
}

func (builder *Builder) openHTMLTag(tag string, style string) {
	fmt.Fprintf(&builder.buf, "<%s %s>", tag, style)
	if builder.opt.prettyOutput {
		builder.buf.WriteByte('\n')
	}
}

func (builder *Builder) closeHTMLTag(tag string) {
	fmt.Fprintf(&builder.buf, "</%s>", tag)
	if builder.opt.prettyOutput {
		builder.buf.WriteByte('\n')
	}
}

func (builder *Builder) outputStyleCSS(format layoutFormat) string {
	builder.styleBuf.Reset()

	builder.styleBuf.WriteString("style=\"")
	for _, f := range format {
		if f == nil {
			continue
		}

		terminateStyle := true

		switch _f := f.(type) {
		case layoutFont:
			fmt.Fprintf(&builder.styleBuf, "font-family: %s", _f.name)
		case layoutTextStyle:
			for i := 0; i < int(layoutTextStyleMAX); i += 1 {
				var mask byte = 1 << i
				if (byte(_f)&mask)>>i == 1 {
					switch layoutTextStyleKind(i) {
					case layoutTextStyleItalic:
						fmt.Fprintf(&builder.styleBuf, "font-style: italic;")
					case layoutTextStyleStrike:
						fmt.Fprintf(&builder.styleBuf, "text-decoration-line: line-through;")
					}
				}
			}
			terminateStyle = false
		case layoutColor:
			fmt.Fprintf(&builder.styleBuf, "color: rgba(%d, %d, %d, %.1f)", _f.r, _f.g, _f.b, float64(_f.a)/255)
		case layoutFontSize:
			fmt.Fprintf(&builder.styleBuf, "font-size: %d", _f)
		case layoutFontWeight:
			fmt.Fprintf(&builder.styleBuf, "font-weight: %s", layoutFontWeightStr[_f])
		case layoutTextAlign:
			fmt.Fprintf(&builder.styleBuf, "text-align: %s", layoutTextAlignStr[_f])
		case layoutTextIndent:
			indentValue := ConvertUnits(_f.value, _f.unit, MeasuringUnitEm)
			if _f.firstLineOffset != 0 {
				firstLineIndentValue := ConvertUnits(_f.firstLineOffset, _f.unit, MeasuringUnitEm)
				fmt.Fprintf(&builder.styleBuf, "padding-left: %dem;", indentValue)
				fmt.Fprintf(&builder.styleBuf, "text-indent: %dem", firstLineIndentValue)
			} else {
				fmt.Fprintf(&builder.styleBuf, "text-indent: %dem", indentValue)
			}
		}

		if terminateStyle {
			builder.styleBuf.WriteByte(';')
		}
	}
	builder.styleBuf.WriteString("\"")

	return builder.styleBuf.String()
}

// func (builder *Builder) buildTextStyleCSS(f layoutTextStyle) {

// }
