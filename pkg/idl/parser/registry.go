package parser

import (
	"fmt"

	"github.com/iamrekas/anchor-go/pkg/idl/types"
)

// Registry holds all available parsers
type Registry struct {
	parsers []Parser
	verbose bool
}

// NewRegistry creates a new parser registry
func NewRegistry(verbose bool) *Registry {
	return &Registry{
		parsers: []Parser{},
		verbose: verbose,
	}
}

// RegisterParser registers a new parser
func (r *Registry) RegisterParser(parser Parser) {
	r.parsers = append(r.parsers, parser)
}

// Parse parses an IDL using the appropriate parser
func (r *Registry) Parse(data []byte) (types.IDL, error) {
	// Detect features
	features, err := DetectFeatures(data)
	if err != nil {
		return nil, fmt.Errorf("failed to detect features: %w", err)
	}

	// Set verbose flag
	features.Verbose = r.verbose

	// Try each parser in order
	for _, parser := range r.parsers {
		if parser.CanParse(data) {
			if r.verbose {
				fmt.Printf("Using parser: %T\n", parser)
			}
			return parser.Parse(data, features)
		}
	}

	// If no parser can handle the data, return an error
	return nil, fmt.Errorf("no parser found for the given IDL")
}
