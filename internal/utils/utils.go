// Package utils provides common utility functions for the anchor-go project.
package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// CreateDirectoryIfNotExists creates a directory if it doesn't exist
func CreateDirectoryIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

// ToCamel converts a string to CamelCase
func ToCamel(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	words := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-' || r == ' '
	})

	for i, word := range words {
		if word == "" {
			continue
		}
		words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
	}

	return strings.Join(words, "")
}

// ToLowerCamel converts a string to lowerCamelCase
func ToLowerCamel(s string) string {
	s = ToCamel(s)
	if s == "" {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}

// ToSnakeCase converts a string to snake_case
func ToSnakeCase(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	var result strings.Builder
	for i, r := range s {
		if i > 0 && (r >= 'A' && r <= 'Z') {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}

	return strings.ToLower(result.String())
}

// FormatSighash formats a byte array as a string representation
func FormatSighash(buf []byte) string {
	elems := make([]string, 0)
	for _, v := range buf {
		elems = append(elems, fmt.Sprintf("%d", v))
	}

	return "[" + strings.Join(elems, ", ") + "]"
}

// FileExists checks if a file exists
func FileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// DirExists checks if a directory exists
func DirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// GetAbsolutePath returns the absolute path
func GetAbsolutePath(path string) (string, error) {
	return filepath.Abs(path)
}

// JoinPaths joins multiple path segments
func JoinPaths(paths ...string) string {
	return filepath.Join(paths...)
}
