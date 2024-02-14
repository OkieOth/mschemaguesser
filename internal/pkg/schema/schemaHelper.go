package schema

import (
	"okieoth/schemaguesser/internal/pkg/mongoHelper"
	"os"
	"text/template"
)

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

func PrintSchema(database string, collection string, mainType *mongoHelper.ComplexType, otherComplexTypes []mongoHelper.ComplexType) {
	tmplFile := "json_schema.tmpl"
	tmpl := template.Must(template.New(tmplFile).Funcs(template.FuncMap{
		"LastIndexProps": lastIndexProps, "LastIndexTypes": lastIndexTypes,
	}).ParseFiles(tmplFile))

	input := TemplateInput{
		MainType:          mainType,
		OtherComplexTypes: otherComplexTypes,
		Database:          database,
		Collection:        collection,
	}
	err := tmpl.Execute(os.Stdout, input)
	if err != nil {
		panic(err)
	}
}
