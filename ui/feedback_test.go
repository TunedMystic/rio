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
		{BadgeNeutral, "bg-[var(--color-text)]/8 text-[var(--color-text)]"},
		{BadgeSuccess, "bg-[var(--color-success)]/12 text-[var(--color-success)]"},
		{BadgeWarning, "bg-[var(--color-warning)]/12 text-[var(--color-warning)]"},
		{BadgeDanger, "bg-[var(--color-danger)]/12 text-[var(--color-danger)]"},
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
