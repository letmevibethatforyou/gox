package idx

import (
	"fmt"
	"strings"
)

// ID represents an AWS-style identifier with environment, type, and object ID components.
// The string format is: environment:type:object_id
// Example: "vibe:user:custom_value"
type ID struct {
	env        string
	objectType Type
	objectID   string
}

// Env returns the environment component of the ID.
func (id ID) Env() string {
	return id.env
}

// Type returns the object type component of the ID.
func (id ID) Type() Type {
	return id.objectType
}

// Value returns the object ID component of the ID.
func (id ID) Value() string {
	return id.objectID
}

// String returns the full string representation of the ID in the format: environment:type:object_id
func (id ID) String() string {
	return fmt.Sprintf("%s:%s:%s", id.env, id.objectType, id.objectID)
}

// ParseID parses a string representation of an ID and returns an ID struct.
// The input must be in the format: environment:type:object_id
// Returns an error if the format is invalid or any component fails validation.
func ParseID(s string) (ID, error) {
	parts := strings.Split(s, ":")
	if len(parts) != 3 {
		return ID{}, fmt.Errorf("invalid ID format: expected 3 parts separated by colons, got %d parts", len(parts))
	}

	environment := parts[0]
	if environment == "" {
		return ID{}, fmt.Errorf("invalid ID: env cannot be empty")
	}

	objectType, err := ParseType(parts[1])
	if err != nil {
		return ID{}, fmt.Errorf("invalid ID: %w", err)
	}

	objectID := parts[2]
	if objectID == "" {
		return ID{}, fmt.Errorf("invalid ID: object ID cannot be empty")
	}

	return ID{
		env:        environment,
		objectType: objectType,
		objectID:   objectID,
	}, nil
}

// Validate checks that all components of the ID are valid.
// Returns an error if any component is invalid or empty.
func (id ID) Validate() error {
	if id.env == "" {
		return fmt.Errorf("env cannot be empty")
	}

	if err := id.objectType.Validate(); err != nil {
		return fmt.Errorf("invalid object type: %w", err)
	}

	if id.objectID == "" {
		return fmt.Errorf("object ID cannot be empty")
	}

	return nil
}
