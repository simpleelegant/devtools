package dataconvert

import (
	"regexp"
	"strings"
)

var indulgentInvalidGoFieldNameChar = regexp.MustCompile("[^a-z0-9]+")

func toGoFieldName(s string) string {
	return strings.Replace(
		strings.Title(
			indulgentInvalidGoFieldNameChar.ReplaceAllString(
				strings.ToLower(s), " ")),
		" ",
		"",
		-1)
}

func multipleString(s string, multiple int) string {
	if multiple <= 0 {
		return ""
	}

	b := ""
	for i := 0; i < multiple; i++ {
		b += s
	}

	return b
}

func insureLen(s string, min int) string {
	return s + multipleString(" ", min-len(s))
}
