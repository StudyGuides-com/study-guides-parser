//go:build !prod

package cleanstring

import (
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected CleanString
	}{
		{
			name:     "empty string",
			input:    "",
			expected: CleanString(""),
		},
		{
			name:     "simple string",
			input:    "hello world",
			expected: CleanString("hello world"),
		},
		{
			name:     "string with spaces",
			input:    "  hello world  ",
			expected: CleanString("  hello world  "),
		},
		{
			name:     "string with special characters",
			input:    "hello\n\t\rworld",
			expected: CleanString("hello\n\t\rworld"),
		},
		{
			name:     "unicode string",
			input:    "你好世界",
			expected: CleanString("你好世界"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := New(tt.input)
			if result != tt.expected {
				t.Errorf("New(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCleanString_String(t *testing.T) {
	tests := []struct {
		name     string
		input    CleanString
		expected string
	}{
		{
			name:     "empty string",
			input:    CleanString(""),
			expected: "",
		},
		{
			name:     "simple string",
			input:    CleanString("hello world"),
			expected: "hello world",
		},
		{
			name:     "string with spaces",
			input:    CleanString("  hello world  "),
			expected: "  hello world  ",
		},
		{
			name:     "unicode string",
			input:    CleanString("你好世界"),
			expected: "你好世界",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.String()
			if result != tt.expected {
				t.Errorf("String() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestCleanString_Clean(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "simple string",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "string with leading/trailing spaces",
			input:    "  hello world  ",
			expected: "hello world",
		},
		{
			name:     "string with tabs and newlines",
			input:    "\t\nhello\tworld\n",
			expected: "hello\tworld",
		},
		{
			name:     "string with invisible characters",
			input:    "hello\u200Bworld\uFEFF", // zero-width space and BOM
			expected: "helloworld",
		},
		{
			name:     "string with unicode spaces",
			input:    "hello\u00A0world\u2000", // non-breaking space and en quad
			expected: "helloworld",
		},
		{
			name:     "string with format characters",
			input:    "hello\u200Dworld\u2060", // zero-width joiner and word joiner
			expected: "helloworld",
		},
		{
			name:     "preserves regular spaces",
			input:    "hello world test",
			expected: "hello world test",
		},
		{
			name:     "mixed case preserved",
			input:    "Hello World",
			expected: "Hello World",
		},
		{
			name:     "unicode string",
			input:    "你好世界",
			expected: "你好世界",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := New(tt.input).Clean()
			if result != tt.expected {
				t.Errorf("Clean() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestCleanString_CleanLower(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "simple string",
			input:    "Hello World",
			expected: "hello world",
		},
		{
			name:     "string with leading/trailing spaces",
			input:    "  Hello World  ",
			expected: "hello world",
		},
		{
			name:     "string with invisible characters",
			input:    "Hello\u200BWorld\uFEFF", // zero-width space and BOM
			expected: "helloworld",
		},
		{
			name:     "already lowercase",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "mixed case with special characters",
			input:    "HeLLo\tWoRld\n",
			expected: "hello\tworld",
		},
		{
			name:     "unicode string",
			input:    "Hello世界",
			expected: "hello世界",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := New(tt.input).CleanLower()
			if result != tt.expected {
				t.Errorf("CleanLower() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestCleanString_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "empty string",
			input:    "",
			expected: true,
		},
		{
			name:     "whitespace only",
			input:    "   ",
			expected: true,
		},
		{
			name:     "tabs and newlines only",
			input:    "\t\n\r",
			expected: true,
		},
		{
			name:     "invisible characters only",
			input:    "\u200B\uFEFF\u2000", // zero-width space, BOM, en quad
			expected: false,                // TrimSpace doesn't remove invisible characters
		},
		{
			name:     "non-empty string",
			input:    "hello",
			expected: false,
		},
		{
			name:     "string with content and spaces",
			input:    "  hello  ",
			expected: false,
		},
		{
			name:     "string with invisible characters and content",
			input:    "\u200Bhello\uFEFF",
			expected: false,
		},
		{
			name:     "unicode string",
			input:    "你好",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := New(tt.input).IsEmpty()
			if result != tt.expected {
				t.Errorf("IsEmpty() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCleanString_HasPrefix(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		prefix   string
		expected bool
	}{
		{
			name:     "exact match",
			input:    "hello world",
			prefix:   "hello",
			expected: true,
		},
		{
			name:     "case insensitive match",
			input:    "Hello World",
			prefix:   "hello",
			expected: true,
		},
		{
			name:     "prefix with spaces",
			input:    "  hello world",
			prefix:   "hello",
			expected: true,
		},
		{
			name:     "prefix with invisible characters",
			input:    "hello\u200Bworld",
			prefix:   "hello",
			expected: true,
		},
		{
			name:     "no match",
			input:    "hello world",
			prefix:   "world",
			expected: false,
		},
		{
			name:     "empty prefix",
			input:    "hello world",
			prefix:   "",
			expected: true,
		},
		{
			name:     "empty input",
			input:    "",
			prefix:   "hello",
			expected: false,
		},
		{
			name:     "both empty",
			input:    "",
			prefix:   "",
			expected: true,
		},
		{
			name:     "unicode prefix",
			input:    "你好世界",
			prefix:   "你好",
			expected: true,
		},
		{
			name:     "unicode prefix case insensitive",
			input:    "你好世界",
			prefix:   "HELLO",
			expected: false,
		},
		{
			name:     "prefix longer than input",
			input:    "hello",
			prefix:   "hello world",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := New(tt.input).HasPrefix(tt.prefix)
			if result != tt.expected {
				t.Errorf("HasPrefix(%q) = %v, want %v", tt.prefix, result, tt.expected)
			}
		})
	}
}

func TestRemoveInvisibleCharacters(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "no invisible characters",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "zero-width space",
			input:    "hello\u200Bworld",
			expected: "helloworld",
		},
		{
			name:     "byte order mark",
			input:    "\uFEFFhello world",
			expected: "hello world",
		},
		{
			name:     "en quad space",
			input:    "hello\u2000world",
			expected: "helloworld",
		},
		{
			name:     "em quad space",
			input:    "hello\u2001world",
			expected: "helloworld",
		},
		{
			name:     "en space",
			input:    "hello\u2002world",
			expected: "helloworld",
		},
		{
			name:     "em space",
			input:    "hello\u2003world",
			expected: "helloworld",
		},
		{
			name:     "three-per-em space",
			input:    "hello\u2004world",
			expected: "helloworld",
		},
		{
			name:     "four-per-em space",
			input:    "hello\u2005world",
			expected: "helloworld",
		},
		{
			name:     "six-per-em space",
			input:    "hello\u2006world",
			expected: "helloworld",
		},
		{
			name:     "figure space",
			input:    "hello\u2007world",
			expected: "helloworld",
		},
		{
			name:     "punctuation space",
			input:    "hello\u2008world",
			expected: "helloworld",
		},
		{
			name:     "thin space",
			input:    "hello\u2009world",
			expected: "helloworld",
		},
		{
			name:     "hair space",
			input:    "hello\u200Aworld",
			expected: "helloworld",
		},
		{
			name:     "zero-width joiner",
			input:    "hello\u200Dworld",
			expected: "helloworld",
		},
		{
			name:     "word joiner",
			input:    "hello\u2060world",
			expected: "helloworld",
		},
		{
			name:     "preserves regular space",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "multiple invisible characters",
			input:    "hello\u200B\uFEFF\u2000world",
			expected: "helloworld",
		},
		{
			name:     "unicode characters preserved",
			input:    "hello\u4F60\u597Dworld", // 你好
			expected: "hello\u4F60\u597Dworld",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := removeInvisibleCharacters(tt.input)
			if result != tt.expected {
				t.Errorf("removeInvisibleCharacters(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestNormalizeText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		toLower  bool
		expected string
	}{
		{
			name:     "empty string, no lowercase",
			input:    "",
			toLower:  false,
			expected: "",
		},
		{
			name:     "empty string, lowercase",
			input:    "",
			toLower:  true,
			expected: "",
		},
		{
			name:     "simple string, no lowercase",
			input:    "  Hello World  ",
			toLower:  false,
			expected: "Hello World",
		},
		{
			name:     "simple string, lowercase",
			input:    "  Hello World  ",
			toLower:  true,
			expected: "hello world",
		},
		{
			name:     "with invisible characters, no lowercase",
			input:    "  Hello\u200BWorld  ",
			toLower:  false,
			expected: "HelloWorld",
		},
		{
			name:     "with invisible characters, lowercase",
			input:    "  Hello\u200BWorld  ",
			toLower:  true,
			expected: "helloworld",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeText(tt.input, tt.toLower)
			if result != tt.expected {
				t.Errorf("normalizeText(%q, %v) = %q, want %q", tt.input, tt.toLower, result, tt.expected)
			}
		})
	}
}
