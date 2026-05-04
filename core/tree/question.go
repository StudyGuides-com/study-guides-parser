package tree

import (
	"github.com/studyguides-com/study-guides-parser/core/idgen"
)

type Question struct {
	InsertID    string   `json:"insert_id,omitempty"`
	Hash        string   `json:"hash,omitempty"`
	Prompt      string   `json:"prompt"`
	Answer      string   `json:"answer"`
	Distractors []string `json:"distractors"`
	LearnMore   string   `json:"learn_more"`
	Order       int      `json:"order"`
}

func NewQuestion(prompt string, answer string, distractors []string, learnMore string, order int) *Question {
	// If distractors is nil, create an empty slice
	if distractors == nil {
		distractors = []string{}
	}

	return &Question{
		InsertID:    idgen.NewCUID(),
		Hash:        idgen.HashFrom(prompt + answer),
		Prompt:      prompt,
		Answer:      answer,
		Distractors: distractors,
		LearnMore:   learnMore,
		Order:       order,
	}
}
