package cmd

import (
	"testing"

	"okieoth/schemaguesser/internal/pkg/mongoHelper"
	testhelper "okieoth/schemaguesser/internal/pkg/testHelper"
	"okieoth/schemaguesser/internal/pkg/utils"
)

func Test_bsonForAllDatabases_IT(t *testing.T) {
	outputDir = "../../../temp"
	defer utils.CleanDirectory(outputDir, false)

	if !testhelper.ValidateEmptyDir(outputDir, t) {
		return
	}

	conStr = "mongodb://{MONGO_USER}:{MONGO_PASSWORD}@{MONGO_HOST}:{MONGO_PORT}/admin"
	useDumps = false
	client, err := mongoHelper.Connect(conStr)
	defer mongoHelper.CloseConnection(client)
	if err != nil {
		t.Errorf("Failed to get client: %v", err)
		return
	}

	bsonForAllDatabases(client, false)

	expected := []string{"admin_system_users.bson", "admin_system_version.bson", "config_system_sessions.bson", "dummy_c1.meta",
		"dummy_c2.meta", "local_startup_log.meta", "admin_system_users.meta", "admin_system_version.meta",
		"dummy_c1.bson", "dummy_c2.bson", "local_startup_log.bson"}

	if !testhelper.ValidateExpectedFiles(outputDir, expected, t) {
		return
	}
	// because for some of the files are there no permissions to read ... so the test is melted down
	// to my own collections
	expected2 := []string{"dummy_c1.meta", "dummy_c2.meta", "dummy_c1.bson", "dummy_c2.bson"}
	_, _ = testhelper.CheckFilesNonZero(outputDir, expected2, t)
}

func Test_bsonForAllCollections_IT(t *testing.T) {
	outputDir = "../../../temp"
	defer utils.CleanDirectory(outputDir, false)

	if !testhelper.ValidateEmptyDir(outputDir, t) {
		return
	}

	conStr = "mongodb://{MONGO_USER}:{MONGO_PASSWORD}@{MONGO_HOST}:{MONGO_PORT}/admin"
	useDumps = false
	client, err := mongoHelper.Connect(conStr)
	defer mongoHelper.CloseConnection(client)
	if err != nil {
		t.Errorf("Failed to get client: %v", err)
		return
	}

	bsonForAllCollections(client, "dummy", false)

	expected := []string{"dummy_c1.bson", "dummy_c2.bson", "dummy_c1.meta", "dummy_c2.meta"}

	if !testhelper.ValidateExpectedFiles(outputDir, expected, t) {
		return
	}
	_, _ = testhelper.CheckFilesNonZero(outputDir, expected, t)
}

func Test_bsonForOneCollection_IT(t *testing.T) {
	outputDir = "../../../temp"
	defer utils.CleanDirectory(outputDir, false)

	if !testhelper.ValidateEmptyDir(outputDir, t) {
		return
	}

	conStr = "mongodb://{MONGO_USER}:{MONGO_PASSWORD}@{MONGO_HOST}:{MONGO_PORT}/admin"
	useDumps = false
	client, err := mongoHelper.Connect(conStr)
	defer mongoHelper.CloseConnection(client)
	if err != nil {
		t.Errorf("Failed to get client: %v", err)
		return
	}

	bsonForOneCollection(client, "dummy", "c2", false, false)

	expected := []string{"dummy_c2.bson", "dummy_c2.meta"}

	if !testhelper.ValidateExpectedFiles(outputDir, expected, t) {
		return
	}
	_, _ = testhelper.CheckFilesNonZero(outputDir, expected, t)
}
