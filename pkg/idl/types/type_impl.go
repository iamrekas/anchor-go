package types

import (
	"fmt"
	"strconv"
	"strings"
)

// BasicType represents a basic type like bool, u8, i64, etc.
type BasicType struct {
	TypeName string
}

// IsArray implements Type.IsArray
func (t *BasicType) IsArray() bool {
	return false
}

// IsOption implements Type.IsOption
func (t *BasicType) IsOption() bool {
	return false
}

// IsVector implements Type.IsVector
func (t *BasicType) IsVector() bool {
	return false
}

// IsDefined implements Type.IsDefined
func (t *BasicType) IsDefined() bool {
	return false
}

// IsBasic implements Type.IsBasic
func (t *BasicType) IsBasic() bool {
	return true
}

// String implements Type.String
func (t *BasicType) String() string {
	return t.TypeName
}

// ArrayType represents an array type like [5]u8
type ArrayType struct {
	ElemType  Type
	ArraySize int
}

// IsArray implements Type.IsArray
func (t *ArrayType) IsArray() bool {
	return true
}

// IsOption implements Type.IsOption
func (t *ArrayType) IsOption() bool {
	return false
}

// IsVector implements Type.IsVector
func (t *ArrayType) IsVector() bool {
	return false
}

// IsDefined implements Type.IsDefined
func (t *ArrayType) IsDefined() bool {
	return false
}

// IsBasic implements Type.IsBasic
func (t *ArrayType) IsBasic() bool {
	return false
}

// String implements Type.String
func (t *ArrayType) String() string {
	return fmt.Sprintf("[%d]%s", t.ArraySize, t.ElemType.String())
}

// ElementType returns the element type of the array
func (t *ArrayType) ElementType() Type {
	return t.ElemType
}

// Size returns the size of the array
func (t *ArrayType) Size() int {
	return t.ArraySize
}

// VectorType represents a vector type like Vec<u8>
type VectorType struct {
	ElementType Type
}

// IsArray implements Type.IsArray
func (t *VectorType) IsArray() bool {
	return false
}

// IsOption implements Type.IsOption
func (t *VectorType) IsOption() bool {
	return false
}

// IsVector implements Type.IsVector
func (t *VectorType) IsVector() bool {
	return true
}

// IsDefined implements Type.IsDefined
func (t *VectorType) IsDefined() bool {
	return false
}

// IsBasic implements Type.IsBasic
func (t *VectorType) IsBasic() bool {
	return false
}

// String implements Type.String
func (t *VectorType) String() string {
	return fmt.Sprintf("Vec<%s>", t.ElementType.String())
}

// OptionType represents an option type like Option<u8>
type OptionType struct {
	ElementType Type
}

// IsArray implements Type.IsArray
func (t *OptionType) IsArray() bool {
	return false
}

// IsOption implements Type.IsOption
func (t *OptionType) IsOption() bool {
	return true
}

// IsVector implements Type.IsVector
func (t *OptionType) IsVector() bool {
	return false
}

// IsDefined implements Type.IsDefined
func (t *OptionType) IsDefined() bool {
	return false
}

// IsBasic implements Type.IsBasic
func (t *OptionType) IsBasic() bool {
	return false
}

// String implements Type.String
func (t *OptionType) String() string {
	return fmt.Sprintf("Option<%s>", t.ElementType.String())
}

// DefinedType represents a user-defined type like MyStruct
type DefinedType struct {
	Name string
}

// IsArray implements Type.IsArray
func (t *DefinedType) IsArray() bool {
	return false
}

// IsOption implements Type.IsOption
func (t *DefinedType) IsOption() bool {
	return false
}

// IsVector implements Type.IsVector
func (t *DefinedType) IsVector() bool {
	return false
}

// IsDefined implements Type.IsDefined
func (t *DefinedType) IsDefined() bool {
	return true
}

// IsBasic implements Type.IsBasic
func (t *DefinedType) IsBasic() bool {
	return false
}

