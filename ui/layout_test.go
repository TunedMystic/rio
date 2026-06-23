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
