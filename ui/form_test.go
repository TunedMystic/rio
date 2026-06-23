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
	opts := []Option{{Value: "us", Label: "USA"}, {Value: "ca", Label: "Canada"}}
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
