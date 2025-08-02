package operations

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/vampire/je/internal/parser"
)

func TestApplyAssignments(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		assignments []parser.Assignment
		expected    map[string]interface{}
		wantErr     bool
	}{
		{
			name:  "string assignment",
			input: "{}",
			assignments: []parser.Assignment{
				{Path: "name", Operator: parser.OpAssignString, Value: "john"},
			},
			expected: map[string]interface{}{"name": "john"},
		},
		{
			name:  "number assignment",
			input: "{}",
			assignments: []parser.Assignment{
				{Path: "age", Operator: parser.OpAssignJSON, Value: "30"},
			},
			expected: map[string]interface{}{"age": float64(30)},
		},
		{
			name:  "boolean assignment",
			input: "{}",
			assignments: []parser.Assignment{
				{Path: "active", Operator: parser.OpAssignJSON, Value: "true"},
			},
			expected: map[string]interface{}{"active": true},
		},
		{
			name:  "null assignment",
			input: `{"value": "something"}`,
			assignments: []parser.Assignment{
				{Path: "value", Operator: parser.OpAssignJSON, Value: "null"},
			},
			expected: map[string]interface{}{"value": nil},
		},
		{
			name:  "nested assignment",
			input: "{}",
			assignments: []parser.Assignment{
				{Path: "user.name", Operator: parser.OpAssignString, Value: "john"},
				{Path: "user.age", Operator: parser.OpAssignJSON, Value: "30"},
			},
			expected: map[string]interface{}{
				"user": map[string]interface{}{
					"name": "john",
					"age":  float64(30),
				},
			},
		},
		{
			name:  "array append",
			input: `{"tags": ["old"]}`,
			assignments: []parser.Assignment{
				{Path: "tags[]", Operator: parser.OpAppendArray, Value: "new"},
			},
			expected: map[string]interface{}{
				"tags": []interface{}{"old", "new"},
			},
		},
		{
			name:  "array append JSON",
			input: `{"ports": [8080]}`,
			assignments: []parser.Assignment{
				{Path: "ports[]", Operator: parser.OpAppendArrayJSON, Value: "9090"},
			},
			expected: map[string]interface{}{
				"ports": []interface{}{float64(8080), float64(9090)},
			},
		},
		{
			name:  "delete key",
			input: `{"name": "john", "age": 30}`,
			assignments: []parser.Assignment{
				{Path: "age", Operator: parser.OpAssignJSON, Value: ""},
			},
			expected: map[string]interface{}{"name": "john"},
		},
		{
			name:  "complex nested structure",
			input: "{}",
			assignments: []parser.Assignment{
				{Path: "config.server.host", Operator: parser.OpAssignString, Value: "localhost"},
				{Path: "config.server.port", Operator: parser.OpAssignJSON, Value: "3000"},
				{Path: "config.debug", Operator: parser.OpAssignJSON, Value: "true"},
			},
			expected: map[string]interface{}{
				"config": map[string]interface{}{
					"server": map[string]interface{}{
						"host": "localhost",
						"port": float64(3000),
					},
					"debug": true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ApplyAssignments([]byte(tt.input), tt.assignments)
			if (err != nil) != tt.wantErr {
				t.Errorf("ApplyAssignments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				var got map[string]interface{}
				if err := json.Unmarshal(result, &got); err != nil {
					t.Fatalf("Failed to unmarshal result: %v", err)
				}

				if !jsonEqual(got, tt.expected) {
					t.Errorf("ApplyAssignments() = %v, want %v", got, tt.expected)
				}
			}
		})
	}
}

func TestFileOperations(t *testing.T) {
	// Create temp file for testing
	tmpfile, err := os.CreateTemp("", "test*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	testContent := "Hello from file"
	if _, err := tmpfile.Write([]byte(testContent)); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	// Test file assignment
	assignments := []parser.Assignment{
		{Path: "content", Operator: parser.OpAssignFile, Value: tmpfile.Name()},
	}

	result, err := ApplyAssignments([]byte("{}"), assignments)
	if err != nil {
		t.Errorf("ApplyAssignments() with file error = %v", err)
		return
	}

	var got map[string]interface{}
	if err := json.Unmarshal(result, &got); err != nil {
		t.Fatalf("Failed to unmarshal result: %v", err)
	}

	if got["content"] != testContent {
		t.Errorf("File content = %v, want %v", got["content"], testContent)
	}
}

// Helper function to compare JSON objects
func jsonEqual(a, b interface{}) bool {
	aJSON, _ := json.Marshal(a)
	bJSON, _ := json.Marshal(b)
	return string(aJSON) == string(bJSON)
}
