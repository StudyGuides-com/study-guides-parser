package tree

import (
	"fmt"

	"github.com/studyguides-com/study-guides-parser/core/config"
	"github.com/studyguides-com/study-guides-parser/core/idgen"
	"github.com/studyguides-com/study-guides-parser/core/ontology"
)

// TagQATarget defines the interface for tags that can be QA'd
type TagQATarget interface {
	GetTagType() ontology.TagType
	GetTitle() string
}

// TreeQAble defines the interface for types that can be QA'd
// The visitor function is called for each tag with its depth
// Traverse must visit all tags in the tree
// GetWarnings/SetWarnings manage QA warnings
// GetQAPassed/SetQAPassed manage QA pass/fail status
type TreeQAble interface {
	Traverse(visitor func(TagQATarget, int))
	GetWarnings() []string
	SetWarnings(warnings []string)
	GetQAPassed() bool
	SetQAPassed(passed bool)
}
// Root represents the file-level container, not an actual content tag
type Root struct {
	Title     string   `json:"title"`
	QAPassed  bool     `json:"qa_passed,omitempty"`
	Warnings  []string `json:"warnings,omitempty"`
	ChildTags []*Tag   `json:"child_tags,omitempty"`
}

func NewRoot() *Root {
	return &Root{
		Title: "Root",
	}
}

func (r *Root) GetChildTags() []*Tag {
	return r.ChildTags
}

func (r *Root) AddChildTag(tag *Tag) {
	r.ChildTags = append(r.ChildTags, tag)
}

type Tree struct {
	Root     *Root            `json:"root"`
	Metadata *config.Metadata `json:"metadata"`
}

