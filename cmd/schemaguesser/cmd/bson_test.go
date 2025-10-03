package cmd

import (
	"context"
	"os"
	"testing"

	"okieoth/schemaguesser/internal/pkg/importHelper"
	"okieoth/schemaguesser/internal/pkg/meta"
	"okieoth/schemaguesser/internal/pkg/mongoHelper"
	testhelper "okieoth/schemaguesser/internal/pkg/testHelper"
	"okieoth/schemaguesser/internal/pkg/utils"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
)

const TEMP_BASE = "../../../temp"

func Test_bsonForAllDatabases_IT(t *testing.T) {
	tmpDir, err := os.MkdirTemp(TEMP_BASE, "mschemag-*")
	require.Nil(t, err)
	defer os.RemoveAll(tmpDir)
	outputDir = tmpDir

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

	bsonForAllDatabases(client, false)

	expected := []string{"admin_system_users.bson", "admin_system_version.bson", "config_system_sessions.bson", "dummy_c1.meta",
		"dummy_c2.meta", "local_startup_log.meta", "admin_system_users.meta", "admin_system_version.meta",
		"dummy_c1.bson", "dummy_c2.bson", "dummy_c3.bson", "dummy_c3.meta", "local_startup_log.bson"}

	if !testhelper.ValidateExpectedFiles(outputDir, expected, t) {
		return
	}
	// because for some of the files are there no permissions to read ... so the test is melted down
	// to my own collections
	expected2 := []string{"dummy_c1.meta", "dummy_c2.meta", "dummy_c1.bson", "dummy_c2.bson"}
	_, _ = testhelper.CheckFilesNonZero(outputDir, expected2, t)

	validateViaReimport(t, outputDir, "dummy_c1.bson", client)
}

func validateViaReimport(t *testing.T, outputDir, desiredFileName string, client *mongo.Client) {
	metaInfos, err := meta.GetAllMetaInfos(outputDir)
	require.Nil(t, err, "couldn't read meta infos")
	var neededMetaInfo *meta.MetaInfo
	for _, m := range metaInfos {
		if m.FileName == desiredFileName {
			neededMetaInfo = &m
			break
		}
	}
	require.NotNil(t, neededMetaInfo, "couldn't find desired mata file")
	newColName := neededMetaInfo.Collection + "_test"
	importFile := utils.GetFileName(outputDir, "bson", neededMetaInfo.Db, neededMetaInfo.Collection)
	ctx := context.Background()
	itemCount, err := importHelper.ImportData(client, importFile, neededMetaInfo.Db, newColName, 100, &ctx)
	require.Nil(t, err, "error while re-import exported bson")
	defer func() {
		// delete new collection
		err = client.Database(neededMetaInfo.Db).Collection(newColName).Drop(ctx)
	}()
	require.Equal(t, neededMetaInfo.ItemCount, itemCount, "wrong number of data sets re-imported")
}

func Test_bsonForAllCollections_IT(t *testing.T) {
	tmpDir, err := os.MkdirTemp(TEMP_BASE, "mschemag-*")
	require.Nil(t, err)
	defer os.RemoveAll(tmpDir)
	outputDir = tmpDir

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

	bsonForAllCollections(client, "dummy", false)

	expected := []string{"dummy_c1.bson", "dummy_c2.bson", "dummy_c3.bson", "dummy_c1.meta", "dummy_c2.meta", "dummy_c3.meta"}

	if !testhelper.ValidateExpectedFiles(outputDir, expected, t) {
		return
	}
	_, _ = testhelper.CheckFilesNonZero(outputDir, expected, t)
}

func Test_bsonForOneCollection_IT(t *testing.T) {
	tmpDir, err := os.MkdirTemp(TEMP_BASE, "mschemag-*")
	require.Nil(t, err)
	defer os.RemoveAll(tmpDir)
	outputDir = tmpDir

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

	bsonForOneCollection(client, "dummy", "c2", false, false)

	expected := []string{"dummy_c2.bson", "dummy_c2.meta"}

	if !testhelper.ValidateExpectedFiles(outputDir, expected, t) {
		return
	}
	_, _ = testhelper.CheckFilesNonZero(outputDir, expected, t)
}
