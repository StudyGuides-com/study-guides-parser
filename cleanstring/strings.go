package cleanstring

import (
	"strings"
	"unicode"
)

// CleanString is a string type that implements CleanStringer.
type CleanString string

func New(s string) CleanString {
	return CleanString(s)
}

func (s CleanString) String() string {
	return string(s)
}

// Clean returns the cleaned version of the string.
func (s CleanString) Clean() string {
	return normalizeText(string(s), false)
}

func (s CleanString) CleanLower() string {
	return normalizeText(string(s), true)
}

func (s CleanString) IsEmpty() bool {
	return strings.TrimSpace(string(s)) == ""
}

func (s CleanString) HasPrefix(prefix string) bool {
	return strings.HasPrefix(
		s.CleanLower(),
		New(prefix).CleanLower(),
	)
}

// NormalizeText performs common text normalization operations:
// 1. Trims whitespace
// 2. Converts to lowercase (if requested)
// 3. Removes invisible characters
func normalizeText(text string, toLower bool) string {
	trimmed := strings.TrimSpace(text)
	if toLower {
		trimmed = strings.ToLower(trimmed)
	}
	return removeInvisibleCharacters(trimmed)
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
