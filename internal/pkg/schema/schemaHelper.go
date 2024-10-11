package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"okieoth/schemaguesser/internal/pkg/mongoHelper"
	"okieoth/schemaguesser/internal/pkg/utils"
	"os"
	"slices"
	"text/template"
	"unicode"
)

type TypeRelation struct {
	Start string
	End   string
}

func NewTypeRelation(start string, end string) TypeRelation {
	return TypeRelation{
		Start: start,
		End:   end,
	}
}

type PumlTemplateInput struct {
	Database          string
	Collection        string
	MainType          *mongoHelper.ComplexType
	Relations         []TypeRelation
	OtherComplexTypes []mongoHelper.ComplexType
}

type TemplateInput struct {
	Database          string
	Collection        string
	MainType          *mongoHelper.ComplexType
	OtherComplexTypes []mongoHelper.ComplexType
}

func lastIndexProps(array []mongoHelper.BasicElemInfo) int {
	return len(array) - 1
}

func lastIndexTypes(array []mongoHelper.ComplexType) int {
	return len(array) - 1
}

func getComplexTypeByName(name string, otherComplexTypes []mongoHelper.ComplexType) (*mongoHelper.ComplexType, error) {
	var complexType mongoHelper.ComplexType
	for _, e := range otherComplexTypes {
		if e.Name == name {
			complexType = e
			return &complexType, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("can't find complex type with name: %s", name))
}

func containsProp(propName string, t *mongoHelper.ComplexType) bool {
	for _, p := range t.Properties {
		if p.AttribName == propName {
			return true
		}
	}
	return false
}

func complexTypesAreNotTheSame(lastTypeInst *mongoHelper.ComplexType, currentTypeInfo *mongoHelper.BasicElemInfo, otherComplexTypes []mongoHelper.ComplexType) bool {
	currentTypeInst, err := getComplexTypeByName(currentTypeInfo.ValueType, otherComplexTypes)
	if err != nil {
		fmt.Sprintf("Error while try to resolve name to complex type: %v\n", err)
		return false
	}
	var missingAttribCount int
	for _, p := range currentTypeInst.Properties {
		if !containsProp(p.AttribName, lastTypeInst) {
			missingAttribCount += 1
		}
	}

	// point where it's decided if the type is the same
	if missingAttribCount > (len(lastTypeInst.Properties) / 4) {
		return true
	} else {
		return false
	}
}

func checkForSameTypesOfAllProps_old(complexType mongoHelper.ComplexType, otherComplexTypes []mongoHelper.ComplexType) bool {
	var lastType string
	var lastTypeInst *mongoHelper.ComplexType
	var err error
	for _, p := range complexType.Properties {
		if lastType == "" {
			lastType = p.ValueType
			if p.IsComplex {
				lastTypeInst, err = getComplexTypeByName(lastType, otherComplexTypes)
				if err != nil {
					log.Printf("Error while try to resolve name (%s) to complex type: %v\n", lastType, err)
					return false
				}
			}
		} else {
			if p.IsComplex {
				if complexTypesAreNotTheSame(lastTypeInst, &p, otherComplexTypes) {
					return false
				}
			} else {
				if lastType != p.ValueType {
					// properties have different types
					return false
				}
			}
		}
	}
	return true
}

func checkForSameTypesOfAllProps(complexType mongoHelper.ComplexType, otherComplexTypes []mongoHelper.ComplexType) bool {
	var lastType string
	for _, p := range complexType.Properties {
		if !p.IsComplex {
			return false
		}
		if lastType == "" {
			lastType = p.ValueType
		} else {
			if lastType != p.ValueType {
				// properties have different types
				return false
			}
		}
	}
	return true
}

func notInTypesToRemove(typeName string, typesToRemove []string) bool {
	for _, t := range typesToRemove {
		if t == typeName {
			return false
		}
	}
	return true
}

func containsPropWithSameType(propToFind *mongoHelper.BasicElemInfo, t *mongoHelper.ComplexType, otherComplexTypes []mongoHelper.ComplexType) bool {
	var c1 *mongoHelper.ComplexType
	var err error
	if propToFind.IsComplex {
		c1, err = getComplexTypeByName(propToFind.ValueType, otherComplexTypes)
		if err != nil {
			log.Printf("containsPropWithSameType: error while resolve complex type (1): %v", err)
			return false
		}
	}
	for _, p := range t.Properties {
		if p.AttribName == propToFind.AttribName {
			if p.BsonType == "null" {
				log.Printf("containsPropWithSameType: skip property for BsonType null: type: %s, property: %v", t.Name, p)
				continue
			}

			if propToFind.IsComplex {
				c2, err := getComplexTypeByName(p.ValueType, otherComplexTypes)
				if err != nil {
					log.Printf("containsPropWithSameType: error while resolve complex type (2): %v, type: %s, property: %v", err, t.Name, p)
					return false
				}
				if !typesAreEqual(c1, c2, otherComplexTypes) {
					return false
				}
				return true
			} else {
				if p.ValueType == propToFind.ValueType {
					return true
				} else {
					return false
				}
			}
		}
	}
	return false
}

func typesAreEqual(t1, t2 *mongoHelper.ComplexType, otherComplexTypes []mongoHelper.ComplexType) bool {
	if t1.IsDictionary && t2.IsDictionary {
		return t1.DictValueType == t2.DictValueType
	} else {
		if t1.IsDictionary || t2.IsDictionary {
			return false
		} else {
			for _, p := range t1.Properties {
				if !containsPropWithSameType(&p, t2, otherComplexTypes) {
					return false
				}
			}
			return true
		}
	}
}

func replaceAllTypeReferences(typeNameToReplace string, typeNameReplacement string, otherComplexTypes []mongoHelper.ComplexType, mainType *mongoHelper.ComplexType) {
	for i, t := range otherComplexTypes {
		for j, p := range t.Properties {
			if p.ValueType == typeNameToReplace {
				(otherComplexTypes)[i].Properties[j].ValueType = typeNameReplacement
			}
		}
	}
	for j, p := range mainType.Properties {
		if p.ValueType == typeNameToReplace {
			mainType.Properties[j].ValueType = typeNameReplacement
		}
	}
}

func removeUnneededTypes(typesToRemove []string, otherComplexTypes []mongoHelper.ComplexType, mainType *mongoHelper.ComplexType) []mongoHelper.ComplexType {
	var ret []mongoHelper.ComplexType
	for _, t := range otherComplexTypes {
		if notInTypesToRemove(t.LongName, typesToRemove) {
			ret = append(ret, t)
		}
	}
	return removeDigitsFromTypeNames(mainType, ret)
}

func removeTrailingDigits(name string) string {
	runes := []rune(name)
	for i := len(runes) - 1; i >= 0; i-- {
		if unicode.IsDigit(runes[i]) {
			runes = runes[:i]
		} else {
			break
		}
	}
	return string(runes)
}

func removeDigitsFromTypeNames(mainType *mongoHelper.ComplexType, complexTypes []mongoHelper.ComplexType) []mongoHelper.ComplexType {
	for i, t := range complexTypes {
		trimmedName := removeTrailingDigits(t.Name)
		if trimmedName == "" {
			trimmedName = mongoHelper.GetNewTypeName("Type", complexTypes)
		}
		if trimmedName != t.Name {
			for j, ct := range complexTypes {
				if ct.Name == t.Name {
					continue
				}
				for k, p := range ct.Properties {
					if p.ValueType == t.Name {
						complexTypes[j].Properties[k].ValueType = trimmedName
					}
				}
			}
			for k, p := range mainType.Properties {
				if p.ValueType == t.Name {
					mainType.Properties[k].ValueType = trimmedName
				}
			}
			complexTypes[i].Name = trimmedName
		}
	}
	return complexTypes
}

func ReduceDoubleTypesByName(otherComplexTypes []mongoHelper.ComplexType) []mongoHelper.ComplexType {
	var typesToRemove []string
	for i, e1 := range otherComplexTypes {
		if e1.TypeReduced {
			continue
		}
		for j, e2 := range otherComplexTypes[i+1:] {
			if e2.TypeReduced {
				continue
			}
			if e1.Name == e2.Name {
				typesToRemove = append(typesToRemove, e2.LongName)
				otherComplexTypes[i+j+1].TypeReduced = true
			}
		}
	}
	var ret []mongoHelper.ComplexType
	for _, t := range otherComplexTypes {
		if notInTypesToRemove(t.LongName, typesToRemove) {
			ret = append(ret, t)
		}
	}
	return ret
}

func ReduceTypes(mainType *mongoHelper.ComplexType, otherComplexTypes []mongoHelper.ComplexType) []mongoHelper.ComplexType {
	var typesToRemove []string
	for i, e1 := range otherComplexTypes {
		if e1.TypeReduced {
			continue
		}
		for j, e2 := range otherComplexTypes[i+1:] {
			if e2.TypeReduced {
				continue
			}
			if typesAreEqual(&e1, &e2, otherComplexTypes) {
				typesToRemove = append(typesToRemove, e2.LongName)
				otherComplexTypes[i+j+1].TypeReduced = true
				replaceAllTypeReferences(e2.Name, e1.Name, otherComplexTypes, mainType)
			}
		}
	}
	return removeUnneededTypes(typesToRemove, otherComplexTypes, mainType)
}

func GuessDicts(otherComplexTypes []mongoHelper.ComplexType) []mongoHelper.ComplexType {
	var typesToRemove []string
	for i := range otherComplexTypes {
		e := otherComplexTypes[i]
		if checkForSameTypesOfAllProps(e, otherComplexTypes) {
			var typeNameToUse string
			for i, p := range e.Properties {
				e.UsedKeys = append(e.UsedKeys, p.AttribName)
				if i == 0 {
					typeNameToUse = p.ValueType
				} else {
					if p.ValueType != typeNameToUse {
						typesToRemove = append(typesToRemove, p.ValueType)
					}
				}
			}
			otherComplexTypes[i].Properties = make([]mongoHelper.BasicElemInfo, 0)
			otherComplexTypes[i].IsDictionary = true
			otherComplexTypes[i].DictValueType = typeNameToUse
		}
	}
	var ret []mongoHelper.ComplexType
	for _, t := range otherComplexTypes {
		if notInTypesToRemove(t.Name, typesToRemove) {
			ret = append(ret, t)
		}
	}
	return ret
}

func PrintSchema(database string, collection string, mainType *mongoHelper.ComplexType, otherComplexTypes []mongoHelper.ComplexType, outputDir string) {
	input := TemplateInput{
		MainType:          mainType,
		OtherComplexTypes: otherComplexTypes,
		Database:          database,

		Collection: collection,
	}
	printTemplateBase("schema.tmpl", schemaTemplateStr, "schema.json", database, collection, &input, outputDir)
}

func PersistSchemaBase(database string, collection string, mainType *mongoHelper.ComplexType, otherComplexTypes []mongoHelper.ComplexType, outputDir string) {
	schemaRaw := mongoHelper.SchemaRaw{
		MainType:          mainType,
		OtherComplexTypes: &otherComplexTypes,
	}

	jsonData, err := json.MarshalIndent(schemaRaw, "", "  ")
	if err != nil {
		log.Printf("[%s:%s] error, failed to marshall schema raw data: %v", database, collection, err)
	}

	if outputDir == "stdout" {
		fmt.Println(jsonData)
	} else {
		outputFile, err := utils.CreateOutputFile(outputDir, "schema-raw.json", database, collection)
		if err != nil {
			panic(err)
		}
		defer outputFile.Close()

		_, err = outputFile.Write(jsonData)
		if err != nil {
			log.Printf("[%s:%s] error, failed to write schema raw data to file: %v", database, collection, err)
		}
	}
}

func WritePlantUml(database string, collection string, mainType *mongoHelper.ComplexType, otherComplexTypes []mongoHelper.ComplexType, outputDir string) {
	typeRelations := make([]TypeRelation, 0)
	for _, e := range mainType.Properties {
		if (e.IsComplex) && (!slices.ContainsFunc(typeRelations, func(v TypeRelation) bool {
			return (v.Start == mainType.Name) && (v.End == e.ValueType)
		})) {
			typeRelations = append(typeRelations, NewTypeRelation(mainType.Name, e.ValueType))
		}
	}

	for _, eo := range otherComplexTypes {
		for _, ep := range eo.Properties {
			if (ep.IsComplex) && (!slices.ContainsFunc(typeRelations, func(v TypeRelation) bool {
				return (v.Start == eo.Name) && (v.End == ep.ValueType)
			})) {
				typeRelations = append(typeRelations, NewTypeRelation(eo.Name, ep.ValueType))
			}
		}
	}

	input := PumlTemplateInput{
		MainType:          mainType,
		OtherComplexTypes: otherComplexTypes,
		Relations:         typeRelations,
		Database:          database,
		Collection:        collection,
	}
	printTemplateBase("plantuml.tmpl", pumlTemplateStr, "schema.puml", database, collection, &input, outputDir)
}

func printTemplateBase(templateName string, templateStr string, fileExt string, database string, collection string, input interface{}, outputDir string) {
	tmpl := template.Must(template.New(templateName).Funcs(template.FuncMap{
		"LastIndexProps": lastIndexProps, "LastIndexTypes": lastIndexTypes,
	}).Parse(templateStr))

	if outputDir == "stdout" {
		err := tmpl.Execute(os.Stdout, input)
		if err != nil {
			panic(err)
		}
	} else {
		outputFile, err := utils.CreateOutputFile(outputDir, fileExt, database, collection)
		if err != nil {
			panic(err)
		}
		defer outputFile.Close()
		err = tmpl.Execute(outputFile, input)
		if err != nil {
			panic(err)
		}
	}
}
