package parser

import (
	"encoding/json"
	"fmt"

	"github.com/iamrekas/anchor-go/pkg/idl/types"
)

// ParseTypes parses types from raw type data
func ParseTypes(rawTypes []json.RawMessage, verbose bool) []types.TypeDef {
	if verbose {
		fmt.Printf("Parsing %d types\n", len(rawTypes))
		fmt.Println("Types:")
	}

	// First, create a map of type name to raw type for lookup
	typeMap := make(map[string]json.RawMessage)
	for _, rawType := range rawTypes {
		var typeDef struct {
			Name string `json:"name"`
		}
		if err := json.Unmarshal(rawType, &typeDef); err == nil && typeDef.Name != "" {
			typeMap[typeDef.Name] = rawType
		}
	}

	typeDefs := make([]types.TypeDef, 0, len(rawTypes))
	for _, rawType := range rawTypes {
		var typeDef struct {
			Name          string   `json:"name"`
			Discriminator []byte   `json:"discriminator,omitempty"`
			Docs          []string `json:"docs,omitempty"`
			Type          struct {
				Kind     string          `json:"kind"`
				Fields   json.RawMessage `json:"fields,omitempty"`
				Variants json.RawMessage `json:"variants,omitempty"`
			} `json:"type"`
		}

		if err := json.Unmarshal(rawType, &typeDef); err != nil {
			if verbose {
				fmt.Printf("Warning: Failed to parse type: %v\n", err)
			}
			continue
		}

		// Create a basic type definition
		def := types.TypeDef{
			Name:          typeDef.Name,
			Kind:          types.TypeDefKind(typeDef.Type.Kind),
			Discriminator: typeDef.Discriminator,
			Docs:          typeDef.Docs,
			Fields:        []types.Field{},
			Variants:      []types.EnumVariant{},
		}

		// Parse fields if it's a struct
		if typeDef.Type.Kind == "struct" {
			var fields []map[string]interface{}
			if err := json.Unmarshal(typeDef.Type.Fields, &fields); err == nil && fields != nil {
				for _, field := range fields {
					name, _ := field["name"].(string)
					docs, _ := field["docs"].([]string)
					optional, _ := field["optional"].(bool)

					// Validate field name is not empty
					if name == "" {
						// Get the raw JSON for the field
						rawField, _ := json.Marshal(field)

						if verbose {
							fmt.Printf("Warning: Field with empty name detected in type '%s'\n", typeDef.Name)
							fmt.Printf("JSON data: %s\n", string(rawField))
						}
						panic(fmt.Sprintf("Error: Field with empty name detected in type '%s'. JSON data: %s. Empty field names are not supported.",
							typeDef.Name, string(rawField)))
					}

					// Validate field type is present
					fieldType, ok := field["type"]
					if !ok || fieldType == nil {
						// Get the raw JSON for the field
						rawField, _ := json.Marshal(field)

						if verbose {
							fmt.Printf("Warning: Field '%s' in type '%s' has no type definition\n", name, typeDef.Name)
							fmt.Printf("JSON data: %s\n", string(rawField))
						}
						panic(fmt.Sprintf("Error: Field '%s' in type '%s' has no type definition. JSON data: %s. Fields must have a type.",
							name, typeDef.Name, string(rawField)))
					}

					// Create a basic field
					fieldDef := types.Field{
						Name:     name,
						Docs:     docs,
						Optional: optional,
						Type:     &types.BasicType{TypeName: "string"}, // Default to string for now
					}

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

					def.Fields = append(def.Fields, fieldDef)
				}
			}
		}

		// Parse variants if it's an enum
		if typeDef.Type.Kind == "enum" {
			var variants []map[string]interface{}
			if err := json.Unmarshal(typeDef.Type.Variants, &variants); err == nil && variants != nil {
				for _, variant := range variants {
					name, _ := variant["name"].(string)
					docs, _ := variant["docs"].([]string)

					// Create a basic variant
					variantDef := types.EnumVariant{
						Name:   name,
						Docs:   docs,
						Fields: []types.Field{},
					}

					// Parse fields if available
					if fieldsRaw, ok := variant["fields"].([]interface{}); ok {
						for _, fieldRaw := range fieldsRaw {
							// Check if the field is a map (has name and type)
							if field, ok := fieldRaw.(map[string]interface{}); ok {
								name, _ := field["name"].(string)
								docs, _ := field["docs"].([]string)
								optional, _ := field["optional"].(bool)

								// Get variant name
								variantName, _ := variant["name"].(string)

								// Check if this is a named field or a direct defined type
								_, hasName := field["name"]
								_, hasType := field["type"]
								definedData, hasDefined := field["defined"]

								// If it's a direct defined type without a name field
								if !hasName && !hasType && hasDefined {
									// It's a defined type directly in the field
									definedMap, ok := definedData.(map[string]interface{})
									if ok {
										typeName, ok := definedMap["name"].(string)
										if !ok || typeName == "" {
											// Get the raw JSON for the field
											rawField, _ := json.Marshal(field)

											if verbose {
												fmt.Printf("Warning: Field in variant '%s' of enum '%s' has empty defined type name\n", variantName, typeDef.Name)
												fmt.Printf("JSON data: %s\n", string(rawField))
											}
											panic(fmt.Sprintf("Error: Field in variant '%s' of enum '%s' has empty defined type name. JSON data: %s. Empty type names are not supported.",
												variantName, typeDef.Name, string(rawField)))
										}

										// Create a field with the defined type but no name
										fieldDef := types.Field{
											Name: "", // No name for this field
											Type: &types.DefinedType{Name: typeName},
										}
										variantDef.Fields = append(variantDef.Fields, fieldDef)
										continue
									}
								}

								// For regular named fields, validate name is not empty
								if hasName && name == "" {
									// Get the raw JSON for the field
									rawField, _ := json.Marshal(field)

									if verbose {
										fmt.Printf("Warning: Field with empty name detected in variant '%s' of enum '%s'\n", variantName, typeDef.Name)
										fmt.Printf("JSON data: %s\n", string(rawField))
									}
									panic(fmt.Sprintf("Error: Field with empty name detected in variant '%s' of enum '%s'. JSON data: %s. Empty field names are not supported.",
										variantName, typeDef.Name, string(rawField)))
								}

								// Validate field type is present for named fields
								fieldType, ok := field["type"]
								if hasName && (!ok || fieldType == nil) {
									// Get the raw JSON for the field
									rawField, _ := json.Marshal(field)

									if verbose {
										fmt.Printf("Warning: Field '%s' in variant '%s' of enum '%s' has no type definition\n", name, variantName, typeDef.Name)
										fmt.Printf("JSON data: %s\n", string(rawField))
									}
									panic(fmt.Sprintf("Error: Field '%s' in variant '%s' of enum '%s' has no type definition. JSON data: %s. Fields must have a type.",
										name, variantName, typeDef.Name, string(rawField)))
								}

								// Create a basic field
								fieldDef := types.Field{
									Name:     name,
									Docs:     docs,
									Optional: optional,
									Type:     &types.BasicType{TypeName: "string"}, // Default to string for now
								}

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

								variantDef.Fields = append(variantDef.Fields, fieldDef)
							} else {
								// Handle the case where the field is a defined type directly
								var definedType struct {
									Defined struct {
										Name string `json:"name"`
									} `json:"defined"`
								}
								fieldRawJSON, _ := json.Marshal(fieldRaw)
								if err := json.Unmarshal(fieldRawJSON, &definedType); err == nil && definedType.Defined.Name != "" {
									// Create a field with the defined type
									fieldDef := types.Field{
										Name: "", // No name for this field
										Type: &types.DefinedType{Name: definedType.Defined.Name},
									}
									variantDef.Fields = append(variantDef.Fields, fieldDef)
								}
							}
						}
					}

					def.Variants = append(def.Variants, variantDef)
				}
			}
		}

		// Check for empty types and handle them appropriately
		if typeDef.Type.Kind == "struct" && len(def.Fields) == 0 {
			// Get the raw JSON for the type
			rawTypeStr, _ := json.Marshal(rawType)

			if verbose {
				fmt.Printf("Warning: Struct type '%s' has no fields\n", typeDef.Name)
				fmt.Printf("JSON data: %s\n", string(rawTypeStr))
			}
			panic(fmt.Sprintf("Error: Struct type '%s' has no fields. JSON data: %s. Empty struct types are not supported.",
				typeDef.Name, string(rawTypeStr)))
		}

		if typeDef.Type.Kind == "enum" && len(def.Variants) == 0 {
			// Get the raw JSON for the type
			rawTypeStr, _ := json.Marshal(rawType)

			if verbose {
				fmt.Printf("Warning: Enum type '%s' has no variants\n", typeDef.Name)
				fmt.Printf("JSON data: %s\n", string(rawTypeStr))
			}
			panic(fmt.Sprintf("Error: Enum type '%s' has no variants. JSON data: %s. Empty enum types are not supported.",
				typeDef.Name, string(rawTypeStr)))
		}

		if verbose {
			// Format the output to match the account parser logs
			if def.Kind == "struct" {
				fmt.Printf("  - %s: %d fields\n", def.Name, len(def.Fields))

				// Log fields
				if len(def.Fields) > 0 {
					fmt.Printf("    Fields:\n")
					for i, field := range def.Fields {
						fmt.Printf("      %d. %s: %s", i+1, field.Name, field.Type.String())

						if field.Optional {
							fmt.Printf(" (optional)")
						}

						fmt.Println()
					}
				}

				fmt.Println()
			} else if def.Kind == "enum" {
				fmt.Printf("  - %s: %d variants\n", def.Name, len(def.Variants))

				// Log variants
				if len(def.Variants) > 0 {
					fmt.Printf("    Variants:\n")
					for i, variant := range def.Variants {
						fmt.Printf("      %d. %s", i+1, variant.Name)

						if len(variant.Fields) > 0 {
							fmt.Printf(": %d fields\n", len(variant.Fields))

							// Log variant fields
							fmt.Printf("        Fields:\n")
							for j, field := range variant.Fields {
								if field.Name == "" {
									// For unnamed fields (like in enum variants), just show the type
									fmt.Printf("          %d. %s", j+1, field.Type.String())
								} else {
									fmt.Printf("          %d. %s: %s", j+1, field.Name, field.Type.String())
								}

								if field.Optional {
									fmt.Printf(" (optional)")
								}

								fmt.Println()
							}
						} else {
							fmt.Println()
						}
					}
				}

				fmt.Println()
			} else {
				// For other types (like aliases)
				fmt.Printf("  - %s: %s\n", def.Name, def.Kind)
				fmt.Println()
			}
		}

		typeDefs = append(typeDefs, def)
	}

	if verbose {
		fmt.Printf("Parsed %d types\n", len(typeDefs))
	}

	return typeDefs
}
