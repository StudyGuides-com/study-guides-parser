package lexer

import (
	"github.com/StudyGuides-com/study-guides-parser/core/cleanstring"
)

// Lexer provides functionality to process individual lines of text, detecting their type and
// parsing their content accordingly.
type Lexer struct {
	classifiers []TokenClassifier
}

// NewLexer creates and returns a new instance of Lexer.
// The lexer is initialized with a set of classifiers that are used to identify different types of lines in the study guide.
//
// The classifiers are executed in the following order:
//  1. Binary content detection
//  2. File header detection (must be first line)
//  3. Empty line detection
//  4. Comment detection
//  5. Question detection
//  6. Header detection
//  7. Passage detection
//  8. Learn More line detection
func NewLexer() *Lexer {
	return &Lexer{
		classifiers: []TokenClassifier{
			isBinary,     // Check for binary content first
			isFileHeader, // Then file headers (must be first line)
			isComment,    // Then comments
			isQuestion,   // Then questions
			isHeader,     // Then headers
			isPassage,    // Then passages
			isLearnMore,  // Then learn more lines
			isEmpty,      // Empty lines last since they're the most generic
		},
	}
}

// ProcessLine processes a single line of text, determining its type and parsing
// its content. It returns a LineInfo struct containing the processed information
// and any errors that occurred during processing.
//
// The function:
//  1. Trims whitespace from the line
//  2. Attempts to detect the line type using registered classifiers
//  3. Falls back to content type if no specific type is detected
//  4. Parses the line using the appropriate parser for its type
//
// Parameters:
//   - line: The text line to process
//   - lineNum: The line number in the source file
//
// Returns:
//   - LineInfo: Information about the processed line
//   - *LexerError: Any error that occurred during processing
func (l *Lexer) ProcessLine(line string, lineNum int) (LineInfo, *LexerError) {
	cleaned := cleanstring.New(line).Clean()
	lineInfo := LineInfo{
		Number: lineNum,
		Text:   line,
	}

	// Try each classifier in order
	var tokenType TokenType
	var classifierErr *LexerError
	for _, classify := range l.classifiers {
		tokenType, classifierErr = classify(cleaned, lineNum)
		if tokenType != "" {
			lineInfo.Type = tokenType
			break
		}
	}

	if tokenType == TokenTypeBinary {
		lineInfo.Type = TokenTypeBinary
		return lineInfo, NewLexerError(CodeBinaryContent, "contains binary or non-printable characters", lineInfo)
	}

	// If no type was detected, it's content
	if tokenType == "" {
		tokenType = TokenTypeContent
	}

	lineInfo.Type = tokenType

	return lineInfo, classifierErr
}
