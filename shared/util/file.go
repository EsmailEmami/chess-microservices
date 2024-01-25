package util

import (
	"os"
)

// ReadFile reads the content of a file and returns it as a string.
func ReadFile(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// WriteFile writes content to a file.
func WriteFile(filename string, content []byte) error {
	return os.WriteFile(filename, content, 0644)
}

// FileExists checks if a file exists at the specified path.
func FileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return !os.IsNotExist(err)
}
