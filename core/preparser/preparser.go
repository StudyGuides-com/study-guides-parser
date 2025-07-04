package preparser

type Preparser struct {
	ParserType string
	Lines      []LineInfo
}

func NewPreparser(lines []LineInfo, parserType string) *Preparser {
	return &Preparser{
		Lines:      lines,
		ParserType: parserType,
	}
}

func (p *Preparser) Parse() ([]ParsedLineInfo, error) {
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
