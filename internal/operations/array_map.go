package operations

import (
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// parseArrayMapPath extracts the base path and property from an array map path.
// Example: "users.[].active" -> basePath: "users", property: "active"
func parseArrayMapPath(path string) (basePath, property string, err error) {
	idx := strings.Index(path, "[].")
	if idx < 0 {
		return "", "", fmt.Errorf("invalid array map path %q: expected format like 'users.[].property'", path)
	}

	basePath = strings.TrimSuffix(path[:idx], ".")
	property = path[idx+3:]

	return basePath, property, nil
}

// validateArrayPath ensures the path exists and is an array.
func validateArrayPath(jsonStr, basePath string) error {
	result := gjson.Get(jsonStr, basePath)
	if !result.Exists() {
		return fmt.Errorf("array path %q does not exist", basePath)
	}
	if !result.IsArray() {
		return fmt.Errorf("path %q is not an array", basePath)
	}
	return nil
}

// applyToArrayElements sets a property on each element of an array.
func applyToArrayElements(jsonStr, basePath, property string, value interface{}) (string, error) {
	result := gjson.Get(jsonStr, basePath)
	array := result.Array()

	for i := range array {
		elementPath := fmt.Sprintf("%s.%d.%s", basePath, i, property)
		var err error
		jsonStr, err = sjson.Set(jsonStr, elementPath, value)
		if err != nil {
			return "", fmt.Errorf("failed to set %s: %w", elementPath, err)
		}
	}

	return jsonStr, nil
}

