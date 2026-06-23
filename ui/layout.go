package ui

import "github.com/tunedmystic/rio/dom"

type Gap int

const (
	GapSm Gap = iota
	GapMd
	GapLg
)

// Container is a max-width centered page wrapper.
func Container(children ...dom.Node) dom.Node {
	return dom.Div(withClass("max-w-7xl mx-auto px-4", children)...)
}

// Section is a vertical spacing band separating page regions.
func Section(children ...dom.Node) dom.Node {
	return dom.Section(withClass("py-12", children)...)
}

// Card is a rounded, bordered raised surface with padding.
func Card(children ...dom.Node) dom.Node {
	return dom.Div(withClass("bg-[var(--color-surface)] border border-[var(--color-border)] rounded-[var(--radius-base)] p-6 shadow-sm", children)...)
}

// Stack is a flex column with a configurable gap.
func Stack(gap Gap, children ...dom.Node) dom.Node {
	return dom.Div(withClass(stackClasses(gap), children)...)
}

func stackClasses(gap Gap) string {
	base := "flex flex-col"
	switch gap {
	case GapMd:
		return Class(base, "gap-4")
	case GapLg:
		return Class(base, "gap-8")
	default:
		return Class(base, "gap-2")
	}
}
