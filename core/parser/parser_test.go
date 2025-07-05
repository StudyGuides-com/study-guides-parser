//go:build !prod

package parser

import (
	"testing"

	"github.com/StudyGuides-com/study-guides-parser/core/lexer"
	"github.com/StudyGuides-com/study-guides-parser/core/preparser"
)

func TestParser(t *testing.T) {
	tests := []struct {
		name       string
		parserType ParserType
		lines      []preparser.ParsedLineInfo
		wantErr    bool
	}{
		{
			name:       "valid AP exam study guide",
			parserType: "APExams",
			lines: []preparser.ParsedLineInfo{
				{
					Number: 1,
					Type:   lexer.TokenTypeFileHeader,
					ParsedValue: preparser.ParsedValue{FileHeader: &preparser.FileHeaderResult{
						Title: "AP® African American Studies",
					}},
				},
				{
					Number: 2,
					Type:   lexer.TokenTypeHeader,
					ParsedValue: preparser.ParsedValue{Header: &preparser.HeaderResult{Parts: []string{
						"Advanced Placement (AP)",
						"AP African American Studies",
						"Origins of the African Diaspora (900 BCE - 16th Century)",
						"Introduction to African American Studies",
					}}},
				},
				{
					Number: 3,
					Type:   lexer.TokenTypeComment,
					ParsedValue: preparser.ParsedValue{Comment: &preparser.CommentResult{
						Text: "This is a comment",
					}},
				},
				{
					Number: 4,
					Type:   lexer.TokenTypeQuestion,
					ParsedValue: preparser.ParsedValue{Question: &preparser.QuestionResult{
						QuestionText: "What distinguishes African American Studies as an interdisciplinary field?",
						AnswerText:   "It integrates history, anthropology, sociology, literature, and political science to analyze the experiences of African-descended peoples.",
					}},
				},
			},
			wantErr: false,
		},
		{
			name:       "valid entrance exam study guide",
			parserType: "EntranceExams",
			lines: []preparser.ParsedLineInfo{
				{
					Number: 1,
					Type:   lexer.TokenTypeFileHeader,
					ParsedValue: preparser.ParsedValue{FileHeader: &preparser.FileHeaderResult{
						Title: "American College Testing (ACT)",
					}},
				},
				{
					Number: 2,
					Type:   lexer.TokenTypeHeader,
					ParsedValue: preparser.ParsedValue{Header: &preparser.HeaderResult{Parts: []string{
						"Entrance Exams",
						"American College Testing (ACT)",
						"English",
					}}},
				},
				{
					Number: 3,
					Type:   lexer.TokenTypePassage,
					ParsedValue: preparser.ParsedValue{Passage: &preparser.PassageResult{
						Text: "The Rise of Digital Libraries",
					}},
				},
				{
					Number: 4,
					Type:   lexer.TokenTypeContent,
					ParsedValue: preparser.ParsedValue{Content: &preparser.ContentResult{
						Text: "In the early 21st century, digital libraries emerged...",
					}},
				},
				{
					Number: 5,
					Type:   lexer.TokenTypeQuestion,
					ParsedValue: preparser.ParsedValue{Question: &preparser.QuestionResult{
						QuestionText: "At first, many scholars and researchers were skeptical...",
						AnswerText:   "No Change",
					}},
				},
			},
			wantErr: false,
		},
		{
			name:       "valid college study guide",
			parserType: "Colleges",
			lines: []preparser.ParsedLineInfo{
				{
					Number: 1,
					Type:   lexer.TokenTypeFileHeader,
					ParsedValue: preparser.ParsedValue{FileHeader: &preparser.FileHeaderResult{
						Title: "Principles of Financial Accounting",
					}},
				},
				{
					Number: 2,
					Type:   lexer.TokenTypeHeader,
					ParsedValue: preparser.ParsedValue{Header: &preparser.HeaderResult{Parts: []string{
						"Colleges",
						"Virginia",
						"Old Dominion University (ODU)",
						"Accounting (ACCT)",
						"ACCT 201 Principles of Financial Accounting",
						"Introduction to Financial Accounting",
					}}},
				},
				{
					Number: 3,
					Type:   lexer.TokenTypeQuestion,
					ParsedValue: preparser.ParsedValue{Question: &preparser.QuestionResult{
						QuestionText: "What is financial accounting?",
						AnswerText:   "Tracks financial transactions.",
					}},
				},
			},
			wantErr: false,
		},
		{
			name:       "invalid document - missing file header",
			parserType: "Colleges",
			lines: []preparser.ParsedLineInfo{
				{
					Number: 1,
					Type:   lexer.TokenTypeHeader,
					ParsedValue: preparser.ParsedValue{Header: &preparser.HeaderResult{Parts: []string{
						"Colleges",
						"Virginia",
						"Old Dominion University (ODU)",
					}}},
				},
			},
			wantErr: true,
		},
		{
			name:       "invalid document - empty",
			parserType: "Colleges",
			lines:      []preparser.ParsedLineInfo{},
			wantErr:    true,
		},
		{
			name:       "valid DOD study guide",
			parserType: "DOD",
			lines: []preparser.ParsedLineInfo{
				{
					Number: 1,
					Type:   lexer.TokenTypeFileHeader,
					ParsedValue: preparser.ParsedValue{FileHeader: &preparser.FileHeaderResult{
						Title: "COMNAVIDFORINST 1550.1, Navy Information Dominance Forces Language Readiness Program (January 2019)",
					}},
				},
				{
					Number: 2,
					Type:   lexer.TokenTypeHeader,
					ParsedValue: preparser.ParsedValue{Header: &preparser.HeaderResult{Parts: []string{
						"Department of Defense (DoD)",
						"United States Navy (USN)",
						"COMNAVIDFORINST",
						"COMNAVIDFORINST 1550.1",
						"COMNAVIDFORINST 1550.1, Navy Information Dominance Forces Language Readiness Program (January 2019)",
					}}},
				},
				{
					Number: 3,
					Type:   lexer.TokenTypeQuestion,
					ParsedValue: preparser.ParsedValue{Question: &preparser.QuestionResult{
						QuestionText: "What is the date of the COMNAVIFORINST 1550.1 instruction?",
						AnswerText:   "4 Jan 2019",
					}},
				},
				{
					Number: 4,
					Type:   lexer.TokenTypeLearnMore,
					ParsedValue: preparser.ParsedValue{LearnMore: &preparser.LearnMoreResult{
						Text: "COMNAVIFORINST 1550.1 4 Jan 2019 (Page 2)",
					}},
				},
				{
					Number: 5,
					Type:   lexer.TokenTypeQuestion,
					ParsedValue: preparser.ParsedValue{Question: &preparser.QuestionResult{
						QuestionText: "Who leads the CTI Community Management Network according to the instruction?",
						AnswerText:   "Commander, NAVIFOR",
					}},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.parserType, tt.lines)
			result, err := parser.Parse()

			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if result == nil {
					t.Error("Parser.Parse() returned nil result when no error expected")
					return
				}

				if result.ParserType != tt.parserType {
					t.Errorf("Parser.Parse() returned parser type = %v, want %v", result.ParserType, tt.parserType)
				}

				if result.Root == nil {
					t.Error("Parser.Parse() returned nil root node")
				}
			}
		})
	}
}
