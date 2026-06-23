// Package ui is a shared, server-rendered themed component library built on
// rio/dom. Component structure and fixed styling live here; per-product
// variation flows through CSS variables emitted by Tokens.StyleVars.
package ui

import (
	"strings"

	"github.com/tunedmystic/rio/dom"
)

// Tokens is the complete configuration surface a product supplies.
type Tokens struct {
	// Typography
	FontFamily   string
	FontSizeBase string
	FontSizeSm   string
	FontSizeLg   string
	FontSizeXl   string
	FontSize2xl  string

	// Colors
	ColorPrimary    string
	OnPrimary       string
	ColorSecondary  string
	OnSecondary     string
	ColorBackground string
	ColorSurface    string
	ColorText       string
	ColorTextMuted  string
	ColorBorder     string

	// Semantic (status) colors — used by Badge / Alert.
	ColorSuccess string
	ColorWarning string
	ColorDanger  string
	ColorInfo    string

	// Structure — a small set of non-color knobs products can tune.
	RadiusBase        string // base corner radius, e.g. "0.5rem"
	FontWeightHeading string // heading font weight, e.g. "700"
}

// StyleVars renders the product's tokens as a :root {...} <style> block.
// Token values are product-controlled compile-time constants, not user input,
// so they are emitted raw (CSS values must not be HTML-escaped).
func (tk Tokens) StyleVars() dom.Node {
	var b strings.Builder
	b.WriteString(":root{")
	writeVar(&b, "--font-family", tk.FontFamily)
	writeVar(&b, "--font-size-sm", tk.FontSizeSm)
	writeVar(&b, "--font-size-base", tk.FontSizeBase)
	writeVar(&b, "--font-size-lg", tk.FontSizeLg)
	writeVar(&b, "--font-size-xl", tk.FontSizeXl)
	writeVar(&b, "--font-size-2xl", tk.FontSize2xl)
	writeVar(&b, "--color-primary", tk.ColorPrimary)
	writeVar(&b, "--color-on-primary", tk.OnPrimary)
	writeVar(&b, "--color-secondary", tk.ColorSecondary)
	writeVar(&b, "--color-on-secondary", tk.OnSecondary)
	writeVar(&b, "--color-background", tk.ColorBackground)
	writeVar(&b, "--color-surface", tk.ColorSurface)
	writeVar(&b, "--color-text", tk.ColorText)
	writeVar(&b, "--color-text-muted", tk.ColorTextMuted)
	writeVar(&b, "--color-border", tk.ColorBorder)
	writeVar(&b, "--color-success", tk.ColorSuccess)
	writeVar(&b, "--color-warning", tk.ColorWarning)
	writeVar(&b, "--color-danger", tk.ColorDanger)
	writeVar(&b, "--color-info", tk.ColorInfo)
	writeVar(&b, "--radius-base", tk.RadiusBase)
	writeVar(&b, "--font-weight-heading", tk.FontWeightHeading)
	b.WriteString("}")
	return dom.StyleEl(dom.Raw(b.String()))
}

func writeVar(b *strings.Builder, name, val string) {
	if val == "" {
		return
	}
	b.WriteString(name)
	b.WriteString(":")
	b.WriteString(val)
	b.WriteString(";")
}

// Class joins class-name parts, trimming whitespace and dropping empties.
// It must NOT transform class names — that would defeat the Tailwind scanner,
// which only emits CSS for classes it finds as complete literal substrings.
func Class(parts ...string) string {
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return strings.Join(out, " ")
}

// withClass prepends a class attribute to a children slice. Shared by the
// components that simply wrap a variadic children tail.
func withClass(class string, children []dom.Node) []dom.Node {
	out := make([]dom.Node, 0, len(children)+1)
	out = append(out, dom.Class(class))
	out = append(out, children...)
	return out
}
