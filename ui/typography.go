package ui

import "github.com/tunedmystic/rio/dom"

type HeadingLevel int

const (
	H1 HeadingLevel = iota + 1
	H2
	H3
	H4
)

type TextTone int

const (
	TextDefault TextTone = iota
	TextMuted
)

// Heading renders an h1–h4 sized from the font-size scale.
func Heading(level HeadingLevel, text string, attrs ...dom.Node) dom.Node {
	children := make([]dom.Node, 0, len(attrs)+2)
	children = append(children, dom.Class(headingClasses(level)))
	children = append(children, attrs...)
	children = append(children, dom.Text(text))
	switch level {
	case H2:
		return dom.H2(children...)
	case H3:
		return dom.H3(children...)
	case H4:
		return dom.H4(children...)
	default:
		return dom.H1(children...)
	}
}

func headingClasses(level HeadingLevel) string {
	base := "[font-weight:var(--font-weight-heading)] tracking-tight leading-tight text-[var(--color-text)]"
	switch level {
	case H2:
		return Class(base, "text-[length:var(--font-size-xl)]")
	case H3:
		return Class(base, "text-[length:var(--font-size-lg)]")
	case H4:
		return Class(base, "text-[length:var(--font-size-base)]")
	default:
		return Class(base, "text-[length:var(--font-size-2xl)]")
	}
}

// Text renders a body paragraph with a default or muted tone.
func Text(tone TextTone, content string, attrs ...dom.Node) dom.Node {
	children := make([]dom.Node, 0, len(attrs)+2)
	children = append(children, dom.Class(textClasses(tone)))
	children = append(children, attrs...)
	children = append(children, dom.Text(content))
	return dom.P(children...)
}

func textClasses(tone TextTone) string {
	base := "text-[length:var(--font-size-base)] leading-relaxed"
	switch tone {
	case TextMuted:
		return Class(base, "text-[var(--color-text-muted)]")
	default:
		return Class(base, "text-[var(--color-text)]")
	}
}

// Link renders an inline anchor in the primary color.
func Link(href, label string, attrs ...dom.Node) dom.Node {
	children := make([]dom.Node, 0, len(attrs)+2)
	children = append(children, dom.Class("font-medium text-[var(--color-primary)] underline decoration-2 underline-offset-2 transition-opacity hover:opacity-80"), dom.Href(href))
	children = append(children, attrs...)
	children = append(children, dom.Text(label))
	return dom.A(children...)
}
