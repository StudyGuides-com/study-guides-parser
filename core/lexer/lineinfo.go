package lexer

type LineInfo struct {
	Number int       // Line number in the file
	Text   string    // The actual text content
	Type   TokenType // The type of line (empty, content, comment, question, header)
}
