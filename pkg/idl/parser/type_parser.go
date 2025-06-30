package parser

import (
	"encoding/json"
	"fmt"

	"github.com/iamrekas/anchor-go/pkg/idl/types"
)

// ParseFieldType parses a field type from raw JSON data
func ParseFieldType(fieldName string, fieldType json.RawMessage, verbose bool) (types.Type, bool) {
	// Try to parse as a basic type first
	var typeStr string
	if err := json.Unmarshal(fieldType, &typeStr); err == nil {
		// Basic type
		return &types.BasicType{TypeName: typeStr}, true
	}

	// Get the raw JSON for error messages
	rawFieldType, _ := json.Marshal(fieldType)

	// Try to parse as a complex type
	var typeObj map[string]interface{}
	if err := json.Unmarshal(fieldType, &typeObj); err == nil {
		// Check if it's an array type
		if arrayData, ok := typeObj["array"]; ok {
			// It's an array type
			arraySlice, ok := arrayData.([]interface{})
			if ok && len(arraySlice) == 2 {
				// First, check if element type is a defined type
				elementTypeRaw, _ := json.Marshal(arraySlice[0])

				// Try to parse as a defined type
				var definedType struct {
					Defined interface{} `json:"defined"`
				}

				if err := json.Unmarshal(elementTypeRaw, &definedType); err == nil && definedType.Defined != nil {
					// It's a defined type
					var typeName string

					// Handle both string and object formats
					switch v := definedType.Defined.(type) {
					case string:
						typeName = v
					case map[string]interface{}:
						if name, ok := v["name"].(string); ok {
							typeName = name
						}
					}

					if typeName != "" {
						// Parse size
						sizeRaw, _ := json.Marshal(arraySlice[1])
						var size int
						if err := json.Unmarshal(sizeRaw, &size); err != nil {
							if verbose {
								fmt.Printf("Warning: Failed to parse array size for field '%s'\n", fieldName)
								fmt.Printf("JSON data: %s\n", string(rawFieldType))
							}
							return nil, false
						}

						// Create array type with a defined type as the element type
						return &types.ArrayType{
							ElemType:  &types.DefinedType{Name: typeName},
							ArraySize: size,
						}, true
					}
				}

				// If not a defined type, try as a basic type
				var elementType string
				if err := json.Unmarshal(elementTypeRaw, &elementType); err == nil {
					// Parse size
					sizeRaw, _ := json.Marshal(arraySlice[1])
					var size int
					if err := json.Unmarshal(sizeRaw, &size); err != nil {
						if verbose {
							fmt.Printf("Warning: Failed to parse array size for field '%s'\n", fieldName)
							fmt.Printf("JSON data: %s\n", string(rawFieldType))
						}
						return nil, false
					}

					// Create array type with a basic type as the element type
					return &types.ArrayType{
						ElemType:  &types.BasicType{TypeName: elementType},
						ArraySize: size,
					}, true
				}
			}
		}

		// Check if it's a defined type
		if definedData, ok := typeObj["defined"]; ok {
			// It's a defined type
			var typeName string

			// Handle both string and object formats
			switch v := definedData.(type) {
			case string:
				typeName = v
			case map[string]interface{}:
				if name, ok := v["name"].(string); ok {
					typeName = name
				}
			}

			if typeName != "" {
				return &types.DefinedType{Name: typeName}, true
			}
		}

		// Check if it's a vector type
		if vecData, ok := typeObj["vec"]; ok {
			// Try to handle vector of defined type
			vecMap, ok := vecData.(map[string]interface{})
			if ok && vecMap["defined"] != nil {
				var typeName string

				// Handle both string and object formats
				switch v := vecMap["defined"].(type) {
				case string:
					typeName = v
				case map[string]interface{}:
					if name, ok := v["name"].(string); ok {
						typeName = name
					}
				}

				if typeName != "" {
					return &types.VectorType{ElementType: &types.DefinedType{Name: typeName}}, true
				}
			}

			// Try as a basic type
			var elemType string
			elemTypeRaw, _ := json.Marshal(vecData)
			if err := json.Unmarshal(elemTypeRaw, &elemType); err == nil {
				return &types.VectorType{ElementType: &types.BasicType{TypeName: elemType}}, true
			}
		}

		// Check if it's an option type
		if optionData, ok := typeObj["option"]; ok {
			// Try to handle option of defined type
			optionMap, ok := optionData.(map[string]interface{})
			if ok && optionMap["defined"] != nil {
				var typeName string

				// Handle both string and object formats
				switch v := optionMap["defined"].(type) {
				case string:
					typeName = v
				case map[string]interface{}:
					if name, ok := v["name"].(string); ok {
						typeName = name
					}
				}

				if typeName != "" {
					return &types.OptionType{ElementType: &types.DefinedType{Name: typeName}}, true
				}
			}

			// Try as a basic type
			var elemType string
			elemTypeRaw, _ := json.Marshal(optionData)
			if err := json.Unmarshal(elemTypeRaw, &elemType); err == nil {
				return &types.OptionType{ElementType: &types.BasicType{TypeName: elemType}}, true
			}
		}
	}

	// Try the old way as a fallback
	var definedType struct {
		Defined struct {
			Name string `json:"name"`
		} `json:"defined"`
	}

	var arrayType struct {
		Array []json.RawMessage `json:"array"`
	}

	if err := json.Unmarshal(fieldType, &definedType); err == nil {
		if definedType.Defined.Name != "" {
			return &types.DefinedType{Name: definedType.Defined.Name}, true
		}
	} else if err := json.Unmarshal(fieldType, &arrayType); err == nil && len(arrayType.Array) == 2 {
		// Try to parse element type as a defined type first
		var definedElemType struct {
			Defined struct {
				Name string `json:"name"`
			} `json:"defined"`
		}

		if err := json.Unmarshal(arrayType.Array[0], &definedElemType); err == nil && definedElemType.Defined.Name != "" {
			// Parse size
			var size int
			if err := json.Unmarshal(arrayType.Array[1], &size); err == nil {
				return &types.ArrayType{
					ElemType:  &types.DefinedType{Name: definedElemType.Defined.Name},
					ArraySize: size,
				}, true
			}
		}

		// If not a defined type, try as a basic type
		var elementType string
		var size int

		// Parse element type
		if err := json.Unmarshal(arrayType.Array[0], &elementType); err == nil {
			// Parse size
			if err := json.Unmarshal(arrayType.Array[1], &size); err == nil {
				// Create array type with a basic type as the element type
				return &types.ArrayType{
					ElemType:  &types.BasicType{TypeName: elementType},
					ArraySize: size,
				}, true
			}
		}
	}

	if verbose {
		fmt.Printf("Warning: Failed to parse type for field '%s'\n", fieldName)
		fmt.Printf("JSON data: %s\n", string(rawFieldType))
	}

	return nil, false
}
