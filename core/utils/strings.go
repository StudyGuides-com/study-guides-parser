package utils

import (
	"strings"
	"unicode"
)

// NormalizeText performs common text normalization operations:
// 1. Trims whitespace
// 2. Converts to lowercase (if requested)
// 3. Removes invisible characters
func NormalizeText(text string, toLower bool) string {
	trimmed := strings.TrimSpace(text)
	if toLower {
		trimmed = strings.ToLower(trimmed)
	}
	return RemoveInvisibleCharacters(trimmed)
}

// IsEmpty checks if a string is empty or contains only whitespace
func IsEmpty(text string) bool {
	return strings.TrimSpace(text) == ""
}

// HasPrefix checks if a string has a specific prefix (case-insensitive)
func HasPrefix(text, prefix string) bool {
	return strings.HasPrefix(strings.ToLower(strings.TrimSpace(text)), strings.ToLower(prefix))
}

// RemoveInvisibleCharacters removes invisible and formatting characters from the input text
// while preserving regular spaces and visible characters.
func RemoveInvisibleCharacters(text string) string {
	var result []rune
	for _, char := range text {
		// Keep regular spaces and visible characters
		if char == ' ' || (!unicode.Is(unicode.Cf, char) && !unicode.Is(unicode.Cs, char) && !unicode.Is(unicode.Zs, char)) {
			result = append(result, char)
		}
	}
	return string(result)
} 