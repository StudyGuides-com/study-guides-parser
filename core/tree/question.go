package tree

import (
	"github.com/studyguides-com/study-guides-parser/core/idgen"
)

type Question struct {
	InsertID   string   `json:"insert_id,omitempty"`
	Hash       string   `json:"hash,omitempty"`
	Prompt     string   `json:"prompt"`
	Answer     string   `json:"answer"`
	Distractor []string `json:"distractor,omitempty"`
	LearnMore  string   `json:"learn_more,omitempty"`
}

func NewQuestion(prompt string, answer string, distractor []string, learnMore string) *Question {
	return &Question{
		InsertID:   idgen.NewCUID(),
		Hash:       idgen.HashFrom(prompt + answer),
		Prompt:     prompt,
		Answer:     answer,
		Distractor: distractor,
		LearnMore:  learnMore,
	}
}
