// Package idl provides functions for working with Anchor Interface Definition Language (IDL) files.
package idl

import (
	"github.com/iamrekas/anchor-go/pkg/idl/parser"
	"github.com/iamrekas/anchor-go/pkg/idl/types"
)

// IDL is an alias for types.IDL
type IDL = types.IDL

// Instruction is an alias for types.Instruction
type Instruction = types.Instruction

// AccountDef is an alias for types.AccountDef
type AccountDef = types.AccountDef

// TypeDef is an alias for types.TypeDef
type TypeDef = types.TypeDef

// TypeDefKind is an alias for types.TypeDefKind
type TypeDefKind = types.TypeDefKind

// Field is an alias for types.Field
type Field = types.Field

// Type is an alias for types.Type
type Type = types.Type

// AccountItem is an alias for types.AccountItem
type AccountItem = types.AccountItem

// Account is an alias for types.Account
type Account = types.Account

// AccountGroup is an alias for types.AccountGroup
type AccountGroup = types.AccountGroup

// PDA is an alias for types.PDA
type PDA = types.PDA

// PDASeed is an alias for types.PDASeed
type PDASeed = types.PDASeed

// EnumVariant is an alias for types.EnumVariant
type EnumVariant = types.EnumVariant

// Event is an alias for types.Event
type Event = types.Event

// ErrorCode is an alias for types.ErrorCode
type ErrorCode = types.ErrorCode

// Constant is an alias for types.Constant
type Constant = types.Constant

// Constants for TypeDefKind
const (
	TypeDefKindStruct = types.TypeDefKindStruct
	TypeDefKindEnum   = types.TypeDefKindEnum
	TypeDefKindAlias  = types.TypeDefKindAlias
)

// ParseFile parses an IDL file
func ParseFile(path string) (IDL, error) {
	return parser.ParseFile(path)
}

// ParseFileWithVerbose parses an IDL file with verbose logging
func ParseFileWithVerbose(path string, verbose bool) (IDL, error) {
	return parser.ParseFileWithVerbose(path, verbose)
}
