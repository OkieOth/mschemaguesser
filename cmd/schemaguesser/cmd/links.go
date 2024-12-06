package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"slices"
	"sync"
	"time"

	linkshelper "okieoth/schemaguesser/internal/pkg/linksHelper"
	"okieoth/schemaguesser/internal/pkg/meta"
	"okieoth/schemaguesser/internal/pkg/progressbar"
	"okieoth/schemaguesser/internal/pkg/utils"

	"github.com/spf13/cobra"
)

var keyValuesDir string
var ignoreSameAttribRefs bool
var whitelist []string

func init() {
	linksCmd.Flags().StringVar(&keyValuesDir, "key_values_dir", "", "Directory where the previously dumped key values of the databases and collections can be found")
	linksCmd.Flags().BoolVar(&ignoreSameAttribRefs, "ignore_same_attrib_refs", false, "Since the algorithm by default also catches foreign-key-2-foreign-key relations and they are often behind the use of the same attrib name, this flag can simplify the situation")
	linksCmd.PersistentFlags().StringSliceVar(&whitelist, "whitelist", []string{}, "Whitelist attribute names to find the relations")
	linksCmd.MarkPersistentFlagRequired("key_values_dir")
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

		colRefs := make([]linkshelper.ColRefs, 0)

		if databaseName == "all" {
			colRefs = linksForAllDatabases(metaInfos, colRefs, true)
		} else {
			if collectionName == "all" {
				colRefs = linksForAllCollections(metaInfos, colRefs, databaseName, true)
			} else {
				colRefs = linksForOneCollection(metaInfos, colRefs, databaseName, collectionName, false, true)
			}
		}
		colRefs = linkshelper.AggregateRefs(colRefs)
		printLinks(colRefs)
	},
}

func printLinks(colRefs []linkshelper.ColRefs) {
	log.Printf("Found references for %d collections\n", len(colRefs))
	if outputDir == "stdout" {
		outputDir = "."
	}
	outputFile, err := utils.CreateOutputFile(outputDir, "json", "db", "references")
	if err != nil {
		log.Fatalf("Error while creating output file: %v", err)
	}
	defer outputFile.Close()
	jsonData, err := json.MarshalIndent(colRefs, "", "  ")
	if err != nil {
		log.Fatalf("Error while marshalling JSON data: %v", err)
	}
	_, err = outputFile.Write(jsonData)
	if err != nil {
		log.Fatalf("Failed to write database references to file: %v", err)
	}
}

func linksForOneCollection(metaInfos []meta.MetaInfo, colRefs []linkshelper.ColRefs, dbName string, collName string, doRecover bool, initProgressBar bool) []linkshelper.ColRefs {
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
	keyValues, err := linkshelper.GetKeyValues(keyValuesDir, dbName, collName, whitelist)
	if err != nil {
		log.Printf("[%s:%s] Error while reading key-values: %v", dbName, collName, err)
		return colRefs
	}

	if len(keyValues) == 0 {
		// there are no key values (maybe they are ignored for the white list) in related files
		return colRefs
	}

	lenMetaInfos := len(metaInfos)
	if lenMetaInfos > 1 {
		var wg sync.WaitGroup
		wg.Add(lenMetaInfos - 1)

		var waitForCollectRefsChannel sync.WaitGroup
		collectRefsChannel := make(chan linkshelper.ColRefs)

		waitForCollectRefsChannel.Add(1)
		go func(chOut <-chan linkshelper.ColRefs) {
			defer func() {
				waitForCollectRefsChannel.Done()
			}()
			for v := range chOut {
				log.Printf("Received Ref: %v", v)
				colRefs = append(colRefs, v)
			}
		}(collectRefsChannel)

		for _, metaInfo := range metaInfos {
			if (metaInfo.Db == dbName) && (metaInfo.Collection == collName) {
				continue
			}
			go func(mf meta.MetaInfo, chIn chan<- linkshelper.ColRefs) {
				defer func() {
					wg.Done()
				}()
				for k, v := range keyValues {
					found, err := linkshelper.FindKeyValues(keyValuesDir, mf.Db, mf.Collection, k, v, dbName, collName, chIn, ignoreSameAttribRefs)
					if err != nil {
						log.Printf("[%s:%s] Error while searching for value (%s) in %s:%s: %v", dbName, collName, k, mf.Db, mf.Collection, err)
					}
					if found {
						// if at least one value found, then the dependency is considered as proved
						log.Printf("[%s:%s] Found reference, skip additional tries", mf.Db, mf.Collection)
						break
					}
				}
			}(metaInfo, collectRefsChannel)
		}
		wg.Wait()
		close(collectRefsChannel)
		waitForCollectRefsChannel.Wait()
	}
	log.Printf("[%s:%s] Links of collection are gathered in %v\n", dbName, collName, time.Since(startTime))
	if initProgressBar {
		progressbar.ProgressOne()
	}
	return colRefs
}

func linksForAllCollections(metaInfos []meta.MetaInfo, colRefs []linkshelper.ColRefs, dbName string, initProgressBar bool) []linkshelper.ColRefs {
	collections := getAllCollectionsOrPanic(nil, keyValuesDir, true, dbName)
	if initProgressBar {
		progressbar.Init(int64(len(collections)), "Links for all collections")
	}

	for _, coll := range collections {
		if slices.Contains(blacklist, coll) {
			log.Printf("[%s:%s] skip blacklisted collection\n", dbName, coll)
			continue
		} else {
			colRefs = linksForOneCollection(metaInfos, colRefs, dbName, coll, true, false)
			if initProgressBar {
				progressbar.ProgressOne()
			}
		}
	}
	return colRefs
}

func linksForAllDatabases(metaInfos []meta.MetaInfo, colRefs []linkshelper.ColRefs, initProgressBar bool) []linkshelper.ColRefs {
	dbs := getAllDatabasesOrPanic(nil, keyValuesDir, true)
	if initProgressBar {
		progressbar.Init(int64(len(dbs)), "Links for all databases")
	}
	for _, db := range dbs {
		if slices.Contains(blacklist, db) {
			log.Printf("[%s] skip blacklisted DB\n", db)
			continue
		}
		startTime := time.Now()
		colRefs = linksForAllCollections(metaInfos, colRefs, db, false)
		if initProgressBar {
			progressbar.ProgressOne()
		}
		log.Printf("[%s] Links for DB in %v\n", db, time.Since(startTime))
	}
	return colRefs
}
