package preparser

import (
	"fmt"
	"regexp"
	"strings"
)

// Regex pattern for list item prefixes (numbers or bullets)
var listItemPrefixRegex = regexp.MustCompile(`^(\d+\.|\*|\-)\s+`)

// Parse implements LineParser for QuestionParser
func (p *QuestionParser) Parse(lineInfo LineInfo) (interface{}, error) {
	// Trim leading spaces before checking for prefix
	trimmedLine := strings.TrimLeft(lineInfo.Text, " ")
	if !listItemPrefixRegex.MatchString(trimmedLine) {
		return nil, NewPreParsingError(CodeValidation, "question must start with a number or bullet point", lineInfo)
	}
	if !strings.Contains(lineInfo.Text, " - ") {
		return nil, NewPreParsingError(CodeValidation, "question must contain answer delimiter ' - '", lineInfo)
	}

	// Split into question and answer using the first occurrence of ' - '
	parts := strings.SplitN(lineInfo.Text, " - ", 2)
	if len(parts) != 2 {
		return nil, NewPreParsingError(CodeValidation, "invalid question format", lineInfo)
	}

	// Remove the prefix from the question
	questionText := strings.TrimSpace(listItemPrefixRegex.ReplaceAllString(parts[0], ""))
	answerText := strings.TrimSpace(parts[1])

	// Sanitize both texts
	questionText = RemoveInvisibleCharacters(questionText)
	answerText = RemoveInvisibleCharacters(answerText)

	return &QuestionResult{
		QuestionText: questionText,
		AnswerText:   answerText,
	}, nil
}

// Parse implements LineParser for HeaderParser
func (p *HeaderParser) Parse(lineInfo LineInfo) (interface{}, error) {
	// Split the line by the colon ":" character
	parts := strings.Split(lineInfo.Text, ":")

	if len(parts) < 3 {
		return nil, NewPreParsingError(CodeValidation, "header must contain at least two colons", lineInfo)
	}

	// Clean up each part by removing invisible characters and trimming whitespace
	cleanedParts := make([]string, len(parts))
	for i, part := range parts {
		cleanedParts[i] = RemoveInvisibleCharacters(strings.TrimSpace(part))
	}

	// Return the cleaned parts as a result
	return cleanedParts, nil
}

// Parse implements LineParser for CommentParser
func (p *CommentParser) Parse(lineInfo LineInfo) (interface{}, error) {
	if !strings.HasPrefix(lineInfo.Text, "#") || strings.HasPrefix(lineInfo.Text, "##") {
		return nil, NewPreParsingError(CodeValidation, "comment must start with exactly one #", lineInfo)
	}
	// Remove the # and sanitize
	text := RemoveInvisibleCharacters(strings.TrimSpace(strings.TrimPrefix(lineInfo.Text, "#")))
	return &CommentResult{
		Text: text,
	}, nil
}

// Parse implements LineParser for EmptyLineParser
func (p *EmptyLineParser) Parse(lineInfo LineInfo) (interface{}, error) {
	if strings.TrimSpace(lineInfo.Text) != "" {
		return nil, NewPreParsingError(CodeValidation, "line must be empty or contain only whitespace", lineInfo)
	}
	return &EmptyLineResult{}, nil
}

// Parse implements LineParser for FileHeaderParser
func (p *FileHeaderParser) Parse(lineInfo LineInfo) (interface{}, error) {
	if lineInfo.Number != 1 {
		return nil, NewPreParsingError(CodeValidation, "file header must be on line 1", lineInfo)
	}
	// If it's a regular header (with multiple colons), it's not a file header
	if strings.Count(lineInfo.Text, ":") >= 2 {
		return nil, NewPreParsingError(CodeValidation, "file header should not be a regular header", lineInfo)
	}
	// Sanitize the title
	title := RemoveInvisibleCharacters(strings.TrimSpace(lineInfo.Text))
	return &FileHeaderResult{
		Title: title,
	}, nil
}

// Parse implements LineParser for PassageParser
func (p *PassageParser) Parse(lineInfo LineInfo) (interface{}, error) {
	lowerLine := strings.ToLower(lineInfo.Text)

	// Find the index of "passage:" in the lowercased line
	passageIdx := strings.Index(lowerLine, "passage:")
	if passageIdx == -1 {
		return nil, NewPreParsingError(CodeValidation, "passage must contain 'Passage:'", lineInfo)
	}

	// Get everything after "passage:"
	rest := lineInfo.Text[passageIdx+len("passage:"):]
	text := RemoveInvisibleCharacters(strings.TrimSpace(rest))
	if text == "" {
		return nil, NewPreParsingError(CodeValidation, "passage must contain text after 'Passage:'", lineInfo)
	}
	return &PassageResult{
		Text: text,
	}, nil
}

// Parse implements LineParser for LearnMoreParser
func (p *LearnMoreParser) Parse(lineInfo LineInfo) (interface{}, error) {
	lowerLine := strings.ToLower(lineInfo.Text)
	if !strings.HasPrefix(lowerLine, "learn more:") {
		return nil, NewPreParsingError(CodeValidation, "learn more line must start with 'Learn More:'", lineInfo)
	}
	colonIdx := len("learn more:")
	rest := lineInfo.Text[colonIdx:]
	text := RemoveInvisibleCharacters(strings.TrimSpace(rest))
	if text == "" {
		return nil, NewPreParsingError(CodeValidation, "learn more line must contain text after 'Learn More:'", lineInfo)
	}
	return &LearnMoreResult{
		Text: text,
	}, nil
}

// Parse implements LineParser for ContentParser
func (p *ContentParser) Parse(lineInfo LineInfo) (interface{}, error) {
	// Content lines have no specific format requirements
	// Just sanitize the text
	text := RemoveInvisibleCharacters(strings.TrimSpace(lineInfo.Text))
	if text == "" {
		return nil, NewPreParsingError(CodeValidation, "content line must not be empty or whitespace only", lineInfo)
	}
	return &ContentResult{
		Text: text,
	}, nil
}

// Parse implements LineParser for BinaryParser
func (p *BinaryParser) Parse(lineInfo LineInfo) (interface{}, error) {
	// Binary lines contain non-printable characters
	// Return the original text as-is
	return &BinaryResult{
		Text: lineInfo.Text,
	}, nil
}

// GetParserForType returns the appropriate parser for a given line type
func GetParserForType(lineType TokenType, lineInfo LineInfo) (LineParser, error) {
	switch lineType {
	case TokenTypeQuestion:
		return &QuestionParser{}, nil
	case TokenTypeHeader:
		return &HeaderParser{}, nil
	case TokenTypeComment:
		return &CommentParser{}, nil
	case TokenTypeEmpty:
		return &EmptyLineParser{}, nil
	case TokenTypeFileHeader:
		return &FileHeaderParser{}, nil
	case TokenTypePassage:
		return &PassageParser{}, nil
	case TokenTypeLearnMore:
		return &LearnMoreParser{}, nil
	case TokenTypeContent:
		return &ContentParser{}, nil
	case TokenTypeBinary:
		return &BinaryParser{}, nil
	default:
		return nil, NewPreParsingError(CodeValidation, fmt.Sprintf("unknown line type: %v", lineType), lineInfo)
	}
}
