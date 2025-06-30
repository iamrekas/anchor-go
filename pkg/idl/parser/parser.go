package parser

import (
	"fmt"
	"io"
	"os"

	"github.com/iamrekas/anchor-go/pkg/idl/types"
)

// Parse reads an IDL file and returns an IDL implementation based on feature detection
func Parse(reader io.Reader, verbose bool) (types.IDL, error) {
	// Read the entire content
	rawBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read IDL: %w", err)
	}

	// Create a parser registry
	registry := NewRegistry(verbose)

	// Register parsers
	registry.RegisterParser(&UniversalParser{})

	// Parse the IDL
	idl, err := registry.Parse(rawBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse IDL: %w", err)
	}

	// Validate the IDL
	if err := idl.Validate(); err != nil {
		return nil, fmt.Errorf("IDL validation failed: %w", err)
	}

	return idl, nil
}

// ParseFile reads an IDL file from a path and returns an IDL implementation
func ParseFile(path string) (types.IDL, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open IDL file: %w", err)
	}
	defer file.Close()

	// Get verbose flag from environment
	verbose := os.Getenv("ANCHOR_GO_VERBOSE") == "1"

	return Parse(file, verbose)
}

// ParseFileWithVerbose reads an IDL file from a path and returns an IDL implementation
// with the specified verbose flag
func ParseFileWithVerbose(path string, verbose bool) (types.IDL, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open IDL file: %w", err)
	}
	defer file.Close()

	return Parse(file, verbose)
}
