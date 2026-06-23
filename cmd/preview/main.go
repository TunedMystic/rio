// Command preview is a dev-only gallery for viewing the rio/ui components
// while developing the library. It uses the TailwindCSS v4 Play CDN (which
// compiles classes from the rendered DOM at runtime) so there is no build
// step. It never ships to production — products use the vendored-binary build.
package main

import (
	"net/http"

	"github.com/tunedmystic/rio/dom"
	"github.com/tunedmystic/rio/ui"
)

// themes are demo token sets so the CSS variables resolve. Swap via ?theme=.
var themes = map[string]ui.Tokens{
	"apron": {
		FontFamily:        `"Inter", ui-sans-serif, system-ui, sans-serif`,
		FontSizeSm:        "16px",
		FontSizeBase:      "18px",
		FontSizeLg:        "20px",
		FontSizeXl:        "24px",
		FontSize2xl:       "30px",
		ColorPrimary:      "#059669",
		OnPrimary:         "#ffffff",
		ColorSecondary:    "#475569",
		OnSecondary:       "#ffffff",
		ColorBackground:   "#f8fafc",
		ColorSurface:      "#ffffff",
		ColorText:         "#0f172a",
		ColorTextMuted:    "#64748b",
		ColorBorder:       "#e2e8f0",
		ColorSuccess:      "#16a34a",
		ColorWarning:      "#d97706",
		ColorDanger:       "#dc2626",
		ColorInfo:         "#2563eb",
		RadiusBase:        "0.5rem",
		FontWeightHeading: "700",
	},
	"teddy": {
		FontFamily:        `"Plus Jakarta Sans", ui-sans-serif, system-ui, sans-serif`,
		FontSizeSm:        "16px",
		FontSizeBase:      "18px",
		FontSizeLg:        "20px",
		FontSizeXl:        "24px",
		FontSize2xl:       "30px",
		ColorPrimary:      "#4f46e5",
		OnPrimary:         "#ffffff",
		ColorSecondary:    "#64748b",
		OnSecondary:       "#ffffff",
		ColorBackground:   "#ffffff",
		ColorSurface:      "#f8fafc",
		ColorText:         "#1e1b2e",
		ColorTextMuted:    "#6b7280",
		ColorBorder:       "#e5e7eb",
		ColorSuccess:      "#16a34a",
		ColorWarning:      "#d97706",
		ColorDanger:       "#dc2626",
		ColorInfo:         "#2563eb",
		RadiusBase:        "1rem",
		FontWeightHeading: "650",
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
				dom.Link(dom.Rel("preconnect"), dom.Href("https://fonts.googleapis.com")),
				dom.Link(dom.Rel("preconnect"), dom.Href("https://fonts.gstatic.com"), dom.Crossorigin("anonymous")),
				dom.Link(dom.Rel("stylesheet"), dom.Href("https://fonts.googleapis.com/css2?family=Inter:wght@400..700&family=Plus+Jakarta+Sans:wght@400..700&display=swap")),
				tokens.StyleVars(),
				dom.Script(dom.Src("https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4")),
			),
			dom.Body(
				dom.Class("bg-[var(--color-background)] text-[var(--color-text)] font-[family-name:var(--font-family)] text-[length:var(--font-size-base)] leading-relaxed antialiased p-8"),
				gallery(),
			),
		))
		_ = page.Render(w)
	})

	_ = http.ListenAndServe(":8080", nil)
}

// row lays out intrinsically-sized items (buttons, badges) horizontally with
// wrapping. A one-off gallery layout, so it uses dom.Div directly per the
// library's "use dom.Div for one-off layout" guidance.
func row(children ...dom.Node) dom.Node {
	return dom.Div(append([]dom.Node{dom.Class("flex flex-wrap items-center gap-3")}, children...)...)
}

// section groups a heading with its demo body and even internal spacing.
func section(title string, body ...dom.Node) dom.Node {
	return dom.Div(append([]dom.Node{
		dom.Class("flex flex-col gap-4"),
		ui.Heading(ui.H2, title),
	}, body...)...)
}

// gallery renders one of every component with all its variants.
func gallery() dom.Node {
	return ui.Container(
		dom.Div(
			dom.Class("flex flex-col gap-12 py-10 max-w-3xl"),
			ui.Heading(ui.H1, "rio/ui component gallery"),

			section("Typography",
				ui.Heading(ui.H3, "Heading level 3"),
				ui.Text(ui.TextDefault, "Default body text in the product font."),
				ui.Text(ui.TextMuted, "Muted secondary text."),
				ui.Link("#", "An inline link"),
			),

			section("Buttons",
				row(
					ui.Button(ui.ButtonPrimary, "Primary"),
					ui.Button(ui.ButtonSecondary, "Secondary"),
					ui.Button(ui.ButtonDanger, "Danger"),
					ui.Button(ui.ButtonGhost, "Ghost"),
					ui.ButtonLink(ui.ButtonPrimary, "#", "Button Link"),
				),
			),

			section("Card",
				ui.Card(
					ui.Heading(ui.H3, "A card"),
					ui.Text(ui.TextMuted, "Raised surface with border and padding."),
				),
			),

			section("Form",
				ui.Card(
					ui.TextField("email", "Email", "", ""),
					ui.TextField("user", "Username", "taken", "That username is taken"),
					ui.Textarea("bio", "Bio", "", ""),
					ui.Select("country", "Country",
						[]ui.Option{{Value: "us", Label: "USA"}, {Value: "ca", Label: "Canada"}}, "ca", ""),
					row(
						ui.Checkbox("agree", "I agree to the terms", true),
						ui.Radio("plan", "Pro plan", "pro", true),
					),
				),
			),

			section("Badges",
				row(
					ui.Badge(ui.BadgeNeutral, "Neutral"),
					ui.Badge(ui.BadgeSuccess, "Success"),
					ui.Badge(ui.BadgeWarning, "Warning"),
					ui.Badge(ui.BadgeDanger, "Danger"),
				),
			),

			section("Alerts",
				ui.Stack(ui.GapMd,
					ui.Alert(ui.AlertInfo, ui.Text(ui.TextDefault, "Informational message.")),
					ui.Alert(ui.AlertSuccess, ui.Text(ui.TextDefault, "It worked.")),
					ui.Alert(ui.AlertWarning, ui.Text(ui.TextDefault, "Careful now.")),
					ui.Alert(ui.AlertError, ui.Text(ui.TextDefault, "Something went wrong.")),
				),
			),
		),
	)
}
