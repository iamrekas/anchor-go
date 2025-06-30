# Anchor-Go

A Go code generator for Solana's Anchor framework IDL files.

## Overview

Anchor-Go generates Go client code from Anchor IDL (Interface Definition Language) files. It creates type-safe Go structures and methods for interacting with Solana programs written in Anchor.

## Features

- **Feature-based IDL parsing**: Automatically detects and adapts to different IDL formats
- **Modular code generation**: Generates code for accounts, instructions, events, errors, etc.
- **Customizable output**: Configure package names, output directories, and more
- **Backward compatibility**: Supports older IDL formats while providing access to newer features

## Project Structure

```
anchor-go/
├── cmd/
│   └── anchor-go/             # Command-line interface
│       └── main.go
├── pkg/
│   ├── idl/                   # IDL parsing and structures
│   │   ├── types.go           # Core IDL types
│   │   ├── type_impl.go       # Type implementations
│   │   └── parser.go          # IDL parsing logic
│   ├── generator/             # Code generation logic
│   │   ├── generator.go       # Generator interface
│   │   ├── accounts.go        # Account code generation
│   │   └── ...                # Other generators
│   └── codegen/               # Code generation helpers
│       └── ...
├── old/                       # Storage for old implementation
│   └── ...
├── examples/                  # Examples of usage
│   └── ...
└── README.md
```

## Installation

```bash
go install github.com/iamrekas/anchor-go/cmd/anchor-go@latest
```

## Usage

```bash
# Basic usage
anchor-go -src path/to/idl.json -dst ./generated

# Multiple IDL files
anchor-go -src path/to/idl1.json -src path/to/idl2.json -dst ./generated

# Custom package name
anchor-go -src path/to/idl.json -dst ./generated -pkg mypackage

# Generate go.mod file
anchor-go -src path/to/idl.json -dst ./generated -mod github.com/myuser/myrepo
```

### Command-line Options

| Flag | Description | Default |
|------|-------------|---------|
| `-src` | Path to source IDL file (can be used multiple times) | (required) |
| `-dst` | Destination folder | `generated` |
| `-pkg` | Package name | (from IDL) |
| `-v` | Verbose mode (show detailed logs) | `false` |
| `-remove-account-suffix` | Remove "Account" suffix from accessors | `false` |
| `-codec` | Choose codec (borsh, bincode) | `borsh` |
| `-type-id` | Choose typeID kind (anchor, uvarint32, uint32, uint8, notype) | `anchor` |
| `-mod` | Generate a go.mod file with dependencies | (none) |

## Feature Detection

Anchor-Go uses feature detection to handle different IDL formats. Instead of relying on strict versioning, it examines the IDL content to determine which features are available. This approach provides better compatibility with custom or modified IDL formats.

Detected features include:
- Discriminators
- Program Derived Addresses (PDAs)
- Nested accounts
- Documentation
- Constants
- Events
- Errors
- Relations

## Development

### Building from Source

```bash
git clone https://github.com/iamrekas/anchor-go.git
cd anchor-go
go build -o anchor-go ./cmd/anchor-go
```

### Adding New Features

To add support for new IDL features:

1. Update the feature detection in `pkg/idl/parser.go`
2. Add the necessary types in `pkg/idl/types.go`
3. Implement the parsing logic in `pkg/idl/parser.go`
4. Add code generation in the appropriate generator

## License

[MIT License](LICENSE)
