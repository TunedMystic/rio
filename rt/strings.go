package rt

import "unicode/utf8"

// StrLen returns the length of s in rune count.
func StrLen(s string) int {
	return utf8.RuneCountInString(s)
}
