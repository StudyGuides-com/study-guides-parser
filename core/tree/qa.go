package tree

// QAResult represents the result of a single QA step
type QAResult struct {
	Name     string   `json:"name"`     // Friendly name for the QA step
	Passed   bool     `json:"passed"`   // Whether this QA step passed
	Warnings []string `json:"warnings"` // Any warnings/errors from this step
}

// NewQAResult creates a new QAResult with default values
func NewQAResult(name string, passed bool) QAResult {
	return QAResult{
		Name:     name,
		Passed:   passed,
		Warnings: []string{}, // Empty slice by default
	}
}

// QAResults represents all QA results for a tree
type QAResults struct {
	OverallPassed bool        `json:"overall_passed"` // Whether all QA steps passed
	Results       []QAResult  `json:"results"`        // Individual QA step results
}

// NewQAResults creates a new QAResults with default values
func NewQAResults() QAResults {
	return QAResults{
		OverallPassed: true,
		Results:       []QAResult{}, // Empty slice by default
	}
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

// GetQAResults returns the QA results from the root
func (t *Tree) GetQAResults() QAResults {
	if t.Root == nil {
		return QAResults{}
	}
	return t.Root.GetQAResults()
}

// SetQAResults sets the QA results on the root
func (t *Tree) SetQAResults(results QAResults) {
	if t.Root == nil {
		return
	}
	t.Root.SetQAResults(results)
} 