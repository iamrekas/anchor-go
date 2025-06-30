package parser

import (
	"encoding/json"
	"fmt"

	"github.com/iamrekas/anchor-go/pkg/idl/types"
)

// ParseEvents parses events from raw event data
func ParseEvents(rawEvents []json.RawMessage, rawTypes []json.RawMessage, verbose bool) []types.Event {
	if verbose {
		fmt.Printf("Parsing %d events\n", len(rawEvents))
	}

	if verbose {
		fmt.Println("Events:")
	}

	// First, parse all types to find event types
	typeMap := make(map[string][]types.Field)
	for _, rawType := range rawTypes {
		var typeDef struct {
			Name string `json:"name"`
			Type struct {
				Kind   string                   `json:"kind"`
				Fields []map[string]interface{} `json:"fields"`
			} `json:"type"`
		}

		if err := json.Unmarshal(rawType, &typeDef); err != nil {
			continue
		}

		// Skip non-struct types
		if typeDef.Type.Kind != "struct" {
			continue
		}

		// Parse fields
		fields := make([]types.Field, 0, len(typeDef.Type.Fields))
		for _, field := range typeDef.Type.Fields {
			name, _ := field["name"].(string)
			docs, _ := field["docs"].([]string)
			optional, _ := field["optional"].(bool)

			// Create a basic field
			fieldDef := types.Field{
				Name:     name,
				Docs:     docs,
				Optional: optional,
				Type:     &types.BasicType{TypeName: "string"}, // Default to string for now
			}

			// Parse the field type
			if fieldType, ok := field["type"]; ok {
				// Use the helper function to parse the type
				rawFieldType, _ := json.Marshal(fieldType)
				parsedType, success := ParseFieldType(name, rawFieldType, verbose)
				if success {
					fieldDef.Type = parsedType

					// Check if it's an option type to set the optional flag
					if optType, ok := parsedType.(*types.OptionType); ok {
						fieldDef.Optional = true
						fieldDef.Type = optType
					}
				} else {
					// If parsing failed, use string as fallback
					if verbose {
						fmt.Printf("Warning: Using string type as fallback for field '%s'\n", name)
					}
					fieldDef.Type = &types.BasicType{TypeName: "string"}
				}
			}

			fields = append(fields, fieldDef)
		}

		// Store fields by type name
		typeMap[typeDef.Name] = fields
	}

	events := make([]types.Event, 0, len(rawEvents))
	for _, rawEvent := range rawEvents {
		var event struct {
			Name          string   `json:"name"`
			Discriminator []byte   `json:"discriminator,omitempty"`
			Docs          []string `json:"docs,omitempty"`
		}

		if err := json.Unmarshal(rawEvent, &event); err != nil {
			if verbose {
				fmt.Printf("Warning: Failed to parse event: %v\n", err)
			}
			continue
		}

		// Create a basic event
		eventDef := types.Event{
			Name:          event.Name,
			Discriminator: event.Discriminator,
			Docs:          event.Docs,
			Fields:        []types.Field{},
		}

		// First try to parse fields directly from the event
		var fields []map[string]interface{}
		if err := json.Unmarshal(rawEvent, &struct {
			Fields *[]map[string]interface{} `json:"fields"`
		}{
			Fields: &fields,
		}); err == nil && fields != nil {
			for _, field := range fields {
				name, _ := field["name"].(string)
				docs, _ := field["docs"].([]string)
				optional, _ := field["optional"].(bool)

				// Create a basic field
				fieldDef := types.Field{
					Name:     name,
					Docs:     docs,
					Optional: optional,
					Type:     &types.BasicType{TypeName: "string"}, // Default to string for now
				}

				// Parse the field type
				if fieldType, ok := field["type"]; ok {
					// Use the helper function to parse the type
					rawFieldType, _ := json.Marshal(fieldType)
					parsedType, success := ParseFieldType(name, rawFieldType, verbose)
					if success {
						fieldDef.Type = parsedType

						// Check if it's an option type to set the optional flag
						if optType, ok := parsedType.(*types.OptionType); ok {
							fieldDef.Optional = true
							fieldDef.Type = optType
						}
					} else {
						// If parsing failed, use string as fallback
						if verbose {
							fmt.Printf("Warning: Using string type as fallback for field '%s'\n", name)
						}
						fieldDef.Type = &types.BasicType{TypeName: "string"}
					}
				}

				eventDef.Fields = append(eventDef.Fields, fieldDef)
			}
		} else {
			// If no fields found directly in the event, look for a matching type
			if fields, ok := typeMap[event.Name]; ok {
				eventDef.Fields = fields
			}
		}

		if verbose {
			fmt.Printf("  - %s: %d fields\n", eventDef.Name, len(eventDef.Fields))

			// Log fields
			if len(eventDef.Fields) > 0 {
				fmt.Printf("    Fields:\n")
				for i, field := range eventDef.Fields {
					fmt.Printf("      %d. %s: %s\n", i+1, field.Name, field.Type.String())
				}
			}

			fmt.Println()
		}

		events = append(events, eventDef)
	}

	if verbose {
		fmt.Printf("Parsed %d events\n", len(events))
	}

	return events
}