// String implements Type.String
func (t *DefinedType) String() string {
	return t.Name
}

// Common basic types
var (
	TypeBool      = &BasicType{TypeName: "bool"}
	TypeU8        = &BasicType{TypeName: "u8"}
	TypeI8        = &BasicType{TypeName: "i8"}
	TypeU16       = &BasicType{TypeName: "u16"}
	TypeI16       = &BasicType{TypeName: "i16"}
	TypeU32       = &BasicType{TypeName: "u32"}
	TypeI32       = &BasicType{TypeName: "i32"}
	TypeU64       = &BasicType{TypeName: "u64"}
	TypeI64       = &BasicType{TypeName: "i64"}
	TypeU128      = &BasicType{TypeName: "u128"}
	TypeI128      = &BasicType{TypeName: "i128"}
	TypeBytes     = &BasicType{TypeName: "bytes"}
	TypeString    = &BasicType{TypeName: "string"}
	TypePublicKey = &BasicType{TypeName: "publicKey"}
	TypeF32       = &BasicType{TypeName: "f32"}
	TypeF64       = &BasicType{TypeName: "f64"}

	// Custom additions
	TypeUnixTimestamp = &BasicType{TypeName: "unixTimestamp"}
	TypeHash          = &BasicType{TypeName: "hash"}
	TypeDuration      = &BasicType{TypeName: "duration"}
)

// ParseType parses a type string into a Type
func ParseType(typeStr string) (Type, error) {
	typeStr = strings.TrimSpace(typeStr)

	// Check for array type
	if strings.HasPrefix(typeStr, "[") && strings.Contains(typeStr, "]") {
		closeBracket := strings.Index(typeStr, "]")
		sizeStr := typeStr[1:closeBracket]
		size, err := strconv.Atoi(sizeStr)
		if err != nil {
			return nil, fmt.Errorf("invalid array size: %s", sizeStr)
		}

		elemTypeStr := typeStr[closeBracket+1:]
		elemType, err := ParseType(elemTypeStr)
		if err != nil {
			return nil, err
		}

		return &ArrayType{
			ElemType:  elemType,
			ArraySize: size,
		}, nil
	}

	// Check for vector type
	if strings.HasPrefix(typeStr, "Vec<") && strings.HasSuffix(typeStr, ">") {
		elemTypeStr := typeStr[4 : len(typeStr)-1]
		elemType, err := ParseType(elemTypeStr)
		if err != nil {
			return nil, err
		}

		return &VectorType{
			ElementType: elemType,
		}, nil
	}

	// Check for option type
	if strings.HasPrefix(typeStr, "Option<") && strings.HasSuffix(typeStr, ">") {
		elemTypeStr := typeStr[7 : len(typeStr)-1]
		elemType, err := ParseType(elemTypeStr)
		if err != nil {
			return nil, err
		}

		return &OptionType{
			ElementType: elemType,
		}, nil
	}

	// Check for basic types
	switch typeStr {
	case "bool":
		return TypeBool, nil
	case "u8":
		return TypeU8, nil
	case "i8":
		return TypeI8, nil
	case "u16":
		return TypeU16, nil
	case "i16":
		return TypeI16, nil
	case "u32":
		return TypeU32, nil
	case "i32":
		return TypeI32, nil
	case "u64":
		return TypeU64, nil
	case "i64":
		return TypeI64, nil
	case "u128":
		return TypeU128, nil
	case "i128":
		return TypeI128, nil
	case "bytes":
		return TypeBytes, nil
	case "string":
		return TypeString, nil
	case "pubkey", "publicKey":
		return TypePublicKey, nil
	case "f32":
		return TypeF32, nil
	case "f64":
		return TypeF64, nil
	case "unixTimestamp":
		return TypeUnixTimestamp, nil
	case "hash":
		return TypeHash, nil
	case "duration":
		return TypeDuration, nil
	default:
		// Assume it's a defined type
		return &DefinedType{
			Name: typeStr,
		}, nil
	}
}
