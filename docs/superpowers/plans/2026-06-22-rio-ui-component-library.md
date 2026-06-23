# rio/ui Component Library Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the `rio/ui` themed component library (~16 Tier-1 components) plus a dev-only preview gallery, all inside the existing `rio` module.

**Architecture:** Free package-level functions (Pattern A) that build `dom.Node` trees. All per-product theming flows through CSS variables emitted by `Tokens.StyleVars()`. Every variant selects among full literal Tailwind class strings via `switch` — class names are never constructed at runtime. A `cmd/preview` HTTP server renders a gallery using the Tailwind v4 Play CDN.

**Tech Stack:** Go 1.24, `github.com/tunedmystic/rio/dom`, `internal/assert`, TailwindCSS v4 (Play CDN for preview only).

## Global Constraints

- Module path: `github.com/tunedmystic/rio` (lowercase). The `ui` package imports `github.com/tunedmystic/rio/dom`.
- **Never construct Tailwind class names at runtime.** No `"bg-" + x`, no `fmt.Sprintf`. Select full literals via `switch`.
- Per-product theming goes only through CSS variables (`bg-[var(--color-primary)]`).
- Components use typed `dom` helpers (`dom.Class`, `dom.Href`, `dom.Type`, `dom.Id`, `dom.Name`, `dom.Value`, `dom.For`, `dom.Map`) — not raw `dom.CreateAttr`.
- Caller contract: extra `attrs ...dom.Node` may carry `id`, `hx-*`, `aria-*` — but NOT `class` (the library owns the class attribute).
- Arbitrary font-size classes use the `length:` hint: `text-[length:var(--font-size-base)]` (disambiguates from color).
- Tests: package `ui`, import `github.com/tunedmystic/rio/internal/assert`. Render to a `bytes.Buffer`; assert HTML substrings with `strings.Contains(...)` + `assert.True(t, ...)`. Assert exact variant class literals.
- `ui` depends only on `rio/dom` and the standard library.
- Run all tests with `go test ./...` from repo root.

---

### Task 1: theme.go — Tokens, Class, StyleVars

**Files:**
- Create: `ui/theme.go`
- Test: `ui/theme_test.go`

**Interfaces:**
- Consumes: `dom.Node`, `dom.StyleEl`, `dom.Raw`, `dom.Class` from `rio/dom`.
- Produces:
  - `type Tokens struct { ... }` (fields below)
  - `func (tk Tokens) StyleVars() dom.Node`
  - `func Class(parts ...string) string`
  - `func withClass(class string, children []dom.Node) []dom.Node` (unexported helper used by later tasks)

- [ ] **Step 1: Write the failing test**

```go
// ui/theme_test.go
package ui

import (
	"strings"
	"testing"

	"github.com/tunedmystic/rio/internal/assert"
)

func Test_Class_TrimsAndJoins(t *testing.T) {
	assert.Equal(t, Class("a", "", "  b  ", "c"), "a b c")
	assert.Equal(t, Class(), "")
	assert.Equal(t, Class("  ", ""), "")
}

func Test_StyleVars_EmitsVariables(t *testing.T) {
	tk := Tokens{
		FontFamily:      `"Inter", sans-serif`,
		FontSizeBase:    "16px",
		FontSize2xl:     "32px",
		ColorPrimary:    "#059669",
		OnPrimary:       "#ffffff",
		ColorText:       "#0f172a",
		ColorTextMuted:  "#64748b",
		ColorBorder:     "#e2e8f0",
		ColorSuccess:    "#16a34a",
	}
	html := render(tk.StyleVars())

	for _, want := range []string{
		"<style>", ":root{",
		"--font-family:\"Inter\", sans-serif;",
		"--font-size-base:16px;",
		"--font-size-2xl:32px;",
		"--color-primary:#059669;",
		"--color-on-primary:#ffffff;",
		"--color-text:#0f172a;",
		"--color-text-muted:#64748b;",
		"--color-border:#e2e8f0;",
		"--color-success:#16a34a;",
		"}</style>",
	} {
		assert.True(t, strings.Contains(html, want))
	}
}

func Test_StyleVars_OmitsEmptyTokens(t *testing.T) {
	html := render(Tokens{ColorPrimary: "#000"}.StyleVars())
	assert.True(t, strings.Contains(html, "--color-primary:#000;"))
	assert.False(t, strings.Contains(html, "--color-secondary"))
}

// render is a shared test helper for the ui package test suite.
func render(n interface{ Render(w interface {
	Write([]byte) (int, error)
}) error }) string {
	var b strings.Builder
	_ = n.Render(&b)
	return b.String()
}
```

