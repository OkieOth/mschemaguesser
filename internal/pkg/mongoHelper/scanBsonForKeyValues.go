package mongoHelper

import (
	"errors"
	"fmt"
	"log"
	"okieoth/schemaguesser/internal/pkg/utils"
	"path/filepath"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

func checkIfStringIsUUIDString(value bson.RawValue) (bool, error) {
	if value.Type != bson.TypeString {
		return false, errors.New("value is not of type string")
	}

	str := value.StringValue()

	uuidRegex := regexp.MustCompile(`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$`)

	if uuidRegex.MatchString(str) {
		return true, nil
	}

	return false, nil
}

func checkIfBinaryIsUUID(value bson.RawValue) (bool, error) {
	// Ensure the value is of binary type
	if value.Type != bson.TypeBinary {
		return false, errors.New("value is not of binary type")
	}

	// Extract the binary data and its subtype
	subtype, _, _, ok := bsoncore.ReadBinary(value.Value)

	if !ok {
		return false, errors.New("checkIfBinaryIsUUID: couldn't read binary")
	}

	// Check if the binary subtype is 3 or 4, indicating a UUID
	if subtype == 0x03 || subtype == 0x04 {
		return true, nil
	}
	return false, nil
}

func GetPersistenceFileName(outputDir string, dbName string, collName string, attribName string) string {
	dirName := utils.GetKeyPersistenceDirName(outputDir, dbName, collName)
	sanitizedAttribName := utils.Sanitize(attribName)

	return filepath.Join(dirName, fmt.Sprintf("%s.keyvalues.txt", sanitizedAttribName))
}

func persistStringValue(value bson.RawValue, dbName string, collName string, attribName string, outputDir string) error {
	// construct the related file name
	// stringify the value
	// write to the file
	return nil // TODO
}

func persistBinaryValue(value bson.RawValue, dbName string, collName string, attribName string, outputDir string) error {
	return nil // TODO
}

func persistObjectIdValue(value bson.RawValue, dbName string, collName string, attribName string, outputDir string) error {
	return nil // TODO
}

func handleStringKeyValue(value bson.RawValue, dbName string, collName string, attribName string, outputDir string) error {
	if b, err := checkIfStringIsUUIDString(value); err != nil {
		log.Printf("Error while checking string value (%v) for uuid format: %v", value, err)
		return err
	} else {
		if b {
			err := persistStringValue(value, dbName, collName, attribName, outputDir)
			if err != nil {
				log.Printf("[%s:%s - %s] Error while writing string value (%v): %v", dbName, collName, attribName, value, err)
				return err
			}
		}
	}
	return nil
}

func handleUuidKeyValue(value bson.RawValue, dbName string, collName string, attribName string, outputDir string) error {
	if b, err := checkIfBinaryIsUUID(value); err != nil {
		log.Printf("Error while checking value (%v) for being uuid: %v", value, err)
	} else {
		if b {
			err := persistBinaryValue(value, dbName, collName, attribName, outputDir)
			if err != nil {
				log.Printf("[%s:%s - %s] Error while writing string value (%v): %v", dbName, collName, attribName, value, err)
			}
		}
	}
	return nil
}

func handleTypeArrayKeyValues(value bson.RawValue, dbName string, collName string, attribName string, outputDir string) error {
	arrayRaw := bson.Raw(value.Value)
	elements, err := arrayRaw.Elements()
	if err != nil {
		return err
	}

	var lastType *bsontype.Type
	for _, elem := range elements {
		if (lastType != nil) && (*lastType != elem.Value().Type) {
			return errors.New(fmt.Sprintf("[%s:%s - %s] array type consists of different types, multiple type arrays are not supported", dbName, collName))
		} else {
			if lastType == nil {
				*lastType = elem.Value().Type
			}
		}
		switch elem.Value().Type {
		case bson.TypeString:
			if err := handleStringKeyValue(value, dbName, collName, attribName, outputDir); err != nil {
				log.Printf("[%s:%s - %s] error while persisting string key value: %v", dbName, collName, attribName, err)
			}
		case bson.TypeEmbeddedDocument:
			if err := handleComplexTypeKeyValues(value, dbName, collName, attribName+"_sub", outputDir); err != nil {
				log.Printf("[%s:%s - %s] error while persisting array key values: %v", dbName, collName, attribName, err)
			}
		case bson.TypeArray:
			if err := handleTypeArrayKeyValues(value, dbName, collName, attribName+"_sub", outputDir); err != nil {
				log.Printf("[%s:%s - %s] error while persisting array key values: %v", dbName, collName, attribName, err)
			}
		case bson.TypeBinary:
			if err := handleUuidKeyValue(value, dbName, collName, attribName, outputDir); err != nil {
				log.Printf("[%s:%s - %s] error while persisting binary key value: %v", dbName, collName, attribName, err)
			}
		case bson.TypeObjectID:
			err := persistObjectIdValue(elem.Value(), dbName, collName, attribName, outputDir)
			if err != nil {
				log.Printf("[%s:%s - %s] Error while writing objectId value (%v): %v", dbName, collName, attribName, elem.Value(), err)
			}
		}
	}
	return nil
}

func handleComplexTypeKeyValues(value bson.RawValue, dbName string, collName string, attribName string, outputDir string) error {
	embeddedDoc := bson.Raw(value.Value)
	elements, err := embeddedDoc.Elements()
	if err != nil {
		return fmt.Errorf("[%s:%s - %s] error while parsing complex type: %v", dbName, collName, attribName, err)
	}
	for _, elem := range elements {
		switch elem.Value().Type {
		case bson.TypeString:
			if err := handleStringKeyValue(elem.Value(), dbName, collName, fmt.Sprintf("%s-%s", attribName, elem.Key()), outputDir); err != nil {
				log.Printf("[%s:%s] error while persisting string key value: %v", dbName, collName, err)
			}
		case bson.TypeEmbeddedDocument:
			if err := handleComplexTypeKeyValues(elem.Value(), dbName, collName, fmt.Sprintf("%s-%s", attribName, elem.Key()), outputDir); err != nil {
				log.Printf("[%s:%s] error while persisting array key values: %v", dbName, collName, err)
			}
		case bson.TypeArray:
			if err := handleTypeArrayKeyValues(elem.Value(), dbName, collName, fmt.Sprintf("%s-%s", attribName, elem.Key()), outputDir); err != nil {
				log.Printf("[%s:%s] error while persisting array key values: %v", dbName, collName, err)
			}
		case bson.TypeBinary:
			if err := handleUuidKeyValue(elem.Value(), dbName, collName, fmt.Sprintf("%s-%s", attribName, elem.Key()), outputDir); err != nil {
				log.Printf("[%s:%s] error while persisting uuid key value: %v", dbName, collName, err)
			}
		case bson.TypeObjectID:
			err := persistObjectIdValue(elem.Value(), dbName, collName, fmt.Sprintf("%s-%s", attribName, elem.Key()), outputDir)
			if err != nil {
				log.Printf("[%s:%s] Error while writing objectId value (%v): %v", dbName, collName, elem.Value(), err)
			}
		}
	}
	return nil // TODO
}

func ScanBsonForKeyValues(doc bson.Raw, dbName string, collName string, outputDir string) error {
	elements, err := doc.Elements()
	if err != nil {
		log.Printf("Error while parsing bson elements: %v", err)
		return err
	}
	for _, elem := range elements {
		switch elem.Value().Type {
		case bson.TypeString:
			if err := handleStringKeyValue(elem.Value(), dbName, collName, elem.Key(), outputDir); err != nil {
				log.Printf("[%s:%s] error while persisting string key value: %v", dbName, collName, err)
			}
		case bson.TypeEmbeddedDocument:
			if err := handleComplexTypeKeyValues(elem.Value(), dbName, collName, elem.Key(), outputDir); err != nil {
				log.Printf("[%s:%s] error while persisting array key values: %v", dbName, collName, err)
			}
		case bson.TypeArray:
			if err := handleTypeArrayKeyValues(elem.Value(), dbName, collName, elem.Key(), outputDir); err != nil {
				log.Printf("[%s:%s] error while persisting array key values: %v", dbName, collName, err)
			}
		case bson.TypeBinary:
			if err := handleUuidKeyValue(elem.Value(), dbName, collName, elem.Key(), outputDir); err != nil {
				log.Printf("[%s:%s] error while persisting uuid key value: %v", dbName, collName, err)
			}
		case bson.TypeObjectID:
			err := persistObjectIdValue(elem.Value(), dbName, collName, elem.Key(), outputDir)
			if err != nil {
				log.Printf("[%s:%s] Error while writing objectId value (%v): %v", dbName, collName, elem.Value(), err)
			}
		}
	}
	return nil
}
