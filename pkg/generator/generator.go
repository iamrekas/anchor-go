// Package generator provides interfaces and implementations for generating Go code from Anchor IDL.
package generator

import (
	"github.com/iamrekas/anchor-go/pkg/idl"
)

// Generator is the interface for code generators
type Generator interface {
	// Generate generates code from an IDL
	Generate(idl idl.IDL) ([]*GeneratedFile, error)

	// SetConfig sets the configuration for the generator
	SetConfig(config *Config)

	// GetConfig returns the current configuration
	GetConfig() *Config
}

// GeneratedFile represents a generated file
type GeneratedFile struct {
	// Name is the name of the file without extension
	Name string

	// Content is the content of the file
	Content []byte

	// Path is the relative path where the file should be written
	Path string
}

// Config holds configuration options for code generation
type Config struct {
	// DstDir is the destination directory for generated files
	DstDir string

	// Package is the package name for generated files
	Package string

	// Verbose enables verbose logging mode
	Verbose bool

	// RemoveAccountSuffix removes "Account" suffix from accessors
	RemoveAccountSuffix bool

	// Encoding is the encoding to use (borsh, bincode, etc.)
	Encoding EncodingType

	// TypeID is the type ID kind to use
	TypeID TypeIDKind

	// ModPath is the module path for generated code
	ModPath string
}

// EncodingType represents the encoding type
type EncodingType string

const (
	// EncodingBorsh is the Borsh encoding
	EncodingBorsh EncodingType = "borsh"

	// EncodingBincode is the Bincode encoding
	EncodingBincode EncodingType = "bincode"
)

// TypeIDKind represents the type ID kind
type TypeIDKind string

const (
	// TypeIDAnchor is the Anchor type ID kind
	TypeIDAnchor TypeIDKind = "anchor"

	// TypeIDUvarint32 is the Uvarint32 type ID kind
	TypeIDUvarint32 TypeIDKind = "uvarint32"

	// TypeIDUint32 is the Uint32 type ID kind
	TypeIDUint32 TypeIDKind = "uint32"

	// TypeIDUint8 is the Uint8 type ID kind
	TypeIDUint8 TypeIDKind = "uint8"

	// TypeIDNoType is the NoType type ID kind
	TypeIDNoType TypeIDKind = "notype"
)

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		DstDir:              "generated",
		Verbose:             false,
		RemoveAccountSuffix: false,
		Encoding:            EncodingBorsh,
		TypeID:              TypeIDAnchor,
	}
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Validation logic would be implemented here
	return nil
}

// CompositeGenerator is a generator that combines multiple generators
type CompositeGenerator struct {
	generators []Generator
	config     *Config
}

// NewCompositeGenerator creates a new composite generator
func NewCompositeGenerator(generators ...Generator) *CompositeGenerator {
	return &CompositeGenerator{
		generators: generators,
		config:     DefaultConfig(),
	}
}

// Generate implements Generator.Generate
func (g *CompositeGenerator) Generate(idl idl.IDL) ([]*GeneratedFile, error) {
	var files []*GeneratedFile

	// Run each generator and collect the results
	for _, generator := range g.generators {
		generator.SetConfig(g.config)

		generatedFiles, err := generator.Generate(idl)
		if err != nil {
			return nil, err
		}

		files = append(files, generatedFiles...)
	}

	return files, nil
}

// SetConfig implements Generator.SetConfig
func (g *CompositeGenerator) SetConfig(config *Config) {
	g.config = config
}

// GetConfig implements Generator.GetConfig
func (g *CompositeGenerator) GetConfig() *Config {
	return g.config
}
