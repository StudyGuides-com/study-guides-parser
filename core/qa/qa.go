package qa

import (
	"fmt"

	"github.com/studyguides-com/study-guides-parser/core/ontology"
	"github.com/studyguides-com/study-guides-parser/core/tree"
)

type TreeQA interface {
	RunQA(tree tree.TreeQAble) []string
}

type DefaultTreeQA struct{}

func NewDefaultTreeQA() *DefaultTreeQA {
	return &DefaultTreeQA{}
}

func (qa *DefaultTreeQA) RunQA(t tree.TreeQAble) []string {
	var warnings []string
	
	t.Traverse(func(tagQAble tree.TagQATarget, depth int) {
		if tagQAble.GetTagType() == ontology.TagTypeNone {
			warnings = append(warnings, 
				fmt.Sprintf("Tag '%s' at depth %d has TagTypeNone", tagQAble.GetTitle(), depth))
		}
	})
	
	return warnings
}

type TreeQARunner struct {
	qa TreeQA
}

func NewTreeQARunner(qa TreeQA) *TreeQARunner {
	return &TreeQARunner{qa: qa}
}

func (runner *TreeQARunner) RunQAAndUpdate(t tree.TreeQAble) {
	warnings := runner.qa.RunQA(t)
	t.SetWarnings(warnings)
	
	// Set QA passed status based on whether there are any warnings
	qaPassed := len(warnings) == 0
	t.SetQAPassed(qaPassed)
} 