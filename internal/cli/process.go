package cli

import (
	"fmt"

	"github.com/vampire/je/internal/json"
	"github.com/vampire/je/internal/operations"
	"github.com/vampire/je/internal/parser"
)

// ProcessResult holds the result of processing a JSON file.
type ProcessResult struct {
	Original []byte
	Modified []byte
	Filename string
}

// ProcessJSONFile applies assignments to a JSON file and returns the result.
func ProcessJSONFile(filename string, assignments []parser.Assignment, createIfMissing bool) (*ProcessResult, error) {
	// Read JSON file
	data, err := ReadJSONFile(filename, createIfMissing)
	if err != nil {
		return nil, err
	}

	// Validate JSON
	if err := json.Validate(data); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	// Apply assignments
	result, err := operations.ApplyAssignments(data, assignments)
	if err != nil {
		return nil, err
	}

	return &ProcessResult{
		Original: data,
		Modified: result,
		Filename: filename,
	}, nil
}

// WriteResult writes the result to the appropriate destination.
func WriteResult(result []byte, filename, outputFile string) error {
	// Determine output destination
	output := filename
	if outputFile != "" {
		output = outputFile
	}

	// Get original file permissions
	perm := GetFilePermissions(filename)

	// Write result
	if err := json.WriteFile(output, result, perm); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	return nil
}