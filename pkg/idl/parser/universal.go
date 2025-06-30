package parser

import (
	"encoding/json"
	"fmt"

	"github.com/iamrekas/anchor-go/pkg/idl/types"
)

// UniversalParser is a parser that can handle all IDL formats
type UniversalParser struct{}

// CanParse implements Parser.CanParse
func (p *UniversalParser) CanParse(data []byte) bool {
	// This parser can handle any IDL format
	return true
}

// Parse implements Parser.Parse
func (p *UniversalParser) Parse(data []byte, features Features) (types.IDL, error) {
	var idlImpl universalIDL
	if err := json.Unmarshal(data, &idlImpl); err != nil {
		return nil, fmt.Errorf("failed to unmarshal IDL: %w", err)
	}

	idlImpl.features = features

	// Parse instructions
	idlImpl.parsedInstructions = ParseInstructions(idlImpl.RawInstructions, features.Verbose)

	// Parse accounts
	idlImpl.parsedAccounts = ParseAccounts(idlImpl.RawAccounts, idlImpl.RawTypes, features.Verbose)

	// Parse types
	idlImpl.parsedTypes = ParseTypes(idlImpl.RawTypes, features.Verbose)

	// Parse events
	idlImpl.parsedEvents = ParseEvents(idlImpl.RawEvents, idlImpl.RawTypes, features.Verbose)

	// Parse errors
	idlImpl.parsedErrors = ParseErrors(idlImpl.RawErrors, features.Verbose)

	// Parse constants
	idlImpl.parsedConstants = ParseConstants(idlImpl.RawConstants, features.Verbose)

	return &idlImpl, nil
}

// universalIDL is a flexible IDL implementation that can handle all features
type universalIDL struct {
	RawVersion      string            `json:"version"`
	Name            string            `json:"name"`
	RawInstructions []rawInstruction  `json:"instructions"`
	RawAccounts     []json.RawMessage `json:"accounts,omitempty"`
	RawTypes        []json.RawMessage `json:"types,omitempty"`
	RawEvents       []json.RawMessage `json:"events,omitempty"`
	RawErrors       []json.RawMessage `json:"errors,omitempty"`
	RawConstants    []json.RawMessage `json:"constants,omitempty"`
	Address         string            `json:"address,omitempty"`
	Metadata        *rawMetadata      `json:"metadata,omitempty"`

	// Parsed data
	parsedInstructions []types.Instruction
	parsedAccounts     []types.AccountDef
	parsedTypes        []types.TypeDef
	parsedEvents       []types.Event
	parsedErrors       []types.ErrorCode
	parsedConstants    []types.Constant

	// Detected features
	features Features
}

// rawInstruction holds the raw instruction data
type rawInstruction struct {
	Name          string          `json:"name"`
	Discriminator []byte          `json:"discriminator,omitempty"`
	Docs          []string        `json:"docs,omitempty"`
	Accounts      json.RawMessage `json:"accounts"`
	Args          []rawField      `json:"args"`
}

// rawField holds the raw field data
type rawField struct {
	Name string          `json:"name"`
	Docs []string        `json:"docs,omitempty"`
	Type json.RawMessage `json:"type"`
}

// rawMetadata holds the raw metadata
type rawMetadata struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

// Version implements types.IDL.Version
func (i *universalIDL) Version() string {
	return i.RawVersion
}

// ProgramName implements types.IDL.ProgramName
func (i *universalIDL) ProgramName() string {
	return i.Name
}

// ProgramID implements types.IDL.ProgramID
func (i *universalIDL) ProgramID() string {
	if i.Address != "" {
		return i.Address
	}
	if i.Metadata != nil && i.Metadata.Address != "" {
		return i.Metadata.Address
	}
	return ""
}

// Instructions implements types.IDL.Instructions
func (i *universalIDL) Instructions() []types.Instruction {
	return i.parsedInstructions
}

// Accounts implements types.IDL.Accounts
func (i *universalIDL) Accounts() []types.AccountDef {
	return i.parsedAccounts
}

// Types implements types.IDL.Types
func (i *universalIDL) Types() []types.TypeDef {
	return i.parsedTypes
}

// Events implements types.IDL.Events
func (i *universalIDL) Events() []types.Event {
	return i.parsedEvents
}

// Errors implements types.IDL.Errors
func (i *universalIDL) Errors() []types.ErrorCode {
	return i.parsedErrors
}

// Constants implements types.IDL.Constants
func (i *universalIDL) Constants() []types.Constant {
	return i.parsedConstants
}

// Validate implements types.IDL.Validate
func (i *universalIDL) Validate() error {
	// For now, we'll be lenient with validation to support various IDL formats

	// If name is empty but metadata name is available, use that
	if i.Name == "" && i.Metadata != nil && i.Metadata.Name != "" {
		i.Name = i.Metadata.Name
	}

	// If name is still empty, use a default name
	if i.Name == "" {
		i.Name = "AnchorProgram"
	}

	// We'll be lenient about instructions for now
	if len(i.RawInstructions) == 0 {
		// Just log this instead of returning an error
		fmt.Println("Warning: IDL has no instructions")
	}

	return nil
}
