package parser

import (
	"reflect"
	"testing"
)

func TestParseAssignments(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected []Assignment
		wantErr  bool
	}{
		{
			name: "string assignment",
			args: []string{"name=john", "city=NYC"},
			expected: []Assignment{
				{Path: "name", Operator: OpAssignString, Value: "john"},
				{Path: "city", Operator: OpAssignString, Value: "NYC"},
			},
		},
		{
			name: "JSON assignment",
			args: []string{"age:=30", "active:=true", "balance:=99.50"},
			expected: []Assignment{
				{Path: "age", Operator: OpAssignJSON, Value: "30"},
				{Path: "active", Operator: OpAssignJSON, Value: "true"},
				{Path: "balance", Operator: OpAssignJSON, Value: "99.50"},
			},
		},
		{
			name: "file assignment",
			args: []string{"cert@cert.pem", "config:@config.json"},
			expected: []Assignment{
				{Path: "cert", Operator: OpAssignFile, Value: "cert.pem"},
				{Path: "config", Operator: OpAssignJSONFile, Value: "config.json"},
			},
		},
		{
			name: "array append",
			args: []string{"tags[]=new", "ports[]:=8080"},
			expected: []Assignment{
				{Path: "tags[]", Operator: OpAppendArray, Value: "new"},
				{Path: "ports[]", Operator: OpAppendArrayJSON, Value: "8080"},
			},
		},
		{
			name: "nested paths",
			args: []string{"user.name=john", "config.server.port:=3000"},
			expected: []Assignment{
				{Path: "user.name", Operator: OpAssignString, Value: "john"},
				{Path: "config.server.port", Operator: OpAssignJSON, Value: "3000"},
			},
		},
		{
			name: "empty values",
			args: []string{"empty=", "delete:="},
			expected: []Assignment{
				{Path: "empty", Operator: OpAssignString, Value: ""},
				{Path: "delete", Operator: OpAssignJSON, Value: ""},
			},
		},
		{
			name:    "invalid assignment",
			args:    []string{"invalid"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseAssignments(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseAssignments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ParseAssignments() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestParsePath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected []string
	}{
		{
			name:     "simple path",
			path:     "user.name",
			expected: []string{"user", "name"},
		},
		{
			name:     "nested path",
			path:     "config.server.port",
			expected: []string{"config", "server", "port"},
		},
		{
			name:     "array index",
			path:     "users.0.name",
			expected: []string{"users", "0", "name"},
		},
		{
			name:     "escaped dot",
			path:     `user\.name.value`,
			expected: []string{"user.name", "value"},
		},
		{
			name:     "single segment",
			path:     "value",
			expected: []string{"value"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParsePath(tt.path)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ParsePath() = %v, want %v", got, tt.expected)
			}
		})
	}
}
