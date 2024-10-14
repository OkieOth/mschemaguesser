package cmd

import (
	"fmt"
	"log"

	"slices"
	"sync"
	"time"

	linkshelper "okieoth/schemaguesser/internal/pkg/linksHelper"
	"okieoth/schemaguesser/internal/pkg/meta"
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
		metaInfos, err := meta.GetAllMetaInfos(keyValuesDir)
		if err != nil {
			panic(fmt.Sprintf("Error while retrieve all available meta infos in: %s - %v", keyValuesDir, err))
		}

		if databaseName == "all" {
			linksForAllDatabases(metaInfos, true)
		} else {
			if collectionName == "all" {
				linksForAllCollections(metaInfos, databaseName, true)
			} else {
				linksForOneCollection(metaInfos, databaseName, collectionName, false, true)
			}
		}

	},
}

func linksForOneCollection(metaInfos []meta.MetaInfo, dbName string, collName string, doRecover bool, initProgressBar bool) {
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
	keyValues, err := linkshelper.GetKeyValues(keyValuesDir, dbName, collName)
	if err != nil {
		log.Printf("[%s:%s] Error while reading key-values: %v", dbName, collName, err)
		return
	}

	lenMetaInfos := len(metaInfos)
	if lenMetaInfos > 1 {
		var wg sync.WaitGroup
		wg.Add(lenMetaInfos - 1)
		for _, metaInfo := range metaInfos {
			if (metaInfo.Db == dbName) && (metaInfo.Collection == collName) {
				continue
			}
			go func(mf meta.MetaInfo) {
				defer func() {
					wg.Done()
				}()
				for k, _ := range keyValues {
					_, err := linkshelper.FoundKeyValue(keyValuesDir, metaInfo.Db, metaInfo.Collection, k)
					if err != nil {
						log.Printf("[%s:%s] Error while searching for value (%s) in %s:%s: %v", dbName, collName, k, metaInfo.Db, metaInfo.Collection, err)
					}
					// if len(deps) > 0 {
					// 	// TODO - do something with the found links
					// }
				}
			}(metaInfo)
		}

		wg.Wait()
	}
	log.Printf("[%s:%s] Links of collection are gathered in %v\n", dbName, collName, time.Since(startTime))
	if initProgressBar {
		progressbar.ProgressOne()
	}
}

func linksForAllCollections(metaInfos []meta.MetaInfo, dbName string, initProgressBar bool) {
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
			linksForOneCollection(metaInfos, dbName, s, true, false)
		}(coll)
	}
}

func linksForAllDatabases(metaInfos []meta.MetaInfo, initProgressBar bool) {
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
			linksForAllCollections(metaInfos, s, false)
		}(db)
	}
}
