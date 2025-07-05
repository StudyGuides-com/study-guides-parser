package preparser

import (
	"fmt"
)

type Preparser struct {
	ParserType string
	Lines      []LineInfo
}

// NewPreparser creates a new Preparser instance with the given lines and parser type.
// The preparser will process each line according to its detected type and return
// parsed results with semantic information.
func NewPreparser(lines []LineInfo, parserType string) *Preparser {
	return &Preparser{
		Lines:      lines,
		ParserType: parserType,
	}
}

// Parse processes all lines in the preparser and returns parsed line information.
// Each line is processed according to its detected type (question, header, content, etc.)
// and the result contains both the original line info and the parsed semantic data.
//
// Returns:
//   - []ParsedLineInfo: Array of parsed line information
//   - *PreParsingError: Any error that occurred during parsing
func (p *Preparser) Parse() ([]ParsedLineInfo, *PreParsingError) {
	parsedLines := make([]ParsedLineInfo, 0, len(p.Lines))
	
	for _, line := range p.Lines {
		result, err := p.parseLine(line)
		if err != nil {
			return nil, err
		}

		info := ParsedLineInfo{
			Number:      line.Number,
			Type:        line.Type,
			Text:        line.Text,
			ParsedValue: result,
		}
		parsedLines = append(parsedLines, info)
	}
	return parsedLines, nil
}

// parseLine handles the parsing of a single line based on its type
func (p *Preparser) parseLine(line LineInfo) (ParsedValue, *PreParsingError) {
	// Use the registry directly - this is the single source of truth
	if parser, exists := parserRegistry[line.Type]; exists {
		// The parser is already the correct type, just call it
		switch p := parser.(type) {
		case LineParser[*QuestionResult]:
			result, err := p.Parse(line)
			if err != nil {
				return ParsedValue{}, err
			}
			return ParsedValue{Question: result}, nil
		case LineParser[*HeaderResult]:
			result, err := p.Parse(line)
			if err != nil {
				return ParsedValue{}, err
			}
			return ParsedValue{Header: result}, nil
		case LineParser[*CommentResult]:
			result, err := p.Parse(line)
			if err != nil {
				return ParsedValue{}, err
			}
			return ParsedValue{Comment: result}, nil
		case LineParser[*EmptyLineResult]:
			result, err := p.Parse(line)
			if err != nil {
				return ParsedValue{}, err
			}
			return ParsedValue{Empty: result}, nil
		case LineParser[*FileHeaderResult]:
			result, err := p.Parse(line)
			if err != nil {
				return ParsedValue{}, err
			}
			return ParsedValue{FileHeader: result}, nil
		case LineParser[*PassageResult]:
			result, err := p.Parse(line)
			if err != nil {
				return ParsedValue{}, err
			}
			return ParsedValue{Passage: result}, nil
		case LineParser[*LearnMoreResult]:
			result, err := p.Parse(line)
			if err != nil {
				return ParsedValue{}, err
			}
			return ParsedValue{LearnMore: result}, nil
		case LineParser[*ContentResult]:
			result, err := p.Parse(line)
			if err != nil {
				return ParsedValue{}, err
			}
			return ParsedValue{Content: result}, nil
		case LineParser[*BinaryResult]:
			result, err := p.Parse(line)
			if err != nil {
				return ParsedValue{}, err
			}
			return ParsedValue{Binary: result}, nil
		}
	}
	
	// If we get here, no parser was found
	return ParsedValue{}, NewPreParsingError(CodeValidation, fmt.Sprintf("unknown line type: %v", line.Type), line)
}
