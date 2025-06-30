package parser

import (
	"encoding/json"
	"fmt"

	"github.com/iamrekas/anchor-go/pkg/idl/types"
)

// ParseAccounts parses accounts from raw account data and raw type data
func ParseAccounts(rawAccounts []json.RawMessage, rawTypes []json.RawMessage, verbose bool) []types.AccountDef {
	if verbose {
		fmt.Printf("Parsing %d accounts\n", len(rawAccounts))
	}

	if verbose {
		fmt.Println("Accounts:")
	}

	// First, parse all types to look up account types
	typeMap := make(map[string]types.TypeDef)
	for _, rawType := range rawTypes {
		var typeDef struct {
			Name string `json:"name"`
			Type struct {
				Kind   string            `json:"kind"`
				Fields []json.RawMessage `json:"fields,omitempty"`
			} `json:"type"`
		}

		if err := json.Unmarshal(rawType, &typeDef); err != nil {
			if verbose {
				fmt.Printf("Warning: Failed to parse type: %v\n", err)
			}
			continue
		}

		// Parse fields
		fields := make([]types.Field, 0, len(typeDef.Type.Fields))
		for _, rawField := range typeDef.Type.Fields {
			var field struct {
				Name string          `json:"name"`
				Type json.RawMessage `json:"type"`
				Docs []string        `json:"docs,omitempty"`
			}

			if err := json.Unmarshal(rawField, &field); err != nil {
				if verbose {
					fmt.Printf("Warning: Failed to parse field: %v\n", err)
				}
				continue
			}

			// Validate field name is not empty
			if field.Name == "" {
				if verbose {
					fmt.Printf("Warning: Field with empty name detected\n")
				}
				panic("Error: Field with empty name detected. Empty fields are not supported.")
			}

			// Create a basic field
			fieldDef := types.Field{
				Name:     field.Name,
				Docs:     field.Docs,
				Optional: false,
				Type:     &types.BasicType{TypeName: "string"}, // Default to string for now
			}

			// Use the helper function to parse the type
			parsedType, success := ParseFieldType(field.Name, field.Type, verbose)
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
					fmt.Printf("Warning: Using string type as fallback for field '%s'\n", field.Name)
				}
				fieldDef.Type = &types.BasicType{TypeName: "string"}
			}

			fields = append(fields, fieldDef)
		}

		// Create a type definition
		def := types.TypeDef{
			Name:   typeDef.Name,
			Kind:   types.TypeDefKind(typeDef.Type.Kind),
			Fields: fields,
		}

		typeMap[typeDef.Name] = def
	}

	// Now parse accounts
	accounts := make([]types.AccountDef, 0, len(rawAccounts))
	for _, rawAccount := range rawAccounts {
		// First, try to unmarshal as a raw JSON to check the structure
		var rawAccountObj map[string]interface{}
		if err := json.Unmarshal(rawAccount, &rawAccountObj); err != nil {
			if verbose {
				fmt.Printf("Warning: Failed to parse account as JSON: %v\n", err)
			}
			continue
		}

		// Check if the account has a "type" field that is an object
		if typeField, ok := rawAccountObj["type"]; ok {
			// Check if the type field is an object
			if _, ok := typeField.(map[string]interface{}); ok {
				// This is the case we need to handle differently
				// Parse the account with a different structure
				var account struct {
					Name          string   `json:"name"`
					Discriminator []byte   `json:"discriminator,omitempty"`
					Docs          []string `json:"docs,omitempty"`
					Type          struct {
						Kind   string            `json:"kind"`
						Fields []json.RawMessage `json:"fields,omitempty"`
					} `json:"type"`
				}

				if err := json.Unmarshal(rawAccount, &account); err != nil {
					if verbose {
						fmt.Printf("Warning: Failed to parse account with complex type: %v\n", err)
					}
					continue
				}

				// Validate account name is not empty
				if account.Name == "" {
					// Get the raw JSON for the account
					rawAccountJSON, _ := json.Marshal(rawAccount)

					if verbose {
						fmt.Printf("Warning: Account with empty name detected\n")
						fmt.Printf("JSON data: %s\n", string(rawAccountJSON))
					}
					continue
				}

				// Parse fields from the type
				fields := make([]types.Field, 0, len(account.Type.Fields))
				for _, rawField := range account.Type.Fields {
					var field struct {
						Name string          `json:"name"`
						Type json.RawMessage `json:"type"`
						Docs []string        `json:"docs,omitempty"`
					}

					if err := json.Unmarshal(rawField, &field); err != nil {
						if verbose {
							fmt.Printf("Warning: Failed to parse field in account type: %v\n", err)
						}
						continue
					}

					// Create a basic field
					fieldDef := types.Field{
						Name:     field.Name,
						Docs:     field.Docs,
						Optional: false,
						Type:     &types.BasicType{TypeName: "string"}, // Default to string for now
					}

					// Use the helper function to parse the type
					parsedType, success := ParseFieldType(field.Name, field.Type, verbose)
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
							fmt.Printf("Warning: Using string type as fallback for field '%s'\n", field.Name)
						}
						fieldDef.Type = &types.BasicType{TypeName: "string"}
					}

					fields = append(fields, fieldDef)
				}

				// Create a type definition for this account
				typeDef := types.TypeDef{
					Name:   account.Name,
					Kind:   types.TypeDefKind(account.Type.Kind),
					Fields: fields,
				}

				// Create an account definition
				accDef := types.AccountDef{
					Name:          account.Name,
					Discriminator: account.Discriminator,
					Docs:          account.Docs,
					Type:          typeDef,
				}

				if verbose {
					fmt.Printf("  - %s: %d fields\n", accDef.Name, len(accDef.Type.Fields))
				}

				accounts = append(accounts, accDef)
				continue
			}
		}

		// If we get here, try the original approach
		var account struct {
			Name          string   `json:"name"`
			Discriminator []byte   `json:"discriminator,omitempty"`
			Type          *string  `json:"type,omitempty"`
			Docs          []string `json:"docs,omitempty"`
		}

		if err := json.Unmarshal(rawAccount, &account); err != nil {
			if verbose {
				fmt.Printf("Warning: Failed to parse account: %v\n", err)
			}
			continue
		}

		// Look up the type definition for this account
		typeDef, ok := typeMap[account.Name]
		if !ok {
			// If no type definition is found, create a basic one
			typeDef = types.TypeDef{
				Name:   account.Name,
				Kind:   types.TypeDefKindStruct,
				Fields: []types.Field{},
			}
		}

		// Validate account name is not empty
		if account.Name == "" {
			// Get the raw JSON for the account
			rawAccount, _ := json.Marshal(rawAccount)

			if verbose {
				fmt.Printf("Warning: Account with empty name detected\n")
				fmt.Printf("JSON data: %s\n", string(rawAccount))
			}
			panic(fmt.Sprintf("Error: Account with empty name detected. JSON data: %s. Empty account names are not supported.",
				string(rawAccount)))
		}

		// Create an account definition
		accDef := types.AccountDef{
			Name:          account.Name,
			Discriminator: account.Discriminator,
			Docs:          account.Docs,
			Type:          typeDef,
		}

		if verbose {
			fmt.Printf("  - %s: %d fields\n", accDef.Name, len(accDef.Type.Fields))

			// Log fields
			if len(accDef.Type.Fields) > 0 {
				fmt.Printf("    Fields:\n")
				for i, field := range accDef.Type.Fields {
					fmt.Printf("      %d. %s: %s\n", i+1, field.Name, field.Type.String())
				}
			}

			fmt.Println()
		}

		accounts = append(accounts, accDef)
	}

	if verbose {
		fmt.Printf("Parsed %d accounts\n", len(accounts))
	}

	return accounts
}
