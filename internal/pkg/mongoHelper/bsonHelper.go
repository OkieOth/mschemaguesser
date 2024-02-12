package mongoHelper

import (
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

func ProcessBson(doc bson.Raw) {
	elements, err := doc.Elements()
	if err != nil {
		log.Fatal("Error while parsing bson elements: %v", err)
		return
	}
	for _, elem := range elements {
		switch elem.Value().Type {
		case bson.TypeDouble:
			fmt.Println("double")
		case bson.TypeEmbeddedDocument:
			fmt.Println("embededDocument")
		case bson.TypeArray:
			fmt.Println("double")
		case bson.TypeBinary:
			fmt.Println("binary")
		case bson.TypeUndefined:
			fmt.Println("undefined")
		case bson.TypeObjectID:
			fmt.Println("ObjectId")
		case bson.TypeBoolean:
			fmt.Println("boolean")
		case bson.TypeDateTime:
			fmt.Println("dateTime")
		case bson.TypeNull:
			fmt.Println("typeNull")
		case bson.TypeRegex:
			fmt.Println("regex")
		case bson.TypeDBPointer:
			fmt.Println("dbPointer")
		case bson.TypeJavaScript:
			fmt.Println("javascript")
		case bson.TypeSymbol:
			fmt.Println("symbol")
		case bson.TypeCodeWithScope:
			fmt.Println("typeCodeWithScope")
		case bson.TypeInt32:
			fmt.Println("int32")
		case bson.TypeInt64:
			fmt.Println("int64")
		case bson.TypeTimestamp:
			fmt.Println("timestamp")
		case bson.TypeDecimal128:
			fmt.Println("decimal128")
		case bson.TypeMinKey:
			fmt.Println("minKey")
		case bson.TypeMaxKey:
			fmt.Println("maxKey")
		}
	}
}
