package testhelper

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"slices"
	"testing"

	"okieoth/schemaguesser/internal/pkg/utils"
)

func ValidateExpectedFiles(dir string, expected []string, t *testing.T) bool {
	if exist, _ := utils.DirExists(dir); !exist {
		e := os.MkdirAll(dir, 0755)
		if e != nil {
			t.Errorf("Error while creating test dir (%s): %v", dir, e)
			return false
		}
	}
	files, err := utils.GetFilesInDir(dir, false)
	if err != nil {
		t.Errorf("Error while getting files via GetFilesInDir (%s): %v", dir, err)
		return false
	}

	lenExpected := len(expected)

	if len(files) != lenExpected {
		t.Errorf("Retrieved wrong number of files in dir (%s). Expected %d, got %d, %v", dir, lenExpected, len(files), files)
		return false
	}
	ret := true
	for _, f := range expected {
		if !slices.Contains(files, f) {
			ret = false
			t.Errorf("Expected file %s not found, expected files: %v", f, files)
		}
	}
	return ret
}

func ValidateEmptyDir(dir string, t *testing.T) bool {
	if exist, _ := utils.DirExists(dir); !exist {
		e := os.MkdirAll(dir, 0755)
		if e != nil {
			t.Errorf("Error while creating test dir (%s): %v", dir, e)
			return false
		}
	}
	files, err := utils.GetFilesInDir(dir, false)

	if err != nil {
		t.Errorf("Error while getting files via GetFilesInDir (%s): %v", dir, err)
		return false
	}

	if len(files) > 0 {
		t.Errorf("Retrieved wrong number of files in dir (%s), got %d, %v", dir, len(files), files)
		return false
	}
	return true
}

func CheckFilesNonZero(dir string, filesToInclude []string, t *testing.T) (bool, error) {
	directory, err := os.Open(dir)
	if err != nil {
		t.Errorf("Error while open directory (%s): %v", dir, err)
		return false, err
	}
	defer directory.Close()

	files, err := directory.Readdir(-1)
	if err != nil {
		t.Errorf("Error while read directory (%s): %v", dir, err)
		return false, err
	}

	ret := true
	for _, file := range files {
		if (!file.IsDir()) && slices.Contains(filesToInclude, file.Name()) {
			if file.Size() == 0 {
				t.Errorf("File has size of 0: %s", file.Name())
				ret = false
			}
		}
	}

	return ret, nil
}

func hashFileMD5(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not open file %s: %w", filePath, err)
	}
	defer file.Close()

	hash := md5.New()

	if _, err := io.Copy(hash, file); err != nil {
		return nil, fmt.Errorf("failed to calculate MD5 hash for file %s: %w", filePath, err)
	}
	return hash.Sum(nil), nil
}

func CompareTwoFiles(file1, file2 string) (bool, error) {
	hash1, err := hashFileMD5(file1)
	if err != nil {
		return false, err
	}
	hash2, err := hashFileMD5(file2)
	if err != nil {
		return false, err
	}
	return bytes.Equal(hash1, hash2), nil
}
