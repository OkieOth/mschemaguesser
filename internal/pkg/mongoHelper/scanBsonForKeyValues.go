package mongoHelper

import (
	"errors"
	"fmt"
	"log"
	"okieoth/schemaguesser/internal/pkg/utils"
	"os"
	"path/filepath"
	"regexp"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

var KeepNullUuids bool

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

func persistStringValue(value bson.RawValue, dbName string, collName string, attribName string, outputFile *os.File) error {
	strValue := value.StringValue()
	return persistString(strValue, dbName, collName, attribName, outputFile)
}

func persistString(strValue string, dbName string, collName string, attribName string, outputFile *os.File) error {
	if KeepNullUuids && (strValue == "00000000-0000-0000-0000-000000000000") {
		// zero uuids are ignored
		return nil
	}
	sanitizedAttribName := utils.Sanitize(attribName)
	_, err := outputFile.Write([]byte(sanitizedAttribName))
	if err != nil {
		return fmt.Errorf("[%s:%s - %s] Error while writing attrib name: %v", dbName, collName, attribName, err)
	}
	_, err = outputFile.Write([]byte(": "))
	if err != nil {
		return fmt.Errorf("[%s:%s - %s] Error while writing separator: %v", dbName, collName, attribName, err)
	}
	_, err = outputFile.Write([]byte(strValue))
	if err != nil {
		return fmt.Errorf("[%s:%s - %s] Error while writing value: %v", dbName, collName, attribName, err)
	}
	_, err = outputFile.Write([]byte("\n"))
	if err != nil {
		return fmt.Errorf("[%s:%s - %s] Error while writing last new line: %v", dbName, collName, attribName, err)
	}
	return nil
}

func persistBinaryValue(value bson.RawValue, dbName string, collName string, attribName string, outputFile *os.File) error {
	subtype, binary := value.Binary()

	if subtype == 4 || subtype == 3 {
		uuidValue, err := uuid.FromBytes(binary)
		if err != nil {
			return fmt.Errorf("[%s:%s - %s] Error converting binary to UUID: %v", dbName, collName, attribName, err)
		}
		err = persistString(uuidValue.String(), dbName, collName, attribName, outputFile)
		if err != nil {
			log.Printf("[%s:%s - %s] Error while writing string value (%v): %v", dbName, collName, attribName, value, err)
			return err
		}
	}
	return nil
}

func persistObjectIdValue(value bson.RawValue, dbName string, collName string, attribName string, outputFile *os.File) error {
	return persistString(value.String(), dbName, collName, attribName, outputFile)
}

func handleStringKeyValue(value bson.RawValue, dbName string, collName string, attribName string, outputFile *os.File) error {
	if b, err := checkIfStringIsUUIDString(value); err != nil {
		log.Printf("Error while checking string value (%v) for uuid format: %v", value, err)
		return err
	} else {
		if b {
			err := persistStringValue(value, dbName, collName, attribName, outputFile)
			if err != nil {
				log.Printf("[%s:%s - %s] Error while writing string value (%v): %v", dbName, collName, attribName, value, err)
				return err
			}
		}
	}
	return nil
}

func handleUuidKeyValue(value bson.RawValue, dbName string, collName string, attribName string, outputFile *os.File) error {
	if b, err := checkIfBinaryIsUUID(value); err != nil {
		log.Printf("Error while checking value (%v) for being uuid: %v", value, err)
	} else {
		if b {
			err := persistBinaryValue(value, dbName, collName, attribName, outputFile)
			if err != nil {
				log.Printf("[%s:%s - %s] Error while writing string value (%v): %v", dbName, collName, attribName, value, err)
			}
		}
	}
	return nil
}

func handleTypeArrayKeyValues(value bson.RawValue, dbName string, collName string, attribName string, outputFile *os.File) error {
	arrayRaw := bson.Raw(value.Value)
	elements, err := arrayRaw.Elements()
	if err != nil {
		return err
	}

	var lastType bsontype.Type
	lastTypeSet := false
	for _, elem := range elements {
		if (lastTypeSet) && (lastType != elem.Value().Type) {
			return errors.New(fmt.Sprintf("[%s:%s - %s] array type consists of different types, multiple type arrays are not supported", dbName, collName, attribName))
		} else {
			if !lastTypeSet {
				lastType = elem.Value().Type
			}
		}
		lastType = elem.Value().Type
		switch elem.Value().Type {
		case bson.TypeString:
			if err := handleStringKeyValue(elem.Value(), dbName, collName, attribName, outputFile); err != nil {
				log.Printf("[%s:%s - %s] error while persisting string key value: %v", dbName, collName, attribName, err)
			}
		case bson.TypeEmbeddedDocument:
			if err := handleComplexTypeKeyValues(elem.Value(), dbName, collName, attribName+"_sub", outputFile); err != nil {
				log.Printf("[%s:%s - %s] error while persisting array key values: %v", dbName, collName, attribName, err)
			}
		case bson.TypeArray:
			if err := handleTypeArrayKeyValues(elem.Value(), dbName, collName, attribName+"_sub", outputFile); err != nil {
				log.Printf("[%s:%s - %s] error while persisting array key values: %v", dbName, collName, attribName, err)
			}
		case bson.TypeBinary:
			if err := handleUuidKeyValue(elem.Value(), dbName, collName, attribName, outputFile); err != nil {
				log.Printf("[%s:%s - %s] error while persisting binary key value: %v", dbName, collName, attribName, err)
			}
		case bson.TypeObjectID:
			err := persistObjectIdValue(elem.Value(), dbName, collName, attribName, outputFile)
			if err != nil {
				log.Printf("[%s:%s - %s] Error while writing objectId value (%v): %v", dbName, collName, attribName, elem.Value(), err)
			}
		}
	}
	return nil
}

func handleComplexTypeKeyValues(value bson.RawValue, dbName string, collName string, attribName string, outputFile *os.File) error {
	embeddedDoc := bson.Raw(value.Value)
	elements, err := embeddedDoc.Elements()
	if err != nil {
		return fmt.Errorf("[%s:%s - %s] error while parsing complex type: %v", dbName, collName, attribName, err)
	}
	for _, elem := range elements {
		switch elem.Value().Type {
		case bson.TypeString:
			if err := handleStringKeyValue(elem.Value(), dbName, collName, fmt.Sprintf("%s-%s", attribName, elem.Key()), outputFile); err != nil {
				log.Printf("[%s:%s] error while persisting string key value: %v", dbName, collName, err)
			}
		case bson.TypeEmbeddedDocument:
			if err := handleComplexTypeKeyValues(elem.Value(), dbName, collName, fmt.Sprintf("%s-%s", attribName, elem.Key()), outputFile); err != nil {
				log.Printf("[%s:%s] error while persisting array key values: %v", dbName, collName, err)
			}
		case bson.TypeArray:
			if err := handleTypeArrayKeyValues(elem.Value(), dbName, collName, fmt.Sprintf("%s-%s", attribName, elem.Key()), outputFile); err != nil {
				log.Printf("[%s:%s] error while persisting array key values: %v", dbName, collName, err)
			}
		case bson.TypeBinary:
			if err := handleUuidKeyValue(elem.Value(), dbName, collName, fmt.Sprintf("%s-%s", attribName, elem.Key()), outputFile); err != nil {
				log.Printf("[%s:%s] error while persisting uuid key value: %v", dbName, collName, err)
			}
		case bson.TypeObjectID:
			err := persistObjectIdValue(elem.Value(), dbName, collName, fmt.Sprintf("%s-%s", attribName, elem.Key()), outputFile)
			if err != nil {
				log.Printf("[%s:%s] Error while writing objectId value (%v): %v", dbName, collName, elem.Value(), err)
			}
		}
	}
	return nil
}

func ScanBsonForKeyValues(doc bson.Raw, dbName string, collName string, outputFile *os.File) error {

	elements, err := doc.Elements()
	if err != nil {
		log.Printf("Error while parsing bson elements: %v", err)
		return err
	}

	for _, elem := range elements {
		switch elem.Value().Type {
		case bson.TypeString:
			if err := handleStringKeyValue(elem.Value(), dbName, collName, elem.Key(), outputFile); err != nil {
				log.Printf("[%s:%s] error while persisting string key value: %v", dbName, collName, err)
			}
		case bson.TypeEmbeddedDocument:
			if err := handleComplexTypeKeyValues(elem.Value(), dbName, collName, elem.Key(), outputFile); err != nil {
				log.Printf("[%s:%s] error while persisting array key values: %v", dbName, collName, err)
			}
		case bson.TypeArray:
			if err := handleTypeArrayKeyValues(elem.Value(), dbName, collName, elem.Key(), outputFile); err != nil {
				log.Printf("[%s:%s] error while persisting array key values: %v", dbName, collName, err)
			}
		case bson.TypeBinary:
			if err := handleUuidKeyValue(elem.Value(), dbName, collName, elem.Key(), outputFile); err != nil {
				log.Printf("[%s:%s] error while persisting uuid key value: %v", dbName, collName, err)
			}
		case bson.TypeObjectID:
			err := persistObjectIdValue(elem.Value(), dbName, collName, elem.Key(), outputFile)
			if err != nil {
				log.Printf("[%s:%s] Error while writing objectId value (%v): %v", dbName, collName, elem.Value(), err)
			}
		}
	}
	return nil
}
