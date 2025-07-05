package lexer

import (
	"strings"

	"github.com/StudyGuides-com/study-guides-parser/cleanstring"
	"github.com/StudyGuides-com/study-guides-parser/core/constants"
	"github.com/StudyGuides-com/study-guides-parser/core/utils"
)

type TokenClassifier func(string, int) (TokenType, *LexerError)

// isQuestion checks if a line is a question and validates its format.
// A valid question must:
//  1. Start with a list item prefix (number or bullet)
//  2. Contain the answer delimiter " - "
//
// Returns:
//   - TokenType: The type of line (Question if valid, empty string if not)
//   - *TokenizerError: Any validation errors found
func isQuestion(line string, lineNum int) (TokenType, *LexerError) {
	if !utils.ListItemPrefixRegex.MatchString(line) {
		return "", nil
	}
	if !strings.Contains(line, constants.AnswerDelimiter) {
		return TokenTypeQuestion, NewLexerError(
			CodeMissingAnswerDelimiter,
			"missing answer delimiter ' - '",
			LineInfo{Number: lineNum, Text: line, Type: TokenTypeQuestion},
		)
	}
	return TokenTypeQuestion, nil
}

// isHeader checks if a line is a header. A line is considered a header if:
//  1. It contains 2 or more colons (e.g., "Subject: Topic: Subtopic")
//  2. It's not a passage, question, or learn more line
//
// Returns:
//   - TokenType: The type of line (Header if valid, empty string if not)
//   - *LexerError: Any validation errors found
func isHeader(line string, lineNum int) (TokenType, *LexerError) {
	if lineType, _ := isPassage(line, lineNum); lineType != "" {
		return "", nil
	}
	if lineType, _ := isQuestion(line, lineNum); lineType != "" {
		return "", nil
	}
	if lineType, _ := isLearnMore(line, lineNum); lineType != "" {
		return "", nil
	}
	parts := strings.Split(line, constants.ColonDelimiter)
	if len(parts) >= constants.MinHeaderParts {
		return TokenTypeHeader, nil
	}
	return "", nil
}

// isComment checks if a line is a comment. A valid comment:
//  1. Starts with exactly one '#' character
//  2. Does not start with multiple '#' characters
//
// Returns:
//   - TokenType: The type of line (Comment if valid, empty string if not)
//   - *LexerError: Any validation errors found
func isComment(line string, lineNum int) (TokenType, *LexerError) {
	trimmed := strings.TrimSpace(line)
	if strings.HasPrefix(trimmed, constants.CommentPrefix) && !strings.HasPrefix(trimmed, constants.CommentDoublePrefix) {
		return TokenTypeComment, nil
	}
	return "", nil
}

// isEmpty checks if a line is empty (contains only whitespace).
//
// Returns:
//   - TokenType: The type of line (Empty if valid, empty string if not)
//   - *LexerError: Any validation errors found
func isEmpty(line string, lineNum int) (TokenType, *LexerError) {
	if cleanstring.New(line).IsEmpty() {
		return TokenTypeEmpty, nil
	}
	return "", nil
}

// isFileHeader checks if a line is a file header. A file header:
//  1. Must be the first line of the file (lineNum == 1)
//  2. Must not be a regular header
//  3. Must not be empty
//
// Returns:
//   - TokenType: The type of line (FileHeader if valid, empty string if not)
//   - *LexerError: Error if line 1 is not a file header, nil otherwise
func isFileHeader(line string, lineNum int) (TokenType, *LexerError) {
	if lineNum != constants.FirstLineNumber {
		return "", nil
	}
	// If it's empty, that's an error
	if cleanstring.New(line).IsEmpty() {
		return TokenTypeEmpty, NewLexerError(
			CodeMissingFileHeader,
			"first line cannot be empty",
			LineInfo{Number: lineNum, Text: line, Type: TokenTypeFileHeader},
		)
	}
	// If it's a regular header, it's not a file header
	if lineType, _ := isHeader(line, lineNum); lineType != "" {
		return TokenTypeHeader, NewLexerError(
			CodeMissingFileHeader,
			"first line must be a file header, not a regular header",
			LineInfo{Number: lineNum, Text: line, Type: TokenTypeHeader},
		)
	}
	return TokenTypeFileHeader, nil
}

// isPassage checks if a line is a passage header. A valid passage header:
//  1. Starts with "Passage:" (case insensitive)
//  2. May optionally be prefixed with "###" or "####"
//
// Returns:
//   - TokenType: The type of line (Passage if valid, empty string if not)
//   - *LexerError: Any validation errors found
func isPassage(line string, lineNum int) (TokenType, *LexerError) {
	if cleanstring.New(line).HasPrefix(constants.PassagePrefix) ||
		cleanstring.New(line).HasPrefix("### "+constants.PassagePrefix) ||
		cleanstring.New(line).HasPrefix("#### "+constants.PassagePrefix) {
		return TokenTypePassage, nil
	}
	return "", nil
}

// isLearnMore checks if a line is a "Learn More" line. A valid learn more line:
//  1. Starts with "Learn More:" (case insensitive)
//  2. May have any amount of whitespace after the colon
//
// Returns:
//   - TokenType: The type of line (LearnMore if valid, empty string if not)
//   - *LexerError: Any validation errors found
func isLearnMore(line string, lineNum int) (TokenType, *LexerError) {
	if cleanstring.New(line).HasPrefix(constants.LearnMorePrefix) {
		return TokenTypeLearnMore, nil
	}
	return "", nil
}

// isBinary checks if a line contains binary content by looking for null bytes
// or other non-printable characters that would indicate binary data.
func isBinary(line string, lineNum int) (TokenType, *LexerError) {
	for _, char := range line {
		if char < 32 && char != '\t' && char != '\n' && char != '\r' {
			return TokenTypeBinary, NewLexerError(
				CodeBinaryContent,
				"contains binary or non-printable characters",
				LineInfo{Number: lineNum, Text: line, Type: TokenTypeBinary},
			)
		}
	}
	return "", nil
}
