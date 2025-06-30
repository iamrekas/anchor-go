package main

import (
	"fmt"
	"os"

	"github.com/iamrekas/anchor-go/pkg/generator"
	"github.com/iamrekas/anchor-go/pkg/idl/parser"
)

func main() {
	// Check if an IDL file path was provided
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <path-to-idl-file>")
		return
	}

	idlPath := os.Args[1]

	// Parse the IDL file
	idlData, err := parser.ParseFileWithVerbose(idlPath, true)
	if err != nil {
		fmt.Printf("Error parsing IDL: %v\n", err)
		return
	}

	// Print events
	fmt.Println("\nEvents from IDL:")
	for _, event := range idlData.Events() {
		fmt.Printf("- %s: %d fields\n", event.Name, len(event.Fields))
		for i, field := range event.Fields {
			fmt.Printf("  %d. %s: %s\n", i+1, field.Name, field.Type.String())
		}
	}

	// Generate event code
	eventsGenerator := generator.NewEventsGenerator()
	config := generator.DefaultConfig()
	config.DstDir = "generated_events"
	config.Package = "events"
	eventsGenerator.SetConfig(config)

	generatedFiles, err := eventsGenerator.Generate(idlData)
	if err != nil {
		fmt.Printf("Error generating events: %v\n", err)
		return
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(config.DstDir, 0755); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		return
	}

	// Write generated files
	for _, file := range generatedFiles {
		filePath := fmt.Sprintf("%s/%s.go", file.Path, file.Name)
		if err := os.WriteFile(filePath, file.Content, 0644); err != nil {
			fmt.Printf("Error writing file %s: %v\n", filePath, err)
			return
		}
		fmt.Printf("Generated file: %s\n", filePath)
	}

	fmt.Println("\nEvent generation completed successfully!")
}
