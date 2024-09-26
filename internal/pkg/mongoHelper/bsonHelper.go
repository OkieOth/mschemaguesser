package mongoHelper

import (
	"errors"
	"fmt"
	"log"
	"unicode"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

const NUMBER = "number"
const STRING = "string"
const OBJECT = "object"
const INT = "integer"

type ComplexType struct {
	Name          string
	LongName      string
	Properties    []BasicElemInfo
	IsDictionary  bool
	DictValueType string
	UsedKeys      []string
	TypeReduced   bool
}

type BasicElemInfo struct {
	AttribName      string
	ValueType       string
	BsonType        string
	Format          string
	IsArray         bool
	ArrayDimensions uint
	Comment         string
	IsComplex       bool
}

func GetNewTypeName(name string, otherComplexTypes *[]ComplexType) string {
	f := func(s string) bool {
		for _, c := range *otherComplexTypes {
			if c.Name == s {
				return true
			}
		}
		return false
	}
	baseName := firstUpperCase(name)
	newName := baseName
	index := 2
	for f(newName) {
		newName = fmt.Sprintf("%s%d", baseName, index)
		index += 1
	}
	return newName
}

func firstUpperCase(s string) string {
	if len(s) == 0 {
		return s
	}

	firstChar := []rune(s)[0]
	upperFirstChar := unicode.ToUpper(firstChar)

	result := string(upperFirstChar) + s[1:]
	return result
}

func addNewOtherComplexType(otherComplexTypes *[]ComplexType, complexType ComplexType) {
	for i, e := range *otherComplexTypes {
		if e.LongName == complexType.LongName {
			(*otherComplexTypes)[i] = complexType
			return
		}
	}
	*otherComplexTypes = append(*otherComplexTypes, complexType)
}

func addNewProperty(properties *[]BasicElemInfo, prop BasicElemInfo) {
	for i, e := range *properties {
		if e.AttribName == prop.AttribName {
			(*properties)[i] = prop
			return
		}
	}
	*properties = append(*properties, prop)
}

func getAlreadyStoredType(otherComplexTypes *[]ComplexType, typeName string) (ComplexType, bool) {
	for i, e := range *otherComplexTypes {
		if e.LongName == typeName {
			return (*otherComplexTypes)[i], true
		}
	}
	return ComplexType{}, false
}

func hasAlreadyProperty(mainType *ComplexType, attribName string) bool {
	for _, p := range mainType.Properties {
		if p.AttribName == attribName {
			return true
		}
	}
	return false
}

func isBasicType(elem bson.RawElement) bool {
	return !((elem.Value().Type == bson.TypeArray) || (elem.Value().Type == bson.TypeEmbeddedDocument))
}

func ProcessBson(doc bson.Raw, collectionName string, mainType *ComplexType, otherComplexTypes *[]ComplexType) error {
	if mainType == nil {
		return errors.New("no mainType given")
	}
	if otherComplexTypes == nil {
		return errors.New("no otherComplexTypes given")
	}
	elements, err := doc.Elements()
	if err != nil {
		log.Fatalf("Error while parsing bson elements: %v", err)
		return err
	}
	if mainType.Name == "" {
		colNameFirstUpper := firstUpperCase(collectionName)
		mainType.Name = colNameFirstUpper
		mainType.LongName = colNameFirstUpper
	}
	for _, elem := range elements {
		isAlreadyThere := hasAlreadyProperty(mainType, elem.Key())
		if isAlreadyThere && isBasicType(elem) {
			continue
		}
		typeInfo := BasicElemInfo{AttribName: elem.Key()}
		typeInfo.AttribName = elem.Key()
		switch elem.Value().Type {
		case bson.TypeString:
			handleTypeString(elem, &typeInfo)
		case bson.TypeDouble:
			handleTypeDouble(elem, &typeInfo)
		case bson.TypeEmbeddedDocument:
			newTypeLongName := firstUpperCase(collectionName) + firstUpperCase(elem.Key())
			newSchemaType, existingOne := getAlreadyStoredType(otherComplexTypes, newTypeLongName)
			var newTypeName string
			if !existingOne {
				newSchemaType = ComplexType{}
				newSchemaType.LongName = newTypeLongName
				newTypeName = GetNewTypeName(elem.Key(), otherComplexTypes)
				newSchemaType.Name = newTypeName
			} else {
				newTypeName = newSchemaType.Name
			}
			typeInfo.ValueType = newTypeName
			handleTypeEmbeddedDocument(elem, &typeInfo, &newSchemaType, otherComplexTypes, newTypeName, true)
			addNewOtherComplexType(otherComplexTypes, newSchemaType)
		case bson.TypeArray:
			handleTypeArray(elem, &typeInfo, otherComplexTypes, mainType.Name)
		case bson.TypeBinary:
			handleTypeBinary(elem, &typeInfo)
		case bson.TypeUndefined:
			handleTypeUndefined(elem, &typeInfo)
		case bson.TypeObjectID:
			handleTypeObjectID(elem, &typeInfo)
		case bson.TypeBoolean:
			handleTypeBoolean(elem, &typeInfo)
		case bson.TypeDateTime:
			handleTypeDateTime(elem, &typeInfo)
		case bson.TypeNull:
			handleTypeNull(elem, &typeInfo)
		case bson.TypeRegex:
			handleTypeRegex(elem, &typeInfo)
		case bson.TypeDBPointer:
			handleTypeDBPointer(elem, &typeInfo)
		case bson.TypeJavaScript:
			handleTypeJavaScript(elem, &typeInfo)
		case bson.TypeSymbol:
			handleTypeSymbol(elem, &typeInfo)
		case bson.TypeCodeWithScope:
			handleTypeCodeWithScope(elem, &typeInfo)
		case bson.TypeInt32:
			handleTypeInt32(elem, &typeInfo)
		case bson.TypeInt64:
			handleTypeInt64(elem, &typeInfo)
		case bson.TypeTimestamp:
			handleTypeTimestamp(elem, &typeInfo)
		case bson.TypeDecimal128:
			handleTypeDecimal128(elem, &typeInfo)
		case bson.TypeMinKey:
			handleTypeMinKey(elem, &typeInfo)
		case bson.TypeMaxKey:
			handleTypeMaxKey(elem, &typeInfo)
		}
		addNewProperty(&mainType.Properties, typeInfo)
	}
	return nil
}

func handleTypeString(elem bson.RawElement, typeInfo *BasicElemInfo) {
	typeInfo.ValueType = STRING
	typeInfo.BsonType = STRING
}

func handleTypeDouble(elem bson.RawElement, typeInfo *BasicElemInfo) {
	typeInfo.ValueType = NUMBER
	typeInfo.BsonType = "double"
}

func handleTypeEmbeddedDocument(elem bson.RawElement, typeInfo *BasicElemInfo, schemaType *ComplexType, otherComplexTypes *[]ComplexType, prefix string, addToOtherSchemas bool) {
	typeInfo.BsonType = "embeddedDocument - unofficial type"
	typeInfo.IsComplex = true

	schemaType.Name = typeInfo.ValueType
	embeddedDoc := bson.Raw(elem.Value().Value)

	elements, err := embeddedDoc.Elements()
	if err != nil {
		typeInfo.Comment = fmt.Sprintf("error while parsing complex type: %v", err)
	} else {
		for _, elem := range elements {
			typeInfo := BasicElemInfo{AttribName: elem.Key()}
			typeInfo.AttribName = elem.Key()
			switch elem.Value().Type {
			case bson.TypeString:
				handleTypeString(elem, &typeInfo)
			case bson.TypeDouble:
				handleTypeDouble(elem, &typeInfo)
			case bson.TypeEmbeddedDocument:
				var newTypeLongName, newTypeName string
				newTypeLongName = schemaType.LongName + firstUpperCase(elem.Key())
				newTypeName = GetNewTypeName(elem.Key(), otherComplexTypes)

				// newTypeLongName = schemaType.LongName + firstUpperCase(elem.Key())
				// newTypeName = getNewTypeName(elem.Key(), otherComplexTypes)

				newSchemaType, existingOne := getAlreadyStoredType(otherComplexTypes, newTypeLongName)
				if !existingOne {
					newSchemaType = ComplexType{}
					newSchemaType.LongName = newTypeLongName
					newSchemaType.Name = newTypeName
				}
				typeInfo.ValueType = newTypeName
				handleTypeEmbeddedDocument(elem, &typeInfo, &newSchemaType, otherComplexTypes, schemaType.Name, true)
				addNewOtherComplexType(otherComplexTypes, newSchemaType)
			case bson.TypeArray:
				handleTypeArray(elem, &typeInfo, otherComplexTypes, schemaType.Name)
			case bson.TypeBinary:
				handleTypeBinary(elem, &typeInfo)
			case bson.TypeUndefined:
				handleTypeUndefined(elem, &typeInfo)
			case bson.TypeObjectID:
				handleTypeObjectID(elem, &typeInfo)
			case bson.TypeBoolean:
				handleTypeBoolean(elem, &typeInfo)
			case bson.TypeDateTime:
				handleTypeDateTime(elem, &typeInfo)
			case bson.TypeNull:
				handleTypeNull(elem, &typeInfo)
			case bson.TypeRegex:
				handleTypeRegex(elem, &typeInfo)
			case bson.TypeDBPointer:
				handleTypeDBPointer(elem, &typeInfo)
			case bson.TypeJavaScript:
				handleTypeJavaScript(elem, &typeInfo)
			case bson.TypeSymbol:
				handleTypeSymbol(elem, &typeInfo)
			case bson.TypeCodeWithScope:
				handleTypeCodeWithScope(elem, &typeInfo)
			case bson.TypeInt32:
				handleTypeInt32(elem, &typeInfo)
			case bson.TypeInt64:
				handleTypeInt64(elem, &typeInfo)
			case bson.TypeTimestamp:
				handleTypeTimestamp(elem, &typeInfo)
			case bson.TypeDecimal128:
				handleTypeDecimal128(elem, &typeInfo)
			case bson.TypeMinKey:
				handleTypeMinKey(elem, &typeInfo)
			case bson.TypeMaxKey:
				handleTypeMaxKey(elem, &typeInfo)
			}
			addNewProperty(&schemaType.Properties, typeInfo)

		}
	}
}

func handleTypeArray(elem bson.RawElement, typeInfo *BasicElemInfo, otherComplexTypes *[]ComplexType, prefix string) {
	arrayRaw := bson.Raw(elem.Value().Value)

	typeInfo.IsArray = true
	typeInfo.ArrayDimensions++
	typeInfo.BsonType = "couldn't be retrieved - no elems"
	typeInfo.ValueType = OBJECT
	newTypeLongName := prefix + firstUpperCase(elem.Key())
	newTypeName := GetNewTypeName(elem.Key(), otherComplexTypes)

	elements, err := arrayRaw.Elements()
	if err != nil {
		typeInfo.Comment = fmt.Sprintf("error while parsing array type: %v", err)
		typeInfo.BsonType = "array type - unofficial type"
		return
	}

	var lastType *bsontype.Type
	var complexArrayType ComplexType
	for _, elem := range elements {
		if (lastType != nil) && (*lastType != elem.Value().Type) {
			typeInfo.Comment = "array type consists of different types, multiple type arrays are not supported"
			typeInfo.BsonType = "array type - unofficial type"
			return
		}

		if lastType == nil {
			switch elem.Value().Type {
			case bson.TypeString:
				handleTypeString(elem, typeInfo)
			case bson.TypeDouble:
				handleTypeDouble(elem, typeInfo)
			case bson.TypeEmbeddedDocument:
				newSchemaType, existingOne := getAlreadyStoredType(otherComplexTypes, newTypeLongName)
				if !existingOne {
					newSchemaType = ComplexType{}
					newSchemaType.LongName = newTypeLongName
					newSchemaType.Name = newTypeName
				}
				typeInfo.ValueType = newTypeName
				handleTypeEmbeddedDocument(elem, typeInfo, &newSchemaType, otherComplexTypes, newTypeName, true)
				addNewOtherComplexType(otherComplexTypes, newSchemaType)
				complexArrayType = newSchemaType
			case bson.TypeArray:
				handleTypeArray(elem, typeInfo, otherComplexTypes, newTypeName)
			case bson.TypeBinary:
				handleTypeBinary(elem, typeInfo)
			case bson.TypeUndefined:
				handleTypeUndefined(elem, typeInfo)
			case bson.TypeObjectID:
				handleTypeObjectID(elem, typeInfo)
			case bson.TypeBoolean:
				handleTypeBoolean(elem, typeInfo)
			case bson.TypeDateTime:
				handleTypeDateTime(elem, typeInfo)
			case bson.TypeNull:
				handleTypeNull(elem, typeInfo)
			case bson.TypeRegex:
				handleTypeRegex(elem, typeInfo)
			case bson.TypeDBPointer:
				handleTypeDBPointer(elem, typeInfo)
			case bson.TypeJavaScript:
				handleTypeJavaScript(elem, typeInfo)
			case bson.TypeSymbol:
				handleTypeSymbol(elem, typeInfo)
			case bson.TypeCodeWithScope:
				handleTypeCodeWithScope(elem, typeInfo)
			case bson.TypeInt32:
				handleTypeInt32(elem, typeInfo)
			case bson.TypeInt64:
				handleTypeInt64(elem, typeInfo)
			case bson.TypeTimestamp:
				handleTypeTimestamp(elem, typeInfo)
			case bson.TypeDecimal128:
				handleTypeDecimal128(elem, typeInfo)
			case bson.TypeMinKey:
				handleTypeMinKey(elem, typeInfo)
			case bson.TypeMaxKey:
				handleTypeMaxKey(elem, typeInfo)
			}
		} else {
			// only complex types needs to be reviewed for additional attributes
			if elem.Value().Type == bson.TypeEmbeddedDocument {
				handleTypeEmbeddedDocument(elem, typeInfo, &complexArrayType, otherComplexTypes, newTypeName, true)
			}

		}
	}

}

func handleTypeBinary(elem bson.RawElement, typeInfo *BasicElemInfo) {
	subtype, _ := elem.Value().Binary()
	typeInfo.ValueType = STRING
	typeInfo.Comment = fmt.Sprintf("Mongodb type binary: subtype=%v", subtype)
	switch subtype {
	case 3:
		typeInfo.Format = "uuid"
	case 4:
		typeInfo.Format = "uuid"
	case 5:
		typeInfo.Format = "md5"
	}
	typeInfo.BsonType = "binData"
}

func handleTypeUndefined(elem bson.RawElement, typeInfo *BasicElemInfo) {
	typeInfo.Comment = "deprecated in mongodb"
	typeInfo.ValueType = OBJECT
	typeInfo.BsonType = "undefined"
}

func handleTypeObjectID(elem bson.RawElement, typeInfo *BasicElemInfo) {
	typeInfo.BsonType = "objectId"
	typeInfo.ValueType = OBJECT
}

func handleTypeBoolean(elem bson.RawElement, typeInfo *BasicElemInfo) {
	typeInfo.ValueType = "boolean"
	typeInfo.BsonType = "bool"
}

func handleTypeDateTime(elem bson.RawElement, typeInfo *BasicElemInfo) {
	typeInfo.ValueType = STRING
	typeInfo.Format = "date-time"
	typeInfo.BsonType = "date"
}

func handleTypeNull(elem bson.RawElement, typeInfo *BasicElemInfo) {
	typeInfo.ValueType = OBJECT
	typeInfo.BsonType = "null"
}

func handleTypeRegex(elem bson.RawElement, typeInfo *BasicElemInfo) {
	typeInfo.ValueType = STRING
	typeInfo.Format = "regex"
	typeInfo.BsonType = "regex"
}

func handleTypeDBPointer(elem bson.RawElement, typeInfo *BasicElemInfo) {
	typeInfo.ValueType = OBJECT
	typeInfo.BsonType = "dbPointer"
	typeInfo.Comment = "deprecated in mongodb"
}

func handleTypeJavaScript(elem bson.RawElement, typeInfo *BasicElemInfo) {
	typeInfo.ValueType = STRING
	typeInfo.BsonType = "javascript"
}

func handleTypeSymbol(elem bson.RawElement, typeInfo *BasicElemInfo) {
	typeInfo.Comment = "deprecated in mongodb"
	typeInfo.ValueType = OBJECT
	typeInfo.BsonType = "symbol"
}

func handleTypeCodeWithScope(elem bson.RawElement, typeInfo *BasicElemInfo) {
	typeInfo.ValueType = STRING
	typeInfo.BsonType = "javascriptWithScope"
	typeInfo.Comment = "deprecated in mongodb"
}

func handleTypeInt32(elem bson.RawElement, typeInfo *BasicElemInfo) {
	typeInfo.ValueType = INT
	typeInfo.Format = "int32"
	typeInfo.BsonType = "int"
}

func handleTypeInt64(elem bson.RawElement, typeInfo *BasicElemInfo) {
	typeInfo.ValueType = INT
	typeInfo.Format = "int64"
	typeInfo.BsonType = "long"
}

func handleTypeTimestamp(elem bson.RawElement, typeInfo *BasicElemInfo) {
	typeInfo.ValueType = STRING
	typeInfo.Format = "time"
	typeInfo.BsonType = "timestamp"
}

func handleTypeDecimal128(elem bson.RawElement, typeInfo *BasicElemInfo) {
	typeInfo.ValueType = NUMBER
	typeInfo.BsonType = "decimal"
}

func handleTypeMinKey(elem bson.RawElement, typeInfo *BasicElemInfo) {
	typeInfo.ValueType = OBJECT
	typeInfo.BsonType = "minKey"
}

func handleTypeMaxKey(elem bson.RawElement, typeInfo *BasicElemInfo) {
	typeInfo.ValueType = OBJECT
	typeInfo.BsonType = "maxKey"
}
