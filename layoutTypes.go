package main

type (
	LayoutNodeKind int

	LayoutNode interface {
		kind() LayoutNodeKind
	}

	LayoutText struct {
		value string
	}
)
