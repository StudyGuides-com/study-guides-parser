package preparser

import (
	"strings"

	"github.com/StudyGuides-com/study-guides-parser/cleanstring"
	"github.com/StudyGuides-com/study-guides-parser/core/constants"
	"github.com/StudyGuides-com/study-guides-parser/core/utils"
)

// ParseQuestion parses question lines
func ParseQuestion(lineInfo LineInfo) (*QuestionResult, *PreParsingError) {
	// Trim leading spaces before checking for prefix
	trimmedLine := strings.TrimLeft(lineInfo.Text, " ")
	if !utils.ListItemPrefixRegex.MatchString(trimmedLine) {
		return nil, NewPreParsingError(CodeValidation, "question must start with a number or bullet point", lineInfo)
	}
	if !strings.Contains(lineInfo.Text, constants.AnswerDelimiter) {
		return nil, NewPreParsingError(CodeValidation, "question must contain answer delimiter ' - '", lineInfo)
	}

	// Split into question and answer using the first occurrence of ' - '
	parts := strings.SplitN(lineInfo.Text, constants.AnswerDelimiter, constants.QuestionAnswerParts)
	if len(parts) != constants.QuestionAnswerParts {
		return nil, NewPreParsingError(CodeValidation, "invalid question format", lineInfo)
	}

	// Remove the prefix from the question
	questionText := strings.TrimSpace(utils.ListItemPrefixRegex.ReplaceAllString(parts[0], ""))
	answerText := strings.TrimSpace(parts[1])

	// Sanitize both texts
	questionText = cleanstring.New(questionText).Clean()
	answerText = cleanstring.New(answerText).Clean()

	return &QuestionResult{
		QuestionText: questionText,
		AnswerText:   answerText,
	}, nil
}

// ParseHeader parses header lines
func ParseHeader(lineInfo LineInfo) (*HeaderResult, *PreParsingError) {
	// Split the line by the colon ":" character
	parts := strings.Split(lineInfo.Clean(), constants.ColonDelimiter)

	if len(parts) < constants.MinHeaderParts {
		return nil, NewPreParsingError(CodeValidation, "header must contain at least two colons", lineInfo)
	}

	// Clean up each part by removing invisible characters and trimming whitespace
	cleanedParts := make([]string, len(parts))
	for i, part := range parts {
		cleanedParts[i] = cleanstring.New(part).Clean()
	}

	// Return the cleaned parts as a result
	return &HeaderResult{Parts: cleanedParts}, nil
}

// ParseComment parses comment lines
func ParseComment(lineInfo LineInfo) (*CommentResult, *PreParsingError) {
	if !strings.HasPrefix(lineInfo.Text, constants.CommentPrefix) || strings.HasPrefix(lineInfo.Text, constants.CommentDoublePrefix) {
		return nil, NewPreParsingError(CodeValidation, "comment must start with exactly one #", lineInfo)
	}
	// Remove the # and sanitize
	text := cleanstring.New(strings.TrimPrefix(lineInfo.Text, constants.CommentPrefix)).Clean()
	return &CommentResult{
		Text: text,
	}, nil
}

// ParseEmptyLine parses empty lines
func ParseEmptyLine(lineInfo LineInfo) (*EmptyLineResult, *PreParsingError) {
	if !cleanstring.New(lineInfo.Text).IsEmpty() {
		return nil, NewPreParsingError(CodeValidation, "line must be empty or contain only whitespace", lineInfo)
	}
	return &EmptyLineResult{}, nil
}

// ParseFileHeader parses file header lines
func ParseFileHeader(lineInfo LineInfo) (*FileHeaderResult, *PreParsingError) {
	if lineInfo.Number != constants.FirstLineNumber {
		return nil, NewPreParsingError(CodeValidation, "file header must be on line 1", lineInfo)
	}
	// If it's a regular header (with multiple colons), it's not a file header
	if strings.Count(lineInfo.Text, constants.ColonDelimiter) >= constants.MinHeaderColons-1 {
		return nil, NewPreParsingError(CodeValidation, "file header should not be a regular header", lineInfo)
	}
	// Sanitize the title
	title := lineInfo.Clean()
	return &FileHeaderResult{
		Title: title,
	}, nil
}

// ParsePassage parses passage lines
func ParsePassage(lineInfo LineInfo) (*PassageResult, *PreParsingError) {
	lowerLine := strings.ToLower(lineInfo.Text)

	// Find the index of "passage:" in the lowercased line
	passageIdx := strings.Index(lowerLine, constants.PassagePrefix)
	if passageIdx == -1 {
		return nil, NewPreParsingError(CodeValidation, "passage must contain 'Passage:'", lineInfo)
	}

	// Get everything after "passage:"
	rest := lineInfo.Text[passageIdx+len(constants.PassagePrefix):]
	text := cleanstring.New(rest).Clean()
	if text == "" {
		return nil, NewPreParsingError(CodeValidation, "passage must contain text after 'Passage:'", lineInfo)
	}
	return &PassageResult{
		Text: text,
	}, nil
}

// ParseLearnMore parses learn more lines
func ParseLearnMore(lineInfo LineInfo) (*LearnMoreResult, *PreParsingError) {
	lowerLine := strings.ToLower(lineInfo.Text)
	if !strings.HasPrefix(lowerLine, constants.LearnMorePrefix) {
		return nil, NewPreParsingError(CodeValidation, "learn more line must start with 'Learn More:'", lineInfo)
	}
	colonIdx := len(constants.LearnMorePrefix)
	rest := lineInfo.Text[colonIdx:]
	text := cleanstring.New(rest).Clean()
	if text == "" {
		return nil, NewPreParsingError(CodeValidation, "learn more line must contain text after 'Learn More:'", lineInfo)
	}
	return &LearnMoreResult{
		Text: text,
	}, nil
}

// ParseContent parses content lines
func ParseContent(lineInfo LineInfo) (*ContentResult, *PreParsingError) {
	// Content lines have no specific format requirements
	// Just sanitize the text
	text := cleanstring.New(lineInfo.Text).Clean()
	if text == "" {
		return nil, NewPreParsingError(CodeValidation, "content line must not be empty or whitespace only", lineInfo)
	}
	return &ContentResult{
		Text: text,
	}, nil
}

// ParseBinary parses binary lines
func ParseBinary(lineInfo LineInfo) (*BinaryResult, *PreParsingError) {
	// Binary lines contain non-printable characters
	// Return the original text as-is
	return &BinaryResult{
		Text: lineInfo.Text,
	}, nil
}

// parserRegistry maps token types to their corresponding parser functions
var parserRegistry = map[TokenType]any{
	TokenTypeQuestion:   ParserFunc[*QuestionResult](ParseQuestion),
	TokenTypeHeader:     ParserFunc[*HeaderResult](ParseHeader),
	TokenTypeComment:    ParserFunc[*CommentResult](ParseComment),
	TokenTypeEmpty:      ParserFunc[*EmptyLineResult](ParseEmptyLine),
	TokenTypeFileHeader: ParserFunc[*FileHeaderResult](ParseFileHeader),
	TokenTypePassage:    ParserFunc[*PassageResult](ParsePassage),
	TokenTypeLearnMore:  ParserFunc[*LearnMoreResult](ParseLearnMore),
	TokenTypeContent:    ParserFunc[*ContentResult](ParseContent),
	TokenTypeBinary:     ParserFunc[*BinaryResult](ParseBinary),
}
