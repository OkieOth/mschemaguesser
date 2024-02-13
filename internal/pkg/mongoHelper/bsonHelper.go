package mongoHelper

import (
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
	Properties []BasicElemInfo
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

type SchemaType struct {
	Name         string
	Properties   []BasicElemInfo
	ComplexTypes []ComplexType
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

func ProcessBson(doc bson.Raw, collectionName string) (*SchemaType, *[]SchemaType, error) {
	var mainType SchemaType
	var otherComplexTypes = make([]SchemaType, 1)
	elements, err := doc.Elements()
	if err != nil {
		log.Fatalf("Error while parsing bson elements: %v", err)
		return &mainType, &otherComplexTypes, err
	}
	mainType.Name = collectionName
	colNameFirstUpper := firstUpperCase(collectionName)
	for _, elem := range elements {
		typeInfo := BasicElemInfo{AttribName: elem.Key()}
		typeInfo.AttribName = elem.Key()
		mainType.Properties = append(mainType.Properties, typeInfo)
		switch elem.Value().Type {
		case bson.TypeDouble:
			handleTypeDouble(elem, &typeInfo)
		case bson.TypeEmbeddedDocument:
			newTypeName := collectionName + firstUpperCase(elem.Key())
			typeInfo.ValueType = newTypeName
			handleTypeEmbeddedDocument(elem, &typeInfo, &mainType, otherComplexTypes, colNameFirstUpper, false)
		case bson.TypeArray:
			handleTypeArray(elem, &typeInfo, otherComplexTypes, colNameFirstUpper)
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
	}
	return &mainType, &otherComplexTypes, nil
}

func handleTypeDouble(elem bson.RawElement, typeInfo *BasicElemInfo) {
	typeInfo.ValueType = NUMBER
	typeInfo.BsonType = "double"
}

func handleTypeEmbeddedDocument(elem bson.RawElement, typeInfo *BasicElemInfo, schemaType *SchemaType, otherComplexTypes []SchemaType, prefix string, addToOtherSchemas bool) {
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
			schemaType.Properties = append(schemaType.Properties, typeInfo)
			switch elem.Value().Type {
			case bson.TypeDouble:
				handleTypeDouble(elem, &typeInfo)
			case bson.TypeEmbeddedDocument:
				var newSchemaType SchemaType
				newTypeName := schemaType.Name + firstUpperCase(elem.Key())
				typeInfo.ValueType = newTypeName
				handleTypeEmbeddedDocument(elem, &typeInfo, &newSchemaType, otherComplexTypes, schemaType.Name, true)
				otherComplexTypes = append(otherComplexTypes, newSchemaType)
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
		}
	}
}

func handleTypeArray(elem bson.RawElement, typeInfo *BasicElemInfo, otherComplexTypes []SchemaType, prefix string) {
	arrayRaw := bson.Raw(elem.Value().Value)

	typeInfo.IsArray = true
	typeInfo.ArrayDimensions++
	newTypeName := prefix + firstUpperCase(elem.Key())

	elements, err := arrayRaw.Elements()
	if err != nil {
		typeInfo.Comment = fmt.Sprintf("error while parsing array type: %v", err)
		typeInfo.BsonType = "array type - unofficial type"
		return
	}

	var lastType *bsontype.Type
	for _, elem := range elements {
		if (lastType != nil) && (*lastType != elem.Value().Type) {
			typeInfo.Comment = "array type consists of different types, multiple type arrays are not supported"
			typeInfo.BsonType = "array type - unofficial type"
			return
		}

		if lastType == nil {
			switch elem.Value().Type {
			case bson.TypeDouble:
				handleTypeDouble(elem, typeInfo)
			case bson.TypeEmbeddedDocument:
				var newSchemaType SchemaType
				typeInfo.ValueType = newTypeName
				handleTypeEmbeddedDocument(elem, typeInfo, &newSchemaType, otherComplexTypes, newTypeName, true)
				otherComplexTypes = append(otherComplexTypes, newSchemaType)
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
		}
	}

}

func handleTypeBinary(elem bson.RawElement, typeInfo *BasicElemInfo) {
	typeInfo.ValueType = STRING
	typeInfo.Comment = "Mongodb type binary"
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
