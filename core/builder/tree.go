package builder

import (
	"github.com/studyguides-com/study-guides-parser/core/config"
	"github.com/studyguides-com/study-guides-parser/core/idgen"
)

type Tree struct {
	Root *Tag `json:"root"`
	Metadata *config.Metadata `json:"metadata"`
}

func NewTree(metadata *config.Metadata) *Tree {
	return &Tree{
		Metadata: metadata,
		Root: NewTag("Root"),
	}
}

type Passage struct {
	Hash string `json:"hash,omitempty"`
	Title string `json:"title"`
	Content string `json:"content,omitempty"`
	Questions []*Question `json:"questions,omitempty"`
}

func NewPassage(title string, content string, questions []*Question) *Passage {
	return &Passage{
		Hash: idgen.HashFrom(title),
		Title: title,
		Content: content,
		Questions: questions,
	}
}

type Question struct {
	Hash string `json:"hash,omitempty"`
	Prompt string `json:"prompt"`
	Answer   string `json:"answer"`
	Distractor []string `json:"distractor,omitempty"`
	LearnMore string `json:"learn_more,omitempty"`
}

func NewQuestion(prompt string, answer string, distractor []string, learnMore string) *Question {
	return &Question{
		Hash: idgen.HashFrom(prompt + answer),
		Prompt: prompt,
		Answer: answer,
		Distractor: distractor,
		LearnMore: learnMore,
	}
}

type Tag struct {
	Title    string `json:"title"`
	Hash     string `json:"hash,omitempty"`
	Questions       []*Question `json:"questions,omitempty"`
	Passages       []*Passage `json:"passages,omitempty"`
	ChildTags     []*Tag `json:"child_tags,omitempty"`
}


func NewTag(title string) *Tag {
	return &Tag{
		Title: title,
		Hash:  idgen.HashFrom(title),
	}
}

// NewTagWithParent creates a new tag with a hash based on its parent's title
func NewTagWithParent(title string, parentTitle string) *Tag {
	return &Tag{
		Title: title,
		Hash:  idgen.HashFrom(parentTitle + title),
	}
}

