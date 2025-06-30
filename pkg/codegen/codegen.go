// Package codegen provides helpers for generating Go code.
package codegen

import (
	"fmt"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/iamrekas/anchor-go/internal/utils"
)

// Common package paths
const (
	PkgSolanaGo       = "github.com/gagliardetto/solana-go"
	PkgSolanaGoText   = "github.com/gagliardetto/solana-go/text"
	PkgBinary         = "github.com/gagliardetto/binary"
	PkgTreeout        = "github.com/gagliardetto/treeout"
	PkgFormat         = "github.com/gagliardetto/solana-go/text/format"
	PkgGoFuzz         = "github.com/gagliardetto/gofuzz"
	PkgTestifyRequire = "github.com/stretchr/testify/require"
)

// File represents a Go source file
type File struct {
	*jen.File
	Name string
	Path string
}

// NewFile creates a new Go source file
func NewFile(name, packageName string) *File {
	return &File{
		File: jen.NewFile(packageName),
		Name: name,
		Path: "",
	}
}

// SetPath sets the file path
func (f *File) SetPath(path string) *File {
	f.Path = path
	return f
}

// AddHeaderComment adds a header comment to the file
func (f *File) AddHeaderComment(comment string) *File {
	f.HeaderComment(comment)
	return f
}

// AddImports adds imports to the file
func (f *File) AddImports(imports map[string]string) *File {
	for alias, path := range imports {
		if alias == "" {
			f.ImportName(path, "")
		} else {
			f.ImportAlias(path, alias)
		}
	}
	return f
}

// AddStruct adds a struct definition to the file
func (f *File) AddStruct(name string, fields []Field) *File {
	structDef := jen.Type().Id(name).Struct()

	for _, field := range fields {
		fieldDef := jen.Id(field.Name)

		if field.IsPointer {
			fieldDef.Op("*")
		}

		fieldDef.Add(field.Type)

		if field.Tag != "" {
			fieldDef.Tag(map[string]string{"bin": field.Tag})
		}

		if field.Comment != "" {
			fieldDef.Comment(field.Comment)
		}

		structDef.Add(fieldDef)
	}

	f.Add(structDef)
	return f
}

// AddMethod adds a method to the file
func (f *File) AddMethod(receiver, name string, isPointer bool, params, results []Field, body ...jen.Code) *File {
	method := jen.Func()

	// Add receiver
	if isPointer {
		method.Params(jen.Id("r").Op("*").Id(receiver))
	} else {
		method.Params(jen.Id("r").Id(receiver))
	}

	method.Id(name)

	// Add parameters
	if len(params) > 0 {
		method.Params(fieldSliceToParams(params)...)
	} else {
		method.Params()
	}

	// Add results
	if len(results) > 0 {
		method.Params(fieldSliceToParams(results)...)
	} else {
		method.Params()
	}

	// Add body
	method.Block(body...)

	f.Add(method)
	return f
}

// AddFunction adds a function to the file
func (f *File) AddFunction(name string, params, results []Field, body ...jen.Code) *File {
	function := jen.Func().Id(name)

	// Add parameters
	if len(params) > 0 {
		function.Params(fieldSliceToParams(params)...)
	} else {
		function.Params()
	}

	// Add results
	if len(results) > 0 {
		function.Params(fieldSliceToParams(results)...)
	} else {
		function.Params()
	}

	// Add body
	function.Block(body...)

	f.Add(function)
	return f
}

// AddConstant adds a constant to the file
func (f *File) AddConstant(name string, value interface{}) *File {
	f.Const().Id(name).Op("=").Lit(value)
	return f
}

// AddConstantGroup adds a group of constants to the file
func (f *File) AddConstantGroup(constants map[string]interface{}) *File {
	constGroup := jen.Const().Parens(jen.Line())

	for name, value := range constants {
		constGroup.Id(name).Op("=").Lit(value).Line()
	}

	f.Add(constGroup)
	return f
}

// AddVar adds a variable to the file
func (f *File) AddVar(name string, typeName jen.Code, value jen.Code) *File {
	if value == nil {
		f.Var().Id(name).Add(typeName)
	} else {
		f.Var().Id(name).Add(typeName).Op("=").Add(value)
	}
	return f
}

// AddVarGroup adds a group of variables to the file
func (f *File) AddVarGroup(vars map[string]jen.Code) *File {
	varGroup := jen.Var().Parens(jen.Line())

	for name, value := range vars {
		varGroup.Id(name).Add(value).Line()
	}

	f.Add(varGroup)
	return f
}

// Field represents a field in a struct or a parameter in a function
type Field struct {
	Name      string
	Type      jen.Code
	Tag       string
	Comment   string
	IsPointer bool
}

// NewField creates a new field
func NewField(name string, typeCode jen.Code) Field {
	return Field{
		Name: name,
		Type: typeCode,
	}
}