> Note: the `render` helper above is intentionally defined once here. If a cleaner signature is preferred, use `func render(n dom.Node) string` importing `dom` and `bytes` — but keep exactly one definition across the `ui` test files.

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./ui/ -run 'Class|StyleVars' -v`
Expected: FAIL — `undefined: Tokens`, `undefined: Class`.

- [ ] **Step 3: Write minimal implementation**

```go
// ui/theme.go
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

	// Semantic (status) colors
	ColorSuccess string
	ColorWarning string
	ColorDanger  string
	ColorInfo    string
}

// StyleVars renders the product's tokens as a :root {...} <style> block.
// Token values are product-controlled compile-time constants, not user
// input, so they are emitted raw (no HTML escaping of CSS values).
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
// It must NOT transform class names (that would defeat the Tailwind scanner).
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
// layout/feedback components that wrap a variadic children tail.
func withClass(class string, children []dom.Node) []dom.Node {
	out := make([]dom.Node, 0, len(children)+1)
	out = append(out, dom.Class(class))
	out = append(out, children...)
	return out
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./ui/ -run 'Class|StyleVars' -v`
Expected: PASS. If the `render` helper signature causes friction, replace it with `func render(n dom.Node) string { var b bytes.Buffer; _ = n.Render(&b); return b.String() }` (import `bytes` and `dom`) and re-run.

- [ ] **Step 5: Commit**

```bash
git add ui/theme.go ui/theme_test.go
git commit -m "feat(ui): add Tokens, StyleVars and Class theme foundation"
```

---

### Task 2: button.go — Button, ButtonLink

**Files:**
- Create: `ui/button.go`
- Test: `ui/button_test.go`

**Interfaces:**
- Consumes: `Class`, `dom.Class`, `dom.Type`, `dom.Href`, `dom.Text`, `dom.Button`, `dom.A`.
- Produces:
  - `type ButtonVariant int` with `ButtonPrimary, ButtonSecondary, ButtonDanger, ButtonGhost`
  - `func Button(variant ButtonVariant, label string, attrs ...dom.Node) dom.Node`
  - `func ButtonLink(variant ButtonVariant, href, label string, attrs ...dom.Node) dom.Node`

- [ ] **Step 1: Write the failing test**

```go
// ui/button_test.go
package ui

import (
	"strings"
	"testing"

	"github.com/tunedmystic/rio/dom"
	"github.com/tunedmystic/rio/internal/assert"
)

func Test_Button_Variants(t *testing.T) {
	cases := []struct {
		variant ButtonVariant
		wantBg  string
	}{
		{ButtonPrimary, "bg-[var(--color-primary)] text-[var(--color-on-primary)]"},
		{ButtonSecondary, "bg-[var(--color-secondary)] text-[var(--color-on-secondary)]"},
		{ButtonDanger, "bg-[var(--color-danger)] text-white"},
		{ButtonGhost, "bg-transparent text-[var(--color-text)]"},
	}
	for _, c := range cases {
		html := render(Button(c.variant, "Go"))
		assert.True(t, strings.HasPrefix(html, "<button "))
		assert.True(t, strings.Contains(html, `type="button"`))
		assert.True(t, strings.Contains(html, c.wantBg))
		assert.True(t, strings.Contains(html, ">Go</button>"))
	}
}

func Test_Button_PassesExtraAttrs(t *testing.T) {
	html := render(Button(ButtonPrimary, "Go", dom.Id("submit")))
	assert.True(t, strings.Contains(html, `id="submit"`))
}

func Test_ButtonLink_RendersAnchor(t *testing.T) {
	html := render(ButtonLink(ButtonPrimary, "/next", "Continue"))
	assert.True(t, strings.HasPrefix(html, "<a "))
	assert.True(t, strings.Contains(html, `href="/next"`))
	assert.True(t, strings.Contains(html, "bg-[var(--color-primary)]"))
	assert.True(t, strings.Contains(html, ">Continue</a>"))
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./ui/ -run Button -v`
Expected: FAIL — `undefined: Button`, `undefined: ButtonVariant`.

- [ ] **Step 3: Write minimal implementation**

```go
// ui/button.go
package ui

import "github.com/tunedmystic/rio/dom"

type ButtonVariant int

const (
	ButtonPrimary ButtonVariant = iota
	ButtonSecondary
	ButtonDanger
	ButtonGhost
)

// Button renders a styled <button>. Pass extra attributes (id, hx-*) via
// attrs; do not pass a class attribute — Button owns the class.
func Button(variant ButtonVariant, label string, attrs ...dom.Node) dom.Node {
	children := make([]dom.Node, 0, len(attrs)+3)
	children = append(children, dom.Class(buttonClasses(variant)), dom.Type("button"))
	children = append(children, attrs...)
	children = append(children, dom.Text(label))
	return dom.Button(children...)
}

// ButtonLink renders an <a> styled identically to Button, for CTAs that
// are navigation.
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
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./ui/ -run Button -v`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add ui/button.go ui/button_test.go
git commit -m "feat(ui): add Button and ButtonLink with variants"
```

---

### Task 3: cmd/preview — Play-CDN gallery skeleton + Makefile target

**Files:**
- Create: `cmd/preview/main.go`
- Modify: `Makefile` (add `preview` target)

**Interfaces:**
- Consumes: `ui.Tokens`, `ui.StyleVars`, `ui.Container`, `ui.Heading`, `ui.Stack`, `ui.Button` (only Button + theme exist so far; later tasks extend `gallery()`).
- Produces: a runnable `cmd/preview` binary serving `:8080`.

> This task has no unit test (it is a dev-only HTTP entry point). Its deliverable is verified by `go build ./cmd/preview` and a manual `make preview` smoke check. The gallery is filled out in Task 8.

- [ ] **Step 1: Write the gallery server**

```go
// cmd/preview/main.go
package main

import (
	"net/http"

	"github.com/tunedmystic/rio/dom"
	"github.com/tunedmystic/rio/ui"
)

// themes are demo token sets so the CSS variables resolve. Swap via ?theme=.
var themes = map[string]ui.Tokens{
	"apron": {
		FontFamily:      `"Source Serif 4", serif`,
		FontSizeSm:      "14px",
		FontSizeBase:    "16px",
		FontSizeLg:      "18px",
		FontSizeXl:      "24px",
		FontSize2xl:     "32px",
		ColorPrimary:    "#059669",
		OnPrimary:       "#ffffff",
		ColorSecondary:  "#475569",
		OnSecondary:     "#ffffff",
		ColorBackground: "#f8fafc",
		ColorSurface:    "#ffffff",
		ColorText:       "#0f172a",
		ColorTextMuted:  "#64748b",
		ColorBorder:     "#e2e8f0",
		ColorSuccess:    "#16a34a",
		ColorWarning:    "#d97706",
		ColorDanger:     "#dc2626",
		ColorInfo:       "#2563eb",
	},
	"teddy": {
		FontFamily:      `"Inter", sans-serif`,
		FontSizeSm:      "13px",
		FontSizeBase:    "15px",
		FontSizeLg:      "17px",
		FontSizeXl:      "22px",
		FontSize2xl:     "30px",
		ColorPrimary:    "#4f46e5",
		OnPrimary:       "#ffffff",
		ColorSecondary:  "#64748b",
		OnSecondary:     "#ffffff",
		ColorBackground: "#ffffff",
		ColorSurface:    "#f8fafc",
		ColorText:       "#1e1b2e",
		ColorTextMuted:  "#6b7280",
		ColorBorder:     "#e5e7eb",
		ColorSuccess:    "#16a34a",
		ColorWarning:    "#d97706",
		ColorDanger:     "#dc2626",
		ColorInfo:       "#2563eb",
	},
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tokens, ok := themes[r.URL.Query().Get("theme")]
		if !ok {
			tokens = themes["apron"]
		}

		page := dom.Doctype(dom.Html(
			dom.Head(
				dom.Meta(dom.Charset("utf-8")),
				dom.Meta(dom.Name("viewport"), dom.Content("width=device-width, initial-scale=1")),
				tokens.StyleVars(),
				dom.Script(dom.Src("https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4")),
			),
			dom.Body(
				dom.Class("bg-[var(--color-background)] text-[var(--color-text)] p-8"),
				gallery(),
			),
		))
		_ = page.Render(w)
	})

	_ = http.ListenAndServe(":8080", nil)
}

// gallery renders one of every component with all its variants. Extended as
// components are added.
func gallery() dom.Node {
	return ui.Container(
		ui.Heading(ui.H1, "rio/ui component gallery"),

		ui.Heading(ui.H2, "Buttons"),
		ui.Stack(ui.GapSm,
			ui.Button(ui.ButtonPrimary, "Primary"),
			ui.Button(ui.ButtonSecondary, "Secondary"),
			ui.Button(ui.ButtonDanger, "Danger"),
			ui.Button(ui.ButtonGhost, "Ghost"),
		),
	)
}
```

> The gallery references `ui.Heading`, `ui.H1/H2`, `ui.Stack`, `ui.GapSm` — these do not exist until Tasks 4–5. To keep this task independently buildable, comment out everything in `gallery()` except `return ui.Container(ui.Button(ui.ButtonPrimary, "Primary"))` for now, and restore the full body in Task 8. Alternatively, sequence this task after Task 5; either is acceptable. The committed state of THIS task must compile.

- [ ] **Step 2: Add the Makefile target**

Add after the existing app-related targets in `Makefile`:

```make
## @(app) - 🎨 Serve the component gallery on :8080
preview:
	@echo "✨📦✨ Serving component gallery on :8080\n"
	@go run ./cmd/preview
```

- [ ] **Step 3: Verify it builds**

Run: `go build ./cmd/preview`
Expected: no output, exit 0. (Ensure `gallery()` only references components that exist at this commit.)

- [ ] **Step 4: Commit**

```bash
git add cmd/preview/main.go Makefile
git commit -m "feat(ui): add cmd/preview Play-CDN gallery and make target"
```

---

### Task 4: typography.go — Heading, Text, Link

**Files:**
- Create: `ui/typography.go`
- Test: `ui/typography_test.go`

**Interfaces:**
- Consumes: `Class`, `dom.Class`, `dom.Href`, `dom.Text`, `dom.H1`..`dom.H4`, `dom.P`, `dom.A`.
- Produces:
  - `type HeadingLevel int` with `H1, H2, H3, H4` (`iota + 1`)
  - `type TextTone int` with `TextDefault, TextMuted`
  - `func Heading(level HeadingLevel, text string, attrs ...dom.Node) dom.Node`
  - `func Text(tone TextTone, content string, attrs ...dom.Node) dom.Node`
  - `func Link(href, label string, attrs ...dom.Node) dom.Node`

- [ ] **Step 1: Write the failing test**

```go
// ui/typography_test.go
package ui

import (
	"strings"
	"testing"

	"github.com/tunedmystic/rio/internal/assert"
)

func Test_Heading_LevelsAndSizes(t *testing.T) {
	cases := []struct {
		level    HeadingLevel
		tag      string
		sizeVar  string
	}{
		{H1, "h1", "text-[length:var(--font-size-2xl)]"},
		{H2, "h2", "text-[length:var(--font-size-xl)]"},
		{H3, "h3", "text-[length:var(--font-size-lg)]"},
		{H4, "h4", "text-[length:var(--font-size-base)]"},
	}
	for _, c := range cases {
		html := render(Heading(c.level, "Title"))
		assert.True(t, strings.HasPrefix(html, "<"+c.tag+" "))
		assert.True(t, strings.Contains(html, c.sizeVar))
		assert.True(t, strings.Contains(html, ">Title</"+c.tag+">"))
	}
}

func Test_Text_Tones(t *testing.T) {
	def := render(Text(TextDefault, "body"))
	assert.True(t, strings.HasPrefix(def, "<p "))
	assert.True(t, strings.Contains(def, "text-[var(--color-text)]"))

	muted := render(Text(TextMuted, "body"))
	assert.True(t, strings.Contains(muted, "text-[var(--color-text-muted)]"))
}

func Test_Link_PrimaryColor(t *testing.T) {
	html := render(Link("/about", "About"))
	assert.True(t, strings.HasPrefix(html, "<a "))
	assert.True(t, strings.Contains(html, `href="/about"`))
	assert.True(t, strings.Contains(html, "text-[var(--color-primary)]"))
	assert.True(t, strings.Contains(html, ">About</a>"))
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./ui/ -run 'Heading|Text|Link' -v`
Expected: FAIL — `undefined: Heading`, etc.

- [ ] **Step 3: Write minimal implementation**

```go
// ui/typography.go
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

// Heading renders an h1–h4 from the font-size scale.
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
	base := "font-bold text-[var(--color-text)]"
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
	children = append(children, dom.Class("text-[var(--color-primary)] underline hover:opacity-80"), dom.Href(href))
	children = append(children, attrs...)
	children = append(children, dom.Text(label))
	return dom.A(children...)
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./ui/ -run 'Heading|Text|Link' -v`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add ui/typography.go ui/typography_test.go
git commit -m "feat(ui): add Heading, Text and Link typography"
```

---

### Task 5: layout.go — Container, Section, Card, Stack

**Files:**
- Create: `ui/layout.go`
- Test: `ui/layout_test.go`

**Interfaces:**
- Consumes: `Class`, `withClass`, `dom.Div`, `dom.Section`.
- Produces:
  - `func Container(children ...dom.Node) dom.Node`
  - `func Section(children ...dom.Node) dom.Node`
  - `func Card(children ...dom.Node) dom.Node`
  - `type Gap int` with `GapSm, GapMd, GapLg`
  - `func Stack(gap Gap, children ...dom.Node) dom.Node`

- [ ] **Step 1: Write the failing test**

```go
// ui/layout_test.go
package ui

import (
	"strings"
	"testing"

	"github.com/tunedmystic/rio/dom"
	"github.com/tunedmystic/rio/internal/assert"
)

func Test_Container(t *testing.T) {
	html := render(Container(dom.Text("inside")))
	assert.True(t, strings.HasPrefix(html, "<div "))
	assert.True(t, strings.Contains(html, "max-w-7xl mx-auto px-4"))
	assert.True(t, strings.Contains(html, ">inside</div>"))
}

func Test_Section(t *testing.T) {
	html := render(Section(dom.Text("band")))
	assert.True(t, strings.HasPrefix(html, "<section "))
	assert.True(t, strings.Contains(html, "py-12"))
}

func Test_Card(t *testing.T) {
	html := render(Card(dom.Text("raised")))
	assert.True(t, strings.Contains(html, "bg-[var(--color-surface)]"))
	assert.True(t, strings.Contains(html, "border-[var(--color-border)]"))
	assert.True(t, strings.Contains(html, "rounded-lg"))
}

func Test_Stack_Gaps(t *testing.T) {
	assert.True(t, strings.Contains(render(Stack(GapSm, dom.Text("x"))), "flex flex-col gap-2"))
	assert.True(t, strings.Contains(render(Stack(GapMd, dom.Text("x"))), "flex flex-col gap-4"))
	assert.True(t, strings.Contains(render(Stack(GapLg, dom.Text("x"))), "flex flex-col gap-8"))
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./ui/ -run 'Container|Section|Card|Stack' -v`
Expected: FAIL — `undefined: Container`, etc.

- [ ] **Step 3: Write minimal implementation**

```go
// ui/layout.go
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
	return dom.Div(withClass("bg-[var(--color-surface)] border border-[var(--color-border)] rounded-lg p-6", children)...)
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
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./ui/ -run 'Container|Section|Card|Stack' -v`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add ui/layout.go ui/layout_test.go
git commit -m "feat(ui): add Container, Section, Card and Stack layout"
```

---

### Task 6: form.go — Label, FieldError, TextField, Textarea, Select, Checkbox, Radio

**Files:**
- Create: `ui/form.go`
- Test: `ui/form_test.go`

**Interfaces:**
- Consumes: `Class`, `dom.Class`, `dom.For`, `dom.Type`, `dom.Id`, `dom.Name`, `dom.Value`, `dom.Selected`, `dom.Checked`, `dom.Text`, `dom.Raw`, `dom.Div`, `dom.Label`, `dom.Input`, `dom.Textarea`, `dom.Select`, `dom.Option`, `dom.P`, `dom.Map`.
- Produces:
  - `type Option struct { Value, Label string }`
  - `func Label(forID, text string, attrs ...dom.Node) dom.Node`
  - `func FieldError(msg string) dom.Node`
  - `func TextField(name, label, value, errMsg string, attrs ...dom.Node) dom.Node`
  - `func Textarea(name, label, value, errMsg string, attrs ...dom.Node) dom.Node`
  - `func Select(name, label string, options []Option, selected, errMsg string) dom.Node`
  - `func Checkbox(name, label string, checked bool, attrs ...dom.Node) dom.Node`
  - `func Radio(name, label, value string, checked bool, attrs ...dom.Node) dom.Node`

- [ ] **Step 1: Write the failing test**

```go
// ui/form_test.go
package ui

import (
	"strings"
	"testing"

	"github.com/tunedmystic/rio/internal/assert"
)

func Test_Label(t *testing.T) {
	html := render(Label("email", "Email"))
	assert.True(t, strings.HasPrefix(html, "<label "))
	assert.True(t, strings.Contains(html, `for="email"`))
	assert.True(t, strings.Contains(html, ">Email</label>"))
}

func Test_FieldError(t *testing.T) {
	assert.Equal(t, render(FieldError("")), "")
	html := render(FieldError("Required"))
	assert.True(t, strings.Contains(html, "text-[var(--color-danger)]"))
	assert.True(t, strings.Contains(html, ">Required</p>"))
}

func Test_TextField(t *testing.T) {
	html := render(TextField("email", "Email", "a@b.com", ""))
	assert.True(t, strings.Contains(html, `for="email"`))
	assert.True(t, strings.Contains(html, `type="text"`))
	assert.True(t, strings.Contains(html, `name="email"`))
	assert.True(t, strings.Contains(html, `value="a@b.com"`))
	// no error => no danger text
	assert.False(t, strings.Contains(html, "text-[var(--color-danger)]"))
}

func Test_TextField_WithError(t *testing.T) {
	html := render(TextField("email", "Email", "", "Required"))
	assert.True(t, strings.Contains(html, "text-[var(--color-danger)]"))
	assert.True(t, strings.Contains(html, ">Required</p>"))
}

func Test_Textarea(t *testing.T) {
	html := render(Textarea("bio", "Bio", "hello", ""))
	assert.True(t, strings.Contains(html, "<textarea "))
	assert.True(t, strings.Contains(html, `name="bio"`))
	assert.True(t, strings.Contains(html, ">hello</textarea>"))
}

func Test_Select_MarksSelected(t *testing.T) {
	opts := []Option{{"us", "USA"}, {"ca", "Canada"}}
	html := render(Select("country", "Country", opts, "ca", ""))
	assert.True(t, strings.Contains(html, "<select "))
	assert.True(t, strings.Contains(html, `<option value="us">USA</option>`))
	assert.True(t, strings.Contains(html, `<option value="ca" selected>Canada</option>`))
}

func Test_Checkbox(t *testing.T) {
	html := render(Checkbox("agree", "I agree", true))
	assert.True(t, strings.Contains(html, `type="checkbox"`))
	assert.True(t, strings.Contains(html, `name="agree"`))
	assert.True(t, strings.Contains(html, " checked"))
	assert.True(t, strings.Contains(html, ">I agree</label>"))
}

func Test_Radio(t *testing.T) {
	on := render(Radio("plan", "Pro", "pro", true))
	assert.True(t, strings.Contains(on, `type="radio"`))
	assert.True(t, strings.Contains(on, `value="pro"`))
	assert.True(t, strings.Contains(on, " checked"))

	off := render(Radio("plan", "Free", "free", false))
	assert.False(t, strings.Contains(off, " checked"))
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./ui/ -run 'Label|FieldError|TextField|Textarea|Select|Checkbox|Radio' -v`
Expected: FAIL — `undefined: Label`, etc.

- [ ] **Step 3: Write minimal implementation**

```go
// ui/form.go
package ui

import "github.com/tunedmystic/rio/dom"

// Option is a single choice in a Select.
type Option struct {
	Value string
	Label string
}

// inputClasses is the shared field styling for text-like inputs.
const inputClasses = "block w-full rounded-md border border-[var(--color-border)] bg-[var(--color-surface)] px-3 py-2 text-[var(--color-text)] focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)]"

// Label renders a form label bound to a field id.
func Label(forID, text string, attrs ...dom.Node) dom.Node {
	children := make([]dom.Node, 0, len(attrs)+3)
	children = append(children, dom.Class("block text-[length:var(--font-size-sm)] font-medium text-[var(--color-text)] mb-1"), dom.For(forID))
	children = append(children, attrs...)
	children = append(children, dom.Text(text))
	return dom.Label(children...)
}

// FieldError renders small error text under a field, or nothing when msg is empty.
func FieldError(msg string) dom.Node {
	if msg == "" {
		return dom.Raw("")
	}
	return dom.P(
		dom.Class("text-[length:var(--font-size-sm)] text-[var(--color-danger)] mt-1"),
		dom.Text(msg),
	)
}

// TextField is a label + text input + optional error.
func TextField(name, label, value, errMsg string, attrs ...dom.Node) dom.Node {
	input := make([]dom.Node, 0, len(attrs)+5)
	input = append(input, dom.Class(inputClasses), dom.Type("text"), dom.Id(name), dom.Name(name), dom.Value(value))
	input = append(input, attrs...)
	return dom.Div(
		dom.Class("mb-4"),
		Label(name, label),
		dom.Input(input...),
		FieldError(errMsg),
	)
}

// Textarea is a label + multi-line input + optional error.
func Textarea(name, label, value, errMsg string, attrs ...dom.Node) dom.Node {
	ta := make([]dom.Node, 0, len(attrs)+4)
	ta = append(ta, dom.Class(inputClasses), dom.Id(name), dom.Name(name))
	ta = append(ta, attrs...)
	ta = append(ta, dom.Text(value))
	return dom.Div(
		dom.Class("mb-4"),
		Label(name, label),
		dom.Textarea(ta...),
		FieldError(errMsg),
	)
}

// Select is a label + dropdown built from options, marking the selected value.
func Select(name, label string, options []Option, selected, errMsg string) dom.Node {
	return dom.Div(
		dom.Class("mb-4"),
		Label(name, label),
		dom.Select(
			dom.Class(inputClasses),
			dom.Id(name),
			dom.Name(name),
			dom.Map(options, func(o Option) dom.Node {
				opt := make([]dom.Node, 0, 3)
				opt = append(opt, dom.Value(o.Value))
				if o.Value == selected {
					opt = append(opt, dom.Selected())
				}
				opt = append(opt, dom.Text(o.Label))
				return dom.Option(opt...)
			}),
		),
		FieldError(errMsg),
	)
}

// Checkbox is a box + label.
func Checkbox(name, label string, checked bool, attrs ...dom.Node) dom.Node {
	return checkable("checkbox", name, label, "", checked, attrs...)
}

// Radio is a circle + label.
func Radio(name, label, value string, checked bool, attrs ...dom.Node) dom.Node {
	return checkable("radio", name, label, value, checked, attrs...)
}

func checkable(kind, name, label, value string, checked bool, attrs ...dom.Node) dom.Node {
	input := make([]dom.Node, 0, len(attrs)+5)
	input = append(input, dom.Class("h-4 w-4 rounded border-[var(--color-border)] text-[var(--color-primary)] focus:ring-[var(--color-primary)]"), dom.Type(kind), dom.Name(name))
	if value != "" {
		input = append(input, dom.Value(value))
	}
	if checked {
		input = append(input, dom.Checked())
	}
	input = append(input, attrs...)
	return dom.Label(
		dom.Class("inline-flex items-center gap-2 text-[var(--color-text)]"),
		dom.Input(input...),
		dom.Text(label),
	)
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./ui/ -run 'Label|FieldError|TextField|Textarea|Select|Checkbox|Radio' -v`
Expected: PASS. (If `Test_Select_MarksSelected` fails on attribute order, inspect the rendered output and adjust the expected substring to match `dom`'s actual attribute ordering — attributes render in child order: `value`, then `selected`.)

- [ ] **Step 5: Commit**

```bash
git add ui/form.go ui/form_test.go
git commit -m "feat(ui): add form components (TextField, Select, Checkbox, etc.)"
```

---

### Task 7: feedback.go — Badge, Alert

**Files:**
- Create: `ui/feedback.go`
- Test: `ui/feedback_test.go`

**Interfaces:**
- Consumes: `Class`, `dom.Class`, `dom.Role`, `dom.Text`, `dom.Span`, `dom.Div`.
- Produces:
  - `type BadgeVariant int` with `BadgeNeutral, BadgeSuccess, BadgeWarning, BadgeDanger`
  - `func Badge(variant BadgeVariant, label string) dom.Node`
  - `type AlertVariant int` with `AlertInfo, AlertSuccess, AlertWarning, AlertError`
  - `func Alert(variant AlertVariant, content ...dom.Node) dom.Node`

- [ ] **Step 1: Write the failing test**

```go
// ui/feedback_test.go
package ui

import (
	"strings"
	"testing"

	"github.com/tunedmystic/rio/dom"
	"github.com/tunedmystic/rio/internal/assert"
)

func Test_Badge_Variants(t *testing.T) {
	cases := []struct {
		variant BadgeVariant
		want    string
	}{
		{BadgeNeutral, "bg-[var(--color-border)] text-[var(--color-text)]"},
		{BadgeSuccess, "bg-[var(--color-success)] text-white"},
		{BadgeWarning, "bg-[var(--color-warning)] text-white"},
		{BadgeDanger, "bg-[var(--color-danger)] text-white"},
	}
	for _, c := range cases {
		html := render(Badge(c.variant, "New"))
		assert.True(t, strings.HasPrefix(html, "<span "))
		assert.True(t, strings.Contains(html, "rounded-full"))
		assert.True(t, strings.Contains(html, c.want))
		assert.True(t, strings.Contains(html, ">New</span>"))
	}
}

func Test_Alert_Variants(t *testing.T) {
	cases := []struct {
		variant AlertVariant
		border  string
	}{
		{AlertInfo, "border-[var(--color-info)]"},
		{AlertSuccess, "border-[var(--color-success)]"},
		{AlertWarning, "border-[var(--color-warning)]"},
		{AlertError, "border-[var(--color-danger)]"},
	}
	for _, c := range cases {
		html := render(Alert(c.variant, dom.Text("heads up")))
		assert.True(t, strings.HasPrefix(html, "<div "))
		assert.True(t, strings.Contains(html, `role="alert"`))
		assert.True(t, strings.Contains(html, c.border))
		assert.True(t, strings.Contains(html, "heads up"))
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./ui/ -run 'Badge|Alert' -v`
Expected: FAIL — `undefined: Badge`, etc.

- [ ] **Step 3: Write minimal implementation**

```go
// ui/feedback.go
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
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./ui/ -run 'Badge|Alert' -v`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add ui/feedback.go ui/feedback_test.go
git commit -m "feat(ui): add Badge and Alert feedback components"
```

---

### Task 8: Fill the gallery + full verification

**Files:**
- Modify: `cmd/preview/main.go` (expand `gallery()`)

**Interfaces:**
- Consumes: every exported `ui` component.

- [ ] **Step 1: Replace `gallery()` with the full version**

```go
// gallery renders one of every component with all its variants.
func gallery() dom.Node {
	return ui.Container(
		ui.Heading(ui.H1, "rio/ui component gallery"),

		ui.Heading(ui.H2, "Typography"),
		ui.Stack(ui.GapSm,
			ui.Heading(ui.H3, "Heading level 3"),
			ui.Text(ui.TextDefault, "Default body text in the product font."),
			ui.Text(ui.TextMuted, "Muted secondary text."),
			ui.Link("#", "An inline link"),
		),

		ui.Heading(ui.H2, "Buttons"),
		ui.Stack(ui.GapSm,
			ui.Button(ui.ButtonPrimary, "Primary"),
			ui.Button(ui.ButtonSecondary, "Secondary"),
			ui.Button(ui.ButtonDanger, "Danger"),
			ui.Button(ui.ButtonGhost, "Ghost"),
			ui.ButtonLink(ui.ButtonPrimary, "#", "Button Link"),
		),

		ui.Heading(ui.H2, "Card"),
		ui.Card(
			ui.Heading(ui.H3, "A card"),
			ui.Text(ui.TextMuted, "Raised surface with border and padding."),
		),

		ui.Heading(ui.H2, "Form"),
		ui.Card(
			ui.TextField("email", "Email", "", ""),
			ui.TextField("user", "Username", "taken", "That username is taken"),
			ui.Textarea("bio", "Bio", "", ""),
			ui.Select("country", "Country",
				[]ui.Option{{Value: "us", Label: "USA"}, {Value: "ca", Label: "Canada"}}, "ca", ""),
			ui.Checkbox("agree", "I agree to the terms", true),
			ui.Radio("plan", "Pro plan", "pro", true),
		),

		ui.Heading(ui.H2, "Badges"),
		ui.Stack(ui.GapSm,
			ui.Badge(ui.BadgeNeutral, "Neutral"),
			ui.Badge(ui.BadgeSuccess, "Success"),
			ui.Badge(ui.BadgeWarning, "Warning"),
			ui.Badge(ui.BadgeDanger, "Danger"),
		),

		ui.Heading(ui.H2, "Alerts"),
		ui.Stack(ui.GapMd,
			ui.Alert(ui.AlertInfo, ui.Text(ui.TextDefault, "Informational message.")),
			ui.Alert(ui.AlertSuccess, ui.Text(ui.TextDefault, "It worked.")),
			ui.Alert(ui.AlertWarning, ui.Text(ui.TextDefault, "Careful now.")),
			ui.Alert(ui.AlertError, ui.Text(ui.TextDefault, "Something went wrong.")),
		),
	)
}
```

- [ ] **Step 2: Build and run the full test suite**

Run: `go build ./... && go vet ./... && go test ./...`
Expected: build succeeds, vet clean, all tests PASS.

- [ ] **Step 3: Manual smoke check (optional but recommended)**

Run: `make preview` then open `http://localhost:8080` and `http://localhost:8080?theme=teddy`.
Expected: gallery renders; switching `?theme=` swaps colors and font. Stop with Ctrl-C.

- [ ] **Step 4: Commit**

```bash
git add cmd/preview/main.go
git commit -m "feat(ui): render full component gallery in preview"
```

---

## Self-Review

**Spec coverage:**
- theme.go (Tokens, StyleVars, Class, withClass) → Task 1 ✓
- Variant enums → defined in their component files (Tasks 2, 4, 5, 7) ✓
- Button/ButtonLink → Task 2 ✓
- cmd/preview + Makefile → Tasks 3, 8 ✓
- Heading/Text/Link → Task 4 ✓
- Container/Section/Card/Stack → Task 5 ✓
- Label/FieldError/TextField/Textarea/Select/Checkbox/Radio + Option → Task 6 ✓
- Badge/Alert → Task 7 ✓
- Hard rule (no constructed class names): every variant uses `switch` over full literals ✓
- CSS-variable theming: all colors via `var(--...)` ✓
- Testing: render-to-buffer + literal class assertions per component ✓
- Out of scope (product wiring, vendored Tailwind build, Tier-2) → not included ✓

**Type consistency:** `ButtonVariant`, `HeadingLevel`, `TextTone`, `Gap`, `BadgeVariant`, `AlertVariant`, `Option` defined once; `render` helper defined once in `ui/theme_test.go`; `inputClasses` const shared within form.go; `withClass` defined in theme.go and used by layout.go. Consistent across tasks.

**Known follow-up:** Task 3's `gallery()` must compile at its own commit — it references components built in Tasks 4–5. Resolved by committing a minimal `gallery()` in Task 3 (Button only) and expanding in Task 8.
