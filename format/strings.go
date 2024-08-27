package format

import "strings"

func Title(s string) string {
	var strs []string
	for _, item := range strings.Split(s, " ") {
		strs = append(strs, TitleFirst(item))
	}
	return strings.Join(strs, " ")
}

func TitleFirst(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}
