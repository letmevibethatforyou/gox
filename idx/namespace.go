package idx

import (
	"fmt"
	"strings"

	"github.com/segmentio/ksuid"
)

// Namespace represents an environment context for creating IDs.
// It encapsulates the environment name and provides methods to create new IDs within that environment.
type Namespace struct {
	environment string
}

// NewNamespace creates a new Namespace with the given environment.
// Special handling: "prd" and empty string environments are normalized to "vibe".
func NewNamespace(environment string) Namespace {
	env := normalizeEnvironment(environment)
	return Namespace{environment: env}
}

// Environment returns the normalized environment name for this namespace.
func (n Namespace) Environment() string {
	return n.environment
}

// NewID creates a new ID within this namespace using the specified object type.
// The object ID component is automatically generated to ensure uniqueness.
// Returns an error if the object type is invalid.
func (n Namespace) NewID(objectType Type) (ID, error) {
	value := ksuid.New()
	return n.NewIDWithValue(objectType, value.String())
}

// NewIDWithValue creates a new ID within this namespace using the specified object type and custom value.
// This allows callers to provide their own object ID value instead of using auto-generation.
// Returns an error if the object type is invalid or the value is empty.
func (n Namespace) NewIDWithValue(objectType Type, value string) (ID, error) {
	if err := objectType.Validate(); err != nil {
		return ID{}, fmt.Errorf("invalid object type: %w", err)
	}

	if value == "" {
		return ID{}, fmt.Errorf("value cannot be empty")
	}

	return ID{
		env:        n.environment,
		objectType: objectType,
		objectID:   value,
	}, nil
}

// normalizeEnvironment applies special transformation rules to environment names.
// Both "prd" and empty string are converted to "vibe" for consistency.
// All other environment names are trimmed of whitespace but otherwise unchanged.
func normalizeEnvironment(env string) string {
	env = strings.TrimSpace(env)

	if env == "" || env == "prd" {
		return "vibe"
	}

	return env
}
