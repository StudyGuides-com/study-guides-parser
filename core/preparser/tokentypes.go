package preparser

import "github.com/StudyGuides-com/study-guides-parser/core/lexer"

// TokenType is imported from the lexer package
type TokenType = lexer.TokenType

// Re-export the constants for convenience
const (
	TokenTypeHeader     = lexer.TokenTypeHeader
	TokenTypeQuestion   = lexer.TokenTypeQuestion
	TokenTypeComment    = lexer.TokenTypeComment
	TokenTypeEmpty      = lexer.TokenTypeEmpty
	TokenTypeContent    = lexer.TokenTypeContent
	TokenTypeFileHeader = lexer.TokenTypeFileHeader
	TokenTypePassage    = lexer.TokenTypePassage
	TokenTypeMisc       = lexer.TokenTypeMisc
	TokenTypeLearnMore  = lexer.TokenTypeLearnMore
	TokenTypeSpacer     = lexer.TokenTypeSpacer
	TokenTypeBinary     = lexer.TokenTypeBinary
)

// HeaderResult represents the parsed result of a header line
type HeaderResult struct {
	Parts []string
}

// QuestionResult represents the parsed result of a question line
type QuestionResult struct {
	QuestionText string
	AnswerText   string
}

// EmptyLineResult represents the parsed result of an empty line
type EmptyLineResult struct{}

// FileHeaderResult represents the parsed result of a file header line
type FileHeaderResult struct {
	Title string
}

// PassageResult represents the parsed result of a passage line
type PassageResult struct {
	Text string
}

// LearnMoreResult represents the parsed result of a learn more line
type LearnMoreResult struct {
	Text string
}

// ContentResult represents the parsed result of a content line
type ContentResult struct {
	Text string
}

// CommentResult represents the parsed result of a comment line
type CommentResult struct {
	Text string
}

// BinaryResult represents the parsed result of a binary line
type BinaryResult struct {
	Text string
}
