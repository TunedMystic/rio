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
	base := "inline-flex items-center rounded-full px-2.5 py-0.5 text-[length:var(--font-size-sm)] font-medium ring-1 ring-inset"
	switch v {
	case BadgeSuccess:
		return Class(base, "bg-[var(--color-success)]/12 text-[var(--color-success)] ring-[var(--color-success)]/25")
	case BadgeWarning:
		return Class(base, "bg-[var(--color-warning)]/12 text-[var(--color-warning)] ring-[var(--color-warning)]/25")
	case BadgeDanger:
		return Class(base, "bg-[var(--color-danger)]/12 text-[var(--color-danger)] ring-[var(--color-danger)]/25")
	default:
		return Class(base, "bg-[var(--color-text)]/8 text-[var(--color-text)] ring-[var(--color-text)]/15")
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
	base := "rounded-[var(--radius-base)] border-l-4 p-4 text-[var(--color-text)]"
	switch v {
	case AlertSuccess:
		return Class(base, "border-[var(--color-success)] bg-[var(--color-success)]/8")
	case AlertWarning:
		return Class(base, "border-[var(--color-warning)] bg-[var(--color-warning)]/8")
	case AlertError:
		return Class(base, "border-[var(--color-danger)] bg-[var(--color-danger)]/8")
	default:
		return Class(base, "border-[var(--color-info)] bg-[var(--color-info)]/8")
	}
}
