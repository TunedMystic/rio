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
