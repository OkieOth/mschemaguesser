package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"okieoth/schemaguesser/internal/pkg/mongoHelper"
	testhelper "okieoth/schemaguesser/internal/pkg/testHelper"
	"okieoth/schemaguesser/internal/pkg/utils"
)

func Test_jsonForAllDatabases_IT(t *testing.T) {
	outputDir = "../../../temp"
	defer utils.CleanDirectory(outputDir, false)

	if !testhelper.ValidateEmptyDir(outputDir, t) {
		return
	}

	conStr := "mongodb://{MONGO_USER}:{MONGO_PASSWORD}@{MONGO_HOST}:{MONGO_PORT}/admin"
	useDumps = false
	client, err := mongoHelper.Connect(conStr)
	defer mongoHelper.CloseConnection(client)
	if err != nil {
		t.Errorf("Failed to get client: %v", err)
		return
	}

	jsonForAllDatabases(client, false)

	expected := []string{"admin_system_users.json", "admin_system_version.json", "config_system_sessions.json",
		"dummy_c1.json", "dummy_c2.json", "local_startup_log.json"}

	if !testhelper.ValidateExpectedFiles(outputDir, expected, t) {
		return
	}
	// because for some of the files are there no permissions to read ... so the test is melted down
	// to my own collections
	expected2 := []string{"dummy_c1.json", "dummy_c2.json"}
	_, _ = testhelper.CheckFilesNonZero(outputDir, expected2, t)
}

func Test_jsonForAllCollections_IT(t *testing.T) {
	outputDir = "../../../temp"
	defer utils.CleanDirectory(outputDir, false)

	if !testhelper.ValidateEmptyDir(outputDir, t) {
		return
	}

	conStr := "mongodb://{MONGO_USER}:{MONGO_PASSWORD}@{MONGO_HOST}:{MONGO_PORT}/admin"
	useDumps = false
	client, err := mongoHelper.Connect(conStr)
	defer mongoHelper.CloseConnection(client)
	if err != nil {
		t.Errorf("Failed to get client: %v", err)
		return
	}

	jsonForAllCollections(client, "dummy", false)

	expected := []string{"dummy_c1.json", "dummy_c2.json"}

	if !testhelper.ValidateExpectedFiles(outputDir, expected, t) {
		return
	}
	_, _ = testhelper.CheckFilesNonZero(outputDir, expected, t)
}

func Test_jsonForOneCollection_IT(t *testing.T) {
	outputDir = "../../../temp"
	defer utils.CleanDirectory(outputDir, false)

	if !testhelper.ValidateEmptyDir(outputDir, t) {
		return
	}

	_ = queryOneFileFromDb(t)
}

func queryOneFileFromDb(t *testing.T) bool {
	conStr := "mongodb://{MONGO_USER}:{MONGO_PASSWORD}@{MONGO_HOST}:{MONGO_PORT}/admin"
	useDumps = false
	client, err := mongoHelper.Connect(conStr)
	defer mongoHelper.CloseConnection(client)
	if err != nil {
		t.Errorf("Failed to get client: %v", err)
		return false
	}

	jsonForOneCollection(client, "dummy", "c2", false, false)

	expected := []string{"dummy_c2.json"}

	if !testhelper.ValidateExpectedFiles(outputDir, expected, t) {
		return false
	}
	_, err = testhelper.CheckFilesNonZero(outputDir, expected, t)
	return err == nil
}

func Test_jsonForAllDatabases(t *testing.T) {
	outputDir = "../../../temp"
	defer utils.CleanDirectory(outputDir, false)

	if !testhelper.ValidateEmptyDir(outputDir, t) {
		return
	}

	useDumps = true
	dumpDir = "../../../resources/bson"

	jsonForAllDatabases(nil, false)

	expected := []string{"dummy_c1.json", "dummy_c2.json"}

	if !testhelper.ValidateExpectedFiles(outputDir, expected, t) {
		return
	}
	_, _ = testhelper.CheckFilesNonZero(outputDir, expected, t)
}

func Test_jsonForAllCollections(t *testing.T) {
	outputDir = "../../../temp"
	defer utils.CleanDirectory(outputDir, false)

	if !testhelper.ValidateEmptyDir(outputDir, t) {
		return
	}

	useDumps = true
	dumpDir = "../../../resources/bson"

	jsonForAllCollections(nil, "dummy", false)

	expected := []string{"dummy_c1.json", "dummy_c2.json"}

	if !testhelper.ValidateExpectedFiles(outputDir, expected, t) {
		return
	}
	_, _ = testhelper.CheckFilesNonZero(outputDir, expected, t)
}

func Test_jsonForOneCollection(t *testing.T) {
	outputDir = "../../../temp"
	defer utils.CleanDirectory(outputDir, false)

	if !testhelper.ValidateEmptyDir(outputDir, t) {
		return
	}

	if !queryOneFileFromDb(t) {
		return
	}
	newFileName := "dummy_c2_from_db.json"
	pathNewFileName := filepath.Join(outputDir, newFileName)
	origFileName := "dummy_c2.json"
	pathOrigFileName := filepath.Join(outputDir, origFileName)
	err := os.Rename(pathOrigFileName, pathNewFileName)
	if err != nil {
		t.Errorf("Error renaming test file (%s -> %s): %v", origFileName, newFileName, err)
		return
	}

	useDumps = true
	dumpDir = "../../../resources/bson"

	jsonForOneCollection(nil, "dummy", "c2", false, false)

	expected := []string{origFileName, newFileName}

	if !testhelper.ValidateExpectedFiles(outputDir, expected, t) {
		return
	}
	_, _ = testhelper.CheckFilesNonZero(outputDir, expected, t)

	testhelper.CompareTwoFiles(pathNewFileName, pathOrigFileName)
}
