package idx

import (
	"strings"
	"testing"
)

func TestType_String(t *testing.T) {
	tests := map[string]struct {
		typ      Type
		expected string
	}{
		"simple type": {
			typ:      Type("user"),
			expected: "user",
		},
		"with underscore": {
			typ:      Type("order_item"),
			expected: "order_item",
		},
		"with numbers": {
			typ:      Type("session123"),
			expected: "session123",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := tt.typ.String()
			if result != tt.expected {
				t.Errorf("String() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestType_Validate(t *testing.T) {
	tests := map[string]struct {
		typ     Type
		wantErr bool
		errMsg  string
	}{
		// Valid cases
		"valid simple type": {
			typ:     Type("user"),
			wantErr: false,
		},
		"valid with underscore": {
			typ:     Type("order_item"),
			wantErr: false,
		},
		"valid with numbers": {
			typ:     Type("session123"),
			wantErr: false,
		},
		"valid mixed case": {
			typ:     Type("API_Key"),
			wantErr: false,
		},
		"valid single letter": {
			typ:     Type("a"),
			wantErr: false,
		},
		"valid max length": {
			typ:     Type("a" + strings.Repeat("b", 31)), // 32 chars total
			wantErr: false,
		},
		// Invalid cases
		"empty type": {
			typ:     Type(""),
			wantErr: true,
			errMsg:  "type cannot be empty",
		},
		"too long": {
			typ:     Type("a" + strings.Repeat("b", 32)), // 33 chars total
			wantErr: true,
			errMsg:  "type cannot be longer than 32 characters",
		},
		"contains colon": {
			typ:     Type("user:item"),
			wantErr: true,
			errMsg:  "type cannot contain colons",
		},
		"starts with number": {
			typ:     Type("1user"),
			wantErr: true,
			errMsg:  "type must start with a letter and contain only letters, numbers, and underscores",
		},
		"starts with underscore": {
			typ:     Type("_user"),
			wantErr: true,
			errMsg:  "type must start with a letter and contain only letters, numbers, and underscores",
		},
		"contains hyphen": {
			typ:     Type("user-item"),
			wantErr: true,
			errMsg:  "type must start with a letter and contain only letters, numbers, and underscores",
		},
		"contains space": {
			typ:     Type("user item"),
			wantErr: true,
			errMsg:  "type must start with a letter and contain only letters, numbers, and underscores",
		},
		"contains special char": {
			typ:     Type("user@item"),
			wantErr: true,
			errMsg:  "type must start with a letter and contain only letters, numbers, and underscores",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := tt.typ.Validate()
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

func TestParseType(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected Type
		wantErr  bool
		errMsg   string
	}{
		// Valid cases
		"valid simple type": {
			input:    "user",
			expected: Type("user"),
			wantErr:  false,
		},
		"valid with underscore": {
			input:    "order_item",
			expected: Type("order_item"),
			wantErr:  false,
		},
		"valid with numbers": {
			input:    "session123",
			expected: Type("session123"),
			wantErr:  false,
		},
		// Invalid cases
		"empty string": {
			input:   "",
			wantErr: true,
			errMsg:  "type cannot be empty",
		},
		"starts with number": {
			input:   "1user",
			wantErr: true,
			errMsg:  "type must start with a letter",
		},
		"contains colon": {
			input:   "user:item",
			wantErr: true,
			errMsg:  "type cannot contain colons",
		},
		"too long": {
			input:   strings.Repeat("a", 33),
			wantErr: true,
			errMsg:  "type cannot be longer than 32 characters",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := ParseType(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseType() expected error but got nil")
					return
				}
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("ParseType() error = %v, want error containing %q", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ParseType() unexpected error = %v", err)
					return
				}
				if result != tt.expected {
					t.Errorf("ParseType() = %q, want %q", result, tt.expected)
				}
			}
		})
	}
}
