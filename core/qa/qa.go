package qa

import (
	"fmt"

	"github.com/studyguides-com/study-guides-parser/core/ontology"
	"github.com/studyguides-com/study-guides-parser/core/tree"
)

// TreeQA interface for running QA checks
type TreeQA interface {
	RunQA(tree tree.TreeQAble) tree.QAResult
}

// TagTypeQA validates that tags have proper tag types assigned
type TagTypeQA struct{}

func NewTagTypeQA() *TagTypeQA {
	return &TagTypeQA{}
}

func (qa *TagTypeQA) RunQA(t tree.TreeQAble) tree.QAResult {
	var warnings []string
	
	t.Traverse(func(tagQAble tree.TagQATarget, depth int) {
		if tagQAble.GetTagType() == ontology.TagTypeNone {
			warnings = append(warnings, 
				fmt.Sprintf("Tag '%s' at depth %d has TagTypeNone", tagQAble.GetTitle(), depth))
		}
	})
	
	result := tree.NewQAResult("Must have TagType", len(warnings) == 0)
	if len(warnings) > 0 {
		result.Warnings = warnings
	}
	// If warnings is empty, result.Warnings remains the empty slice from NewQAResult
	return result
}

// ContextTypeQA validates that tags have proper context types assigned
type ContextTypeQA struct{}

func NewContextTypeQA() *ContextTypeQA {
	return &ContextTypeQA{}
}

func (qa *ContextTypeQA) RunQA(t tree.TreeQAble) tree.QAResult {
	var warnings []string
	
	t.Traverse(func(tagQAble tree.TagQATarget, depth int) {
		if tagQAble.GetContext() == ontology.ContextTypeNone {
			warnings = append(warnings, 
				fmt.Sprintf("Tag '%s' at depth %d has ContextTypeNone", tagQAble.GetTitle(), depth))
		}
	})
	
	result := tree.NewQAResult("Must have ContextType", len(warnings) == 0)
	if len(warnings) > 0 {
		result.Warnings = warnings
	}
	// If warnings is empty, result.Warnings remains the empty slice from NewQAResult
	return result
}

// DefaultTreeQA is kept for backward compatibility but now just runs TagTypeQA
type DefaultTreeQA struct{}

func NewDefaultTreeQA() *DefaultTreeQA {
	return &DefaultTreeQA{}
}

func (qa *DefaultTreeQA) RunQA(t tree.TreeQAble) tree.QAResult {
	tagTypeQA := NewTagTypeQA()
	return tagTypeQA.RunQA(t)
}

type TreeQARunner struct {
	qaSteps []TreeQA
}

func NewTreeQARunner(qaSteps ...TreeQA) *TreeQARunner {
	return &TreeQARunner{qaSteps: qaSteps}
}

func (runner *TreeQARunner) RunQAAndUpdate(t tree.TreeQAble) {
	var results []tree.QAResult
	overallPassed := true
	
	for _, qaStep := range runner.qaSteps {
		result := qaStep.RunQA(t)
		results = append(results, result)
		if !result.Passed {
			overallPassed = false
		}
	}
	
	qaResults := tree.NewQAResults()
	qaResults.OverallPassed = overallPassed
	qaResults.Results = results
	
	t.SetQAResults(qaResults)
} 