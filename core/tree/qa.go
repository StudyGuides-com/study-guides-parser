package tree

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