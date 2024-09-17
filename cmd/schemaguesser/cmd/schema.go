package cmd

import (
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"

	"okieoth/schemaguesser/internal/pkg/mongoHelper"
	"okieoth/schemaguesser/internal/pkg/schema"

	"github.com/spf13/cobra"
)

var dbName string

var colName string

var outputDir string

var itemCount int32

var plantUml bool
var plantUmlDest string

var raw bool
var rawDest string

var data bool
var dataDest string

var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "functions around the schemas",
	Long:  "With this command you can create schemas out of mongodb collection",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := mongoHelper.Connect(mongoHelper.ConStr)
		if err != nil {
			msg := fmt.Sprintf("Failed to connect to db: %v", err)
			panic(msg)
		}
		defer mongoHelper.CloseConnection(client)

		if dbName == "all" {
			printSchemasForAllDatabases(client)
		} else {
			if colName == "all" {
				printSchemasForAllCollections(client, dbName)
			} else {
				printSchemaForOneCollection(client, dbName, colName, false)
			}
		}

	},
}

func init() {
	schemaCmd.Flags().StringVar(&dbName, "database", "all", "Database to query existing collections")

	schemaCmd.Flags().StringVar(&colName, "collection", "all", "Name of the collection to show the indexes")

	schemaCmd.Flags().StringVar(&outputDir, "output", "stdout", "stdout or the directory to write the created schema file, default is 'stdout'")

	schemaCmd.Flags().Int32Var(&itemCount, "item_count", 100, "Number of collection entries used to build the schema")

	schemaCmd.Flags().BoolVar(&plantUml, "plantuml", false, "If this flag is set, are to the schema additional plant uml diagrams created")

	schemaCmd.Flags().StringVar(&plantUmlDest, "pumldir", "", "This flag specifies the target directory for the PlantUml diagrams")

	schemaCmd.Flags().BoolVar(&raw, "raw", false, "If this flag is set, the collected and aggregated schema data are persisted too")

	schemaCmd.Flags().StringVar(&rawDest, "rawdir", "", "Destination directory for the raw schema data")

	schemaCmd.Flags().BoolVar(&data, "data", false, "If this flag is set, the read mongo data are persisted too")

	schemaCmd.Flags().StringVar(&dataDest, "datadir", "", "Destination directory for the read database content")
}

func printSchemaForOneCollection(client *mongo.Client, dbName string, collName string, doRecover bool) {
	defer func() {
		if doRecover {
			if r := recover(); r != nil {
				fmt.Printf("Recovered while handling collection (db: %s, collection: %s): %v", dbName, collName, r)
			}
		}
	}()
	bsonRaw, err := mongoHelper.QueryCollection(client, dbName, collName, int(itemCount))
	if err != nil {
		msg := fmt.Sprintf("Error while reading data for collection (%s.%s): \n%v\n", dbName, collName, err)
		panic(msg)
	}
	var otherComplexTypes []mongoHelper.ComplexType
	var mainType mongoHelper.ComplexType
	for _, b := range bsonRaw {
		err = mongoHelper.ProcessBson(b, collName, &mainType, &otherComplexTypes)
		if err != nil {
			fmt.Printf("Error while processing bson for schema: %v", err)
		}
	}
	if len(bsonRaw) > 0 {
		schema.ReduceTypes(&mainType, &otherComplexTypes)
		//schema.GuessDicts(&otherComplexTypes)
		schema.PrintSchema(dbName, collName, &mainType, &otherComplexTypes, outputDir)
	} else {
		fmt.Printf("No data for database: %s, collection: %s\n", dbName, collName)
	}
}

func printSchemasForAllCollections(client *mongo.Client, dbName string) {
	collections := mongoHelper.ReadCollectionsOrPanic(client, dbName)
	var wg sync.WaitGroup
	for _, coll := range *collections {
		wg.Add(1)
		go func(s string) {
			defer wg.Done()
			printSchemaForOneCollection(client, dbName, s, true)
		}(coll)
	}
	wg.Wait()
}

func printSchemasForAllDatabases(client *mongo.Client) {
	dbs := mongoHelper.ReadDatabasesOrPanic(client)
	var wg sync.WaitGroup
	for _, db := range *dbs {
		wg.Add(1)
		go func(s string) {
			defer wg.Done()
			printSchemasForAllCollections(client, s)
		}(db)
	}
	wg.Wait()
}
