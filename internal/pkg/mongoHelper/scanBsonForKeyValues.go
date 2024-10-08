package mongoHelper

import (
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

func ScanBsonForKeyValues(doc bson.Raw, dbName string, collName string, outputDir string) error {
	elements, err := doc.Elements()
	if err != nil {
		log.Printf("Error while parsing bson elements: %v", err)
		return err
	}
	for _, elem := range elements {
		switch elem.Value().Type {
		case bson.TypeString:
			// TODO
		case bson.TypeEmbeddedDocument:
			// TODO
		case bson.TypeArray:
			// TODO
		case bson.TypeBinary:
			// TODO
		case bson.TypeObjectID:
			// TODO
		}
	}
	return nil
}
