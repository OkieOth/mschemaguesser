package utils

import (
	"fmt"
	"log"
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
		if (!includeDotFiles) && strings.HasPrefix(file.Name(), ".") {
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
		if (!includeDotFiles) && strings.HasPrefix(file.Name(), ".") {
			continue
		}
		ret = append(ret, file.Name())
	}
	return ret, nil
}

func DirExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return info.IsDir(), nil
}

func RemoveDirectory(path string) {
	if exist, err := DirExists(path); err != nil {
		log.Printf("Error while checking for directory (%s) existence: %v", path, err)
	} else {
		if exist {
			err := os.RemoveAll(path)
			if err != nil {
				log.Printf("Error while deleting directory (%s) existence: %v", path, err)
			}
		}
	}
}

func PrepareDirStructure(outputDir string, dbName string, collName string) (string, error) {
	dbSanitized := Sanitize(dbName)
	collSanitized := Sanitize(collName)
	dirPath := filepath.Join(outputDir, filepath.Join(dbSanitized, collSanitized))
	if exist, err := DirExists(dirPath); err != nil {
		return dirPath, err
	} else if !exist {
		err := os.MkdirAll(dirPath, 0755)
		return dirPath, err
	}
	return dirPath, nil
}
