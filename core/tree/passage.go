package tree

import (
	"github.com/studyguides-com/study-guides-parser/core/idgen"
)

type Passage struct {
	Hash      string      `json:"hash,omitempty"`
	Title     string      `json:"title"`
	Content   string      `json:"content,omitempty"`
	Questions []*Question `json:"questions,omitempty"`
}

func NewPassage(title string, content string, questions []*Question) *Passage {
	return &Passage{
		Hash:      idgen.HashFrom(title),
		Title:     title,
		Content:   content,
		Questions: questions,
	}
} 