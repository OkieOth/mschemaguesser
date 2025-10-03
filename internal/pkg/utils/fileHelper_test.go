package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_GetFilesInDir(t *testing.T) {
	baseTempDir := "../../../temp"
	tempDir, err := os.MkdirTemp(baseTempDir, "mschemag-*")
	require.Nil(t, err)
	defer os.RemoveAll(tempDir)
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
	err = os.Mkdir(filepath.Join(tempDir, "Test_GetFilesInDir_3.txt"), 0755)
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

func Test_PrepareDirStructureGetFilesInDir(t *testing.T) {
	testDir := "../../../temp"
	defer RemoveDirectory(filepath.Join(testDir, "testDb"))

	path := filepath.Join(testDir, "testDb")
	exist, err := DirExists(path)
	if err != nil {
		t.Errorf("Error while checking patch (%s): %v", path, err)
		return
	}
	if exist {
		t.Errorf("Directory already exist: %s", path)
		return
	}

	path2 := filepath.Join(testDir, filepath.Join("testDb", "coll1"))
	exist, err = DirExists(path2)

	if err != nil {
		t.Errorf("Error while checking path (%s): %v", path2, err)
		return
	}
	if exist {
		t.Errorf("Directory already exist: %s", path2)
		return
	}

	p1, err := PrepareDirStructure(testDir, "testDb", "coll1")

	if err != nil {
		t.Errorf("Error while prepare dir structure (%s): %v", p1, err)
		return
	}
	if p1 != path2 {
		t.Errorf("Created path is different than the expected: got: %s, expected: %s", p1, path2)
		return
	}

	exist, err = DirExists(path)
	if err != nil {
		t.Errorf("Error while checking patch-2 (%s): %v", path, err)
		return
	}
	if !exist {
		t.Errorf("Directory doesn't exist (2): %s", path)
		return
	}

	p2, err := PrepareDirStructure(testDir, "testDb", "coll1")

	if err != nil {
		t.Errorf("Error while prepare dir structure (%s): %v", p1, err)
		return
	}
	if p2 != p1 {
		t.Errorf("Created path is different than the expected: got: %s, expected: %s", p2, p1)
		return
	}
}

func TestGetKeyPersistenceDirName(t *testing.T) {
	outputDir := "/base/output"

	// Define test cases
	tests := []struct {
		dbName   string
		collName string
		expected string
	}{
		{
			dbName:   "myDB",
			collName: "myCollection",
			expected: filepath.Join(outputDir, "myDB", "myCollection"),
		},
		{
			dbName:   "my DB",    // space in dbName
			collName: "coll@123", // special chars in collName
			expected: filepath.Join(outputDir, "my_DB", "coll_123"),
		},
		{
			dbName:   "dbWithÜnicode", // non-ASCII character in dbName
			collName: "normalColl",
			expected: filepath.Join(outputDir, "dbWith_nicode", "normalColl"),
		},
		{
			dbName:   " ", // dbName with only a space
			collName: "collection",
			expected: filepath.Join(outputDir, "_", "collection"),
		},
		{
			dbName:   "name$with#special&chars!", // multiple special characters in dbName
			collName: "coll$chars",
			expected: filepath.Join(outputDir, "name_with_special_chars_", "coll_chars"),
		},
		{
			dbName:   "123中文", // non-ASCII and numeric characters
			collName: "coll_測試",
			expected: filepath.Join(outputDir, "123__", "coll___"),
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.dbName+"_"+test.collName, func(t *testing.T) {
			result := GetKeyPersistenceDirName(outputDir, test.dbName, test.collName)
			if result != test.expected {
				t.Errorf("Expected %s, but got %s", test.expected, result)
			}
		})
	}
}
