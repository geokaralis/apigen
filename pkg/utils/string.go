package utils

import (
	"strings"
)

func ToCamelCase(input string) string {
	parts := strings.FieldsFunc(input, func(r rune) bool {
		return r == '_' || r == '-'
	})

	for i, part := range parts {
		if i == 0 {
			parts[i] = strings.ToLower(part)
		} else if len(part) > 0 {
			parts[i] = strings.ToUpper(part[0:1]) + strings.ToLower(part[1:])
		}
	}

	return strings.Join(parts, "")
}

func ToPascalCase(input string) string {
	parts := strings.FieldsFunc(input, func(r rune) bool {
		return r == '_' || r == '-'
	})

	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[0:1]) + strings.ToLower(part[1:])
		}
	}

	return strings.Join(parts, "")
}
