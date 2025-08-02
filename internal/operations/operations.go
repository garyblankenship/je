package operations

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/tidwall/sjson"
	"github.com/vampire/je/internal/parser"
)

// ApplyAssignments applies a list of assignments to JSON data
func ApplyAssignments(data []byte, assignments []parser.Assignment) ([]byte, error) {
	jsonStr := string(data)

	for _, assignment := range assignments {
		var err error
		jsonStr, err = applyAssignment(jsonStr, assignment)
		if err != nil {
			return nil, fmt.Errorf("failed to apply %s: %w", assignment.Path, err)
		}
	}

	return []byte(jsonStr), nil
}

func applyAssignment(jsonStr string, assignment parser.Assignment) (string, error) {
	switch assignment.Operator {
	case parser.OpAssignString:
		return applyStringAssignment(jsonStr, assignment.Path, assignment.Value)

	case parser.OpAssignJSON:
		return applyJSONAssignment(jsonStr, assignment.Path, assignment.Value)

	case parser.OpAssignFile:
		content, err := os.ReadFile(assignment.Value)
		if err != nil {
			return "", fmt.Errorf("failed to read file %s: %w", assignment.Value, err)
		}
		return applyStringAssignment(jsonStr, assignment.Path, string(content))

	case parser.OpAssignJSONFile:
		content, err := os.ReadFile(assignment.Value)
		if err != nil {
			return "", fmt.Errorf("failed to read file %s: %w", assignment.Value, err)
		}
		// Validate JSON
		var v interface{}
		if err := json.Unmarshal(content, &v); err != nil {
			return "", fmt.Errorf("invalid JSON in file %s: %w", assignment.Value, err)
		}
		return applyJSONAssignment(jsonStr, assignment.Path, string(content))

	case parser.OpAppendArray:
		return applyArrayAppend(jsonStr, assignment.Path, assignment.Value, false)

	case parser.OpAppendArrayJSON:
		return applyArrayAppend(jsonStr, assignment.Path, assignment.Value, true)

	case parser.OpArrayMap:
		return applyArrayMap(jsonStr, assignment.Path, assignment.Value, false)

	case parser.OpArrayMapJSON:
		return applyArrayMap(jsonStr, assignment.Path, assignment.Value, true)

	default:
		return "", fmt.Errorf("unknown operator type: %d", assignment.Operator)
	}
}

func applyStringAssignment(jsonStr, path, value string) (string, error) {
	result, err := sjson.Set(jsonStr, path, value)
	if err != nil {
		return "", err
	}
	return result, nil
}

func applyJSONAssignment(jsonStr, path, value string) (string, error) {
	// Handle special case: empty value means delete
	if value == "" {
		result, err := sjson.Delete(jsonStr, path)
		if err != nil {
			return "", err
		}
		return result, nil
	}

	// Parse the value as JSON
	v, err := parseJSONValue(value)
	if err != nil {
		return "", fmt.Errorf("invalid JSON value for %q: %w", path, err)
	}

	result, err := sjson.Set(jsonStr, path, v)
	if err != nil {
		return "", err
	}
	return result, nil
}

func applyArrayAppend(jsonStr, path, value string, isJSON bool) (string, error) {
	// Remove [] from path
	basePath := strings.TrimSuffix(path, "[]")

	// Ensure array exists
	var err error
	jsonStr, err = prepareArrayPath(jsonStr, basePath)
	if err != nil {
		return "", err
	}

	// Prepare value to append
	var appendValue interface{}
	if isJSON {
		appendValue, err = parseJSONValue(value)
		if err != nil {
			return "", fmt.Errorf("invalid JSON value for array append: %w", err)
		}
	} else {
		appendValue = value
	}

	// Append to array
	return appendToArray(jsonStr, basePath, appendValue)
}

func applyArrayMap(jsonStr, path, value string, isJSON bool) (string, error) {
	// Parse the array map path
	basePath, property, err := parseArrayMapPath(path)
	if err != nil {
		return "", err
	}

	// Validate the array exists
	if err := validateArrayPath(jsonStr, basePath); err != nil {
		return "", err
	}

	// Prepare the value
	var setValue interface{}
	if isJSON {
		setValue, err = parseJSONValue(value)
		if err != nil {
			return "", fmt.Errorf("invalid JSON value for array map: %w", err)
		}
	} else {
		setValue = value
	}

	// Apply to each array element
	return applyToArrayElements(jsonStr, basePath, property, setValue)
}
