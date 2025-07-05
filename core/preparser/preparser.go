package preparser

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
		lineParser, err := GetParserForType(line.Type, line)
		if err != nil {
			return nil, err
		}
		result, err := lineParser.Parse(line)
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
