package qa

import (
	"testing"

	"github.com/studyguides-com/study-guides-parser/core/config"
	"github.com/studyguides-com/study-guides-parser/core/ontology"
	"github.com/studyguides-com/study-guides-parser/core/tree"
)

func TestTreeQAWarnsOnTagTypeNone(t *testing.T) {
	metadata := &config.Metadata{
		Type:        "test",
		ContextType: ontology.ContextTypeAPExams, // Use APExams which only supports up to depth 6
	}
	treeObj := tree.NewTree(metadata)

	// Add tags with depth 7 structure (which exceeds APExams ontology)
	tag1 := tree.NewTag("Category")
	tag2 := tree.NewTag("AP Exam")
	tag3 := tree.NewTag("Module")
	tag4 := tree.NewTag("Module")
	tag5 := tree.NewTag("Module")
	tag6 := tree.NewTag("Topic")
	tag7 := tree.NewTag("Unclassified") // This will remain TagTypeNone

	tag1.AddChildTag(tag2)
	tag2.AddChildTag(tag3)
	tag3.AddChildTag(tag4)
	tag4.AddChildTag(tag5)
	tag5.AddChildTag(tag6)
	tag6.AddChildTag(tag7)
	treeObj.Root.AddChildTag(tag1)

	// Assign tag types (tag7 will not get a type, as depth 7 is not in ontology for APExams)
	err := treeObj.AssignTagTypes(ontology.ContextTypeAPExams)
	t.Logf("AssignTagTypes error: %v", err) // This should fail for depth 7

	// Debug: Print tag types before QA
	t.Log("Before QA:")
	treeObj.Traverse(func(tag tree.TagQATarget, depth int) {
		t.Logf("Tag: %s, Depth: %d, Type: %s", tag.GetTitle(), depth, tag.GetTagType())
	})

	// Run QA using the new QA package
	qaRunner := NewTreeQARunner(NewDefaultTreeQA())
	qaRunner.RunQAAndUpdate(treeObj)

	t.Logf("Warnings found: %v", treeObj.Root.Warnings)

	if len(treeObj.Root.Warnings) == 0 {
		t.Fatal("Expected at least one warning for TagTypeNone, got none")
	}

	found := false
	for _, w := range treeObj.Root.Warnings {
		if w == "Tag 'Unclassified' at depth 7 has TagTypeNone" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected warning for 'Unclassified' tag at depth 7, got: %v", treeObj.Root.Warnings)
	}
} 