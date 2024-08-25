package format

import "strings"

func Title(s string) string {
	var strs []string
	for _, item := range strings.Split(s, " ") {
		strs = append(strs, strings.ToUpper(item[:1])+item[1:])
	}
	return strings.Join(strs, " ")
}
