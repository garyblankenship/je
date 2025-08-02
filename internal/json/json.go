package json

import (
	"encoding/json"
	"io"
	"os"
)

// ReadFile reads JSON from a file
func ReadFile(path string) ([]byte, error) {
	if path == "-" {
		return io.ReadAll(os.Stdin)
	}
	return os.ReadFile(path)
}

// WriteFile writes JSON to a file atomically
func WriteFile(path string, data []byte, perm os.FileMode) error {
	if path == "-" {
		_, err := os.Stdout.Write(data)
		return err
	}

	// Write to temp file first
	tempFile := path + ".tmp"
	if err := os.WriteFile(tempFile, data, perm); err != nil {
		return err
	}

	// Rename atomically
	return os.Rename(tempFile, path)
}

// Validate checks if data is valid JSON
func Validate(data []byte) error {
	var v interface{}
	return json.Unmarshal(data, &v)
}

// Format formats JSON with optional pretty printing
func Format(data []byte, pretty bool, compact bool) ([]byte, error) {
	if !pretty && !compact {
		return data, nil
	}

	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, err
	}

	if compact {
		return json.Marshal(v)
	}

	// Pretty print
	return json.MarshalIndent(v, "", "  ")
}

// GetFileInfo gets file permissions
func GetFileInfo(path string) (os.FileMode, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0644, err
	}
	return info.Mode(), nil
}
