// Package types provides type definitions for Anchor Interface Definition Language (IDL) files.
package types

// IDL represents a parsed Interface Definition Language (IDL) for an Anchor program
type IDL interface {
	// Version returns the IDL specification version
	Version() string

	// ProgramName returns the name of the program
	ProgramName() string

	// ProgramID returns the program's public key as a string
	ProgramID() string

	// Instructions returns all instructions defined in the IDL
	Instructions() []Instruction

	// Accounts returns all account structures defined in the IDL
	Accounts() []AccountDef

	// Types returns all custom types defined in the IDL
	Types() []TypeDef

	// Events returns all events defined in the IDL
	Events() []Event

	// Errors returns all custom error codes defined in the IDL
	Errors() []ErrorCode

	// Constants returns all constants defined in the IDL (may not be supported in all versions)
	Constants() []Constant

	// Validate verifies the IDL is properly formed
	Validate() error
}

// Instruction represents a program instruction
type Instruction struct {
	Name          string
	Discriminator []byte
	Docs          []string
	Accounts      []AccountItem
	Args          []Field
}

// AccountDef represents an account structure definition
type AccountDef struct {
	Name          string
	Type          TypeDef
	Discriminator []byte
	Docs          []string
}

// TypeDef represents a custom type definition
type TypeDef struct {
	Name          string
	Kind          TypeDefKind
	Fields        []Field
	Variants      []EnumVariant
	Discriminator []byte
	Docs          []string
}

// TypeDefKind represents the kind of a type definition
type TypeDefKind string

const (
	TypeDefKindStruct TypeDefKind = "struct"
	TypeDefKindEnum   TypeDefKind = "enum"
	TypeDefKindAlias  TypeDefKind = "alias"
)

// Field represents a named field with a type
type Field struct {
	Name     string
	Type     Type
	Docs     []string
	Optional bool
}

// Type represents a data type
type Type interface {
	// IsArray returns true if the type is an array
	IsArray() bool

	// IsOption returns true if the type is an option
	IsOption() bool

	// IsVector returns true if the type is a vector
	IsVector() bool

	// IsDefined returns true if the type is a defined type
	IsDefined() bool

	// IsBasic returns true if the type is a basic type
	IsBasic() bool

	// String returns a string representation of the type
	String() string
}

// AccountItem represents an account or a group of accounts
type AccountItem interface {
	// IsAccount returns true if this is a single account
	IsAccount() bool

	// IsAccountGroup returns true if this is a group of accounts
	IsAccountGroup() bool

	// GetName returns the name of the account or account group
	GetName() string

	// GetDocs returns the documentation for the account or account group
	GetDocs() []string
}

// Account represents a single account
type Account struct {
	Name      string
	Signer    bool
	Writable  bool
	Optional  bool
	Address   string
	PDA       *PDA
	Relations []string
	Docs      []string
}

// IsAccount implements AccountItem.IsAccount
func (a *Account) IsAccount() bool {
	return true
}

// IsAccountGroup implements AccountItem.IsAccountGroup
func (a *Account) IsAccountGroup() bool {
	return false
}

// GetName implements AccountItem.GetName
func (a *Account) GetName() string {
	return a.Name
}

// GetDocs implements AccountItem.GetDocs
func (a *Account) GetDocs() []string {
	return a.Docs
}

// AccountGroup represents a group of accounts
type AccountGroup struct {
	Name     string
	Accounts []AccountItem
	Docs     []string
}

// IsAccount implements AccountItem.IsAccount
func (g *AccountGroup) IsAccount() bool {
	return false
}

// IsAccountGroup implements AccountItem.IsAccountGroup
func (g *AccountGroup) IsAccountGroup() bool {
	return true
}

// GetName implements AccountItem.GetName
func (g *AccountGroup) GetName() string {
	return g.Name
}

// GetDocs implements AccountItem.GetDocs
func (g *AccountGroup) GetDocs() []string {
	return g.Docs
}

// PDA represents a Program Derived Address
type PDA struct {
	Seeds   []PDASeed
	Program *PDASeed
}

// PDASeed represents a seed for a Program Derived Address
type PDASeed struct {
	Kind  string // "const" or "account"
	Value []byte // For "const" kind
	Path  string // For "account" kind
}

// EnumVariant represents a variant in an enum
type EnumVariant struct {
	Name   string
	Fields []Field
	Docs   []string
}

// Event represents an event that can be emitted by a program
type Event struct {
	Name          string
	Discriminator []byte
	Fields        []Field
	Docs          []string
}

// ErrorCode represents a custom error code
type ErrorCode struct {
	Code int
	Name string
	Msg  string
}

// Constant represents a constant value defined in the IDL
type Constant struct {
	Name  string
	Type  Type
	Value string
}
