package qa

import (
	"encoding/json"
	"strings"
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
	qaRunner := NewTreeQARunner(
		NewTagTypeQA(),
		NewContextTypeQA(),
	)
	qaRunner.RunQAAndUpdate(treeObj)

	qaResults := treeObj.GetQAResults()
	t.Logf("QA Results: %+v", qaResults)

	if qaResults.OverallPassed {
		t.Fatal("Expected QA to fail due to TagTypeNone, but it passed")
	}

	if len(qaResults.Results) != 2 {
		t.Fatal("Expected 2 QA results, got", len(qaResults.Results))
	}

	// Check the Tag Type Validation result
	tagTypeResult := qaResults.Results[0]
	if tagTypeResult.Name != "Must have TagType" {
		t.Errorf("Expected QA step name 'Must have TagType', got '%s'", tagTypeResult.Name)
	}

	if tagTypeResult.Passed {
		t.Fatal("Expected Tag Type Validation to fail, but it passed")
	}

	if len(tagTypeResult.Warnings) == 0 {
		t.Fatal("Expected at least one warning for TagTypeNone, got none")
	}

	found := false
	for _, w := range tagTypeResult.Warnings {
		if w == "Tag 'Unclassified' at depth 7 has TagTypeNone" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected warning for 'Unclassified' tag at depth 7, got: %v", tagTypeResult.Warnings)
	}

	// Check the Context Type Validation result
	contextTypeResult := qaResults.Results[1]
	if contextTypeResult.Name != "Must have ContextType" {
		t.Errorf("Expected QA step name 'Must have ContextType', got '%s'", contextTypeResult.Name)
	}
}

func TestQAResultsJSONMarshaling(t *testing.T) {
	// Create a tree with valid tags that should pass QA
	metadata := &config.Metadata{
		Type:        "test",
		ContextType: ontology.ContextTypeAPExams,
	}
	treeObj := tree.NewTree(metadata)

	// Add a valid tag structure for APExams depth 3: Category -> AP Exam -> Topic
	tag1 := tree.NewTag("Category")
	tag2 := tree.NewTag("AP Exam")
	tag3 := tree.NewTag("Topic")
	tag1.AddChildTag(tag2)
	tag2.AddChildTag(tag3)
	treeObj.Root.AddChildTag(tag1)

	// Assign tag types
	treeObj.AssignTagTypes(ontology.ContextTypeAPExams)

	// Run QA
	qaRunner := NewTreeQARunner(
		NewTagTypeQA(),
		NewContextTypeQA(),
	)
	qaRunner.RunQAAndUpdate(treeObj)

	// Get the QA results
	qaResults := treeObj.GetQAResults()

	// Marshal to JSON
	jsonData, err := json.Marshal(qaResults)
	if err != nil {
		t.Fatalf("Failed to marshal QA results to JSON: %v", err)
	}

	jsonStr := string(jsonData)
	t.Logf("JSON output: %s", jsonStr)

	// Check that warnings arrays are empty arrays, not null
	if !strings.Contains(jsonStr, `"warnings":[]`) {
		t.Errorf("Expected 'warnings':[] in JSON, but got: %s", jsonStr)
	}

	// Check that we don't have null warnings
	if strings.Contains(jsonStr, `"warnings":null`) {
		t.Errorf("Found 'warnings':null in JSON, expected empty arrays: %s", jsonStr)
	}
}
