package operations

import (
	"testing"

	"github.com/vampire/je/internal/parser"
)

func TestArrayMapEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		assignments []parser.Assignment
		wantErr     bool
		errContains string
	}{
		{
			name:  "array map on non-existent path",
			input: `{"users": []}`,
			assignments: []parser.Assignment{
				{Path: "nonexistent.[].name", Operator: parser.OpArrayMap, Value: "test"},
			},
			wantErr:     true,
			errContains: "does not exist",
		},
		{
			name:  "array map on non-array",
			input: `{"users": "not an array"}`,
			assignments: []parser.Assignment{
				{Path: "users.[].name", Operator: parser.OpArrayMap, Value: "test"},
			},
			wantErr:     true,
			errContains: "not an array",
		},
		{
			name:  "array map with invalid path format",
			input: `{"users": []}`,
			assignments: []parser.Assignment{
				{Path: "users.[]", Operator: parser.OpArrayMap, Value: "test"},
			},
			wantErr:     true,
			errContains: "invalid array map path",
		},
		{
			name:  "array map on empty array",
			input: `{"users": []}`,
			assignments: []parser.Assignment{
				{Path: "users.[].name", Operator: parser.OpArrayMap, Value: "test"},
			},
			wantErr: false, // Should succeed but do nothing
		},
		{
			name:  "deeply nested array map",
			input: `{"data": {"items": [{"tags": [{"name": "old"}]}]}}`,
			assignments: []parser.Assignment{
				{Path: "data.items.0.tags.[].name", Operator: parser.OpArrayMap, Value: "new"},
			},
			wantErr: false,
		},
		{
			name:  "array map with special characters in property",
			input: `{"items": [{"user.name": "old"}]}`,
			assignments: []parser.Assignment{
				{Path: "items.[].user\\.name", Operator: parser.OpArrayMap, Value: "new"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ApplyAssignments([]byte(tt.input), tt.assignments)

			if (err != nil) != tt.wantErr {
				t.Errorf("ApplyAssignments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.errContains != "" {
				if !contains(err.Error(), tt.errContains) {
					t.Errorf("Error message %q should contain %q", err.Error(), tt.errContains)
				}
			}
		})
	}
}

func TestArrayAppendEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		assignments []parser.Assignment
		wantErr     bool
		errContains string
	}{
		{
			name:  "append to non-array",
			input: `{"tags": "not an array"}`,
			assignments: []parser.Assignment{
				{Path: "tags[]", Operator: parser.OpAppendArray, Value: "new"},
			},
			wantErr:     true,
			errContains: "cannot append to non-array",
		},
		{
			name:  "append with invalid JSON value",
			input: `{"numbers": [1, 2]}`,
			assignments: []parser.Assignment{
				{Path: "numbers[]", Operator: parser.OpAppendArrayJSON, Value: "{invalid json}"},
			},
			wantErr:     true,
			errContains: "invalid JSON value",
		},
		{
			name:  "append to deeply nested array",
			input: `{"a": {"b": {"c": [1, 2]}}}`,
			assignments: []parser.Assignment{
				{Path: "a.b.c[]", Operator: parser.OpAppendArrayJSON, Value: "3"},
			},
			wantErr: false,
		},
		{
			name:  "append to array with null values",
			input: `{"items": [null, null]}`,
			assignments: []parser.Assignment{
				{Path: "items[]", Operator: parser.OpAppendArray, Value: "value"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ApplyAssignments([]byte(tt.input), tt.assignments)

			if (err != nil) != tt.wantErr {
				t.Errorf("ApplyAssignments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.errContains != "" {
				if !contains(err.Error(), tt.errContains) {
					t.Errorf("Error message %q should contain %q", err.Error(), tt.errContains)
				}
			}
		})
	}
}

func TestInvalidPaths(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		assignments []parser.Assignment
		wantErr     bool
		errContains string
	}{
		{
			name:  "path with multiple array notations",
			input: `{"data": [[]]}`,
			assignments: []parser.Assignment{
				{Path: "data.[].[].value", Operator: parser.OpArrayMap, Value: "test"},
			},
			wantErr:     true,
			errContains: "failed to set", // sjson error when trying to set invalid path
		},
		{
			name:  "empty path component",
			input: `{"data": {"": {"value": "old"}}}`,
			assignments: []parser.Assignment{
				{Path: "data..value", Operator: parser.OpAssignString, Value: "new"},
			},
			wantErr: false, // gjson handles this
		},
		{
			name:  "array index out of bounds",
			input: `{"items": [1, 2]}`,
			assignments: []parser.Assignment{
				{Path: "items.10", Operator: parser.OpAssignJSON, Value: "3"},
			},
			wantErr: false, // sjson creates intermediate elements
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ApplyAssignments([]byte(tt.input), tt.assignments)

			if (err != nil) != tt.wantErr {
				t.Errorf("ApplyAssignments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.errContains != "" {
				if !contains(err.Error(), tt.errContains) {
					t.Errorf("Error message %q should contain %q", err.Error(), tt.errContains)
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr ||
		len(s) > len(substr) && containsHelper(s[1:], substr)
}

func containsHelper(s, substr string) bool {
	if len(s) < len(substr) {
		return false
	}
	if s[:len(substr)] == substr {
		return true
	}
	return containsHelper(s[1:], substr)
}

