package utils_test

import (
	"os"
	"testing"

	"okieoth/schemaguesser/internal/pkg/utils"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestSanitize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello_world", "hello_world"},
		{"hello world!", "hello_world_"},
		{"123-ABC_xyz!", "123-ABC_xyz_"},
		{"@special*chars#", "_special_chars_"},
	}

	for _, test := range tests {
		result := utils.Sanitize(test.input)
		assert.Equal(t, test.expected, result, "Expected %s but got %s", test.expected, result)
	}
}

func TestGetFileName(t *testing.T) {
	tests := []struct {
		dir      string
		ext      string
		dbName   string
		colName  string
		expected string
	}{
		{"/tmp", "json", "myDB", "myCollection", "/tmp/myDB_myCollection.json"},
		{"/home", "txt", "db-test", "collection_test", "/home/db-test_collection_test.txt"},
	}

	for _, test := range tests {
		result := utils.GetFileName(test.dir, test.ext, test.dbName, test.colName)
		assert.Equal(t, test.expected, result, "Expected %s but got %s", test.expected, result)
	}
}

func TestCreateOutputFile(t *testing.T) {
	outputDir := "../../../temp"
	fileExt := "txt"
	dbName := "testDB"
	colName := "testCollection"

	file, err := utils.CreateOutputFile(outputDir, fileExt, dbName, colName)
	assert.NoError(t, err)
	defer os.Remove(file.Name()) // Cleanup
	defer file.Close()

	expectedFilePath := "../../../temp/testDB_testCollection.txt"
	assert.Equal(t, expectedFilePath, file.Name(), "Expected file path %s but got %s", expectedFilePath, file.Name())
}

func TestDumpBsonCollectionData(t *testing.T) {
	// Creating a temporary file
	dumpFile, err := os.CreateTemp("../../../temp", "test_bson_dump_*.bson")
	assert.NoError(t, err)
	defer os.Remove(dumpFile.Name()) // Cleanup
	defer dumpFile.Close()

	// Example BSON data
	doc := bson.D{{"name", "test"}, {"age", 30}}
	raw, err := bson.Marshal(doc)
	assert.NoError(t, err)

	// Call the function and check for errors
	err = utils.DumpBsonCollectionData(raw, dumpFile)
	assert.NoError(t, err)

	// Read the contents of the file to verify
	data, err := os.ReadFile(dumpFile.Name())
	assert.NoError(t, err)
	assert.Equal(t, raw, data, "Dumped data does not match BSON content")
}

func TestDumpBytesToFile(t *testing.T) {
	// Create a temporary file
	dumpFile, err := os.CreateTemp("../../../temp", "test_dump_bytes_*.txt")
	assert.NoError(t, err)
	defer os.Remove(dumpFile.Name()) // Cleanup
	defer dumpFile.Close()

	// Sample data to dump
	data := []byte("Hello, World!")

	// Call the function and check for errors
	err = utils.DumpBytesToFile(data, dumpFile)
	assert.NoError(t, err)

	// Read the contents of the file to verify
	fileData, err := os.ReadFile(dumpFile.Name())
	assert.NoError(t, err)
	assert.Equal(t, data, fileData, "Expected file data to be %s but got %s", data, fileData)
}
