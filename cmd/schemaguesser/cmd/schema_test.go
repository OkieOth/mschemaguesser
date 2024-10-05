package cmd

import (
	"testing"

	"okieoth/schemaguesser/internal/pkg/mongoHelper"
	testhelper "okieoth/schemaguesser/internal/pkg/testHelper"
	"okieoth/schemaguesser/internal/pkg/utils"
)

func Test_getDocumentCount_IT(t *testing.T) {
	conStr := "mongodb://{MONGO_USER}:{MONGO_PASSWORD}@{MONGO_HOST}:{MONGO_PORT}/admin"
	client, err := mongoHelper.Connect(conStr)
	defer mongoHelper.CloseConnection(client)

	if err != nil {
		t.Errorf("Failed to get client: %v", err)
		return
	}

	var mt mongoHelper.ComplexType
	getDocumentCount(client, "dummy", "c1", &mt)

	if !mt.Count.IsSet {
		t.Error("IsSet not true")
		return
	}

	if mt.Count.Value != 4 {
		t.Errorf("Value != 4, got: %d", mt.Count.Value)
	}

}

func Test_replaceUuidValues(t *testing.T) {
	jsonStr := `{"Category":7,"CharacterDetails":null,"DefaultOffProgramId":{"Subtype":4,"Data":"AAAAAAAAAAAAAAAAAAAAAA=="},"Description":"Vehicle Activated Sign","FullMatrixDetails":null,"Name":"Vehicle Activated Sign","PrismDetails":null,"TenantId":{"Subtype":4,"Data":"BWvPWOF+QrqBhvJf+96LNQ=="},"TenantType":0,"_id":{"Subtype":4,"Data":"AN0UjlWcSeKLPZa8F59Xog=="}}`

	convertedStr, err := replaceUuidValues(4, jsonStr)

	if err != nil {
		t.Errorf("Fail to replace uuids in replaceUuidValues: %v", err)
		return
	}

	expected := `{"Category":7,"CharacterDetails":null,"DefaultOffProgramId":"00000000-0000-0000-0000-000000000000","Description":"Vehicle Activated Sign","FullMatrixDetails":null,"Name":"Vehicle Activated Sign","PrismDetails":null,"TenantId":"056bcf58-e17e-42ba-8186-f25ffbde8b35","TenantType":0,"_id":"00dd148e-559c-49e2-8b3d-96bc179f57a2"}`
	if convertedStr != expected {
		t.Errorf("Got wrong jsonString: expected: %s\ngot: %s", expected, convertedStr)
		return
	}
}

// func Test_printSchemaForAllDatabases_IT(t *testing.T) {
// 	outputDir = "../../../temp"
// 	defer utils.CleanDirectory(outputDir, false)

// 	if !testhelper.ValidateEmptyDir(outputDir, t) {
// 		return
// 	}

// 	conStr := "mongodb://{MONGO_USER}:{MONGO_PASSWORD}@{MONGO_HOST}:{MONGO_PORT}/admin"
// 	useDumps = false
// 	client, err := mongoHelper.Connect(conStr)
// 	defer mongoHelper.CloseConnection(client)
// 	if err != nil {
// 		t.Errorf("Failed to get client: %v", err)
// 		return
// 	}

// 	printSchemasForAllDatabases(client, false)

// 	expected := []string{"dummy_c1.schema.json", "dummy_c2.schema.json"}

// 	if !testhelper.ValidateExpectedFiles(outputDir, expected, t) {
// 		return
// 	}
// 	_, _ = testhelper.CheckFilesNonZero(outputDir, expected, t)
// }

func Test_printSchemaForOneCollection_IT(t *testing.T) {
	outputDir = "../../../temp"
	defer utils.CleanDirectory(outputDir, false)

	if !testhelper.ValidateEmptyDir(outputDir, t) {
		return
	}

	_ = schemaForOneFileFromDb(t)
}

func schemaForOneFileFromDb(t *testing.T) bool {
	conStr := "mongodb://{MONGO_USER}:{MONGO_PASSWORD}@{MONGO_HOST}:{MONGO_PORT}/admin"
	useDumps = false
	client, err := mongoHelper.Connect(conStr)
	defer mongoHelper.CloseConnection(client)
	if err != nil {
		t.Errorf("Failed to get client: %v", err)
		return false
	}

	printSchemaForOneCollection(client, "dummy", "c2", false, false)

	expected := []string{"dummy_c2.schema.json"}

	if !testhelper.ValidateExpectedFiles(outputDir, expected, t) {
		return false
	}
	_, err = testhelper.CheckFilesNonZero(outputDir, expected, t)
	return err == nil
}

func Test_printSchemaForAllCollections_IT(t *testing.T) {
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

	printSchemasForAllCollections(client, "dummy", false)

	expected := []string{"dummy_c1.schema.json", "dummy_c2.schema.json"}

	if !testhelper.ValidateExpectedFiles(outputDir, expected, t) {
		return
	}
	_, _ = testhelper.CheckFilesNonZero(outputDir, expected, t)
}
