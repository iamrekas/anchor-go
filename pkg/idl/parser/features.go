package parser

import (
	"encoding/json"
)

// Features represents the features detected in an IDL
type Features struct {
	HasVersion        bool
	HasConstants      bool
	HasDiscriminators bool
	HasPDA            bool
	HasEvents         bool
	HasErrors         bool
	HasMetadata       bool
	HasDirectAddress  bool
	HasNestedAccounts bool
	HasDocs           bool
	HasRelations      bool
	Verbose           bool // For verbose logging
}

// DetectFeatures detects features from a raw IDL map
func DetectFeatures(data []byte) (Features, error) {
	// First, try to unmarshal into a generic map to inspect the structure
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return Features{}, err
	}

	features := Features{}

	// Check for version
	if _, ok := raw["version"]; ok {
		features.HasVersion = true
	}

	// Check for constants
	if constants, ok := raw["constants"]; ok && constants != nil {
		features.HasConstants = true
	}

	// Check for events
	if events, ok := raw["events"]; ok && events != nil {
		features.HasEvents = true
	}

	// Check for errors
	if errors, ok := raw["errors"]; ok && errors != nil {
		features.HasErrors = true
	}

	// Check for metadata
	if metadata, ok := raw["metadata"]; ok && metadata != nil {
		features.HasMetadata = true
	}

	// Check for direct address
	if _, ok := raw["address"]; ok {
		features.HasDirectAddress = true
	}

	// Check for discriminators in instructions
	if instructions, ok := raw["instructions"].([]interface{}); ok && len(instructions) > 0 {
		if instr, ok := instructions[0].(map[string]interface{}); ok {
			if _, ok := instr["discriminator"]; ok {
				features.HasDiscriminators = true
			}
		}
	}

	// Check for PDA in accounts
	if accounts, ok := raw["accounts"].([]interface{}); ok && len(accounts) > 0 {
		if acc, ok := accounts[0].(map[string]interface{}); ok {
			if _, ok := acc["pda"]; ok {
				features.HasPDA = true
			}
		}
	}

	// Check for nested accounts
	if instructions, ok := raw["instructions"].([]interface{}); ok && len(instructions) > 0 {
		if instr, ok := instructions[0].(map[string]interface{}); ok {
			if accs, ok := instr["accounts"].([]interface{}); ok && len(accs) > 0 {
				for _, acc := range accs {
					if accMap, ok := acc.(map[string]interface{}); ok {
						if _, ok := accMap["accounts"]; ok {
							features.HasNestedAccounts = true
							break
						}
					}
				}
			}
		}
	}

	// Check for docs
	if instructions, ok := raw["instructions"].([]interface{}); ok && len(instructions) > 0 {
		if instr, ok := instructions[0].(map[string]interface{}); ok {
			if _, ok := instr["docs"]; ok {
				features.HasDocs = true
			}
		}
	}

	// Check for relations
	if instructions, ok := raw["instructions"].([]interface{}); ok && len(instructions) > 0 {
		if instr, ok := instructions[0].(map[string]interface{}); ok {
			if accs, ok := instr["accounts"].([]interface{}); ok && len(accs) > 0 {
				for _, acc := range accs {
					if accMap, ok := acc.(map[string]interface{}); ok {
						if _, ok := accMap["relations"]; ok {
							features.HasRelations = true
							break
						}
					}
				}
			}
		}
	}

	return features, nil
}
