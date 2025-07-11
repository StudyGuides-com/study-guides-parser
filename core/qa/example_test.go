package qa

import (
	"fmt"
	"testing"

	"github.com/studyguides-com/study-guides-parser/core/config"
	"github.com/studyguides-com/study-guides-parser/core/ontology"
	"github.com/studyguides-com/study-guides-parser/core/tree"
)

func ExampleQAResults() {
	// Create a tree with some tags
	metadata := &config.Metadata{
		Type:        "test",
		ContextType: ontology.ContextTypeAPExams,
	}
	treeObj := tree.NewTree(metadata)

	// Add a tag that will have TagTypeNone (invalid)
	tag := tree.NewTag("Invalid Tag")
	treeObj.Root.AddChildTag(tag)

	// Run QA
	qaRunner := NewTreeQARunner(
		NewTagTypeQA(),
		NewContextTypeQA(),
	)
	qaRunner.RunQAAndUpdate(treeObj)

	// Get the QA results
	qaResults := treeObj.GetQAResults()

	// Print the results
	fmt.Printf("Overall QA Passed: %t\n", qaResults.OverallPassed)
	fmt.Printf("Number of QA steps: %d\n", len(qaResults.Results))

	for i, result := range qaResults.Results {
		fmt.Printf("QA Step %d: %s\n", i+1, result.Name)
		fmt.Printf("  Passed: %t\n", result.Passed)
		fmt.Printf("  Warnings: %d\n", len(result.Warnings))
		for j, warning := range result.Warnings {
			fmt.Printf("    %d. %s\n", j+1, warning)
		}
	}

	// Output:
	// Overall QA Passed: false
	// Number of QA steps: 2
	// QA Step 1: Must have TagType
	//   Passed: false
	//   Warnings: 1
	//     1. Tag 'Invalid Tag' at depth 1 has TagTypeNone
	// QA Step 2: Must have ContextType
	//   Passed: false
	//   Warnings: 1
	//     1. Tag 'Invalid Tag' at depth 1 has ContextTypeNone
}

func TestQAResultsStructure(t *testing.T) {
	// Create a tree with valid tags
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

	// Verify the structure
	if !qaResults.OverallPassed {
		t.Errorf("Expected QA to pass for valid tags, but it failed")
	}

	if len(qaResults.Results) != 2 {
		t.Errorf("Expected 2 QA results, got %d", len(qaResults.Results))
	}

	// Check the Tag Type Validation result
	tagTypeResult := qaResults.Results[0]
	if tagTypeResult.Name != "Must have TagType" {
		t.Errorf("Expected QA step name 'Must have TagType', got '%s'", tagTypeResult.Name)
	}

	if !tagTypeResult.Passed {
		t.Errorf("Expected Tag Type Validation to pass, but it failed")
	}

	if len(tagTypeResult.Warnings) != 0 {
		t.Errorf("Expected no warnings for valid tags, got %d", len(tagTypeResult.Warnings))
	}

	// Check the Context Type Validation result
	contextTypeResult := qaResults.Results[1]
	if contextTypeResult.Name != "Must have ContextType" {
		t.Errorf("Expected QA step name 'Must have ContextType', got '%s'", contextTypeResult.Name)
	}

	if !contextTypeResult.Passed {
		t.Errorf("Expected Context Type Validation to pass, but it failed")
	}

	if len(contextTypeResult.Warnings) != 0 {
		t.Errorf("Expected no warnings for valid tags, got %d", len(contextTypeResult.Warnings))
	}
}
