package ui

import "github.com/tunedmystic/rio/dom"

type ButtonVariant int

const (
	ButtonPrimary ButtonVariant = iota
	ButtonSecondary
	ButtonDanger
	ButtonGhost
)

// Button renders a styled <button>. Pass extra attributes (id, hx-*, aria-*)
// via attrs; do not pass a class attribute — Button owns the class.
func Button(variant ButtonVariant, label string, attrs ...dom.Node) dom.Node {
	children := make([]dom.Node, 0, len(attrs)+3)
	children = append(children, dom.Class(buttonClasses(variant)), dom.Type("button"))
	children = append(children, attrs...)
	children = append(children, dom.Text(label))
	return dom.Button(children...)
}

// ButtonLink renders an <a> styled identically to Button, for CTAs that are
// navigation. Pass extra attributes via attrs; do not pass a class.
func ButtonLink(variant ButtonVariant, href, label string, attrs ...dom.Node) dom.Node {
	children := make([]dom.Node, 0, len(attrs)+3)
	children = append(children, dom.Class(buttonClasses(variant)), dom.Href(href))
	children = append(children, attrs...)
	children = append(children, dom.Text(label))
	return dom.A(children...)
}

func buttonClasses(v ButtonVariant) string {
	base := "inline-flex items-center justify-center rounded-md px-4 py-2 font-medium transition-colors"
	switch v {
	case ButtonPrimary:
		return Class(base, "bg-[var(--color-primary)] text-[var(--color-on-primary)] hover:opacity-90")
	case ButtonSecondary:
		return Class(base, "bg-[var(--color-secondary)] text-[var(--color-on-secondary)] hover:opacity-90")
	case ButtonDanger:
		return Class(base, "bg-[var(--color-danger)] text-white hover:opacity-90")
	case ButtonGhost:
		return Class(base, "bg-transparent text-[var(--color-text)] hover:bg-[var(--color-border)]")
	default:
		return Class(base, "bg-[var(--color-primary)] text-[var(--color-on-primary)] hover:opacity-90")
	}
}
