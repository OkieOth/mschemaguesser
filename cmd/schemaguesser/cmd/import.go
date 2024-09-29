package cmd

import (
	"context"
	"fmt"
	"log"
	"slices"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"okieoth/schemaguesser/internal/pkg/mongoHelper"
	"okieoth/schemaguesser/internal/pkg/progressbar"

	"okieoth/schemaguesser/internal/pkg/utils"

	"github.com/spf13/cobra"
)

var inputDir string

var bulkSize int16

var importCmd = &cobra.Command{
	Use:   "get",
	Short: "Import data to mongodb",
	Long:  `Based on a given mongodb connection you can import data from before stored BSON persistent files.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := mongoHelper.Connect(mongoHelper.ConStr)
		if err != nil {
			msg := fmt.Sprintf("Failed to connect to db: %v", err)
			panic(msg)
		}
		defer mongoHelper.CloseConnection(client)

		if databaseName == "all" {
			importAllDatabases(client, true)
		} else {
			if collectionName == "all" {
				importAllCollections(client, databaseName, true)
			} else {
				importOneCollection(client, databaseName, collectionName, false, true)
			}
		}
	},
}

func init() {

	getCmd.Flags().StringVarP(&databaseName, "database", "d", "all", "Database for the import")

	getCmd.Flags().StringVarP(&collectionName, "collection", "c", "all", "Name of the collection to import")

	getCmd.Flags().StringVar(&inputDir, "input", "", "The directory where the BSON exports can be found. This files need to be created with the 'get bson' commands of this tool")
}

func importOneCollection(client *mongo.Client, dbName string, collName string, doRecover bool, initProgressBar bool) {
	defer func() {
		if doRecover {
			if r := recover(); r != nil {
				log.Printf("Recovered while handling collection (db: %s, collection: %s): %v", dbName, collName, r)
			}
		}
	}()

	outputFile, err := utils.CreateOutputFile(outputDir, "bson", dbName, collName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	var ctx context.Context
	if timeout > 0 {
		c, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
		ctx = c
		defer cancel()
	} else {
		ctx = context.Background()
	}

	mongoHelper.DumpCollectionToFile(ctx, outputFile, client, dbName, collName, itemCount, useAggregation, mongoV44)
}

func importAllCollections(client *mongo.Client, dbName string, initProgressBar bool) {
	collections := mongoHelper.ReadCollectionsOrPanic(client, dbName)
	var wg sync.WaitGroup
	if initProgressBar {
		progressbar.Init(int64(len(*collections)), "BSON export for all collections")
	}

	for _, coll := range *collections {
		if slices.Contains(blacklist, coll) {
			log.Printf("[%s:%s] skip blacklisted collection\n", dbName, coll)
			continue
		}
		wg.Add(1)
		go func(s string) {
			startTime := time.Now()
			defer func() {
				log.Printf("[%s:%s] BSON export of collection in %v\n", dbName, s, time.Since(startTime))
				wg.Done()
				if initProgressBar {
					progressbar.ProgressOne()
				}
			}()
			bsonForOneCollection(client, dbName, s, true, false)
		}(coll)
	}
	wg.Wait()
}

func importAllDatabases(client *mongo.Client, initProgressBar bool) {
	dbs := mongoHelper.ReadDatabasesOrPanic(client)
	var wg sync.WaitGroup
	if initProgressBar {
		progressbar.Init(int64(len(*dbs)), "BSON export for all databases")
	}
	for _, db := range *dbs {
		if slices.Contains(blacklist, db) {
			log.Printf("[%s] skip blacklisted DB\n", db)
			continue
		}
		wg.Add(1)
		go func(s string) {
			startTime := time.Now()
			defer func() {
				log.Printf("[%s] BSON exported from DB in %v\n", s, time.Since(startTime))
				wg.Done()
				if initProgressBar {
					progressbar.ProgressOne()
				}
			}()
			bsonForAllCollections(client, s, false)
		}(db)
	}
	wg.Wait()
}
