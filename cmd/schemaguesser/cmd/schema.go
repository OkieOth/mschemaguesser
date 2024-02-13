package cmd

import (
	"fmt"

	"okieoth/schemaguesser/internal/pkg/mongoHelper"

	"github.com/spf13/cobra"
)

var dbName string

var colName string

var outputDir string

var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "functions around the schemas",
	Long:  "With this command you can create schemas out of mongodb collection",
	Run: func(cmd *cobra.Command, args []string) {
		bsonRaw, err := mongoHelper.QueryCollection(mongoHelper.ConStr, dbName, colName, 100)
		if err != nil {
			msg := fmt.Sprintf("Error while reading indexes for collection (%s.%s): \n%v\n", databaseName, collectionName, err)
			panic(msg)
		}
		var otherComplexTypes = make([]mongoHelper.ComplexType, 1)
		var mainType mongoHelper.ComplexType
		for i, b := range bsonRaw {
			if i > 10 {
				break
			}
			err = mongoHelper.ProcessBson(b, collectionName, &mainType, &otherComplexTypes)
			fmt.Println(b)
		}
	},
}

func init() {
	schemaCmd.Flags().StringVar(&dbName, "database", "", "Database to query existing collections")
	schemaCmd.MarkFlagRequired("database")

	schemaCmd.Flags().StringVar(&colName, "collection", "", "Name of the collection to show the indexes")
	schemaCmd.MarkFlagRequired("collection")

	schemaCmd.Flags().StringVar(&outputDir, "output_dir", "", "Directory to write the created schema file")
}
