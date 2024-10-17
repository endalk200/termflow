package utils

import "strings"

func ToSnakeCase(str string) string {
	var result []rune
	for i, char := range str {
		if i > 0 && IsUpper(char) {
			result = append(result, '_')
		}
		result = append(result, char)
	}
	return strings.ToLower(string(result))
}

func IsUpper(char rune) bool {
	return char >= 'A' && char <= 'Z'
}
