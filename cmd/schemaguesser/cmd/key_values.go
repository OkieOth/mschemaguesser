package cmd

import (
	"fmt"
	"log"
	"slices"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"okieoth/schemaguesser/internal/pkg/meta"
	"okieoth/schemaguesser/internal/pkg/mongoHelper"
	"okieoth/schemaguesser/internal/pkg/progressbar"
	"okieoth/schemaguesser/internal/pkg/utils"

	"github.com/spf13/cobra"
)

var keyValuesCmd = &cobra.Command{
	Use:   "key_values",
	Short: "dump the values of assumed key field to a text file",
	Long:  "With this command you can dump the data of considered key fields from the collections. Potential key fields are '_id', UUIDs or string in the UUID format. The received data are stored in a folder structure by database and collection. Every collection folder contains then the files with the field data (new line separated)",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := mongoHelper.Connect(mongoHelper.ConStr)
		if err != nil {
			msg := fmt.Sprintf("Failed to connect to db: %v", err)
			panic(msg)
		}
		defer mongoHelper.CloseConnection(client)

		if databaseName == "all" {
			keyValuesForAllDatabases(client, true)
		} else {
			if collectionName == "all" {
				keyValuesForAllCollections(client, databaseName, true)
			} else {
				keyValuesForOneCollection(client, databaseName, collectionName, false, true)
			}
		}

	},
}

func keyValuesForOneCollection(client *mongo.Client, dbName string, collName string, doRecover bool, initProgressBar bool) {
	// defer func() {
	// 	if doRecover {
	// 		if r := recover(); r != nil {
	// 			log.Printf("Recovered while handling collection (db: %s, collection: %s): %v", dbName, collName, r)
	// 		}
	// 	}
	// }()
	if initProgressBar {
		descr := fmt.Sprintf("JSON export of %s:%s", dbName, collName)
		progressbar.Init(1, descr)
	}

	outputFile, err := utils.CreateOutputFile(outputDir, "key-values.json", dbName, collName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	startTime := time.Now()
	count := uint64(0)
	err = queryCollection(client, dbName, collName, func(data bson.Raw) error {
		mongoHelper.ScanBsonForKeyValues(data, dbName, collName, outputFile)
		if err != nil {
			log.Printf("[%s:%s] Error while scanning for key values: %v", dbName, collName, err)
			return err
		}
		count++
		return nil
	})
	if err := meta.WriteMetaInfo(outputDir, dbName, collName, count, "", nil); err != nil {
		panic(err)
	}
	log.Printf("[%s:%s] Key values persisted (count = %d) in %v\n", dbName, collName, count, time.Since(startTime))
}

func keyValuesForAllCollections(client *mongo.Client, dbName string, initProgressBar bool) {
	collections := getAllCollectionsOrPanic(client, dbName)
	var wg sync.WaitGroup
	if initProgressBar {
		progressbar.Init(int64(len(collections)), "Key values export for all collections")
	}

	for _, coll := range collections {
		if slices.Contains(blacklist, coll) {
			log.Printf("[%s:%s] skip blacklisted collection\n", dbName, coll)
			continue
		}
		wg.Add(1)
		go func(s string) {
			startTime := time.Now()
			defer func() {
				log.Printf("[%s:%s] Key values export of collection in %v\n", dbName, s, time.Since(startTime))
				wg.Done()
				if initProgressBar {
					progressbar.ProgressOne()
				}
			}()
			keyValuesForOneCollection(client, dbName, s, true, false)
		}(coll)
	}
	wg.Wait()
}

func keyValuesForAllDatabases(client *mongo.Client, initProgressBar bool) {
	dbs := getAllDatabasesOrPanic(client)
	var wg sync.WaitGroup
	if initProgressBar {
		progressbar.Init(int64(len(dbs)), "Key values export for all databases")
	}
	for _, db := range dbs {
		if slices.Contains(blacklist, db) {
			log.Printf("[%s] skip blacklisted DB\n", db)
			continue
		}
		wg.Add(1)
		go func(s string) {
			startTime := time.Now()
			defer func() {
				log.Printf("[%s] Key values exported from DB in %v\n", s, time.Since(startTime))
				wg.Done()
				if initProgressBar {
					progressbar.ProgressOne()
				}
			}()
			keyValuesForAllCollections(client, s, false)
		}(db)
	}
	wg.Wait()
}
