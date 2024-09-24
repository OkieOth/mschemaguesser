package utils

import (
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson"
)

func DumpCollectionData(b bson.Raw, dataDumpFile *os.File) error {
	return dumpToFile([]byte(b), dataDumpFile)
	// var bsonMap map[string]interface{}
	// err := bson.Unmarshal(b, &bsonMap)
	// if err != nil {
	// 	log.Fatal("Error unmarshalling BSON:", err)
	// 	return err
	// }

	// jsonData, err := json.MarshalIndent(bsonMap, "", "  ")
	// if err != nil {
	// 	log.Fatal("Error marshalling JSON:", err)
	// }

	// return dumpToFile(jsonData, dataDumpFile)
}

func dumpToFile(b []byte, dumpFile *os.File) error {
	_, err := dumpFile.Write(b)
	if err != nil {
		return fmt.Errorf("error writing to file: %w\n", err)
	} else {
		return nil
	}
}
