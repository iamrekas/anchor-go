package parser

import (
	"encoding/json"
	"fmt"

	"github.com/iamrekas/anchor-go/pkg/idl/types"
)

// ParseConstants parses constants from raw constant data
func ParseConstants(rawConstants []json.RawMessage, verbose bool) []types.Constant {
	if verbose {
		fmt.Printf("Parsing %d constants\n", len(rawConstants))
	}

	if verbose {
		fmt.Println("Constants:")
	}

	constants := make([]types.Constant, 0, len(rawConstants))
	for _, rawConstant := range rawConstants {
		var constant struct {
			Name  string          `json:"name"`
			Type  json.RawMessage `json:"type"`
			Value string          `json:"value"`
		}

		if err := json.Unmarshal(rawConstant, &constant); err != nil {
			if verbose {
				fmt.Printf("Warning: Failed to parse constant: %v\n", err)
			}
			continue
		}

		// Create a basic constant
		constDef := types.Constant{
			Name:  constant.Name,
			Type:  &types.BasicType{TypeName: "string"}, // Default to string for now
			Value: constant.Value,
		}

		// Parse the type
		parsedType, success := ParseFieldType(constant.Name, constant.Type, verbose)
		if success {
			constDef.Type = parsedType
		} else {
			// If parsing failed, use string as fallback
			if verbose {
				fmt.Printf("Warning: Using string type as fallback for constant '%s'\n", constant.Name)
			}
		}

		if verbose {
			fmt.Printf("  - %s: %s = %s\n", constDef.Name, constDef.Type.String(), constDef.Value)
		}

		constants = append(constants, constDef)
	}

	if verbose {
		fmt.Printf("Parsed %d constants\n", len(constants))
	}

	return constants
}
