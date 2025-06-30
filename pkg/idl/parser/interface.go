package parser

import (
	"github.com/iamrekas/anchor-go/pkg/idl/types"
)

// Parser is the interface for IDL parsers
type Parser interface {
	// Parse parses an IDL from raw data
	Parse(data []byte, features Features) (types.IDL, error)

	// CanParse checks if this parser can parse the given data
	CanParse(data []byte) bool
}
