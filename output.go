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
	fmt.Fprintf(&builder.buf, "<%s>", tag)
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

		switch _f := f.(type) {
		case layoutColor:
			fmt.Fprintf(&builder.styleBuf, "color: rgba(%d, %d, %d, %.1f);", _f.r, _f.g, _f.b, float64(_f.a)/255)
		case layoutFontSize:
			fmt.Fprintf(&builder.styleBuf, "font-size: %d;", _f)
		}
	}
	builder.styleBuf.WriteString("\"")

	return builder.styleBuf.String()
}
