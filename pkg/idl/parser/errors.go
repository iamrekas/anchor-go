package parser

import (
	"encoding/json"
	"fmt"

	"github.com/iamrekas/anchor-go/pkg/idl/types"
)

// ParseErrors parses error codes from raw error data
func ParseErrors(rawErrors []json.RawMessage, verbose bool) []types.ErrorCode {
	if verbose {
		fmt.Printf("Parsing %d errors\n", len(rawErrors))
	}

	if verbose {
		fmt.Println("Errors:")
	}

	errors := make([]types.ErrorCode, 0, len(rawErrors))
	for _, rawError := range rawErrors {
		var errorDef struct {
			Code int    `json:"code"`
			Name string `json:"name"`
			Msg  string `json:"msg"`
		}

		if err := json.Unmarshal(rawError, &errorDef); err != nil {
			if verbose {
				fmt.Printf("Warning: Failed to parse error: %v\n", err)
			}
			continue
		}

		// Create a basic error code
		errCode := types.ErrorCode{
			Code: errorDef.Code,
			Name: errorDef.Name,
			Msg:  errorDef.Msg,
		}

		if verbose {
			fmt.Printf("  - %s (code %d): %s\n", errCode.Name, errCode.Code, errCode.Msg)
		}

		errors = append(errors, errCode)
	}

	if verbose {
		fmt.Printf("Parsed %d errors\n", len(errors))
	}

	return errors
}
