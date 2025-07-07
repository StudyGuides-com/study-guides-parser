package preparser

// ParsedValue represents all possible parsed result types
type ParsedValue struct {
	Question   *QuestionResult `json:"question,omitempty"`
	Header     *HeaderResult   `json:"header,omitempty"`
	Comment    *CommentResult  `json:"comment,omitempty"`
	Empty      *EmptyLineResult `json:"empty,omitempty"`
	FileHeader *FileHeaderResult `json:"file_header,omitempty"`
	Passage    *PassageResult   `json:"passage,omitempty"`
	LearnMore  *LearnMoreResult `json:"learn_more,omitempty"`
	Content    *ContentResult   `json:"content,omitempty"`
	Binary     *BinaryResult    `json:"binary,omitempty"`
}

// GetQuestion returns the QuestionResult if this is a question, nil otherwise
func (pv ParsedValue) GetQuestion() *QuestionResult {
	return pv.Question
}

// GetHeader returns the HeaderResult if this is a header, nil otherwise
func (pv ParsedValue) GetHeader() *HeaderResult {
	return pv.Header
}

// GetComment returns the CommentResult if this is a comment, nil otherwise
func (pv ParsedValue) GetComment() *CommentResult {
	return pv.Comment
}

// GetEmpty returns the EmptyLineResult if this is an empty line, nil otherwise
func (pv ParsedValue) GetEmpty() *EmptyLineResult {
	return pv.Empty
}

// GetFileHeader returns the FileHeaderResult if this is a file header, nil otherwise
func (pv ParsedValue) GetFileHeader() *FileHeaderResult {
	return pv.FileHeader
}

// GetPassage returns the PassageResult if this is a passage, nil otherwise
func (pv ParsedValue) GetPassage() *PassageResult {
	return pv.Passage
}

// GetLearnMore returns the LearnMoreResult if this is a learn more, nil otherwise
func (pv ParsedValue) GetLearnMore() *LearnMoreResult {
	return pv.LearnMore
}

// GetContent returns the ContentResult if this is content, nil otherwise
func (pv ParsedValue) GetContent() *ContentResult {
	return pv.Content
}

// GetBinary returns the BinaryResult if this is binary, nil otherwise
func (pv ParsedValue) GetBinary() *BinaryResult {
	return pv.Binary
}

// IsQuestion returns true if this contains a QuestionResult
func (pv ParsedValue) IsQuestion() bool {
	return pv.Question != nil
}

// IsHeader returns true if this contains a HeaderResult
func (pv ParsedValue) IsHeader() bool {
	return pv.Header != nil
}

// IsComment returns true if this contains a CommentResult
func (pv ParsedValue) IsComment() bool {
	return pv.Comment != nil
}

// IsEmpty returns true if this contains an EmptyLineResult
func (pv ParsedValue) IsEmpty() bool {
	return pv.Empty != nil
}

// IsFileHeader returns true if this contains a FileHeaderResult
func (pv ParsedValue) IsFileHeader() bool {
	return pv.FileHeader != nil
}

// IsPassage returns true if this contains a PassageResult
func (pv ParsedValue) IsPassage() bool {
	return pv.Passage != nil
}

// IsLearnMore returns true if this contains a LearnMoreResult
func (pv ParsedValue) IsLearnMore() bool {
	return pv.LearnMore != nil
}

// IsContent returns true if this contains a ContentResult
func (pv ParsedValue) IsContent() bool {
	return pv.Content != nil
}

// IsBinary returns true if this contains a BinaryResult
func (pv ParsedValue) IsBinary() bool {
	return pv.Binary != nil
}

type ParsedLineInfo struct {
	Number      int         // Line number in the file
	Text        string      // The actual text content
	Type        TokenType   // The type of line (empty, content, comment, question, header)
	ParsedValue ParsedValue // The parsed value of the line, type depends on TokenType
}
