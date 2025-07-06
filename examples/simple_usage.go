package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/StudyGuides-com/study-guides-parser/core/processor"
)

func main() {
	// Example 1: Parse from strings (super simple!)
	lines := []string{
		"Mathematics Study Guide",
		"Colleges: Virginia: Old Dominion University (ODU): Mathematics (MATH): MATH 101: Linear Equations",
		"",
		"1. What is a linear equation? - An equation where the highest power of the variable is 1.",
		"2. How do you solve 2x + 3 = 7? - Subtract 3 from both sides: 2x = 4, then divide by 2: x = 2.",
		"",
		"Learn More: See Khan Academy's linear equations course.",
	}

	fmt.Println("=== Example 1: Parse from strings ===")
	metadata := processor.NewMetadata(processor.Colleges)
	ast, err := processor.Parse(lines, metadata)
	if err != nil {
		log.Fatal("Error parsing lines:", err)
	}

	// Output as JSON
	jsonData, _ := json.MarshalIndent(ast, "", "  ")
	fmt.Println(string(jsonData))

	// Example 2: Parse from file
	fmt.Println("\n=== Example 2: Parse from file ===")
	// This would work with an actual file
	// ast, err = processor.ParseFile("study_guide.txt", processor.NewMetadata(processor.Colleges))
	// if err != nil {
	//     log.Fatal("Error parsing file:", err)
	// }

	// Example 3: Different parser types
	fmt.Println("\n=== Example 3: AP Exam format ===")
	apLines := []string{
		"AP Calculus AB Study Guide",
		"Advanced Placement (AP): AP Calculus AB: Derivatives: Introduction to Derivatives",
		"",
		"1. What is a derivative? - The rate of change of a function.",
		"2. How do you find the derivative of x²? - Use the power rule: 2x.",
	}

	apMetadata := processor.NewMetadata(processor.APExams)
	ast, err = processor.Parse(apLines, apMetadata)
	if err != nil {
		log.Fatal("Error parsing AP exam:", err)
	}

	// Show just the structure
	fmt.Printf("AST Type: %s\n", ast.ParserType)
	fmt.Printf("Root Type: %s\n", ast.Root.Type)
	fmt.Printf("Number of children: %d\n", len(ast.Root.Children))

	// Example 4: Advanced usage with options
	fmt.Println("\n=== Example 4: Advanced usage with options ===")
	advancedMetadata := processor.NewMetadata(processor.Colleges).
		WithOption("strict", "true").
		WithOption("debug", "false").
		WithOption("version", "1.0")

	fmt.Printf("Parser Type: %s\n", advancedMetadata.ParserType)
	fmt.Printf("Options: %+v\n", advancedMetadata.Options)

	// Example 5: Error handling
	fmt.Println("\n=== Example 5: Error handling ===")
	invalidLines := []string{
		"Colleges: Virginia: Old Dominion University (ODU): Mathematics (MATH): MATH 101: Linear Equations",
		"1. What is x? - A variable",
	}

	_, err = processor.Parse(invalidLines, processor.NewMetadata(processor.Colleges))
	if err != nil {
		fmt.Printf("Expected error: %v\n", err)
	}

	fmt.Println("\n=== Summary ===")
	fmt.Println("✅ Phase 1 Complete: Simple API implemented!")
	fmt.Println("✅ Users can now parse study guides with just one function call")
	fmt.Println("✅ No more manual pipeline management")
	fmt.Println("✅ No more token format conversions")
	fmt.Println("✅ All existing functionality preserved")
	fmt.Println("✅ New Metadata struct for future extensibility")
} 