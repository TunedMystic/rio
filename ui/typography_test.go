package ui

import (
	"strings"
	"testing"

	"github.com/tunedmystic/rio/internal/assert"
)

func Test_Heading_LevelsAndSizes(t *testing.T) {
	cases := []struct {
		level   HeadingLevel
		tag     string
		sizeVar string
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
