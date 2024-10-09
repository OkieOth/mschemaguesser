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
	bsonRaw := make([]bson.Raw, 0)
	i := 0
	err := queryCollection(client, dbName, collName, func(data bson.Raw) error {
		i++
		bsonRaw = append(bsonRaw, data)
		return nil
	})

	if i > 0 {
		if err != nil {
			msg := fmt.Sprintf("Error while reading data of collection (%s.%s): \n%v\n", dbName, collName, err)
			panic(msg)
		}
		startTime := time.Now()

		for _, b := range bsonRaw {
			mongoHelper.ScanBsonForKeyValues(b, dbName, collName, outputDir)
			if err != nil {
				log.Printf("[%s:%s] Error while scanning for key values: %v", dbName, collName, err)
			}
		}
		log.Printf("[%s:%s] Key values persisted in %v\n", dbName, collName, time.Since(startTime))
	} else {
		log.Printf("No data for database: %s, collection: %s\n", dbName, collName)
	}

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
