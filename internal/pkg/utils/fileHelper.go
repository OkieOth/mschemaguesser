package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func CleanDirectory(dir string, includeDotFiles bool) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	if len(files) == 0 {
		return nil
	}

	for _, file := range files {
		if includeDotFiles && strings.HasPrefix(file.Name(), ".") {
			continue
		}
		filePath := filepath.Join(dir, file.Name())
		err = os.RemoveAll(filePath)
		if err != nil {
			return fmt.Errorf("failed to delete file: %w", err)
		}
	}
	return nil
}

func GetFilesInDir(dir string, includeDotFiles bool) ([]string, error) {
	ret := make([]string, 0)
	files, err := os.ReadDir(dir)
	if err != nil {
		return ret, fmt.Errorf("failed to read directory: %w", err)
	}

	if len(files) == 0 {
		return ret, nil
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if includeDotFiles && strings.HasPrefix(file.Name(), ".") {
			continue
		}
		ret = append(ret, file.Name())
	}
	return ret, nil
}
