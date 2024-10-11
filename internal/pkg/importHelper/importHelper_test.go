package importHelper

import (
	"slices"
	"testing"
)

// Test for AllDatabases function
func TestAllDatabases(t *testing.T) {
	// Define the input directory containing the .meta files
	inputDir := "../../../resources/meta"

	// Call the function
	databases, err := AllDatabases(inputDir)

	// Check if an error occurred
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Define the expected databases (as present in the .meta files)
	expectedDatabases := []string{"tms-sub-system", "tms-process-data", "tms-strategy"} // Adjust according to the expected data

	l1 := len(databases)
	l2 := len(expectedDatabases)

	if l1 != l2 {
		t.Errorf("Expected array has a different length: len(result): %d, len(expected): %d", l1, l2)
	}

	// Check if the returned databases match the expected ones
	for _, db := range expectedDatabases {
		if !slices.Contains(databases, db) {
			t.Errorf("Expected database %s, but it was not found", db)
		}
	}
}

// Test for AllCollectionsForDb function
func TestAllCollectionsForDb1(t *testing.T) {
	inputDir := "../../../resources/meta"
	dbName := "tms-strategy"

	collections, err := AllCollectionsForDb(inputDir, dbName)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expectedCollections := []string{"strategy"}

	l1 := len(collections)
	l2 := len(expectedCollections)

	if l1 != l2 {
		t.Errorf("Expected array has a different length: len(result): %d, len(expected): %d", l1, l2)
	}

	for _, coll := range expectedCollections {
		if !slices.Contains(collections, coll) {
			t.Errorf("Expected collection %s for database %s, but it was not found", coll, dbName)
		}
	}
}

func TestAllCollectionsForDb2(t *testing.T) {
	inputDir := "../../../resources/meta"
	dbName := "tms-process-data"

	collections, err := AllCollectionsForDb(inputDir, dbName)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expectedCollections := []string{"PublicTransport", "LastDeviceDetectorData", "DetectorsConfigurations", "DataMigration"} // Adjust according to test data

	l1 := len(collections)
	l2 := len(expectedCollections)

	if l1 != l2 {
		t.Errorf("Expected array has a different length: len(result): %d, len(expected): %d", l1, l2)
	}

	for _, coll := range expectedCollections {
		if !slices.Contains(collections, coll) {
			t.Errorf("Expected collection %s for database %s, but it was not found", coll, dbName)
		}
	}
}
