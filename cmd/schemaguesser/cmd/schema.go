package cmd

import (
	"fmt"
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"

	"okieoth/schemaguesser/internal/pkg/mongoHelper"
	"okieoth/schemaguesser/internal/pkg/progressbar"
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
		client, err := mongoHelper.Connect(mongoHelper.ConStr)
		if err != nil {
			msg := fmt.Sprintf("Failed to connect to db: %v", err)
			panic(msg)
		}
		defer mongoHelper.CloseConnection(client)

		if dbName == "all" {
			printSchemasForAllDatabases(client, true)
		} else {
			if colName == "all" {
				printSchemasForAllCollections(client, dbName, true)
			} else {
				printSchemaForOneCollection(client, dbName, colName, false, true)
			}
		}

	},
}

func init() {
	schemaCmd.Flags().StringVar(&dbName, "database", "all", "Database to query existing collections")

	schemaCmd.Flags().StringVar(&colName, "collection", "all", "Name of the collection to show the indexes")

	schemaCmd.Flags().StringVar(&outputDir, "output", "stdout", "stdout or the directory to write the created schema file, default is 'stdout'")

	schemaCmd.Flags().Int32Var(&itemCount, "item_count", 100, "Number of collection entries used to build the schema")
}

func printSchemaForOneCollection(client *mongo.Client, dbName string, collName string, doRecover bool, initProgressBar bool) {
	defer func() {
		if doRecover {
			if r := recover(); r != nil {
				log.Printf("Recovered while handling collection (db: %s, collection: %s): %v", dbName, collName, r)
			}
		}
	}()
	if initProgressBar {
		progressbar.Init(1, "Schema for one collection")
	}
	bsonRaw, err := mongoHelper.QueryCollectionWithAggregation(client, dbName, collName, int(itemCount))
	if err != nil {
		msg := fmt.Sprintf("Error while reading data for collection (%s.%s): \n%v\n", dbName, collName, err)
		panic(msg)
	}
	var otherComplexTypes []mongoHelper.ComplexType
	var mainType mongoHelper.ComplexType
	for _, b := range bsonRaw {
		err = mongoHelper.ProcessBson(b, collName, &mainType, &otherComplexTypes)
		if err != nil {
			log.Printf("Error while processing bson for schema: %v", err)
		}
	}
	if len(bsonRaw) > 0 {
		schema.ReduceTypes(&mainType, &otherComplexTypes)
		//schema.GuessDicts(&otherComplexTypes)
		schema.PrintSchema(dbName, collName, &mainType, &otherComplexTypes, outputDir)
	} else {
		log.Printf("No data for database: %s, collection: %s\n", dbName, collName)
	}
	if initProgressBar {
		progressbar.ProgressOne()
	}
}

func printSchemasForAllCollections(client *mongo.Client, dbName string, initProgressBar bool) {
	collections := mongoHelper.ReadCollectionsOrPanic(client, dbName)
	var wg sync.WaitGroup
	if initProgressBar {
		progressbar.Init(int64(len(*collections)), "Schema for all collections")
	}

	for _, coll := range *collections {
		wg.Add(1)
		go func(s string) {
			defer func() {
				wg.Done()
				if initProgressBar {
					progressbar.ProgressOne()
				}
			}()
			printSchemaForOneCollection(client, dbName, s, true, false)
		}(coll)
	}
	wg.Wait()
}

func printSchemasForAllDatabases(client *mongo.Client, initProgressBar bool) {
	dbs := mongoHelper.ReadDatabasesOrPanic(client)
	var wg sync.WaitGroup
	if initProgressBar {
		progressbar.Init(int64(len(*dbs)), "Schema for all databases")
	}
	for _, db := range *dbs {
		wg.Add(1)
		go func(s string) {
			defer func() {
				wg.Done()
				if initProgressBar {
					progressbar.ProgressOne()
				}
			}()
			printSchemasForAllCollections(client, s, false)
		}(db)
	}
	wg.Wait()
}
