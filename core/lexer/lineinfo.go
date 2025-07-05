package lexer

import "github.com/StudyGuides-com/study-guides-parser/core/utils"

type LineInfo struct {
	Number int       // Line number in the file
	Text   string    // The actual text content
	Type   TokenType // The type of line (empty, content, comment, question, header)
}

// Clean returns the cleaned version of the Text field.
func (li LineInfo) Clean() string {
	return utils.NormalizeText(li.Text, false)
}
