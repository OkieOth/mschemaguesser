package cmd

import (
	"context"
	"fmt"
	"log"
	"slices"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"okieoth/schemaguesser/internal/pkg/importHelper"
	"okieoth/schemaguesser/internal/pkg/mongoHelper"
	"okieoth/schemaguesser/internal/pkg/progressbar"

	"okieoth/schemaguesser/internal/pkg/utils"

	"github.com/spf13/cobra"
)

var inputDir string

var chunkSize int64

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import data to mongodb",
	Long:  `Based on a given mongodb connection you can import data from before stored BSON persistent files.`,
	Run: func(cmd *cobra.Command, args []string) {
		var client *mongo.Client
		var err error
		if !useDumps {
			client, err = mongoHelper.Connect(mongoHelper.ConStr)
			if err != nil {
				msg := fmt.Sprintf("Failed to connect to db: %v", err)
				panic(msg)
			}
			defer mongoHelper.CloseConnection(client)
		}

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

	importCmd.Flags().StringVarP(&databaseName, "database", "d", "all", "Database for the import")

	importCmd.Flags().StringVarP(&collectionName, "collection", "c", "all", "Name of the collection to import")

	importCmd.Flags().StringVar(&inputDir, "input", "", "The directory where the BSON exports can be found. This files need to be created with the 'get bson' commands of this tool")

	importCmd.Flags().StringSliceVarP(&blacklist, "blacklist", "b", []string{}, "Blacklist names to skip")

	importCmd.Flags().Int64Var(&chunkSize, "chunk_size", 100, "Chunk size to use for the imports, default is 100")

}

func importOneCollection(client *mongo.Client, dbName string, collName string, doRecover bool, initProgressBar bool) {
	defer func() {
		if doRecover {
			if r := recover(); r != nil {
				log.Printf("Recovered while handling collection (db: %s, collection: %s): %v", dbName, collName, r)
			}
		}
	}()

	importFile := utils.GetFileName(inputDir, "bson", dbName, collName)

	var ctx context.Context
	if timeout > 0 {
		c, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
		ctx = c
		defer cancel()
	} else {
		ctx = context.Background()
	}

	itemCount, err := importHelper.ImportData(client, importFile, dbName, collName, chunkSize, &ctx)
	if err != nil {
		log.Printf("[%s:%s] Error while importing data: %v\n", dbName, collName, err)
	} else {
		log.Printf("[%s:%s] %d items successfully imported\n", dbName, collName, itemCount)
	}
}

func importAllCollections(client *mongo.Client, dbName string, initProgressBar bool) {
	collections, err := importHelper.AllCollectionsForDb(inputDir, dbName)
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	if initProgressBar {
		progressbar.Init(int64(len(collections)), "BSON import for all collections")
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
				log.Printf("[%s:%s] BSON import of collection in %v\n", dbName, s, time.Since(startTime))
				wg.Done()
				if initProgressBar {
					progressbar.ProgressOne()
				}
			}()
			importOneCollection(client, dbName, s, true, false)
		}(coll)
	}
	wg.Wait()
}

func importAllDatabases(client *mongo.Client, initProgressBar bool) {
	dbs, err := importHelper.AllDatabases(inputDir)
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	if initProgressBar {
		progressbar.Init(int64(len(dbs)), "BSON import for all databases")
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
				log.Printf("[%s] BSON imported into DB in %v\n", s, time.Since(startTime))
				wg.Done()
				if initProgressBar {
					progressbar.ProgressOne()
				}
			}()
			importAllCollections(client, s, false)
		}(db)
	}
	wg.Wait()
}
