package main

import (
	"fmt"

	"github.com/iamrekas/anchor-go/pkg/idl/parser"
)

func main() {
	// Parse the events example IDL
	idlData, err := parser.ParseFileWithVerbose("../../idl/examples/events.json", true)
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
}
