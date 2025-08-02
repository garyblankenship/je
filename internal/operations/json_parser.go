package operations

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// parseJSONValue attempts to parse a string as a JSON value.
// It handles special cases like null, true, false, and numbers.
func parseJSONValue(value string) (interface{}, error) {
	// Try parsing as valid JSON first
	var v interface{}
	if err := json.Unmarshal([]byte(value), &v); err == nil {
		return v, nil
	}

	// Try parsing as simple values
	switch strings.ToLower(value) {
	case "null":
		return nil, nil
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		// Try as number
		if f, err := strconv.ParseFloat(value, 64); err == nil {
			return f, nil
		}
		return nil, fmt.Errorf("invalid JSON value: %q", value)
	}
}

