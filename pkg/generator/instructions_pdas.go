package generator

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/iamrekas/anchor-go/pkg/idl"
)

// generatePDACalculations generates code to calculate and set PDAs in the New<Instruction>Instruction function
func (g *InstructionsGenerator) generatePDACalculations(content *bytes.Buffer, instr idl.Instruction, requiredAccounts []idl.AccountItem, additionalParams map[string]bool) error {
	// First, identify PDAs that depend on other PDAs
	pdaDependencies := make(map[string][]string) // map[pdaName][]dependencyNames
	pdaSeeds := make(map[string][]string)        // map[pdaName][]seedNames
	pdaAccounts := make(map[string]*idl.Account) // map[pdaName]*idl.Account

	// First pass: identify dependencies and collect seeds
	for _, acc := range instr.Accounts {
		if !acc.IsAccount() {
			continue
		}

		account, ok := acc.(*idl.Account)
		if !ok {
			continue
		}

		// Skip non-PDA accounts
		if account.PDA == nil || len(account.PDA.Seeds) == 0 {
			continue
		}

		pdaAccounts[account.Name] = account

		// Collect seeds for this PDA
		seeds := make([]string, 0)
		for _, seed := range account.PDA.Seeds {
			if seed.Kind == "account" {
				seedPath := seed.Path
				if strings.Contains(seedPath, ".") {
					// This is a dependency on another PDA
					parts := strings.SplitN(seedPath, ".", 2)
					pdaDependencies[account.Name] = append(pdaDependencies[account.Name], parts[0])
					paramName := strings.ReplaceAll(seedPath, ".", "_")
					seeds = append(seeds, paramName)
				} else {
					seeds = append(seeds, seedPath)
				}
			}
		}
		pdaSeeds[account.Name] = seeds
	}

	// Track calculated PDAs to use as seeds for other PDAs
	calculatedPDAs := make(map[string]string)

	// 1. First, process PDAs with constant seeds (no input required)
	content.WriteString("\t// Step 1: Calculate PDAs with constant seeds (no input required)\n")
	for pdaName, account := range pdaAccounts {
		seeds, exists := pdaSeeds[pdaName]
		if !exists || len(seeds) == 0 {
			// This is a constant seed PDA
			accName := toCamelCase(account.Name)
			content.WriteString(fmt.Sprintf("\t// Calculate and set %s PDA\n", account.Name))
			content.WriteString(fmt.Sprintf("\t%sPDA, _, err := nd.Find%sAddress()\n", account.Name, accName))
			content.WriteString("\tif err != nil {\n")
			content.WriteString(fmt.Sprintf("\t\treturn nil, fmt.Errorf(\"failed to find %s address: %%w\", err)\n", account.Name))
			content.WriteString("\t}\n")
			content.WriteString(fmt.Sprintf("\tnd.Set%sAccount(%sPDA)\n\n", accName, account.Name))

			// Add this PDA to the calculated PDAs map
			calculatedPDAs[account.Name] = fmt.Sprintf("%sPDA", account.Name)
		}
	}

	// 2. Next, process PDAs that use only arguments as seeds
	content.WriteString("\t// Step 2: Calculate PDAs that use arguments as seeds\n")
	processedInStep2 := make(map[string]bool)
	for pdaName, account := range pdaAccounts {
		// Skip PDAs already processed in step 1
		if calculatedPDAs[account.Name] != "" {
			continue
		}

		seeds, exists := pdaSeeds[pdaName]
		if !exists {
			continue
		}

		// Check if all seeds are available as arguments or required accounts
		allSeedsFromArgs := true
		seedParams := make([]string, 0)
		for _, seedName := range seeds {
			seedFound := false

			// Check if it's an additional parameter
			if strings.Contains(seedName, "_") && !strings.HasPrefix(seedName, "program_") && additionalParams[seedName] {
				seedParams = append(seedParams, seedName)
				seedFound = true
			} else {
				// Check if it's a required account
				for _, reqAcc := range requiredAccounts {
					reqAccount, ok := reqAcc.(*idl.Account)
					if ok && reqAccount.Name == seedName {
						seedParams = append(seedParams, seedName)
						seedFound = true
						break
					}
				}
			}

			if !seedFound {
				allSeedsFromArgs = false
				break
			}
		}

		if allSeedsFromArgs {
			// All seeds are available as arguments
			accName := toCamelCase(account.Name)
			content.WriteString(fmt.Sprintf("\t// Calculate and set %s PDA\n", account.Name))
			content.WriteString(fmt.Sprintf("\t%sPDA, _, err := nd.Find%sAddress(%s)\n", account.Name, accName, strings.Join(seedParams, ", ")))
			content.WriteString("\tif err != nil {\n")
			content.WriteString(fmt.Sprintf("\t\treturn nil, fmt.Errorf(\"failed to find %s address: %%w\", err)\n", account.Name))
			content.WriteString("\t}\n")
			content.WriteString(fmt.Sprintf("\tnd.Set%sAccount(%sPDA)\n\n", accName, account.Name))

			// Add this PDA to the calculated PDAs map
			calculatedPDAs[account.Name] = fmt.Sprintf("%sPDA", account.Name)
			processedInStep2[pdaName] = true
		}
	}

	// 3. Finally, process PDAs that depend on previously generated PDAs
	content.WriteString("\t// Step 3: Calculate PDAs that depend on previously generated PDAs\n")

	// Keep trying to process remaining PDAs until no more can be processed
	remainingPDAs := true
	for remainingPDAs {
		remainingPDAs = false

		for pdaName, account := range pdaAccounts {
			// Skip PDAs already processed
			if calculatedPDAs[account.Name] != "" {
				continue
			}

			seeds, exists := pdaSeeds[pdaName]
			if !exists {
				continue
			}

			// Collect available seeds and check if we can calculate this PDA
			seedParams := make([]string, 0)
			missingSeeds := make([]string, 0)

			for _, seedName := range seeds {
				seedFound := false

				// Check if this seed is a calculated PDA
				if pdaValue, ok := calculatedPDAs[seedName]; ok {
					seedParams = append(seedParams, pdaValue)
					seedFound = true
				} else if strings.Contains(seedName, "_") && !strings.HasPrefix(seedName, "program_") && additionalParams[seedName] {
					// Check if it's an additional parameter
					seedParams = append(seedParams, seedName)
					seedFound = true
				} else {
					// Check if it's a required account
					for _, reqAcc := range requiredAccounts {
						reqAccount, ok := reqAcc.(*idl.Account)
						if ok && reqAccount.Name == seedName {
							seedParams = append(seedParams, seedName)
							seedFound = true
							break
						}
					}
				}

				if !seedFound {
					missingSeeds = append(missingSeeds, seedName)
				}
			}

			// Check if all seeds are available
			if len(missingSeeds) == 0 {
				// All seeds are available, so we can calculate the PDA
				accName := toCamelCase(account.Name)
				content.WriteString(fmt.Sprintf("\t// Calculate and set %s PDA\n", account.Name))
				content.WriteString(fmt.Sprintf("\t%sPDA, _, err := nd.Find%sAddress(%s)\n", account.Name, accName, strings.Join(seedParams, ", ")))
				content.WriteString("\tif err != nil {\n")
				content.WriteString(fmt.Sprintf("\t\treturn nil, fmt.Errorf(\"failed to find %s address: %%w\", err)\n", account.Name))
				content.WriteString("\t}\n")
				content.WriteString(fmt.Sprintf("\tnd.Set%sAccount(%sPDA)\n\n", accName, account.Name))

				// Add this PDA to the calculated PDAs map
				calculatedPDAs[account.Name] = fmt.Sprintf("%sPDA", account.Name)
				remainingPDAs = true // We processed a PDA, so we need to try again
			}
		}
	}

	// Add comments for PDAs that couldn't be calculated
	for pdaName, account := range pdaAccounts {
		// Skip PDAs already processed
		if calculatedPDAs[account.Name] != "" {
			continue
		}

		seeds, exists := pdaSeeds[pdaName]
		if !exists {
			continue
		}

		// Collect available seeds and missing seeds
		seedParams := make([]string, 0)
		missingSeeds := make([]string, 0)

		for _, seedName := range seeds {
			seedFound := false

			// Check if this seed is a calculated PDA
			if pdaValue, ok := calculatedPDAs[seedName]; ok {
				seedParams = append(seedParams, pdaValue)
				seedFound = true
			} else if strings.Contains(seedName, "_") && !strings.HasPrefix(seedName, "program_") && additionalParams[seedName] {
				// Check if it's an additional parameter
				seedParams = append(seedParams, seedName)
				seedFound = true
			} else {
				// Check if it's a required account
				for _, reqAcc := range requiredAccounts {
					reqAccount, ok := reqAcc.(*idl.Account)
					if ok && reqAccount.Name == seedName {
						seedParams = append(seedParams, seedName)
						seedFound = true
						break
					}
				}
			}

			if !seedFound {
				missingSeeds = append(missingSeeds, seedName)
			}
		}

		// Log missing seeds for debugging
		if g.config.Verbose {
			fmt.Printf("Warning: Cannot calculate PDA for %s, missing seeds: %v\n", account.Name, missingSeeds)
		}

		accName := toCamelCase(account.Name)

		// For PDAs with missing seeds, add a comment explaining that these seeds need to be provided
		content.WriteString(fmt.Sprintf("\t// Note: To calculate %s PDA, the following seeds need to be provided: %v\n", account.Name, missingSeeds))

		// Suggest how to handle these missing seeds
		content.WriteString("\t// These seeds can be:\n")
		content.WriteString("\t// 1. Passed as additional arguments to this instruction function\n")
		content.WriteString("\t// 2. Derived from other PDAs that are calculated earlier\n")
		content.WriteString("\t// 3. Fetched from the blockchain before calling this function\n")

		// Add example code for manually calculating and setting the PDA
		content.WriteString(fmt.Sprintf("\t// Example of manually calculating and setting this PDA:\n"))
		content.WriteString(fmt.Sprintf("\t// %sPDA, _, err := nd.Find%sAddress(", account.Name, accName))

		// Add placeholders for missing seeds
		for i, seed := range missingSeeds {
			if i > 0 {
				content.WriteString(", ")
			}
			content.WriteString(seed)
		}

		// Add existing seeds if any
		for _, seed := range seedParams {
			if len(missingSeeds) > 0 || seedParams[0] != seed {
				content.WriteString(", ")
			}
			content.WriteString(seed)
		}

		content.WriteString(")\n")
		content.WriteString("\t// if err != nil {\n")
		content.WriteString(fmt.Sprintf("\t// \treturn nil, fmt.Errorf(\"failed to find %s address: %%w\", err)\n", account.Name))
		content.WriteString("\t// }\n")
		content.WriteString(fmt.Sprintf("\t// nd.Set%sAccount(%sPDA)\n\n", accName, account.Name))
	}

	return nil
}
