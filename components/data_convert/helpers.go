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

func parseKeyValueLines(s string, each func(key, value string)) {
	for _, line := range strings.Split(s, "\n") {
		// skip empty line
		if strings.TrimSpace(line) == "" {
			continue
		}

		kv := strings.SplitN(line, "=", 2)
		if len(kv) != 2 {
			continue
		}

		kv[0] = strings.TrimSpace(kv[0])
		if kv[0] == "" {
			continue
		}

		each(kv[0], strings.TrimSpace(kv[1]))
	}
}
