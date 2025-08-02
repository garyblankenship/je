package cli

import (
	"fmt"
	"os"

	"github.com/vampire/je/internal/json"
)

// ReadJSONFile reads JSON from a file or stdin, creating an empty object if --create is set.
func ReadJSONFile(filename string, createIfMissing bool) ([]byte, error) {
	data, err := json.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) && createIfMissing {
			return []byte("{}"), nil
		}
		if filename == "-" {
			return nil, fmt.Errorf("failed to read from stdin: %w", err)
		}
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return data, nil
}

// GetFilePermissions returns the file permissions, defaulting to 0644 if not found.
func GetFilePermissions(filename string) os.FileMode {
	if perm, err := json.GetFileInfo(filename); err == nil {
		return perm
	}
	return 0644
}

