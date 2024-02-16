package schema

import (
	"fmt"
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

func PrintSchema(database string, collection string, mainType *mongoHelper.ComplexType, otherComplexTypes []mongoHelper.ComplexType, outputDir string) {
	tmplFile := "json_schema.tmpl"
	tmpl := template.Must(template.New("json_schema.tmpl").Funcs(template.FuncMap{
		"LastIndexProps": lastIndexProps, "LastIndexTypes": lastIndexTypes,
	}).ParseFiles("resources/" + tmplFile))

	input := TemplateInput{
		MainType:          mainType,
		OtherComplexTypes: otherComplexTypes,
		Database:          database,

		Collection: collection,
	}

	if outputDir == "stdout" {
		err := tmpl.Execute(os.Stdout, input)
		if err != nil {
			panic(err)
		}
	} else {
		outputFile, err := os.Create(outputDir + string(os.PathSeparator) + fmt.Sprintf("%s_%s.json", database, collection))
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
