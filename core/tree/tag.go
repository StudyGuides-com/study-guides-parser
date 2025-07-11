package tree

import (
	"github.com/studyguides-com/study-guides-parser/core/idgen"
	"github.com/studyguides-com/study-guides-parser/core/ontology"
)

// TagQATarget defines the interface for tags that can be QA'd
type TagQATarget interface {
	GetTagType() ontology.TagType
	GetTitle() string
	GetContext() ontology.ContextType
}

// TagContainer interface for types that can contain child tags
type TagContainer interface {
	GetChildTags() []*Tag
	AddChildTag(*Tag)
}

// TagTypeAssignable interface for tags that can have their types assigned
type TagTypeAssignable interface {
	SetTagType(tagType ontology.TagType)
	SetContext(contextType ontology.ContextType)
	GetTagType() ontology.TagType
	GetContext() ontology.ContextType
}

type Tag struct {
	Title     string               `json:"title"`
	TagType   ontology.TagType     `json:"tag_type,omitempty"`
	InsertID  string               `json:"insert_id,omitempty"`
	Context   ontology.ContextType `json:"context,omitempty"`
	Hash      string               `json:"hash,omitempty"`
	Questions []*Question          `json:"questions,omitempty"`
	Passages  []*Passage           `json:"passages,omitempty"`
	ChildTags []*Tag               `json:"child_tags,omitempty"`
}

func NewTag(title string) *Tag {
	return &Tag{
		InsertID: idgen.NewCUID(),
		Title:    title,
		Hash:     idgen.HashFrom(title),
		TagType:  ontology.TagTypeNone,
		Context:  ontology.ContextTypeNone,
	}
}

// NewTagWithParent creates a new tag with a hash based on its parent's title
func NewTagWithParent(title string, parentTitle string) *Tag {
	return &Tag{
		InsertID: idgen.NewCUID(),
		Title:    title,
		Hash:     idgen.HashFrom(parentTitle + title),
		TagType:  ontology.TagTypeNone,
		Context:  ontology.ContextTypeNone,
	}
}

func (t *Tag) GetChildTags() []*Tag {
	return t.ChildTags
}

func (t *Tag) AddChildTag(tag *Tag) {
	t.ChildTags = append(t.ChildTags, tag)
}

// SetTagType sets the tag type
func (t *Tag) SetTagType(tagType ontology.TagType) {
	t.TagType = tagType
}

// SetContext sets the context type
func (t *Tag) SetContext(contextType ontology.ContextType) {
	t.Context = contextType
}

// GetTagType returns the current tag type
func (t *Tag) GetTagType() ontology.TagType {
	return t.TagType
}

// GetContext returns the current context type
func (t *Tag) GetContext() ontology.ContextType {
	return t.Context
}

// GetTitle implements TagQATarget
func (t *Tag) GetTitle() string {
	return t.Title
}
