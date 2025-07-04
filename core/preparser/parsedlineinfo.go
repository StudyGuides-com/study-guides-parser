package preparser

type ParsedLineInfo struct {
	Number      int         // Line number in the file
	Text        string      // The actual text content
	Type        TokenType   // The type of line (empty, content, comment, question, header)
	ParsedValue interface{} // The parsed value of the line, type depends on TokenType
}
