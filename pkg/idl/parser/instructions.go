package parser

import (
	"encoding/json"
	"fmt"

	"github.com/iamrekas/anchor-go/pkg/idl/types"
)

// ParseInstructions parses instructions from raw instruction data
func ParseInstructions(rawInstructions []rawInstruction, verbose bool) []types.Instruction {
	if verbose {
		fmt.Printf("Parsing %d instructions\n", len(rawInstructions))
	}

	instructions := make([]types.Instruction, 0, len(rawInstructions))

	if verbose {
		fmt.Println("Instructions:")
	}
	for _, rawInstr := range rawInstructions {
		// Create a basic instruction
		instr := types.Instruction{
			Name:          rawInstr.Name,
			Discriminator: rawInstr.Discriminator,
			Docs:          rawInstr.Docs,
			Accounts:      []types.AccountItem{},
			Args:          []types.Field{},
		}

		// Parse accounts
		if rawInstr.Accounts != nil {
			var accounts []map[string]interface{}
			if err := json.Unmarshal(rawInstr.Accounts, &accounts); err == nil {
				for _, acc := range accounts {
					name, _ := acc["name"].(string)
					signer, _ := acc["signer"].(bool)
					writable, _ := acc["writable"].(bool)
					isMut, _ := acc["isMut"].(bool)
					optional, _ := acc["optional"].(bool)
					address, _ := acc["address"].(string)
					docs, _ := acc["docs"].([]string)

					// Create a basic account
					account := &types.Account{
						Name:     name,
						Signer:   signer,
						Writable: writable || isMut, // Set writable to true if either writable or isMut is true
						Optional: optional,
						Address:  address,
						Docs:     docs,
					}

					// Check for PDA
					if pdaData, hasPDA := acc["pda"].(map[string]interface{}); hasPDA {
						if verbose {
							fmt.Printf("      Account %s is a PDA\n", name)
						}

						pda := &types.PDA{
							Seeds: []types.PDASeed{},
						}

						// Parse seeds
						if seedsData, hasSeeds := pdaData["seeds"].([]interface{}); hasSeeds {
							for _, seedData := range seedsData {
								if seedMap, ok := seedData.(map[string]interface{}); ok {
									kind, _ := seedMap["kind"].(string)

									if kind == "const" {
										// Parse constant seed
										if valueData, hasValue := seedMap["value"].([]interface{}); hasValue {
											valueBytes := make([]byte, len(valueData))
											for i, v := range valueData {
												if b, ok := v.(float64); ok {
													valueBytes[i] = byte(b)
												}
											}

											pda.Seeds = append(pda.Seeds, types.PDASeed{
												Kind:  "const",
												Value: valueBytes,
											})
										}
									} else if kind == "account" {
										// Parse account seed
										if path, hasPath := seedMap["path"].(string); hasPath {
											pda.Seeds = append(pda.Seeds, types.PDASeed{
												Kind: "account",
												Path: path,
											})
										}
									}
								}
							}
						}

						// Parse program seed
						if programData, hasProgram := pdaData["program"].(map[string]interface{}); hasProgram {
							kind, _ := programData["kind"].(string)

							if kind == "const" {
								// Parse constant program seed
								if valueData, hasValue := programData["value"].([]interface{}); hasValue {
									valueBytes := make([]byte, len(valueData))
									for i, v := range valueData {
										if b, ok := v.(float64); ok {
											valueBytes[i] = byte(b)
										}
									}

									pda.Program = &types.PDASeed{
										Kind:  "const",
										Value: valueBytes,
									}
								}
							} else if kind == "account" {
								// Parse account program seed
								if path, hasPath := programData["path"].(string); hasPath {
									pda.Program = &types.PDASeed{
										Kind: "account",
										Path: path,
									}
								}
							}
						}

						account.PDA = pda
					}

					instr.Accounts = append(instr.Accounts, account)
				}
			}
		}

		// Parse args
		for _, arg := range rawInstr.Args {
			// Create a basic field
			field := types.Field{
				Name:     arg.Name,
				Docs:     arg.Docs,
				Optional: false,
				Type:     &types.BasicType{TypeName: "string"}, // Default to string for now
			}

			// Parse type
			if arg.Type != nil {
				// Use the helper function to parse the type
				parsedType, success := ParseFieldType(arg.Name, arg.Type, verbose)
				if success {
					field.Type = parsedType

					// Check if it's an option type to set the optional flag
					if optType, ok := parsedType.(*types.OptionType); ok {
						field.Optional = true
						field.Type = optType
					}
				} else {
					// If parsing failed, use string as fallback
					if verbose {
						fmt.Printf("Warning: Using string type as fallback for arg '%s'\n", arg.Name)
					}
					field.Type = &types.BasicType{TypeName: "string"}
				}
			}

			instr.Args = append(instr.Args, field)
		}

		if verbose {
			fmt.Printf("  - %s: %d accounts, %d args\n", instr.Name, len(instr.Accounts), len(instr.Args))

			// Log accounts
			if len(instr.Accounts) > 0 {
				fmt.Printf("    Accounts:\n")
				for i, acc := range instr.Accounts {
					if acc.IsAccount() {
						account := acc.(*types.Account)
						fmt.Printf("      %d. %s (signer: %t, writable: %t)\n",
							i+1, account.GetName(), account.Signer, account.Writable)
					}
				}
			}

			// Log args
			if len(instr.Args) > 0 {
				fmt.Printf("    Args:\n")
				for i, arg := range instr.Args {
					fmt.Printf("      %d. %s: %s\n", i+1, arg.Name, arg.Type.String())
				}
			}

			fmt.Println()
		}

		instructions = append(instructions, instr)
	}

	if verbose {
		fmt.Printf("Parsed %d instructions\n", len(instructions))
	}

	return instructions
}
