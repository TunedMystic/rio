package rt

import (
	"html/template"
	"unicode/utf8"
)

// StrLen returns the length of s in rune count.
func StrLen(s string) int {
	return utf8.RuneCountInString(s)
}

// DisplaySafeHTML converts a string into an HTML fragment, so that
// it can be rendered verbatim in the template.
func SafeHtml(content string) template.HTML {
	return template.HTML(content)
}
