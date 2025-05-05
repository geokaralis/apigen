package ts

import (
	"regexp"
	"strings"
)

type Formatter struct {
	indentWidth    int
	removeSpaces   bool
	prettyNewlines bool
}

func NewFormatter(indentWidth int, removeSpaces, prettyNewlines bool) *Formatter {
	return &Formatter{
		indentWidth:    indentWidth,
		removeSpaces:   removeSpaces,
		prettyNewlines: prettyNewlines,
	}
}

func (f *Formatter) Format(code string) string {
	indent := strings.Repeat(" ", f.indentWidth)

	fourSpace := regexp.MustCompile(`(?m)^(    )+`)
	code = fourSpace.ReplaceAllStringFunc(code, func(s string) string {
		count := len(s) / 4
		return strings.Repeat(indent, count)
	})

	if f.removeSpaces {
		code = strings.ReplaceAll(code, " : ", ": ")

		code = strings.ReplaceAll(code, "{ }", "{}")

		code = strings.ReplaceAll(code, "{ ", "{")
		code = strings.ReplaceAll(code, " }", "}")

		code = strings.ReplaceAll(code, "? :", "?:")

		code = strings.ReplaceAll(code, ", ", ",")
		code = strings.ReplaceAll(code, ",", ", ")
	}

	if f.prettyNewlines {
		types := regexp.MustCompile(`}\n\nexport type`)
		code = types.ReplaceAllString(code, "}\n\nexport type")

		interfaces := regexp.MustCompile(`}\n\nexport interface`)
		code = interfaces.ReplaceAllString(code, "}\n\nexport interface")
	}

	return code
}

func (f *Formatter) FormatAsync(code string) string {
	if f.removeSpaces {
		asyncMethods := regexp.MustCompile(`async\s+([a-zA-Z0-9_]+)\s*\(`)
		code = asyncMethods.ReplaceAllString(code, "async $1(")

		promiseReturn := regexp.MustCompile(`\):\s+Promise<`)
		code = promiseReturn.ReplaceAllString(code, "): Promise<")
	}

	return code
}
