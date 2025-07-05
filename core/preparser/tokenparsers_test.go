//go:build !prod

package preparser

import (
	"testing"
)

func TestLineQuestionParser(t *testing.T) {
	tests := []struct {
		name     string
		lineInfo LineInfo
		want     *QuestionResult
		wantErr  bool
	}{
		{
			name: "valid question with number",
			lineInfo: LineInfo{
				Number: 1,
				Type:   TokenTypeQuestion,
				Text:   "1. What is Go? - A programming language",
			},
			want: &QuestionResult{
				QuestionText: "What is Go?",
				AnswerText:   "A programming language",
			},
			wantErr: false,
		},
		{
			name: "valid question with bullet",
			lineInfo: LineInfo{
				Number: 1,
				Type:   TokenTypeQuestion,
				Text:   "* What is Go? - A programming language",
			},
			want: &QuestionResult{
				QuestionText: "What is Go?",
				AnswerText:   "A programming language",
			},
			wantErr: false,
		},
		{
			name: "valid question with dash",
			lineInfo: LineInfo{
				Number: 1,
				Type:   TokenTypeQuestion,
				Text:   "- What is Go? - A programming language",
			},
			want: &QuestionResult{
				QuestionText: "What is Go?",
				AnswerText:   "A programming language",
			},
			wantErr: false,
		},
		{
			name: "valid question with extra spaces",
			lineInfo: LineInfo{
				Number: 1,
				Type:   TokenTypeQuestion,
				Text:   "1.   What is Go?   -   A programming language",
			},
			want: &QuestionResult{
				QuestionText: "What is Go?",
				AnswerText:   "A programming language",
			},
			wantErr: false,
		},
		{
			name: "invalid question no prefix",
			lineInfo: LineInfo{
				Number: 1,
				Type:   TokenTypeQuestion,
				Text:   "What is Go? - A programming language",
			},
			wantErr: true,
		},
		{
			name: "invalid question no delimiter",
			lineInfo: LineInfo{
				Number: 1,
				Type:   TokenTypeQuestion,
				Text:   "1. What is Go?",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseQuestion(tt.lineInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseQuestion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.QuestionText != tt.want.QuestionText {
					t.Errorf("ParseQuestion() question = %v, want %v", got.QuestionText, tt.want.QuestionText)
				}
				if got.AnswerText != tt.want.AnswerText {
					t.Errorf("ParseQuestion() answer = %v, want %v", got.AnswerText, tt.want.AnswerText)
				}
			}
		})
	}
}

func TestLineHeaderParser(t *testing.T) {
	tests := []struct {
		name     string
		lineInfo LineInfo
		want     []string
		wantErr  bool
	}{
		{
			name: "valid header with two colons",
			lineInfo: LineInfo{
				Number: 1,
				Type:   TokenTypeHeader,
				Text:   "Section: Chapter 1: Introduction",
			},
			want:    []string{"Section", "Chapter 1", "Introduction"},
			wantErr: false,
		},
		{
			name: "valid header with three colons",
			lineInfo: LineInfo{
				Number: 1,
				Type:   TokenTypeHeader,
				Text:   "Section: Chapter 1: Part 1: Introduction",
			},
			want:    []string{"Section", "Chapter 1", "Part 1", "Introduction"},
			wantErr: false,
		},
		{
			name: "valid header with extra spaces",
			lineInfo: LineInfo{
				Number: 1,
				Type:   TokenTypeHeader,
				Text:   "Section : Chapter 1 : Introduction",
			},
			want:    []string{"Section", "Chapter 1", "Introduction"},
			wantErr: false,
		},
		{
			name: "invalid header with one colon",
			lineInfo: LineInfo{
				Number: 1,
				Type:   TokenTypeHeader,
				Text:   "Section: Chapter 1",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseHeader(tt.lineInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseHeader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(got.Parts) != len(tt.want) {
					t.Errorf("ParseHeader() got %d parts, want %d", len(got.Parts), len(tt.want))
					return
				}
				for i, part := range got.Parts {
					if part != tt.want[i] {
						t.Errorf("ParseHeader() part[%d] = %v, want %v", i, part, tt.want[i])
					}
				}
			}
		})
	}
}

