package idx

import (
	"strings"
	"testing"
)

func TestID_Env(t *testing.T) {
	id := ID{
		env:        "test-env",
		objectType: Type("user"),
		objectID:   "123",
	}

	result := id.Env()
	expected := "test-env"

	if result != expected {
		t.Errorf("Env() = %q, want %q", result, expected)
	}
}

func TestID_Type(t *testing.T) {
	expectedType := Type("user")
	id := ID{
		env:        "test-env",
		objectType: expectedType,
		objectID:   "123",
	}

	result := id.Type()

	if result != expectedType {
		t.Errorf("Type() = %q, want %q", result, expectedType)
	}
}

func TestID_Value(t *testing.T) {
	id := ID{
		env:        "test-env",
		objectType: Type("user"),
		objectID:   "test-value-123",
	}

	result := id.Value()
	expected := "test-value-123"

	if result != expected {
		t.Errorf("Value() = %q, want %q", result, expected)
	}
}

func TestID_String(t *testing.T) {
	tests := map[string]struct {
		id       ID
		expected string
	}{
		"basic ID": {
			id: ID{
				env:        "dev",
				objectType: Type("user"),
				objectID:   "123",
			},
			expected: "dev:user:123",
		},
		"vibe environment": {
			id: ID{
				env:        "vibe",
				objectType: Type("session"),
				objectID:   "abc123",
			},
			expected: "vibe:session:abc123",
		},
		"complex object ID": {
			id: ID{
				env:        "staging",
				objectType: Type("order_item"),
				objectID:   "ord_12345_item_67890",
			},
			expected: "staging:order_item:ord_12345_item_67890",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := tt.id.String()
			if result != tt.expected {
				t.Errorf("String() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestParseID(t *testing.T) {
	tests := map[string]struct {
		input      string
		expectedID ID
		wantErr    bool
		errMsg     string
	}{
		// Valid cases
		"basic valid ID": {
			input: "dev:user:123",
			expectedID: ID{
				env:        "dev",
				objectType: "user",
				objectID:   "123",
			},
			wantErr: false,
		},
		"vibe environment": {
			input: "vibe:session:abc123",
			expectedID: ID{
				env:        "vibe",
				objectType: "session",
				objectID:   "abc123",
			},
			wantErr: false,
		},
		"complex type and value": {
			input: "staging:order_item:ord_12345_item_67890",
			expectedID: ID{
				env:        "staging",
				objectType: "order_item",
				objectID:   "ord_12345_item_67890",
			},
			wantErr: false,
		},
		"long KSUID-like value": {
			input: "prod:user:2B5E5fLHQjw1234567890123456",
			expectedID: ID{
				env:        "prod",
				objectType: "user",
				objectID:   "2B5E5fLHQjw1234567890123456",
			},
			wantErr: false,
		},
		// Invalid cases
		"too few parts": {
			input:   "dev:user",
			wantErr: true,
			errMsg:  "invalid ID format: expected 3 parts separated by colons, got 2 parts",
		},
		"too many parts": {
			input:   "dev:user:123:extra",
			wantErr: true,
			errMsg:  "invalid ID format: expected 3 parts separated by colons, got 4 parts",
		},
		"empty environment": {
			input:   ":user:123",
			wantErr: true,
			errMsg:  "invalid ID: env cannot be empty",
		},
		"empty type": {
			input:   "dev::123",
			wantErr: true,
			errMsg:  "invalid ID: type cannot be empty",
		},
		"empty object ID": {
			input:   "dev:user:",
			wantErr: true,
			errMsg:  "invalid ID: object ID cannot be empty",
		},
		"invalid type - starts with number": {
			input:   "dev:1user:123",
			wantErr: true,
			errMsg:  "invalid ID: type must start with a letter",
		},
		"invalid type - contains special char": {
			input:   "dev:user@item:123",
			wantErr: true,
			errMsg:  "invalid ID: type must start with a letter",
		},
		"invalid type - too long": {
			input:   "dev:" + strings.Repeat("a", 33) + ":123",
			wantErr: true,
			errMsg:  "invalid ID: type cannot be longer than 32 characters",
		},
		"no colons": {
			input:   "devuser123",
			wantErr: true,
			errMsg:  "invalid ID format: expected 3 parts separated by colons, got 1 parts",
		},
		"empty string": {
			input:   "",
			wantErr: true,
			errMsg:  "invalid ID format: expected 3 parts separated by colons, got 1 parts",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := ParseID(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseID() expected error but got nil")
					return
				}
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("ParseID() error = %v, want error containing %q", err, tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("ParseID() unexpected error = %v", err)
				return
			}

			// Verify all components
			if result.env != tt.expectedID.env {
				t.Errorf("ParseID().env = %q, want %q", result.env, tt.expectedID.env)
			}
			if result.objectType != tt.expectedID.objectType {
				t.Errorf("ParseID().objectType = %q, want %q", result.objectType, tt.expectedID.objectType)
			}
			if result.objectID != tt.expectedID.objectID {
				t.Errorf("ParseID().objectID = %q, want %q", result.objectID, tt.expectedID.objectID)
			}
		})
	}
}

func TestID_Validate(t *testing.T) {
	tests := map[string]struct {
		id      ID
		wantErr bool
		errMsg  string
	}{
		// Valid cases
		"valid ID": {
			id: ID{
				env:        "dev",
				objectType: Type("user"),
				objectID:   "123",
			},
			wantErr: false,
		},
		"valid with vibe env": {
			id: ID{
				env:        "vibe",
				objectType: Type("session"),
				objectID:   "abc123",
			},
			wantErr: false,
		},
		"valid complex type": {
			id: ID{
				env:        "staging",
				objectType: Type("order_item"),
				objectID:   "complex_value_123",
			},
			wantErr: false,
		},
		// Invalid cases
		"empty environment": {
			id: ID{
				env:        "",
				objectType: Type("user"),
				objectID:   "123",
			},
			wantErr: true,
			errMsg:  "env cannot be empty",
		},
		"invalid type - empty": {
			id: ID{
				env:        "dev",
				objectType: Type(""),
				objectID:   "123",
			},
			wantErr: true,
			errMsg:  "invalid object type: type cannot be empty",
		},
		"invalid type - starts with number": {
			id: ID{
				env:        "dev",
				objectType: Type("1user"),
				objectID:   "123",
			},
			wantErr: true,
			errMsg:  "invalid object type: type must start with a letter",
		},
		"invalid type - contains colon": {
			id: ID{
				env:        "dev",
				objectType: Type("user:item"),
				objectID:   "123",
			},
			wantErr: true,
			errMsg:  "invalid object type: type cannot contain colons",
		},
		"empty object ID": {
			id: ID{
				env:        "dev",
				objectType: Type("user"),
				objectID:   "",
			},
			wantErr: true,
			errMsg:  "object ID cannot be empty",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := tt.id.Validate()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Validate() expected error but got nil")
					return
				}
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Validate() error = %v, want error containing %q", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Validate() unexpected error = %v", err)
				}
			}
		})
	}
}

