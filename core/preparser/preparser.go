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
	switch line.Type {
	case TokenTypeQuestion:
		result, err := ParseQuestion(line)
		if err != nil {
			return ParsedValue{}, err
		}
		return ParsedValue{Question: result}, nil

	case TokenTypeHeader:
		result, err := ParseHeader(line)
		if err != nil {
			return ParsedValue{}, err
		}
		return ParsedValue{Header: result}, nil

	case TokenTypeComment:
		result, err := ParseComment(line)
		if err != nil {
			return ParsedValue{}, err
		}
		return ParsedValue{Comment: result}, nil

	case TokenTypeEmpty:
		result, err := ParseEmptyLine(line)
		if err != nil {
			return ParsedValue{}, err
		}
		return ParsedValue{Empty: result}, nil

	case TokenTypeFileHeader:
		result, err := ParseFileHeader(line)
		if err != nil {
			return ParsedValue{}, err
		}
		return ParsedValue{FileHeader: result}, nil

	case TokenTypePassage:
		result, err := ParsePassage(line)
		if err != nil {
			return ParsedValue{}, err
		}
		return ParsedValue{Passage: result}, nil

	case TokenTypeLearnMore:
		result, err := ParseLearnMore(line)
		if err != nil {
			return ParsedValue{}, err
		}
		return ParsedValue{LearnMore: result}, nil

	case TokenTypeContent:
		result, err := ParseContent(line)
		if err != nil {
			return ParsedValue{}, err
		}
		return ParsedValue{Content: result}, nil

	case TokenTypeBinary:
		result, err := ParseBinary(line)
		if err != nil {
			return ParsedValue{}, err
		}
		return ParsedValue{Binary: result}, nil

	default:
		return ParsedValue{}, NewPreParsingError(CodeValidation, fmt.Sprintf("unknown line type: %v", line.Type), line)
	}
}