// WithTag adds a tag to the field
func (f Field) WithTag(tag string) Field {
	f.Tag = tag
	return f
}

// WithComment adds a comment to the field
func (f Field) WithComment(comment string) Field {
	f.Comment = comment
	return f
}

// AsPointer marks the field as a pointer
func (f Field) AsPointer() Field {
	f.IsPointer = true
	return f
}

// fieldSliceToParams converts a slice of fields to a slice of jen.Code
func fieldSliceToParams(fields []Field) []jen.Code {
	params := make([]jen.Code, len(fields))

	for i, field := range fields {
		param := jen.Id(field.Name)

		if field.IsPointer {
			param.Op("*")
		}

		param.Add(field.Type)
		params[i] = param
	}

	return params
}

// FormatFieldName formats a field name according to Go conventions
func FormatFieldName(name string) string {
	return utils.ToCamel(name)
}

// FormatMethodName formats a method name according to Go conventions
func FormatMethodName(name string) string {
	return utils.ToCamel(name)
}

// FormatPackageName formats a package name according to Go conventions
func FormatPackageName(name string) string {
	return strings.ToLower(utils.ToSnakeCase(name))
}

// FormatFileName formats a file name according to Go conventions
func FormatFileName(name string) string {
	return strings.ToLower(utils.ToSnakeCase(name))
}

// FormatConstantName formats a constant name according to Go conventions
func FormatConstantName(name string) string {
	return strings.ToUpper(utils.ToSnakeCase(name))
}

// FormatVariableName formats a variable name according to Go conventions
func FormatVariableName(name string) string {
	return utils.ToLowerCamel(name)
}

// FormatInterfaceName formats an interface name according to Go conventions
func FormatInterfaceName(name string) string {
	return utils.ToCamel(name)
}

// FormatStructName formats a struct name according to Go conventions
func FormatStructName(name string) string {
	return utils.ToCamel(name)
}

// FormatEnumName formats an enum name according to Go conventions
func FormatEnumName(name string) string {
	return utils.ToCamel(name)
}

// FormatEnumValueName formats an enum value name according to Go conventions
func FormatEnumValueName(enumName, valueName string) string {
	return fmt.Sprintf("%s%s", enumName, utils.ToCamel(valueName))
}

// BasicType returns a jen.Code for a basic type
func BasicType(typeName string) jen.Code {
	switch typeName {
	case "bool":
		return jen.Bool()
	case "u8":
		return jen.Uint8()
	case "i8":
		return jen.Int8()
	case "u16":
		return jen.Uint16()
	case "i16":
		return jen.Int16()
	case "u32":
		return jen.Uint32()
	case "i32":
		return jen.Int32()
	case "u64":
		return jen.Uint64()
	case "i64":
		return jen.Int64()
	case "u128":
		return jen.Qual(PkgBinary, "Uint128")
	case "i128":
		return jen.Qual(PkgBinary, "Int128")
	case "bytes":
		return jen.Index().Byte()
	case "string":
		return jen.String()
	case "publicKey", "pubkey":
		return jen.Qual(PkgSolanaGo, "PublicKey")
	case "f32":
		return jen.Float32()
	case "f64":
		return jen.Float64()
	case "unixTimestamp":
		return jen.Qual(PkgSolanaGo, "UnixTimeSeconds")
	case "hash":
		return jen.Qual(PkgSolanaGo, "Hash")
	case "duration":
		return jen.Qual(PkgSolanaGo, "DurationSeconds")
	default:
		return jen.Interface()
	}
}

// ArrayType returns a jen.Code for an array type
func ArrayType(size int, elemType jen.Code) jen.Code {
	return jen.Index(jen.Lit(size)).Add(elemType)
}

// SliceType returns a jen.Code for a slice type
func SliceType(elemType jen.Code) jen.Code {
	return jen.Index().Add(elemType)
}

// PointerType returns a jen.Code for a pointer type
func PointerType(elemType jen.Code) jen.Code {
	return jen.Op("*").Add(elemType)
}

// MapType returns a jen.Code for a map type
func MapType(keyType, valueType jen.Code) jen.Code {
	return jen.Map(keyType).Add(valueType)
}

// ChannelType returns a jen.Code for a channel type
func ChannelType(elemType jen.Code) jen.Code {
	return jen.Chan().Add(elemType)
}

// FuncType returns a jen.Code for a function type
func FuncType(params, results []jen.Code) jen.Code {
	return jen.Func().Params(params...).Params(results...)
}

// InterfaceType returns a jen.Code for an interface type
func InterfaceType(methods ...jen.Code) jen.Code {
	return jen.Interface(methods...)
}

// StructType returns a jen.Code for a struct type
func StructType(fields ...jen.Code) jen.Code {
	return jen.Struct(fields...)
}
