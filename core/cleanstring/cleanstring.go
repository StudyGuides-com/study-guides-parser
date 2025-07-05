// Package cleanstring provides utilities for cleaning and normalizing text strings.
// It offers a CleanString type that wraps string operations with common text
// normalization features such as whitespace trimming, case conversion, and
// removal of invisible characters.
package cleanstring

import (
	"strings"
	"unicode"
)

// CleanString is a string type that implements text cleaning operations.
// It wraps a string value and provides methods for normalizing and cleaning text.
type CleanString string

// New creates a new CleanString from the given string value.
// This is the primary constructor for creating CleanString instances.
func New(s string) CleanString {
	return CleanString(s)
}

// String returns the underlying string value of the CleanString.
// This implements the fmt.Stringer interface.
func (s CleanString) String() string {
	return string(s)
}

// Clean returns the cleaned version of the string.
// This performs text normalization including whitespace trimming and
// removal of invisible characters while preserving the original case.
func (s CleanString) Clean() string {
	return normalizeText(string(s), false)
}

// CleanLower returns the cleaned version of the string converted to lowercase.
// This performs the same normalization as Clean() but also converts the text
// to lowercase, useful for case-insensitive comparisons.
func (s CleanString) CleanLower() string {
	return normalizeText(string(s), true)
}

// IsEmpty checks if the string is empty after trimming whitespace.
// Returns true if the string contains only whitespace characters or is empty.
func (s CleanString) IsEmpty() bool {
	return strings.TrimSpace(string(s)) == ""
}

// HasPrefix checks if the cleaned string starts with the given prefix.
// The comparison is case-insensitive and uses cleaned versions of both strings.
// This is useful for matching text that may have formatting differences.
func (s CleanString) HasPrefix(prefix string) bool {
	return strings.HasPrefix(
		s.CleanLower(),
		New(prefix).CleanLower(),
	)
}

// normalizeText performs common text normalization operations:
// 1. Trims leading and trailing whitespace
// 2. Converts to lowercase (if requested)
// 3. Removes invisible characters while preserving regular spaces
//
// Parameters:
//   - text: The input string to normalize
//   - toLower: Whether to convert the text to lowercase
//
// Returns the normalized string.
func normalizeText(text string, toLower bool) string {
	trimmed := strings.TrimSpace(text)
	if toLower {
		trimmed = strings.ToLower(trimmed)
	}
	return removeInvisibleCharacters(trimmed)
}

// removeInvisibleCharacters removes invisible and formatting characters from the input text
// while preserving regular spaces and visible characters.
//
// This function filters out:
// - Unicode Cf (Format) characters (e.g., zero-width spaces, directional markers)
// - Unicode Cs (Surrogate) characters (invalid UTF-16 surrogate pairs)
// - Unicode Zs (Space_Separator) characters (various types of spaces except regular space)
//
// It preserves:
// - Regular space characters (' ')
// - All visible characters
// - Other printable characters
//
// Parameters:
//   - text: The input string to clean
//
// Returns the cleaned string with invisible characters removed.
func removeInvisibleCharacters(text string) string {
	var result []rune
	for _, char := range text {
		// Keep regular spaces and visible characters
		// Filter out format, surrogate, and space separator characters
		if char == ' ' || (!unicode.Is(unicode.Cf, char) && !unicode.Is(unicode.Cs, char) && !unicode.Is(unicode.Zs, char)) {
			result = append(result, char)
		}
	}
	return string(result)
}
