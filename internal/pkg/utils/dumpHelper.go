package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
)

func sanitize(input string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9-]`)
	return re.ReplaceAllString(input, "_")
}

func GetFileName(dir string, fileExt string, dbName string, colName string) string {
	safeDbName := sanitize(dbName)
	safeColName := sanitize(colName)
	return filepath.Join(dir, fmt.Sprintf("%s_%s.%s", safeDbName, safeColName, fileExt))
}

func CreateOutputFile(outputDir string, fileExt string, dbName string, colName string) (*os.File, error) {
	filePath := GetFileName(outputDir, fileExt, dbName, colName)
	return os.Create(filePath)
}

func DumpBsonCollectionData(b bson.Raw, dataDumpFile *os.File) error {
	return DumpBytesToFile([]byte(b), dataDumpFile)
}

func DumpBytesToFile(b []byte, dumpFile *os.File) error {
	_, err := dumpFile.Write(b)
	if err != nil {
		return fmt.Errorf("error writing to file: %w\n", err)
	} else {
		return nil
	}
}