func TestLineCommentParser(t *testing.T) {
	tests := []struct {
		name     string
		lineInfo LineInfo
		want     *CommentResult
		wantErr  bool
	}{
		{
			name: "valid comment",
			lineInfo: LineInfo{
				Number: 1,
				Type:   TokenTypeComment,
				Text:   "# This is a comment",
			},
			want: &CommentResult{
				Text: "This is a comment",
			},
			wantErr: false,
		},
		{
			name: "valid comment with extra spaces",
			lineInfo: LineInfo{
				Number: 1,
				Type:   TokenTypeComment,
				Text:   "#   This is a comment   ",
			},
			want: &CommentResult{
				Text: "This is a comment",
			},
			wantErr: false,
		},
		{
			name: "invalid comment with double hash",
			lineInfo: LineInfo{
				Number: 1,
				Type:   TokenTypeComment,
				Text:   "## This is not a valid comment",
			},
			wantErr: true,
		},
		{
			name: "invalid comment without hash",
			lineInfo: LineInfo{
				Number: 1,
				Type:   TokenTypeComment,
				Text:   "This is not a comment",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseComment(tt.lineInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseComment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Text != tt.want.Text {
					t.Errorf("ParseComment() text = %v, want %v", got.Text, tt.want.Text)
				}
			}
		})
	}
}

func TestLinePassageParser(t *testing.T) {
	tests := []struct {
		name     string
		lineInfo LineInfo
		want     *PassageResult
		wantErr  bool
	}{
		{
			name: "valid passage",
			lineInfo: LineInfo{
				Number: 1,
				Type:   TokenTypePassage,
				Text:   "Passage: This is a passage text",
			},
			want: &PassageResult{
				Text: "This is a passage text",
			},
			wantErr: false,
		},
		{
			name: "valid passage with different case",
			lineInfo: LineInfo{
				Number: 1,
				Type:   TokenTypePassage,
				Text:   "PASSAGE: This is a passage text",
			},
			want: &PassageResult{
				Text: "This is a passage text",
			},
			wantErr: false,
		},
		{
			name: "valid passage with extra spaces",
			lineInfo: LineInfo{
				Number: 1,
				Type:   TokenTypePassage,
				Text:   "Passage:   This is a passage text   ",
			},
			want: &PassageResult{
				Text: "This is a passage text",
			},
			wantErr: false,
		},
		{
			name: "invalid passage no keyword",
			lineInfo: LineInfo{
				Number: 1,
				Type:   TokenTypePassage,
				Text:   "This is not a passage",
			},
			wantErr: true,
		},
		{
			name: "invalid passage empty text",
			lineInfo: LineInfo{
				Number: 1,
				Type:   TokenTypePassage,
				Text:   "Passage:",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePassage(tt.lineInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePassage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Text != tt.want.Text {
					t.Errorf("ParsePassage() text = %v, want %v", got.Text, tt.want.Text)
				}
			}
		})
	}
}

func TestLineLearnMoreParser(t *testing.T) {
	tests := []struct {
		name     string
		lineInfo LineInfo
		want     *LearnMoreResult
		wantErr  bool
	}{
		{
			name: "valid learn more",
			lineInfo: LineInfo{
				Number: 1,
				Type:   TokenTypeLearnMore,
				Text:   "Learn More: Additional information here",
			},
			want: &LearnMoreResult{
				Text: "Additional information here",
			},
			wantErr: false,
		},
		{
			name: "valid learn more with different case",
			lineInfo: LineInfo{
				Number: 1,
				Type:   TokenTypeLearnMore,
				Text:   "LEARN MORE: Additional information here",
			},
			want: &LearnMoreResult{
				Text: "Additional information here",
			},
			wantErr: false,
		},
		{
			name: "valid learn more with extra spaces",
			lineInfo: LineInfo{
				Number: 1,
				Type:   TokenTypeLearnMore,
				Text:   "Learn More:   Additional information here   ",
			},
			want: &LearnMoreResult{
				Text: "Additional information here",
			},
			wantErr: false,
		},
		{
			name: "invalid learn more no keyword",
			lineInfo: LineInfo{
				Number: 1,
				Type:   TokenTypeLearnMore,
				Text:   "Additional information here",
			},
			wantErr: true,
		},
		{
			name: "invalid learn more empty text",
			lineInfo: LineInfo{
				Number: 1,
				Type:   TokenTypeLearnMore,
				Text:   "Learn More:",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseLearnMore(tt.lineInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseLearnMore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Text != tt.want.Text {
					t.Errorf("ParseLearnMore() text = %v, want %v", got.Text, tt.want.Text)
				}
			}
		})
	}
}