func NewTree(metadata *config.Metadata) *Tree {
	return &Tree{
		Metadata: metadata,
		Root:     NewRoot(),
	}
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

type Question struct {
	Hash       string   `json:"hash,omitempty"`
	Prompt     string   `json:"prompt"`
	Answer     string   `json:"answer"`
	Distractor []string `json:"distractor,omitempty"`
	LearnMore  string   `json:"learn_more,omitempty"`
}

func NewQuestion(prompt string, answer string, distractor []string, learnMore string) *Question {
	return &Question{
		Hash:       idgen.HashFrom(prompt + answer),
		Prompt:     prompt,
		Answer:     answer,
		Distractor: distractor,
		LearnMore:  learnMore,
	}
}

type Tag struct {
	Title     string              `json:"title"`
	TagType   ontology.TagType    `json:"tag_type,omitempty"`
	InsertID  string              `json:"insert_id,omitempty"`
	Context   ontology.ContextType `json:"context,omitempty"`
	Hash      string              `json:"hash,omitempty"`
	Questions []*Question         `json:"questions,omitempty"`
	Passages  []*Passage          `json:"passages,omitempty"`
	ChildTags []*Tag              `json:"child_tags,omitempty"`
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

// TreeTraverser defines the interface for traversing tree structures
type TreeTraverser interface {
	// Traverse performs a depth-first traversal of the tree
	// The visitor function is called for each tag with its depth
	Traverse(visitor func(TagTypeAssignable, int))
	
	// TraverseWithContext performs traversal with additional context
	TraverseWithContext(visitor func(TagTypeAssignable, int, ontology.ContextType))
}

// TagTypeAssigner defines the interface for assigning tag types
type TagTypeAssigner interface {
	AssignTagTypes(contextType ontology.ContextType) error
}

// TraverseForTagTypes implements TreeTraverser interface for tag type assignment
func (t *Tree) TraverseForTagTypes(visitor func(TagTypeAssignable, int)) {
	if t.Root == nil {
		return
	}
	
	var traverse func(*Tag, int)
	traverse = func(tag *Tag, depth int) {
		if tag == nil {
			return
		}
		
		// Visit current tag
		visitor(tag, depth)
		
		// Recursively visit children
		for _, child := range tag.ChildTags {
			traverse(child, depth+1)
		}
	}
	
	// Start traversal from root-level tags
	for _, child := range t.Root.ChildTags {
		traverse(child, 1)
	}
}

// Traverse implements TreeQAble interface for QA
func (t *Tree) Traverse(visitor func(TagQATarget, int)) {
	if t.Root == nil {
		return
	}
	var traverse func(*Tag, int)
	traverse = func(tag *Tag, depth int) {
		if tag == nil {
			return
		}
		visitor(tag, depth)
		for _, child := range tag.ChildTags {
			traverse(child, depth+1)
		}
	}
	for _, child := range t.Root.ChildTags {
		traverse(child, 1)
	}
}

// TraverseWithContext implements TreeTraverser interface with context
func (t *Tree) TraverseWithContext(visitor func(TagTypeAssignable, int, ontology.ContextType)) {
	if t.Root == nil {
		return
	}
	
	// Get context from metadata
	contextType := ontology.ContextTypeNone
	if t.Metadata != nil {
		contextType = t.Metadata.ContextType
	}
	
	var traverse func(*Tag, int)
	traverse = func(tag *Tag, depth int) {
		if tag == nil {
			return
		}
		
		// Visit current tag with context
		visitor(tag, depth, contextType)
		
		// Recursively visit children
		for _, child := range tag.ChildTags {
			traverse(child, depth+1)
		}
	}
	
	// Start traversal from root-level tags
	for _, child := range t.Root.ChildTags {
		traverse(child, 1)
	}
}

// AssignTagTypes implements TagTypeAssigner interface
func (t *Tree) AssignTagTypes(contextType ontology.ContextType) error {
	// First, determine the maximum depth of the tree
	maxDepth := t.getMaxDepth()
	
	// Find the ontology entry for this total depth
	tagOntology := ontology.FindTagOntology(contextType, maxDepth)
	if tagOntology == nil {
		return fmt.Errorf("no ontology found for context type '%s' with depth %d", contextType, maxDepth)
	}
	
	// Now traverse and assign types based on individual tag depths
	t.TraverseForTagTypes(func(tag TagTypeAssignable, depth int) {
		assignTagTypeFromOntology(tag, contextType, depth, tagOntology)
	})
	
	return nil
}

// getMaxDepth calculates the maximum depth of the tree
func (t *Tree) getMaxDepth() int {
	if t.Root == nil {
		return 0
	}
	
	var maxDepth int
	var traverse func(*Tag, int)
	traverse = func(tag *Tag, depth int) {
		if depth > maxDepth {
			maxDepth = depth
		}
		for _, child := range tag.ChildTags {
			traverse(child, depth+1)
		}
	}
	
	for _, child := range t.Root.ChildTags {
		traverse(child, 1)
	}
	
	return maxDepth
}

// assignTagTypeFromOntology assigns the appropriate tag type based on context and depth within a known ontology
func assignTagTypeFromOntology(tag TagTypeAssignable, contextType ontology.ContextType, depth int, tagOntology *ontology.TagOntology) {
	if depth <= len(tagOntology.TagTypes) {
		tag.SetTagType(tagOntology.TagTypes[depth-1]) // depth is 1-indexed, slice is 0-indexed
		tag.SetContext(contextType)
	}
}

// GetTitle implements TagQATarget
func (t *Tag) GetTitle() string {
	return t.Title
}

// GetWarnings returns the current warnings
func (r *Root) GetWarnings() []string {
	return r.Warnings
}

// SetWarnings sets the warnings
func (r *Root) SetWarnings(warnings []string) {
	r.Warnings = warnings
}

// GetQAPassed returns whether QA passed
func (r *Root) GetQAPassed() bool {
	return r.QAPassed
}

// SetQAPassed sets whether QA passed
func (r *Root) SetQAPassed(passed bool) {
	r.QAPassed = passed
}

// GetWarnings returns the current warnings from the root
func (t *Tree) GetWarnings() []string {
	if t.Root == nil {
		return nil
	}
	return t.Root.GetWarnings()
}

// SetWarnings sets the warnings on the root
func (t *Tree) SetWarnings(warnings []string) {
	if t.Root == nil {
		return
	}
	t.Root.SetWarnings(warnings)
}

// GetQAPassed returns whether QA passed
func (t *Tree) GetQAPassed() bool {
	if t.Root != nil {
		return t.Root.GetQAPassed()
	}
	return false
}

// SetQAPassed sets whether QA passed
func (t *Tree) SetQAPassed(passed bool) {
	if t.Root != nil {
		t.Root.SetQAPassed(passed)
	}
}
