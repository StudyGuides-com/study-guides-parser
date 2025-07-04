package lexer

type TokenType string

const (
	// TokenTypeHeader represents a header line (e.g., "Unit 1: Primitive Types")
	TokenTypeHeader TokenType = "header"
	// TokenTypeQuestion represents a question line (e.g., "1. What is a variable? - A named storage location")
	TokenTypeQuestion TokenType = "question"
	// TokenTypeComment represents a comment line (e.g., "# This is a comment")
	TokenTypeComment TokenType = "comment"
	// TokenTypeEmpty represents an empty line
	TokenTypeEmpty TokenType = "empty"
	// TokenTypeContent represents a content line (not a header, question, comment, or empty)
	TokenTypeContent TokenType = "content"
	// TokenTypeFileHeader represents the first line of the file (e.g., "Advanced Placement (AP): AP Computer Science A Exam")
	TokenTypeFileHeader TokenType = "file_header"
	// TokenTypePassage represents a passage line (e.g., "Passage: ...")
	TokenTypePassage TokenType = "passage"
	// TokenTypeMisc represents miscellaneous lines (e.g., "Questions")
	TokenTypeMisc TokenType = "misc"
	// TokenTypeLearnMore represents a "Learn More" line
	TokenTypeLearnMore TokenType = "learn_more"
	// TokenTypeSpacer represents a spacer line (e.g., "---")
	TokenTypeSpacer TokenType = "spacer"
	// TokenTypeBinary represents a line containing binary or non-text content
	TokenTypeBinary TokenType = "binary"
)