func TestLineContentParser(t *testing.T) {
	tests := []struct {
		name     string
		lineInfo LineInfo
		want     *ContentResult
		wantErr  bool
	}{
		{
			name: "valid content",
			lineInfo: LineInfo{
				Number: 1,
				Type:   TokenTypeContent,
				Text:   "This is some content",
			},
			want: &ContentResult{
				Text: "This is some content",
			},
			wantErr: false,
		},
		{
			name: "valid content with extra spaces",
			lineInfo: LineInfo{
				Number: 1,
				Type:   TokenTypeContent,
				Text:   "   This is some content   ",
			},
			want: &ContentResult{
				Text: "This is some content",
			},
			wantErr: false,
		},
		{
			name: "invalid content empty",
			lineInfo: LineInfo{
				Number: 1,
				Type:   TokenTypeContent,
				Text:   "",
			},
			wantErr: true,
		},
		{
			name: "invalid content whitespace only",
			lineInfo: LineInfo{
				Number: 1,
				Type:   TokenTypeContent,
				Text:   "   \t  ",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseContent(tt.lineInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Text != tt.want.Text {
					t.Errorf("ParseContent() text = %v, want %v", got.Text, tt.want.Text)
				}
			}
		})
	}
}

func TestLineEmptyLineParser(t *testing.T) {
	tests := []struct {
		name     string
		lineInfo LineInfo
		wantErr  bool
	}{
		{
			name: "valid empty line",
			lineInfo: LineInfo{
				Number: 1,
				Type:   TokenTypeEmpty,
				Text:   "",
			},
			wantErr: false,
		},
		{
			name: "valid whitespace line",
			lineInfo: LineInfo{
				Number: 1,
				Type:   TokenTypeEmpty,
				Text:   "   \t  ",
			},
			wantErr: false,
		},
		{
			name: "invalid non-empty line",
			lineInfo: LineInfo{
				Number: 1,
				Type:   TokenTypeEmpty,
				Text:   "not empty",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseEmptyLine(tt.lineInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseEmptyLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// EmptyLineResult has no fields to check, just verify it's not nil
				if got == nil {
					t.Error("ParseEmptyLine() returned nil result")
				}
			}
		})
	}
}

func TestLineFileHeaderParser(t *testing.T) {
	tests := []struct {
		name     string
		lineInfo LineInfo
		want     *FileHeaderResult
		wantErr  bool
	}{
		{
			name: "valid file header",
			lineInfo: LineInfo{
				Number: 1,
				Type:   TokenTypeFileHeader,
				Text:   "Study Guide",
			},
			want: &FileHeaderResult{
				Title: "Study Guide",
			},
			wantErr: false,
		},
		{
			name: "valid file header with single colon",
			lineInfo: LineInfo{
				Number: 1,
				Type:   TokenTypeFileHeader,
				Text:   "Study Guide: Go Programming",
			},
			want: &FileHeaderResult{
				Title: "Study Guide: Go Programming",
			},
			wantErr: false,
		},
		{
			name: "valid file header with extra spaces",
			lineInfo: LineInfo{
				Number: 1,
				Type:   TokenTypeFileHeader,
				Text:   "   Study Guide   ",
			},
			want: &FileHeaderResult{
				Title: "Study Guide",
			},
			wantErr: false,
		},
		{
			name: "invalid file header not line 1",
			lineInfo: LineInfo{
				Number: 2,
				Type:   TokenTypeFileHeader,
				Text:   "Study Guide",
			},
			wantErr: true,
		},
		{
			name: "invalid file header is regular header",
			lineInfo: LineInfo{
				Number: 1,
				Type:   TokenTypeFileHeader,
				Text:   "Section: Chapter 1: Introduction",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFileHeader(tt.lineInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFileHeader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Title != tt.want.Title {
					t.Errorf("ParseFileHeader() title = %v, want %v", got.Title, tt.want.Title)
				}
			}
		})
	}
}
