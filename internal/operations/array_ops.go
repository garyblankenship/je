package operations

import (
	"fmt"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// prepareArrayPath ensures the path contains an array, creating one if necessary.
// Returns the updated JSON string and an error if the path exists but is not an array.
func prepareArrayPath(jsonStr, basePath string) (string, error) {
	result := gjson.Get(jsonStr, basePath)
	if !result.Exists() {
		// Create new array
		return sjson.Set(jsonStr, basePath, []interface{}{})
	}
	if !result.IsArray() {
		return "", fmt.Errorf("cannot append to non-array at path %q", basePath)
	}
	return jsonStr, nil
}

// appendToArray adds a value to the end of an array at the given path.
func appendToArray(jsonStr, basePath string, value interface{}) (string, error) {
	result := gjson.Get(jsonStr, basePath)
	currentArray := result.Array()
	newIndex := len(currentArray)
	newPath := fmt.Sprintf("%s.%d", basePath, newIndex)

	return sjson.Set(jsonStr, newPath, value)
}

