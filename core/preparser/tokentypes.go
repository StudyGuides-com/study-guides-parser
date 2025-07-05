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

// ParseResult represents all possible parsing result types
type ParseResult interface {
	*QuestionResult | HeaderResult | *CommentResult | *EmptyLineResult | 
	*FileHeaderResult | *PassageResult | *LearnMoreResult | *ContentResult | *BinaryResult
}

type LineParser[T ParseResult] interface {
	Parse(lineInfo LineInfo) (T, *PreParsingError)
}

// HeaderResult represents the parsed result of a header line
type HeaderResult = []string

// QuestionParser parses question lines
type QuestionParser struct{}

// QuestionResult represents the parsed result of a question line
type QuestionResult struct {
	QuestionText string
	AnswerText   string
}

// HeaderParser parses header lines
type HeaderParser struct{}

// CommentParser parses comment lines
type CommentParser struct{}

// EmptyLineParser parses empty lines
type EmptyLineParser struct{}

// EmptyLineResult represents the parsed result of an empty line
type EmptyLineResult struct{}

// FileHeaderParser parses file header lines
type FileHeaderParser struct{}

// FileHeaderResult represents the parsed result of a file header line
type FileHeaderResult struct {
	Title string
}

// PassageParser parses passage lines
type PassageParser struct{}

// PassageResult represents the parsed result of a passage line
type PassageResult struct {
	Text string
}

// LearnMoreParser parses learn more lines
type LearnMoreParser struct{}

// LearnMoreResult represents the parsed result of a learn more line
type LearnMoreResult struct {
	Text string
}

// ContentParser parses content lines
type ContentParser struct{}

// ContentResult represents the parsed result of a content line
type ContentResult struct {
	Text string
}

// CommentResult represents the parsed result of a comment line
type CommentResult struct {
	Text string
}

// BinaryParser parses binary lines
type BinaryParser struct{}

// BinaryResult represents the parsed result of a binary line
type BinaryResult struct {
	Text string
}
