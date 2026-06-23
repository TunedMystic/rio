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
