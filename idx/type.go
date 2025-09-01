// Copyright (c) 2025 letmevibethatforyou
// SPDX-License-Identifier: MIT

// Package idx provides AWS-style identifiers with flexible object ID generation.
//
// The package implements a namespace-based ID system where IDs follow the format:
// environment:type:object_id
//
// Example usage:
//
//	ns := idx.NewNamespace("prd") // becomes "vibe"
//	id, err := ns.NewID(idx.Type("user"))
//	fmt.Println(id.String()) // "vibe:user:auto_generated_value"
//
//	// Or with custom value:
//	customID, err := ns.NewIDWithValue(idx.Type("user"), "custom123")
//	fmt.Println(customID.String()) // "vibe:user:custom123"
package idx

import (
	"fmt"
	"regexp"
	"strings"
)

// Type represents an object type identifier used in IDs.
// Types must start with a letter and contain only letters, numbers, and underscores.
// Maximum length is 32 characters.
type Type string

// typeRegex validates that types start with a letter and contain only alphanumeric characters and underscores.
var typeRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`)

// String returns the string representation of the Type.
func (t Type) String() string {
	return string(t)
}

// Validate checks if the Type meets all requirements:
// - Not empty
// - Maximum 32 characters
// - No colons (to avoid conflicts with ID format)
// - Must start with letter and contain only letters, numbers, and underscores
func (t Type) Validate() error {
	str := string(t)

	if str == "" {
		return fmt.Errorf("type cannot be empty")
	}

	if len(str) > 32 {
		return fmt.Errorf("type cannot be longer than 32 characters")
	}

	if strings.Contains(str, ":") {
		return fmt.Errorf("type cannot contain colons")
	}

	if !typeRegex.MatchString(str) {
		return fmt.Errorf("type must start with a letter and contain only letters, numbers, and underscores")
	}

	return nil
}

// ParseType creates a Type from a string and validates it.
// Returns an error if the string doesn't meet Type requirements.
func ParseType(s string) (Type, error) {
	t := Type(s)
	if err := t.Validate(); err != nil {
		return "", err
	}
	return t, nil
}
