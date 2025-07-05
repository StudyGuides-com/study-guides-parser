package regexes

import (
	"testing"
)

func TestListItemPrefixRegex_MatchString(t *testing.T) {
	tests := []struct {
		input   string
		matches bool
	}{
		{"1. Question text", true},
		{"12. Another question", true},
		{"* Bullet point", true},
		{"- Dash bullet", true},
		{" 1. Not a match (leading space)", false},
		{"No prefix here", false},
		{"1.2. Not a match", false}, // Only matches the first prefix, not subsequent ones
		{"*Not a match (no space)", false},
		{"-Not a match (no space)", false},
		{"2) Not a match (wrong symbol)", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := ListItemPrefixRegex.MatchString(tt.input); got != tt.matches {
				t.Errorf("MatchString(%q) = %v, want %v", tt.input, got, tt.matches)
			}
		})
	}
}

func TestListItemPrefixRegex_ReplaceAllString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1. Question text", "Question text"},
		{"12. Another question", "Another question"},
		{"* Bullet point", "Bullet point"},
		{"- Dash bullet", "Dash bullet"},
		{"No prefix here", "No prefix here"},
		{" 1. Not a match (leading space)", " 1. Not a match (leading space)"},
		{"*Not a match (no space)", "*Not a match (no space)"},
		{"-Not a match (no space)", "-Not a match (no space)"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := ListItemPrefixRegex.ReplaceAllString(tt.input, ""); got != tt.expected {
				t.Errorf("ReplaceAllString(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
