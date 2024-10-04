package cmd

import (
	"fmt"
	"log"
	"slices"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"okieoth/schemaguesser/internal/pkg/mongoHelper"
	"okieoth/schemaguesser/internal/pkg/progressbar"
	"okieoth/schemaguesser/internal/pkg/schema"

	"github.com/spf13/cobra"
)

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

		if databaseName == "all" {
			printSchemasForAllDatabases(client, true)
		} else {
			if collectionName == "all" {
				printSchemasForAllCollections(client, databaseName, true)
			} else {
				printSchemaForOneCollection(client, databaseName, collectionName, false, true)
			}
		}
	},
}

var includeCount bool
var commentPotentialKeyFields bool
var persistKeyValues bool
var persistKeyValuesDir string
var keyUuid bool
var keyUuidString bool
var useZeroKeyUuid bool

func init() {
	schemaCmd.Flags().BoolVar(&includeCount, "include_count", false, "If set it includes the current number of elements of the collection into schema comments")
	schemaCmd.Flags().BoolVar(&commentPotentialKeyFields, "key_fields", false, "If set it annotates potential key fields in the schema with a comment. Without additional flags only the fields of type 'objectId' are considered as keys")
	schemaCmd.Flags().BoolVar(&persistKeyValues, "persist_key_values", false, "If set the unique key values are extracted from the sample data and stored in separate files")
	schemaCmd.Flags().StringVar(&persistKeyValuesDir, "key_values_dir", "", "Optional output dir to store the files with the key values. If 'persist_key_values' is set and this flag is empty, then the output dir is used")

	schemaCmd.Flags().BoolVar(&keyUuid, "uuid_keys", false, "If set, binary uuid fields are considered as key, too")
	schemaCmd.Flags().BoolVar(&keyUuid, "uuid_str_keys", false, "If set, uuids in string format (e.g. '056bcf58-e17e-42ba-8186-f25ffbde8b35') are considered as key, too")
	schemaCmd.Flags().BoolVar(&keyUuid, "zero_uuid_keys", false, "Per default zero uuids (e.g. '00000000-0000-0000-0000-000000000000') are ignored, use the switch to integrate them as values when found")
}

func getDocumentCount(client *mongo.Client, dbName string, collName string, mt *mongoHelper.ComplexType) {
	startTime := time.Now()
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered while counting collection (db: %s, collection: %s): %v", dbName, collName, r)
		}
	}()
	if mt == nil {
		log.Printf("[%s:%s] getDocumentCount: mainType pointer is nil, skip count", dbName, collName)
		return
	}
	defer func() {
		log.Printf("[%s:%s] count call finished in %v\n", dbName, collName, time.Since(startTime))
	}()
	log.Printf("[%s:%s] count ...\n", dbName, collName)
	count, err := mongoHelper.CountCollection(client, dbName, collName)
	if err != nil {
		msg := fmt.Sprintf("[%s:%s] error while count elements: %v", dbName, collName, err)
		log.Println(msg)
		mt.Comments = append(mt.Comments, msg)
	} else {
		mt.Count.Set(count)
	}
}

func printSchemaForOneCollection(client *mongo.Client, dbName string, collName string, doRecover bool, initProgressBar bool) {
	defer func() {
		if doRecover {
			if r := recover(); r != nil {
				log.Printf("Recovered while handling collection (db: %s, collection: %s): %v", dbName, collName, r)
			}
		}
		if initProgressBar {
			progressbar.ProgressOne()
		}
	}()
	if initProgressBar {
		descr := fmt.Sprintf("Schema for %s:%s", dbName, collName)
		progressbar.Init(1, descr)
	}
	var otherComplexTypes []mongoHelper.ComplexType
	var mainType mongoHelper.ComplexType

	if includeCount {
		getDocumentCount(client, dbName, collName, &mainType)
	}

	bsonRaw := make([]bson.Raw, 0)
	err := mongoHelper.QueryCollection(client, dbName, collName, int(itemCount), useAggregation, mongoV44, func(data bson.Raw) error {
		bsonRaw = append(bsonRaw, data)
		return nil
	})

	if err != nil {
		msg := fmt.Sprintf("Error while reading data for collection (%s.%s): \n%v\n", dbName, collName, err)
		panic(msg)
	}
	startTime := time.Now()
	for _, b := range bsonRaw {
		err = mongoHelper.ProcessBson(b, collName, &mainType, &otherComplexTypes)
		if err != nil {
			log.Printf("Error while processing bson for schema: %v", err)
		}
	}
	log.Printf("[%s:%s] Mongodb data processed for collection in %v\n", dbName, collName, time.Since(startTime))
	if len(bsonRaw) > 0 {
		schema.ReduceTypes(&mainType, &otherComplexTypes)
		//schema.GuessDicts(&otherComplexTypes)
		schema.PrintSchema(dbName, collName, &mainType, &otherComplexTypes, outputDir)
		log.Printf("[%s:%s] Schema printed in %v\n", dbName, collName, time.Since(startTime))
	} else {
		log.Printf("No data for database: %s, collection: %s\n", dbName, collName)
	}
}

func printSchemasForAllCollections(client *mongo.Client, dbName string, initProgressBar bool) {
	collections := mongoHelper.ReadCollectionsOrPanic(client, dbName)
	var wg sync.WaitGroup
	if initProgressBar {
		progressbar.Init(int64(len(collections)), "Schema for all collections")
	}

	if includeCount {
		wg.Add(len(collections))
	}
	for _, coll := range collections {
		if slices.Contains(blacklist, coll) {
			log.Printf("[%s:%s] skip blacklisted collection\n", dbName, coll)
			continue
		}
		go func(s string) {
			startTime := time.Now()
			defer func() {
				log.Printf("[%s:%s] Schema created for collection in %v\n", dbName, s, time.Since(startTime))
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
		progressbar.Init(int64(len(dbs)), "Schema for all databases")
	}
	wg.Add(len(dbs))
	for _, db := range dbs {
		if slices.Contains(blacklist, db) {
			log.Printf("[%s] skip blacklisted DB\n", db)
			continue
		}
		go func(s string) {
			startTime := time.Now()
			defer func() {
				log.Printf("[%s] Schemas created for DB in %v\n", s, time.Since(startTime))
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
