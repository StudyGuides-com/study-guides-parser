package constants

// File structure constants
const (
	// FirstLineNumber is the expected line number for file headers
	FirstLineNumber = 1

	// MinHeaderColons is the minimum number of colons required for a header
	MinHeaderColons = 3

	// MinHeaderParts is the minimum number of parts after splitting by colons
	MinHeaderParts = 3

	// QuestionAnswerParts is the expected number of parts when splitting by " - "
	QuestionAnswerParts = 2
)

// String constants
const (
	// AnswerDelimiter is the delimiter used to separate questions from answers
	AnswerDelimiter = " - "

	// ColonDelimiter is the delimiter used in headers
	ColonDelimiter = ":"

	// CommentPrefix is the prefix for comment lines
	CommentPrefix = "#"

	// CommentDoublePrefix is the prefix that indicates a non-comment line
	CommentDoublePrefix = "##"

	// PassagePrefix is the prefix for passage lines
	PassagePrefix = "passage:"

	// LearnMorePrefix is the prefix for learn more lines
	LearnMorePrefix = "learn more:"
)
