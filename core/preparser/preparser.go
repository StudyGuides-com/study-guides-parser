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
	parsedLines := []ParsedLineInfo{}
	for _, line := range p.Lines {
		// Use type-safe parsing based on line type
		var result any
		var err *PreParsingError
		
		switch line.Type {
		case TokenTypeQuestion:
			parser, parseErr := GetParserForType[*QuestionResult](line.Type, line)
			if parseErr != nil {
				err = parseErr
			} else {
				result, err = parser.Parse(line)
			}
		case TokenTypeHeader:
			parser, parseErr := GetParserForType[HeaderResult](line.Type, line)
			if parseErr != nil {
				err = parseErr
			} else {
				result, err = parser.Parse(line)
			}
		case TokenTypeComment:
			parser, parseErr := GetParserForType[*CommentResult](line.Type, line)
			if parseErr != nil {
				err = parseErr
			} else {
				result, err = parser.Parse(line)
			}
		case TokenTypeEmpty:
			parser, parseErr := GetParserForType[*EmptyLineResult](line.Type, line)
			if parseErr != nil {
				err = parseErr
			} else {
				result, err = parser.Parse(line)
			}
		case TokenTypeFileHeader:
			parser, parseErr := GetParserForType[*FileHeaderResult](line.Type, line)
			if parseErr != nil {
				err = parseErr
			} else {
				result, err = parser.Parse(line)
			}
		case TokenTypePassage:
			parser, parseErr := GetParserForType[*PassageResult](line.Type, line)
			if parseErr != nil {
				err = parseErr
			} else {
				result, err = parser.Parse(line)
			}
		case TokenTypeLearnMore:
			parser, parseErr := GetParserForType[*LearnMoreResult](line.Type, line)
			if parseErr != nil {
				err = parseErr
			} else {
				result, err = parser.Parse(line)
			}
		case TokenTypeContent:
			parser, parseErr := GetParserForType[*ContentResult](line.Type, line)
			if parseErr != nil {
				err = parseErr
			} else {
				result, err = parser.Parse(line)
			}
		case TokenTypeBinary:
			parser, parseErr := GetParserForType[*BinaryResult](line.Type, line)
			if parseErr != nil {
				err = parseErr
			} else {
				result, err = parser.Parse(line)
			}
		default:
			return nil, NewPreParsingError(CodeValidation, fmt.Sprintf("unknown line type: %v", line.Type), line)
		}
		
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