// Test roundtrip: create ID, convert to string, parse back
func TestID_Roundtrip(t *testing.T) {
	tests := map[string]struct {
		id ID
	}{
		"basic roundtrip": {
			id: ID{
				env:        "dev",
				objectType: Type("user"),
				objectID:   "123",
			},
		},
		"vibe environment roundtrip": {
			id: ID{
				env:        "vibe",
				objectType: Type("session"),
				objectID:   "abc123",
			},
		},
		"complex roundtrip": {
			id: ID{
				env:        "staging",
				objectType: Type("order_item"),
				objectID:   "ord_12345_item_67890",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Convert to string
			idString := tt.id.String()

			// Parse back
			parsed, err := ParseID(idString)
			if err != nil {
				t.Errorf("ParseID() unexpected error = %v", err)
				return
			}

			// Compare all components
			if parsed.env != tt.id.env {
				t.Errorf("roundtrip env mismatch: got %q, want %q", parsed.env, tt.id.env)
			}
			if parsed.objectType != tt.id.objectType {
				t.Errorf("roundtrip objectType mismatch: got %q, want %q", parsed.objectType, tt.id.objectType)
			}
			if parsed.objectID != tt.id.objectID {
				t.Errorf("roundtrip objectID mismatch: got %q, want %q", parsed.objectID, tt.id.objectID)
			}

			// Ensure the string representation is the same
			if parsed.String() != idString {
				t.Errorf("roundtrip string mismatch: got %q, want %q", parsed.String(), idString)
			}
		})
	}
}

// Test integration with Namespace
func TestID_Integration(t *testing.T) {
	ns := NewNamespace("prd") // should become "vibe"
	objectType := Type("user")

	// Test auto-generated ID
	autoID, err := ns.NewID(objectType)
	if err != nil {
		t.Fatalf("NewID() unexpected error = %v", err)
	}

	if autoID.Env() != "vibe" {
		t.Errorf("NewID().Env() = %q, want %q", autoID.Env(), "vibe")
	}
	if autoID.Type() != objectType {
		t.Errorf("NewID().Type() = %q, want %q", autoID.Type(), objectType)
	}
	if autoID.Value() == "" {
		t.Errorf("NewID().Value() should not be empty")
	}

	// Test parsing the auto-generated ID
	parsed, err := ParseID(autoID.String())
	if err != nil {
		t.Fatalf("ParseID() unexpected error = %v", err)
	}

	if parsed.String() != autoID.String() {
		t.Errorf("ParseID roundtrip failed: got %q, want %q", parsed.String(), autoID.String())
	}

	// Test custom value ID
	customID, err := ns.NewIDWithValue(objectType, "custom_123")
	if err != nil {
		t.Fatalf("NewIDWithValue() unexpected error = %v", err)
	}

	expectedString := "vibe:user:custom_123"
	if customID.String() != expectedString {
		t.Errorf("NewIDWithValue().String() = %q, want %q", customID.String(), expectedString)
	}
}
