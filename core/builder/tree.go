package builder

import "github.com/studyguides-com/study-guides-parser/core/config"

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
	Title string `json:"title"`
	Content string `json:"content,omitempty"`
	Questions []*Question `json:"questions,omitempty"`
}

type Question struct {
	Prompt string `json:"prompt"`
	Answer   string `json:"answer"`
	Distractor []string `json:"distractor,omitempty"`
	LearnMore string `json:"learn_more,omitempty"`
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
	}
}

