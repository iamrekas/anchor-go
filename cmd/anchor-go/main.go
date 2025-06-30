// Package main provides the command-line interface for anchor-go.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/iamrekas/anchor-go/pkg/generator"
	"github.com/iamrekas/anchor-go/pkg/idl"
)

// StringArray is a flag.Value that collects multiple string values
type StringArray []string

func (s *StringArray) String() string {
	return strings.Join(*s, ", ")
}

func (s *StringArray) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func main() {
	// Define command-line flags
	var (
		filenames           StringArray
		dstDir              string
		packageName         string
		verbose             bool
		removeAccountSuffix bool
		encoding            string
		typeID              string
		modPath             string
		onlyParse           bool
	)

	// Set up flags
	flag.Var(&filenames, "src", "Path to source IDL file; can use multiple times.")
	flag.StringVar(&dstDir, "dst", "generated", "Destination folder")
	flag.StringVar(&packageName, "pkg", "", "Set package name to generate, default value is metadata.name of the source IDL.")
	flag.BoolVar(&verbose, "v", false, "Verbose mode (show detailed logs)")
	flag.BoolVar(&removeAccountSuffix, "remove-account-suffix", false, "Remove \"Account\" suffix from accessors")
	flag.StringVar(&encoding, "codec", "borsh", "Choose codec (borsh, bincode)")
	flag.StringVar(&typeID, "type-id", "anchor", "Choose typeID kind (anchor, uvarint32, uint32, uint8, notype)")
	flag.StringVar(&modPath, "mod", "", "Generate a go.mod file with the necessary dependencies, and this module")
	flag.BoolVar(&onlyParse, "only-parse", false, "Only parse IDL files without generating or writing code")
	flag.Parse()

	// Validate flags
	if len(filenames) == 0 {
		fmt.Println("Error: No IDL files provided")
		fmt.Println("Usage: anchor-go -src path/to/idl.json [-src path/to/another.json] [options]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Create configuration
	config := &generator.Config{
		DstDir:              dstDir,
		Package:             packageName,
		Verbose:             verbose,
		RemoveAccountSuffix: removeAccountSuffix,
		Encoding:            generator.EncodingType(encoding),
		TypeID:              generator.TypeIDKind(typeID),
		ModPath:             modPath,
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		fmt.Printf("Error: Invalid configuration: %v\n", err)
		os.Exit(1)
	}

	// Create generators
	accountsGen := generator.NewAccountsGenerator()
	instructionsGen := generator.NewInstructionsGenerator()
	typesGen := generator.NewTypesGenerator()
	eventsGen := generator.NewEventsGenerator()
	errorsGen := generator.NewErrorsGenerator()

	// Create composite generator
	gen := generator.NewCompositeGenerator(
		accountsGen,
		instructionsGen,
		typesGen,
		eventsGen,
		errorsGen,
	)
	gen.SetConfig(config)

	// Process each IDL file
	for _, idlFilepath := range filenames {
		if verbose {
			fmt.Printf("Generating client from IDL: %s\n", idlFilepath)
		} else {
			fmt.Printf("Processing: %s\n", idlFilepath)
		}

		// Set environment variable for verbose mode
		if verbose {
			os.Setenv("ANCHOR_GO_VERBOSE", "1")
		} else {
			os.Setenv("ANCHOR_GO_VERBOSE", "0")
		}

		// Parse IDL
		parsedIDL, err := idl.ParseFileWithVerbose(idlFilepath, verbose)
		if err != nil {
			fmt.Printf("Error parsing IDL file %s: %v\n", idlFilepath, err)
			continue
		}

		// If package name is not provided, use the IDL name
		if config.Package == "" {
			config.Package = parsedIDL.ProgramName()
		}

		// Generate code
		files, err := gen.Generate(parsedIDL)
		if err != nil {
			fmt.Printf("Error generating code for %s: %v\n", idlFilepath, err)
			continue
		}

		// Write generated files if not in parse-only mode
		if onlyParse {
			if verbose {
				fmt.Printf("Parsed IDL successfully: %s (skipping file generation due to --only-parse flag)\n", idlFilepath)
			}
		} else {
			for _, file := range files {
				// Create directory if it doesn't exist
				fullPath := filepath.Join(file.Path, file.Name+".go")
				dir := filepath.Dir(fullPath)

				if verbose {
					fmt.Printf("Creating directory: %s\n", dir)
				}

				if err := os.MkdirAll(dir, 0755); err != nil {
					fmt.Printf("Error creating directory %s: %v\n", dir, err)
					continue
				}

				// Write file
				if verbose {
					fmt.Printf("Writing file: %s\n", fullPath)
				}

				if err := os.WriteFile(fullPath, file.Content, 0644); err != nil {
					fmt.Printf("Error writing file %s: %v\n", fullPath, err)
					continue
				}

				if verbose {
					fmt.Printf("Generated: %s\n", fullPath)
				}
			}
		}
	}

	if onlyParse {
		fmt.Println("IDL parsing complete (no files were generated)")
	} else {
		fmt.Println("Code generation complete")
	}
}
