package cmd

import (
	"slices"
	"testing"

	"okieoth/schemaguesser/internal/pkg/mongoHelper"
)

func Test_getAllDatabasesOrPanic_1(t *testing.T) {
	useDumps := true
	dumpDir := "../../../resources/meta"

	expectedDatabases := []string{"tms-sub-system", "tms-process-data", "tms-strategy"} // Adjust according to the expected data

	databases := getAllDatabasesOrPanic(nil, dumpDir, useDumps)

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

func Test_getAllDatabasesOrPanic_2_IT(t *testing.T) {
	conStr := "mongodb://{MONGO_USER}:{MONGO_PASSWORD}@{MONGO_HOST}:{MONGO_PORT}/admin"
	useDumps := false
	client, err := mongoHelper.Connect(conStr)
	defer mongoHelper.CloseConnection(client)
	if err != nil {
		t.Errorf("Failed to get client: %v", err)
		return
	}

	databases := getAllDatabasesOrPanic(client, "", useDumps)
	expectedDbs := []string{"admin", "config", "dummy", "local"}

	l1 := len(databases)
	l2 := len(expectedDbs)

	if l1 != l2 {
		t.Errorf("Expected array has a different length: len(result): %d, len(expected): %d, received: %v", l1, l2, databases)
	}

	for _, db := range expectedDbs {
		if !slices.Contains(databases, db) {
			t.Errorf("Expected database %s, but it was not found", db)
		}
	}
}

func Test_getAllCollectionsOrPanic_1(t *testing.T) {
	useDumps := true
	dumpDir := "../../../resources/meta"
	dbName := "tms-process-data"

	expectedCollections := []string{"PublicTransport", "LastDeviceDetectorData", "DetectorsConfigurations", "DataMigration"} // Adjust according to test data

	collections := getAllCollectionsOrPanic(nil, dumpDir, useDumps, dbName)

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

func Test_getAllCollectionsOrPanic_2_IT(t *testing.T) {
	conStr := "mongodb://{MONGO_USER}:{MONGO_PASSWORD}@{MONGO_HOST}:{MONGO_PORT}/admin"
	useDumps := false
	client, err := mongoHelper.Connect(conStr)
	defer mongoHelper.CloseConnection(client)

	if err != nil {
		t.Errorf("Failed to get client: %v", err)
		return
	}

	collections := getAllCollectionsOrPanic(client, "", useDumps, "dummy")
	expected := []string{"c1", "c2", "c3"}

	l1 := len(collections)
	l2 := len(expected)

	if l1 != l2 {
		t.Errorf("Expected array has a different length: len(result): %d, len(expected): %d, received: %v", l1, l2, collections)
	}

	for _, db := range expected {
		if !slices.Contains(collections, db) {
			t.Errorf("Expected database %s, but it was not found", db)
		}
	}
}
