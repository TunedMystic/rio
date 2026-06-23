# rio/ui — Themed Component Library — Design

**Date:** 2026-06-22
**Status:** Approved, ready for implementation planning
**Source spec:** `rio-ui-spec.md` (provided by user)

## Goal

A shared, server-rendered UI component library living as the `ui` package inside the
`rio` module. Built on `rio/dom`, styled with TailwindCSS v4. Multiple products
configure it with their own design tokens (colors, font scale). Structure and fixed
styling live in the library; per-product variation flows through **CSS variables**.

## Scope (this effort)

Build-order steps 1–7 from the source spec, all self-contained in the `rio` repo:

- The `ui` package: all ~16 Tier-1 components, `theme.go` foundation.
- `cmd/preview`: a dev-only Tailwind Play-CDN gallery for visual development.
- Tests for every component file.
- A `preview` Makefile target.

**Explicitly out of scope** (deferred):

- Product wiring (source spec step 8): tokens + Layout + a real screen in a product.
- The vendored-binary production Tailwind build (source spec §8) — this is product-side.
- All Tier-2 components (source spec §5: Table, Hero, Navbar, Footer, etc.) and generic
  `Box`/`Flex` wrappers. These stay inline in products until the rule of three.

## Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Module path | `github.com/tunedmystic/rio/ui` | Real `go.mod` module is lowercase `tunedmystic`; source spec's `TunedMystic` casing is cosmetic. |
| Theme injection | **Pattern A**: free package-level functions + CSS vars | Least machinery; theming flows entirely through `StyleVars()`. Source spec §6 recommendation. |
| Token color classes | **Raw arbitrary values** `bg-[var(--color-primary)]` | Works on Tailwind v4 with zero extra config. No `@theme` indirection layer. |
| Variant enum location | In each component's file | Source spec §3.4 note ("kept in the relevant file"). `theme.go` holds only `Tokens`, `StyleVars`, `Class`. |
| Testing | Render-to-buffer + `internal/assert` substring assertions | Mirrors existing `dom_test.go` style. Assert exact variant class literals. |

## Hard rules (inherited from source spec, do not break)

1. **Never construct Tailwind class names at runtime.** Always *select* among full
   literal strings via `switch`/maps; never *build* them (`"bg-" + c` is forbidden).
   The Tailwind scanner only emits CSS for classes it finds as complete literals.
2. **Per-product theming goes through CSS variables** (`bg-[var(--color-primary)]`).
3. **Only add a component when it carries token styling OR has variants.** Otherwise
   call `rio/dom` directly.

## Architecture

### `theme.go` — foundation

- **`Tokens`** struct — verbatim from source spec §3.1: typography (`FontFamily`,
  `FontSizeBase/Sm/Lg/Xl/2xl`) and colors (primary/secondary/background/surface/text/
  text-muted/border + on-primary/on-secondary + semantic success/warning/danger/info).
- **`Class(parts ...string) string`** — verbatim: `TrimSpace` each part, drop empties,
  `strings.Join` with a space. Must NOT transform class names.
- **`StyleVars() dom.Node`** — renders a `<style>` block defining all CSS variables in
  `:root`. Implementation: build the CSS body with a `strings.Builder` mapping each
  `Tokens` field to its `--var` name (source spec §3.2 naming convention), wrapped as
  `dom.StyleEl(dom.Raw(":root{…}"))`. Uses `dom.Raw` (not `dom.Text`) so CSS values are
  not HTML-escaped — font-family quotes and the like must pass through verbatim.
  - **Trust boundary:** tokens are product-controlled compile-time constants, not user
    input. No sanitization of token values; a comment documents this assumption.

CSS variable names emitted:

```
--font-family
--font-size-sm | --font-size-base | --font-size-lg | --font-size-xl | --font-size-2xl
--color-primary | --color-on-primary
--color-secondary | --color-on-secondary
--color-background | --color-surface
--color-text | --color-text-muted | --color-border
--color-success | --color-warning | --color-danger | --color-info
```

### The variant → class pattern (the core of the library)

Every variant maps to **full literal class strings via a `switch`**. Each component file
has an unexported `<component>Classes(variant) string` helper. Canonical shape:

```go
func buttonClasses(v ButtonVariant) string {
    base := "inline-flex items-center justify-center rounded-md px-4 py-2 font-medium transition-colors"
    switch v {
    case ButtonPrimary:   return Class(base, "bg-[var(--color-primary)] text-[var(--color-on-primary)] hover:opacity-90")
    case ButtonSecondary: return Class(base, "bg-[var(--color-secondary)] text-[var(--color-on-secondary)] hover:opacity-90")
    case ButtonDanger:    return Class(base, "bg-[var(--color-danger)] text-white hover:opacity-90")
    case ButtonGhost:     return Class(base, "bg-transparent text-[var(--color-text)] hover:bg-[var(--color-border)]")
    default:              return buttonClasses(ButtonPrimary)
    }
}
```

### Component construction convention

Components are package-level functions using the typed `dom` helpers (`dom.Class`,
`dom.Href`, `dom.Type`, `dom.Id`, `dom.Map`, etc.) rather than raw `dom.CreateAttr`.
The variadic `attrs ...dom.Node` tail is spliced into the element's children.

- Arg ordering between the library's own attrs and caller `attrs` is safe:
  `htmlElement.Render` emits all `HtmlAttributer` children before non-attribute
  children regardless of position.
- **Caller contract:** pass extra attributes (`id`, `hx-*`, `aria-*`) via `attrs`, but
  **not** `class` — a second `class` attribute would be emitted as a duplicate. The
  library owns the class. Documented on each component.

