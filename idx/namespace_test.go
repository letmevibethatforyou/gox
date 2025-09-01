package idx

import (
	"strings"
	"testing"
)

func TestNewNamespace(t *testing.T) {
	tests := map[string]struct {
		environment string
		expected    string
	}{
		"prd becomes vibe": {
			environment: "prd",
			expected:    "vibe",
		},
		"empty string becomes vibe": {
			environment: "",
			expected:    "vibe",
		},
		"whitespace only becomes vibe": {
			environment: "   ",
			expected:    "vibe",
		},
		"development stays development": {
			environment: "development",
			expected:    "development",
		},
		"staging stays staging": {
			environment: "staging",
			expected:    "staging",
		},
		"custom env stays custom": {
			environment: "my-custom-env",
			expected:    "my-custom-env",
		},
		"trimmed whitespace": {
			environment: "  dev  ",
			expected:    "dev",
		},
		"prd with whitespace becomes vibe": {
			environment: "  prd  ",
			expected:    "vibe",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ns := NewNamespace(tt.environment)
			if ns.Environment() != tt.expected {
				t.Errorf("NewNamespace(%q).Environment() = %q, want %q", tt.environment, ns.Environment(), tt.expected)
			}
		})
	}
}

func TestNamespace_Environment(t *testing.T) {
	ns := NewNamespace("test-env")
	result := ns.Environment()
	expected := "test-env"

	if result != expected {
		t.Errorf("Environment() = %q, want %q", result, expected)
	}
}

func TestNamespace_NewID(t *testing.T) {
	tests := map[string]struct {
		environment string
		objectType  Type
		wantErr     bool
		errMsg      string
	}{
		"valid type": {
			environment: "dev",
			objectType:  Type("user"),
			wantErr:     false,
		},
		"valid type with vibe env": {
			environment: "prd",
			objectType:  Type("session"),
			wantErr:     false,
		},
		"invalid type - empty": {
			environment: "dev",
			objectType:  Type(""),
			wantErr:     true,
			errMsg:      "invalid object type: type cannot be empty",
		},
		"invalid type - starts with number": {
			environment: "dev",
			objectType:  Type("1user"),
			wantErr:     true,
			errMsg:      "invalid object type: type must start with a letter",
		},
		"invalid type - contains colon": {
			environment: "dev",
			objectType:  Type("user:item"),
			wantErr:     true,
			errMsg:      "invalid object type: type cannot contain colons",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ns := NewNamespace(tt.environment)
			id, err := ns.NewID(tt.objectType)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewID() expected error but got nil")
					return
				}
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("NewID() error = %v, want error containing %q", err, tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("NewID() unexpected error = %v", err)
				return
			}

			// Verify the ID has correct components
			expectedEnv := tt.environment
			if tt.environment == "prd" || strings.TrimSpace(tt.environment) == "" {
				expectedEnv = "vibe"
			}

			if id.Env() != expectedEnv {
				t.Errorf("NewID().Env() = %q, want %q", id.Env(), expectedEnv)
			}

			if id.Type() != tt.objectType {
				t.Errorf("NewID().Type() = %q, want %q", id.Type(), tt.objectType)
			}

			if id.Value() == "" {
				t.Errorf("NewID().Value() should not be empty")
			}

			// Verify string format
			parts := strings.Split(id.String(), ":")
			if len(parts) != 3 {
				t.Errorf("NewID().String() should have 3 parts separated by colons, got %d parts", len(parts))
			}
		})
	}
}

func TestNamespace_NewIDWithValue(t *testing.T) {
	tests := map[string]struct {
		environment string
		objectType  Type
		value       string
		wantErr     bool
		errMsg      string
	}{
		"valid custom value": {
			environment: "dev",
			objectType:  Type("user"),
			value:       "custom123",
			wantErr:     false,
		},
		"valid with prd env": {
			environment: "prd",
			objectType:  Type("order"),
			value:       "ord_12345",
			wantErr:     false,
		},
		"empty value": {
			environment: "dev",
			objectType:  Type("user"),
			value:       "",
			wantErr:     true,
			errMsg:      "value cannot be empty",
		},
		"invalid type": {
			environment: "dev",
			objectType:  Type("1invalid"),
			value:       "custom123",
			wantErr:     true,
			errMsg:      "invalid object type: type must start with a letter",
		},
		"whitespace only value": {
			environment: "dev",
			objectType:  Type("user"),
			value:       "   ",
			wantErr:     false, // whitespace is allowed in values
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ns := NewNamespace(tt.environment)
			id, err := ns.NewIDWithValue(tt.objectType, tt.value)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewIDWithValue() expected error but got nil")
					return
				}
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("NewIDWithValue() error = %v, want error containing %q", err, tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("NewIDWithValue() unexpected error = %v", err)
				return
			}

			// Verify the ID has correct components
			expectedEnv := tt.environment
			if tt.environment == "prd" || strings.TrimSpace(tt.environment) == "" {
				expectedEnv = "vibe"
			}

			if id.Env() != expectedEnv {
				t.Errorf("NewIDWithValue().Env() = %q, want %q", id.Env(), expectedEnv)
			}

			if id.Type() != tt.objectType {
				t.Errorf("NewIDWithValue().Type() = %q, want %q", id.Type(), tt.objectType)
			}

			if id.Value() != tt.value {
				t.Errorf("NewIDWithValue().Value() = %q, want %q", id.Value(), tt.value)
			}

			// Verify string format
			expectedString := expectedEnv + ":" + string(tt.objectType) + ":" + tt.value
			if id.String() != expectedString {
				t.Errorf("NewIDWithValue().String() = %q, want %q", id.String(), expectedString)
			}
		})
	}
}

func TestNormalizeEnvironment(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected string
	}{
		"prd becomes vibe": {
			input:    "prd",
			expected: "vibe",
		},
		"empty becomes vibe": {
			input:    "",
			expected: "vibe",
		},
		"whitespace only becomes vibe": {
			input:    "   ",
			expected: "vibe",
		},
		"prd with whitespace becomes vibe": {
			input:    "  prd  ",
			expected: "vibe",
		},
		"dev stays dev": {
			input:    "dev",
			expected: "dev",
		},
		"staging stays staging": {
			input:    "staging",
			expected: "staging",
		},
		"trimmed whitespace": {
			input:    "  test-env  ",
			expected: "test-env",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := normalizeEnvironment(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeEnvironment(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// Test that NewID generates unique values on multiple calls
func TestNamespace_NewID_Uniqueness(t *testing.T) {
	ns := NewNamespace("test")
	objectType := Type("user")

	// Generate multiple IDs and ensure they're unique
	ids := make(map[string]bool)
	for i := 0; i < 10; i++ {
		id, err := ns.NewID(objectType)
		if err != nil {
			t.Fatalf("NewID() unexpected error = %v", err)
		}

		idString := id.String()
		if ids[idString] {
			t.Errorf("NewID() generated duplicate ID: %s", idString)
		}
		ids[idString] = true

		// Ensure the value part is unique
		if ids[id.Value()] {
			t.Errorf("NewID() generated duplicate value: %s", id.Value())
		}
		ids[id.Value()] = true
	}
}
