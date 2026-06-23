package ui

import "github.com/tunedmystic/rio/dom"

type BadgeVariant int

const (
	BadgeNeutral BadgeVariant = iota
	BadgeSuccess
	BadgeWarning
	BadgeDanger
)

type AlertVariant int

const (
	AlertInfo AlertVariant = iota
	AlertSuccess
	AlertWarning
	AlertError
)

// Badge renders a small status/tag pill.
func Badge(variant BadgeVariant, label string) dom.Node {
	return dom.Span(dom.Class(badgeClasses(variant)), dom.Text(label))
}

func badgeClasses(v BadgeVariant) string {
	base := "inline-flex items-center rounded-full px-2.5 py-0.5 text-[length:var(--font-size-sm)] font-medium"
	switch v {
	case BadgeSuccess:
		return Class(base, "bg-[var(--color-success)] text-white")
	case BadgeWarning:
		return Class(base, "bg-[var(--color-warning)] text-white")
	case BadgeDanger:
		return Class(base, "bg-[var(--color-danger)] text-white")
	default:
		return Class(base, "bg-[var(--color-border)] text-[var(--color-text)]")
	}
}

// Alert renders a callout message box.
func Alert(variant AlertVariant, content ...dom.Node) dom.Node {
	children := make([]dom.Node, 0, len(content)+2)
	children = append(children, dom.Class(alertClasses(variant)), dom.Role("alert"))
	children = append(children, content...)
	return dom.Div(children...)
}

func alertClasses(v AlertVariant) string {
	base := "rounded-md border-l-4 p-4 text-[var(--color-text)] bg-[var(--color-surface)]"
	switch v {
	case AlertSuccess:
		return Class(base, "border-[var(--color-success)]")
	case AlertWarning:
		return Class(base, "border-[var(--color-warning)]")
	case AlertError:
		return Class(base, "border-[var(--color-danger)]")
	default:
		return Class(base, "border-[var(--color-info)]")
	}
}
