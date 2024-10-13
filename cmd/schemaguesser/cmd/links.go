package cmd

import (
	"fmt"
	"log"

	"slices"
	"time"

	linkshelper "okieoth/schemaguesser/internal/pkg/linksHelper"
	"okieoth/schemaguesser/internal/pkg/progressbar"

	"github.com/spf13/cobra"
)

var keyValuesDir string

func init() {
	linksCmd.Flags().StringVar(&keyValuesDir, "key_values_dir", "", "Directory where the previously dumped key values of the databases and collections can be found")
}

var linksCmd = &cobra.Command{
	Use:   "links",
	Short: "Search for ID links between collections in before persisted key values",
	Long:  "With this command you can search for collection links between ID fields (objectId, uuid or strings in uuid format).",
	Run: func(cmd *cobra.Command, args []string) {
		if databaseName == "all" {
			linksForAllDatabases(true)
		} else {
			if collectionName == "all" {
				linksForAllCollections(databaseName, true)
			} else {
				linksForOneCollection(databaseName, collectionName, false, true)
			}
		}

	},
}

func linksForOneCollection(dbName string, collName string, doRecover bool, initProgressBar bool) {
	defer func() {
		if doRecover {
			if r := recover(); r != nil {
				log.Printf("Recovered while handling collection (db: %s, collection: %s): %v", dbName, collName, r)
			}
		}
	}()
	if initProgressBar {
		descr := fmt.Sprintf("Links for Collection %s:%s", dbName, collName)
		progressbar.Init(1, descr)
	}

	// TODO - output format not decided yet
	// outputFile, err := utils.CreateOutputFile(outputDir, "json", dbName, collName)
	// if err != nil {
	// 	panic(err)
	// }
	// defer outputFile.Close()

	startTime := time.Now()
	_, err := linkshelper.GetKeyValues(keyValuesDir, dbName, collName)
	if err != nil {
		log.Printf("[%s:%s] Error while reading key-values: %v", dbName, collName, err)
		return
	}

	// i := 0

	// utils.DumpBytesToFile([]byte("["), outputFile)
	// err = queryCollection(client, dbName, collName, func(data bson.Raw) error {
	// 	bytes, err := getJsonBytes(&data)
	// 	if err != nil {
	// 		log.Printf("Error while converting to JSON: %v", err)
	// 		return err
	// 	}
	// 	if i > 0 {
	// 		utils.DumpBytesToFile([]byte(","), outputFile)
	// 	}
	// 	utils.DumpBytesToFile(bytes, outputFile)
	// 	utils.DumpBytesToFile([]byte("\n"), outputFile)
	// 	i++
	// 	return nil // TODO
	// })
	// utils.DumpBytesToFile([]byte("]"), outputFile)

	// if err != nil {
	// 	msg := fmt.Sprintf("Error while reading data for collection (%s.%s): \n%v\n", dbName, collName, err)
	// 	panic(msg)
	// }
	log.Printf("[%s:%s] Links of collection are gathered in %v\n", dbName, collName, time.Since(startTime))
	if initProgressBar {
		progressbar.ProgressOne()
	}
}

func linksForAllCollections(dbName string, initProgressBar bool) {
	collections := getAllCollectionsOrPanic(nil, keyValuesDir, true, dbName)
	if initProgressBar {
		progressbar.Init(int64(len(collections)), "Links for all collections")
	}

	for _, coll := range collections {
		if slices.Contains(blacklist, coll) {
			log.Printf("[%s:%s] skip blacklisted collection\n", dbName, coll)
			continue
		}
		go func(s string) {
			defer func() {
				if initProgressBar {
					progressbar.ProgressOne()
				}
			}()
			linksForOneCollection(dbName, s, true, false)
		}(coll)
	}
}

func linksForAllDatabases(initProgressBar bool) {
	dbs := getAllDatabasesOrPanic(nil, keyValuesDir, true)
	if initProgressBar {
		progressbar.Init(int64(len(dbs)), "Links for all databases")
	}
	for _, db := range dbs {
		if slices.Contains(blacklist, db) {
			log.Printf("[%s] skip blacklisted DB\n", db)
			continue
		}
		go func(s string) {
			startTime := time.Now()
			defer func() {
				log.Printf("[%s] Links for DB in %v\n", s, time.Since(startTime))
				if initProgressBar {
					progressbar.ProgressOne()
				}
			}()
			linksForAllCollections(s, false)
		}(db)
	}
}
