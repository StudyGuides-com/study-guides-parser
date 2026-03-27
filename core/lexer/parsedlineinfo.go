package lexer

type ParsedLineInfo struct {
	Number      int         `json:"number"`       // Line number in the file
	Text        string      `json:"text"`         // The actual text content
	Type        TokenType   `json:"type"`         // The type of line (empty, content, comment, question, header)
	ParsedValue interface{} `json:"parsed_value"` // The parsed value of the line, type depends on TokenType
}
