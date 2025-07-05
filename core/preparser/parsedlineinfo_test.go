//go:build !prod

package preparser

import (
	"reflect"
	"testing"
)

func TestParsedValueGetters(t *testing.T) {
	tests := []struct {
		name     string
		value    ParsedValue
		getter   func(ParsedValue) interface{}
		expected interface{}
	}{
		{
			name: "GetQuestion returns QuestionResult",
			value: ParsedValue{
				Question: &QuestionResult{
					QuestionText: "What is Go?",
					AnswerText:   "A programming language",
				},
			},
			getter: func(pv ParsedValue) interface{} { return pv.GetQuestion() },
			expected: &QuestionResult{
				QuestionText: "What is Go?",
				AnswerText:   "A programming language",
			},
		},
		{
			name: "GetHeader returns HeaderResult",
			value: ParsedValue{
				Header: &HeaderResult{
					Parts: []string{"Section", "Chapter 1"},
				},
			},
			getter: func(pv ParsedValue) interface{} { return pv.GetHeader() },
			expected: &HeaderResult{
				Parts: []string{"Section", "Chapter 1"},
			},
		},
		{
			name: "GetComment returns CommentResult",
			value: ParsedValue{
				Comment: &CommentResult{
					Text: "This is a comment",
				},
			},
			getter: func(pv ParsedValue) interface{} { return pv.GetComment() },
			expected: &CommentResult{
				Text: "This is a comment",
			},
		},
		{
			name: "GetEmpty returns EmptyLineResult",
			value: ParsedValue{
				Empty: &EmptyLineResult{},
			},
			getter:   func(pv ParsedValue) interface{} { return pv.GetEmpty() },
			expected: &EmptyLineResult{},
		},
		{
			name: "GetFileHeader returns FileHeaderResult",
			value: ParsedValue{
				FileHeader: &FileHeaderResult{
					Title: "Study Guide",
				},
			},
			getter: func(pv ParsedValue) interface{} { return pv.GetFileHeader() },
			expected: &FileHeaderResult{
				Title: "Study Guide",
			},
		},
		{
			name: "GetPassage returns PassageResult",
			value: ParsedValue{
				Passage: &PassageResult{
					Text: "This is a passage",
				},
			},
			getter: func(pv ParsedValue) interface{} { return pv.GetPassage() },
			expected: &PassageResult{
				Text: "This is a passage",
			},
		},
		{
			name: "GetLearnMore returns LearnMoreResult",
			value: ParsedValue{
				LearnMore: &LearnMoreResult{
					Text: "Additional information",
				},
			},
			getter: func(pv ParsedValue) interface{} { return pv.GetLearnMore() },
			expected: &LearnMoreResult{
				Text: "Additional information",
			},
		},
		{
			name: "GetContent returns ContentResult",
			value: ParsedValue{
				Content: &ContentResult{
					Text: "This is content",
				},
			},
			getter: func(pv ParsedValue) interface{} { return pv.GetContent() },
			expected: &ContentResult{
				Text: "This is content",
			},
		},
		{
			name: "GetBinary returns BinaryResult",
			value: ParsedValue{
				Binary: &BinaryResult{
					Text: "binary data",
				},
			},
			getter: func(pv ParsedValue) interface{} { return pv.GetBinary() },
			expected: &BinaryResult{
				Text: "binary data",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.getter(tt.value)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("getter() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestParsedValueGettersReturnNil(t *testing.T) {
	emptyValue := ParsedValue{}

	tests := []struct {
		name   string
		getter func(ParsedValue) interface{}
	}{
		{"GetQuestion returns nil", func(pv ParsedValue) interface{} { return pv.GetQuestion() }},
		{"GetHeader returns nil", func(pv ParsedValue) interface{} { return pv.GetHeader() }},
		{"GetComment returns nil", func(pv ParsedValue) interface{} { return pv.GetComment() }},
		{"GetEmpty returns nil", func(pv ParsedValue) interface{} { return pv.GetEmpty() }},
		{"GetFileHeader returns nil", func(pv ParsedValue) interface{} { return pv.GetFileHeader() }},
		{"GetPassage returns nil", func(pv ParsedValue) interface{} { return pv.GetPassage() }},
		{"GetLearnMore returns nil", func(pv ParsedValue) interface{} { return pv.GetLearnMore() }},
		{"GetContent returns nil", func(pv ParsedValue) interface{} { return pv.GetContent() }},
		{"GetBinary returns nil", func(pv ParsedValue) interface{} { return pv.GetBinary() }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.getter(emptyValue)
			if result != nil {
				v := reflect.ValueOf(result)
				if v.Kind() == reflect.Ptr && v.IsNil() {
					return // ok
				}
				t.Errorf("getter() = %v (type %T), want nil", result, result)
			}
		})
	}
}

func TestParsedValueTypeCheckers(t *testing.T) {
	tests := []struct {
		name     string
		value    ParsedValue
		checker  func(ParsedValue) bool
		expected bool
	}{
		{
			name:     "IsQuestion returns true when Question is set",
			value:    ParsedValue{Question: &QuestionResult{}},
			checker:  func(pv ParsedValue) bool { return pv.IsQuestion() },
			expected: true,
		},
		{
			name:     "IsHeader returns true when Header is set",
			value:    ParsedValue{Header: &HeaderResult{}},
			checker:  func(pv ParsedValue) bool { return pv.IsHeader() },
			expected: true,
		},
		{
			name:     "IsComment returns true when Comment is set",
			value:    ParsedValue{Comment: &CommentResult{}},
			checker:  func(pv ParsedValue) bool { return pv.IsComment() },
			expected: true,
		},
		{
			name:     "IsEmpty returns true when Empty is set",
			value:    ParsedValue{Empty: &EmptyLineResult{}},
			checker:  func(pv ParsedValue) bool { return pv.IsEmpty() },
			expected: true,
		},
		{
			name:     "IsFileHeader returns true when FileHeader is set",
			value:    ParsedValue{FileHeader: &FileHeaderResult{}},
			checker:  func(pv ParsedValue) bool { return pv.IsFileHeader() },
			expected: true,
		},
		{
			name:     "IsPassage returns true when Passage is set",
			value:    ParsedValue{Passage: &PassageResult{}},
			checker:  func(pv ParsedValue) bool { return pv.IsPassage() },
			expected: true,
		},
		{
			name:     "IsLearnMore returns true when LearnMore is set",
			value:    ParsedValue{LearnMore: &LearnMoreResult{}},
			checker:  func(pv ParsedValue) bool { return pv.IsLearnMore() },
			expected: true,
		},
		{
			name:     "IsContent returns true when Content is set",
			value:    ParsedValue{Content: &ContentResult{}},
			checker:  func(pv ParsedValue) bool { return pv.IsContent() },
			expected: true,
		},
		{
			name:     "IsBinary returns true when Binary is set",
			value:    ParsedValue{Binary: &BinaryResult{}},
			checker:  func(pv ParsedValue) bool { return pv.IsBinary() },
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.checker(tt.value)
			if result != tt.expected {
				t.Errorf("checker() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestParsedValueTypeCheckersReturnFalse(t *testing.T) {
	emptyValue := ParsedValue{}

	tests := []struct {
		name    string
		checker func(ParsedValue) bool
	}{
		{"IsQuestion returns false", func(pv ParsedValue) bool { return pv.IsQuestion() }},
		{"IsHeader returns false", func(pv ParsedValue) bool { return pv.IsHeader() }},
		{"IsComment returns false", func(pv ParsedValue) bool { return pv.IsComment() }},
		{"IsEmpty returns false", func(pv ParsedValue) bool { return pv.IsEmpty() }},
		{"IsFileHeader returns false", func(pv ParsedValue) bool { return pv.IsFileHeader() }},
		{"IsPassage returns false", func(pv ParsedValue) bool { return pv.IsPassage() }},
		{"IsLearnMore returns false", func(pv ParsedValue) bool { return pv.IsLearnMore() }},
		{"IsContent returns false", func(pv ParsedValue) bool { return pv.IsContent() }},
		{"IsBinary returns false", func(pv ParsedValue) bool { return pv.IsBinary() }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.checker(emptyValue)
			if result != false {
				t.Errorf("checker() = %v, want false", result)
			}
		})
	}
}
