package ui

import (
	"bytes"
	"strings"
	"testing"

	"github.com/tunedmystic/rio/dom"
	"github.com/tunedmystic/rio/internal/assert"
)

// render is the shared test helper for the ui package: it renders a node to a
// string. Defined once here and reused across the ui test files.
func render(n dom.Node) string {
	var b bytes.Buffer
	_ = n.Render(&b)
	return b.String()
}

func Test_Class_TrimsAndJoins(t *testing.T) {
	assert.Equal(t, Class("a", "", "  b  ", "c"), "a b c")
	assert.Equal(t, Class(), "")
	assert.Equal(t, Class("  ", ""), "")
}

func Test_StyleVars_EmitsVariables(t *testing.T) {
	tk := Tokens{
		FontFamily:        `"Inter", sans-serif`,
		FontSizeBase:      "16px",
		FontSize2xl:       "32px",
		ColorPrimary:      "#059669",
		OnPrimary:         "#ffffff",
		ColorText:         "#0f172a",
		ColorTextMuted:    "#64748b",
		ColorBorder:       "#e2e8f0",
		ColorSuccess:      "#16a34a",
		RadiusBase:        "0.5rem",
		FontWeightHeading: "700",
	}
	html := render(tk.StyleVars())

	for _, want := range []string{
		"<style>", ":root{",
		`--font-family:"Inter", sans-serif;`,
		"--font-size-base:16px;",
		"--font-size-2xl:32px;",
		"--color-primary:#059669;",
		"--color-on-primary:#ffffff;",
		"--color-text:#0f172a;",
		"--color-text-muted:#64748b;",
		"--color-border:#e2e8f0;",
		"--color-success:#16a34a;",
		"--radius-base:0.5rem;",
		"--font-weight-heading:700;",
		"}</style>",
	} {
		if !strings.Contains(html, want) {
			t.Errorf("StyleVars output missing %q\ngot: %s", want, html)
		}
	}
}

func Test_StyleVars_OmitsEmptyTokens(t *testing.T) {
	html := render(Tokens{ColorPrimary: "#000"}.StyleVars())
	assert.True(t, strings.Contains(html, "--color-primary:#000;"))
	assert.False(t, strings.Contains(html, "--color-secondary"))
}
