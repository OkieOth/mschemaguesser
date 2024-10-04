package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func Test_GetFilesInDir(t *testing.T) {
	tempDir := "../../../temp"
	defer CleanDirectory(tempDir, false)

	validateEmptyDir := func() bool {
		files, err := GetFilesInDir(tempDir, false)

		if err != nil {
			t.Errorf("Error while getting file via GetFilesInDir: %v", err)
			return false
		}

		if len(files) > 0 {
			t.Errorf("Retrieved wrong number of files. Expected 0, got %d: %v", len(files), files)
			return false
		}
		return true
	}

	validateEmptyDir()

	for i := 0; i < 3; i++ {
		fName := filepath.Join(tempDir, fmt.Sprintf("Test_GetFilesInDir_%d.txt", i))
		_, err := os.Create(fName)
		if err != nil {
			t.Errorf("Error while creating file (%s): %v", fName, err)
			return
		}
	}

	validate3Files := func() bool {
		files, err := GetFilesInDir(tempDir, false)
		if err != nil {
			t.Errorf("Error (2) while getting file via GetFilesInDir: %v", err)
			return false
		}
		expected := []string{"Test_GetFilesInDir_0.txt", "Test_GetFilesInDir_1.txt", "Test_GetFilesInDir_2.txt"}

		lenExpected := len(expected)

		if len(files) != lenExpected {
			t.Errorf("Retrieved wrong number of files. Expected %d, got %d, %v", lenExpected, len(files), files)
			return false
		}
		ret := true
		for _, f := range expected {
			if !slices.Contains(files, f) {
				ret = false
				t.Errorf("Expected file %s not found", f)
			}
		}
		return ret
	}

	if !validate3Files() { // test that we got 3 files
		return
	}
	err := os.Mkdir(filepath.Join(tempDir, "Test_GetFilesInDir_3.txt"), 0755)
	if err != nil {
		t.Errorf("Error (3) while creating test dir: %v", err)
		return
	}

	if !validate3Files() { // test that we still got 3 files
		return
	}

	CleanDirectory(tempDir, false)
	validateEmptyDir()

}