### Components (Tier-1, source spec §4)

`layout.go`
- `Container(children ...dom.Node) dom.Node` — `max-w-7xl mx-auto px-4`.
- `Section(children ...dom.Node) dom.Node` — vertical rhythm band (`py-12`/`py-16`).
- `Card(children ...dom.Node) dom.Node` — surface bg + border + radius + padding.
- `Stack(gap Gap, children ...dom.Node) dom.Node` — flex column; `Gap` enum (`GapSm/Md/Lg`).

`typography.go`
- `Heading(level HeadingLevel, text string, attrs ...dom.Node) dom.Node` — `HeadingLevel`
  enum `H1..H4` selects both the tag (`dom.H1..dom.H4`) and a font-size class from the
  scale (`text-[var(--font-size-2xl)]` … down to base).
- `Text(tone TextTone, content string, attrs ...dom.Node) dom.Node` — `TextTone` enum
  `TextDefault`/`TextMuted` selects `text-[var(--color-text)]` vs `--color-text-muted`.
- `Link(href, label string, attrs ...dom.Node) dom.Node` — primary color + hover.

`button.go`
- `ButtonVariant` enum: `ButtonPrimary/Secondary/Danger/Ghost`.
- `Button(variant ButtonVariant, label string, attrs ...dom.Node) dom.Node`.
- `ButtonLink(variant ButtonVariant, href, label string, attrs ...dom.Node) dom.Node` —
  an `<a>` styled identically (shares `buttonClasses`).

`form.go`
- `Option struct { Value, Label string }`.
- `Label(forID, text string, attrs ...dom.Node) dom.Node` — shared primitive.
- `FieldError(msg string) dom.Node` — renders an empty node when `msg == ""`.
- `TextField(name, label, value, errMsg string, attrs ...dom.Node) dom.Node` — Label +
  input + optional FieldError. The workhorse.
- `Textarea(name, label, value, errMsg string, attrs ...dom.Node) dom.Node` — same shape.
- `Select(name, label string, options []Option, selected, errMsg string) dom.Node` —
  builds `<option>`s with `dom.Map`; marks `selected`.
- `Checkbox(name, label string, checked bool, attrs ...dom.Node) dom.Node`.
- `Radio(name, label, value string, checked bool, attrs ...dom.Node) dom.Node` — may
  share impl with Checkbox via an internal variant.

`feedback.go`
- `BadgeVariant` enum: `BadgeNeutral/Success/Warning/Danger`.
- `Badge(variant BadgeVariant, label string) dom.Node` — small status pill.
- `AlertVariant` enum: `AlertInfo/Success/Warning/Error`.
- `Alert(variant AlertVariant, content ...dom.Node) dom.Node` — callout box.

### `cmd/preview` — dev-only gallery

- `main.go`: HTTP server on `:8080`. Dev-only, never ships to production. The Tailwind
  v4 **Play CDN** (`<script src="https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4">`)
  compiles classes from the rendered DOM at runtime — no `@source` scanning, no build
  step, no watcher.
- `themes` map with two demo token sets: `apron` (emerald / Source Serif 4) and `teddy`
  (indigo / Inter). `?theme=` query param swaps tokens; default `apron`. This proves the
  CSS-variable theming actually swaps across products.
- Page = `dom.Doctype(dom.Html(dom.Head(tokens.StyleVars(), <cdn script>), dom.Body(...,
  gallery())))`.
- `gallery()` renders one of every component with all variants, driven from a slice of
  demo sections so adding a component is ~one line.
- **Boundary:** the CDN lives *only* in `cmd/preview`. Products would use the §8
  vendored-binary build (out of scope here). The two never mix.

### Makefile

Add one target (mirrors existing self-documenting style):

```make
## @(app) - 🎨 Serve the component gallery on :8080
preview:
	@echo "✨📦✨ Serving component gallery on :8080\n"
	@go run ./cmd/preview
```

## Testing

- Each component is a `dom.Node`; test by rendering to a `bytes.Buffer` and asserting on
  the HTML string with `internal/assert` (mirrors `dom/dom_test.go`).
- Assert that each variant `switch` produces the expected **full class literal** — this
  guards against accidentally introducing computed class names.
- One `_test.go` per source file. Cover: `StyleVars` emits each `--var`; `Class` trims/
  drops/joins; every variant of every component; `FieldError("")` renders nothing;
  `Select` marks the selected option.

## Build order (steps 1–7)

1. `theme.go` — `Tokens`, `StyleVars`, `Class` (+ test).
2. `button.go` — `Button`, `ButtonLink`, `ButtonVariant`, `buttonClasses` (+ test).
   Exercises the variant + token pattern end to end.
3. `cmd/preview` — stand up the Play-CDN gallery skeleton + `preview` Makefile target, so
   every subsequent component is visible the moment it's written.
4. `typography.go` — `Heading`, `Text`, `Link` (+ test).
5. `layout.go` — `Container`, `Section`, `Card`, `Stack` (+ test).
6. `form.go` — `Label`, `FieldError`, `TextField`, then the rest (+ test).
7. `feedback.go` — `Badge`, `Alert` (+ test).
8. Fill `gallery()` with all components; run `go test ./...` and `make preview` as a
   smoke check.

## Success criteria

- `go test ./...` passes; every variant's class literal is asserted.
- `make preview` serves a gallery rendering all ~16 components and all variants.
- `?theme=apron` vs `?theme=teddy` visibly swaps colors and font — confirming the
  CSS-variable token approach works before any real product exists.
- No runtime-constructed Tailwind class names anywhere in `ui`.
