package builder

import (
	"github.com/studyguides-com/study-guides-parser/core/config"
	"github.com/studyguides-com/study-guides-parser/core/idgen"
)

type ContextType string

const (
	ContextTypeColleges             ContextType = "Colleges"
	ContextTypeCertifications       ContextType = "Certifications"
	ContextTypeEntranceExams        ContextType = "EntranceExams"
	ContextTypeAPExams              ContextType = "APExams"
	ContextTypeUserGeneratedContent ContextType = "UserGeneratedContent"
	ContextTypeDoD                  ContextType = "DoD"
	ContextTypeNone                 ContextType = "None"
)

type TagType string

const (
	TagTypeCategory         TagType = "Category"
	TagTypeSubCategory      TagType = "SubCategory"
	TagTypeUniversity       TagType = "University"
	TagTypeRegion           TagType = "Region"
	TagTypeDepartment       TagType = "Department"
	TagTypeCourse           TagType = "Course"
	TagTypeTopic            TagType = "Topic"
	TagTypeUserFolder       TagType = "UserFolder"
	TagTypeUserTopic        TagType = "UserTopic"
	TagTypeCertifyingAgency TagType = "Certifying_Agency"
	TagTypeCertification    TagType = "Certification"
	TagTypeDomain           TagType = "Domain"
	TagTypeModule           TagType = "Module"
	TagTypeEntranceExam     TagType = "Entrance_Exam"
	TagTypeAPExam           TagType = "AP_Exam"
	TagTypeUserContent      TagType = "UserContent"
	TagTypeBranch           TagType = "Branch"
	TagTypeInstructionType  TagType = "Instruction_Type"
	TagTypeInstructionGroup TagType = "Instruction_Group"
	TagTypeInstruction      TagType = "Instruction"
	TagTypeChapter          TagType = "Chapter"
	TagTypeSection          TagType = "Section"
	TagTypePart             TagType = "Part"
	TagTypeNone             TagType = "None"
)

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

// Root represents the file-level container, not an actual content tag
type Root struct {
	Title     string `json:"title"`
	ChildTags []*Tag `json:"child_tags,omitempty"`
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
	Title     string      `json:"title"`
	TagType   TagType     `json:"tag_type,omitempty"`
	InsertID  string      `json:"insert_id,omitempty"`
	Context   ContextType `json:"context,omitempty"`
	Hash      string      `json:"hash,omitempty"`
	Questions []*Question `json:"questions,omitempty"`
	Passages  []*Passage  `json:"passages,omitempty"`
	ChildTags []*Tag      `json:"child_tags,omitempty"`
}

func NewTag(title string) *Tag {
	return &Tag{
		InsertID: idgen.NewCUID(),
		Title:    title,
		Hash:     idgen.HashFrom(title),
		TagType:  TagTypeNone,
		Context:  ContextTypeNone,
	}
}

// NewTagWithParent creates a new tag with a hash based on its parent's title
func NewTagWithParent(title string, parentTitle string) *Tag {
	return &Tag{
		InsertID: idgen.NewCUID(),
		Title:    title,
		Hash:     idgen.HashFrom(parentTitle + title),
		TagType:  TagTypeNone,
		Context:  ContextTypeNone,
	}
}

func (t *Tag) GetChildTags() []*Tag {
	return t.ChildTags
}

func (t *Tag) AddChildTag(tag *Tag) {
	t.ChildTags = append(t.ChildTags, tag)
}
