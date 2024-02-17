package cmd

import (
	"fmt"

	"okieoth/schemaguesser/internal/pkg/mongoHelper"
	"okieoth/schemaguesser/internal/pkg/schema"

	"github.com/spf13/cobra"
)

var dbName string

var colName string

var outputDir string

var itemCount int32

var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "functions around the schemas",
	Long:  "With this command you can create schemas out of mongodb collection",
	Run: func(cmd *cobra.Command, args []string) {
		bsonRaw, err := mongoHelper.QueryCollection(mongoHelper.ConStr, dbName, colName, int(itemCount))
		if err != nil {
			msg := fmt.Sprintf("Error while reading indexes for collection (%s.%s): \n%v\n", databaseName, colName, err)
			panic(msg)
		}
		var otherComplexTypes []mongoHelper.ComplexType
		var mainType mongoHelper.ComplexType
		for _, b := range bsonRaw {
			err = mongoHelper.ProcessBson(b, colName, &mainType, &otherComplexTypes)
			if err != nil {
				fmt.Printf("Error while processing bson for schema: %v", err)
			}
		}
		schema.PrintSchema(dbName, colName, &mainType, otherComplexTypes, outputDir)
	},
}

func init() {
	schemaCmd.Flags().StringVar(&dbName, "database", "all", "Database to query existing collections")

	schemaCmd.Flags().StringVar(&colName, "collection", "all", "Name of the collection to show the indexes")

	schemaCmd.Flags().StringVar(&outputDir, "output", "stdout", "stdout or the directory to write the created schema file, default is 'stdout'")

	schemaCmd.Flags().Int32Var(&itemCount, "item_count", 100, "Number of collection entries used to build the schema")
}
