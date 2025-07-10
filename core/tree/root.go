package tree

// Root represents the file-level container, not an actual content tag
type Root struct {
	Title     string     `json:"title"`
	QAResults QAResults  `json:"qa_results,omitempty"`
	Warnings  []string   `json:"warnings,omitempty"`
	ChildTags []*Tag     `json:"child_tags,omitempty"`
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

// GetQAPassed returns whether QA passed (for backward compatibility)
func (r *Root) GetQAPassed() bool {
	return r.QAResults.OverallPassed
}

// SetQAPassed sets whether QA passed (for backward compatibility)
func (r *Root) SetQAPassed(passed bool) {
	// For backward compatibility, create a simple QA result
	qaResults := NewQAResults()
	qaResults.OverallPassed = passed
	r.QAResults = qaResults
}

// GetQAResults returns the QA results
func (r *Root) GetQAResults() QAResults {
	return r.QAResults
}

// SetQAResults sets the QA results
func (r *Root) SetQAResults(results QAResults) {
	r.QAResults = results
} 