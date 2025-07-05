package utils

import (
	"strings"
	"unicode"
)

// CleanStringer is an interface for types that can return a cleaned string representation.
type CleanStringer interface {
	Clean() string
}

// CleanString is a string type that implements CleanStringer.
type CleanString string

// Clean returns the cleaned version of the string.
func (cs CleanString) Clean() string {
	return NormalizeText(string(cs), false)
}

// NormalizeText performs common text normalization operations:
// 1. Trims whitespace
// 2. Converts to lowercase (if requested)
// 3. Removes invisible characters
func NormalizeText(text string, toLower bool) string {
	trimmed := strings.TrimSpace(text)
	if toLower {
		trimmed = strings.ToLower(trimmed)
	}
	return removeInvisibleCharacters(trimmed)
}

// IsEmpty checks if a string is empty or contains only whitespace
func IsEmpty(text string) bool {
	return strings.TrimSpace(text) == ""
}

// HasPrefix checks if a string has a specific prefix (case-insensitive)
func HasPrefix(text, prefix string) bool {
	return strings.HasPrefix(strings.ToLower(strings.TrimSpace(text)), strings.ToLower(prefix))
}

// removeInvisibleCharacters removes invisible and formatting characters from the input text
// while preserving regular spaces and visible characters.
func removeInvisibleCharacters(text string) string {
	var result []rune
	for _, char := range text {
		// Keep regular spaces and visible characters
		if char == ' ' || (!unicode.Is(unicode.Cf, char) && !unicode.Is(unicode.Cs, char) && !unicode.Is(unicode.Zs, char)) {
			result = append(result, char)
		}
	}
	return string(result)
}
