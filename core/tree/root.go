package tree

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