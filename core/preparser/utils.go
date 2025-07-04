package preparser

import "unicode"

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
