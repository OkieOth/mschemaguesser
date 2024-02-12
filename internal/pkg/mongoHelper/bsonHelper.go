package mongoHelper

import (
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

const NUMBER = "number"
const STRING = "string"
const OBJECT = "object"
const INT = "integer"

type ComplexType struct {
	Properties []BasicElemInfo
}

type BasicElemInfo struct {
	AttribName string
	ValueType  string
	Format     string
	IsArray    bool
	Comment    string
}

type SchemaType struct {
	Properties   []BasicElemInfo
	complexTypes []ComplexType
}

func ProcessBson(doc bson.Raw) (*SchemaType, error) {
	elements, err := doc.Elements()
	if err != nil {
		log.Fatal("Error while parsing bson elements: %v", err)
		return nil, err
	}
	var ret SchemaType
	for _, elem := range elements {
		typeInfo := BasicElemInfo{AttribName: elem.Key()}
		typeInfo.AttribName = elem.Key()
		switch elem.Value().Type {
		case bson.TypeDouble:
			handleTypeDouble(elem, &typeInfo)
		case bson.TypeEmbeddedDocument:
			handleTypeEmbeddedDocument(elem, &typeInfo, &ret)
		case bson.TypeArray:
			handleTypeArray(elem, &typeInfo, &ret)
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
	return &ret, nil
}

func handleTypeDouble(elem bson.RawElement, typeInfo *BasicElemInfo) {
	typeInfo.ValueType = NUMBER
	typeInfo.Comment = "Mongodb type double"
}

func handleTypeEmbeddedDocument(elem bson.RawElement, typeInfo *BasicElemInfo, schemaType *SchemaType) {
	// TODO
}

func handleTypeArray(elem bson.RawElement, typeInfo *BasicElemInfo, schemaType *SchemaType) {
	// TODO
}

func handleTypeBinary(elem bson.RawElement, typeInfo *BasicElemInfo) {
	typeInfo.ValueType = STRING
	typeInfo.Comment = "Mongodb type binary"
}

func handleTypeUndefined(elem bson.RawElement, typeInfo *BasicElemInfo) {
	typeInfo.Comment = "Mongodb type Undefined"
	typeInfo.ValueType = OBJECT
}

func handleTypeObjectID(elem bson.RawElement, typeInfo *BasicElemInfo) {
	typeInfo.Comment = "Mongodb type objectId"
	// TODO
}

func handleTypeBoolean(elem bson.RawElement, typeInfo *BasicElemInfo) {
	typeInfo.ValueType = "boolean"
}

func handleTypeDateTime(elem bson.RawElement, typeInfo *BasicElemInfo) {
	typeInfo.ValueType = STRING
	typeInfo.Format = "date-time"
}

func handleTypeNull(elem bson.RawElement, typeInfo *BasicElemInfo) {
	typeInfo.Comment = "Mongodb type null"
	typeInfo.ValueType = OBJECT
}

func handleTypeRegex(elem bson.RawElement, typeInfo *BasicElemInfo) {
	typeInfo.ValueType = STRING
	typeInfo.Format = "regex"
}

func handleTypeDBPointer(elem bson.RawElement, typeInfo *BasicElemInfo) {
	typeInfo.Comment = "Mongodb type db pointer"
	typeInfo.ValueType = OBJECT
}

func handleTypeJavaScript(elem bson.RawElement, typeInfo *BasicElemInfo) {
	typeInfo.Comment = "Mongodb type javascript"
	typeInfo.ValueType = STRING
}

func handleTypeSymbol(elem bson.RawElement, typeInfo *BasicElemInfo) {
	typeInfo.Comment = "Mongodb type symbol"
	typeInfo.ValueType = OBJECT
}

func handleTypeCodeWithScope(elem bson.RawElement, typeInfo *BasicElemInfo) {
	typeInfo.Comment = "Mongodb type code with scope"
	typeInfo.ValueType = STRING

}

func handleTypeInt32(elem bson.RawElement, typeInfo *BasicElemInfo) {
	typeInfo.ValueType = INT
	typeInfo.Format = "int32"
}

func handleTypeInt64(elem bson.RawElement, typeInfo *BasicElemInfo) {
	typeInfo.ValueType = INT
	typeInfo.Format = "int64"
}

func handleTypeTimestamp(elem bson.RawElement, typeInfo *BasicElemInfo) {
	typeInfo.ValueType = STRING
	typeInfo.Format = "time"
}

func handleTypeDecimal128(elem bson.RawElement, typeInfo *BasicElemInfo) {
	typeInfo.ValueType = NUMBER
	typeInfo.Comment = "Mongodb type Decimal128"
}

func handleTypeMinKey(elem bson.RawElement, typeInfo *BasicElemInfo) {
	typeInfo.ValueType = OBJECT
	typeInfo.Comment = "Mongodb type minKey"
}

func handleTypeMaxKey(elem bson.RawElement, typeInfo *BasicElemInfo) {
	typeInfo.ValueType = OBJECT
	typeInfo.Comment = "Mongodb type maxKey"
}
