package cmd

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"slices"

	"okieoth/schemaguesser/internal/pkg/importHelper"
	"okieoth/schemaguesser/internal/pkg/mongoHelper"
	"okieoth/schemaguesser/internal/pkg/utils"

	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"
)

var databaseName string

var collectionName string

var outputDir string

var itemCount int64

var blacklist []string

var useAggregation bool

var mongoV44 bool

var timeout int64

var useDumps bool

var dumpDir string

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Retrieve information out of mongodb",
	Long: `Based on a given mongodb connection you can extract data and create
                different outputs out of it.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("use this command with: help|schema|bson|json, puml")
	},
}

func init() {
	getCmd.AddCommand(schemaCmd)
	getCmd.AddCommand(bsonCmd)
	getCmd.AddCommand(jsonCmd)
	getCmd.AddCommand(keyValuesCmd)
	getCmd.AddCommand(linksCmd)

	getCmd.PersistentFlags().StringVarP(&databaseName, "database", "d", "all", "Database to query existing collections")

	getCmd.PersistentFlags().StringVarP(&collectionName, "collection", "c", "all", "Name of the collection to handle")

	getCmd.PersistentFlags().StringVarP(&outputDir, "output", "o", "stdout", "The directory to write the created schema file")

	getCmd.PersistentFlags().Int64VarP(&itemCount, "item_count", "i", 100, "Number of collection entries used to build the schema")

	getCmd.PersistentFlags().Int64VarP(&timeout, "timeout", "t", 30, "Timeout seconds of database queries. The default is 30s. In case you don't want any timeout, set the value to 0")

	getCmd.PersistentFlags().StringSliceVarP(&blacklist, "blacklist", "b", []string{}, "Blacklist names to skip")

	getCmd.PersistentFlags().BoolVar(&useAggregation, "use_aggregation", false, "Use an aggregation pipeline to query the collections, this allows to enable the disk use for sorting also in mongo < 4.4")

	getCmd.PersistentFlags().BoolVar(&mongoV44, "mongo_v44", false, "The connection is to a mongodb newer than 4.4, enables additional driver features")

	getCmd.PersistentFlags().BoolVar(&useDumps, "use_dumps", false, "This flag allows to use before dumped bson data (by the use of the `get bson`command)")

	getCmd.PersistentFlags().StringVar(&dumpDir, "dump_dir", "", "The directory where the dumps to use, can be found")

}

func getAllDatabasesOrPanic(client *mongo.Client, dirWithMetaFiles string, useMetaFiles bool) []string {
	if useMetaFiles {
		if dirWithMetaFiles == "" {
			panic("no directory given to find the needed meta files for offline processing, so no idea from where to get the data")
		}
		ret, err := importHelper.AllDatabases(dirWithMetaFiles)
		if err != nil {
			log.Fatalf("error while reading dbs from 'dump_dir': %v", err)
		}
		return ret
	} else {
		if client == nil {
			panic("mongo client not initialized to query databases")
		}
		return mongoHelper.ReadDatabasesOrPanic(client)
	}
}

func getAllCollectionsOrPanic(client *mongo.Client, dirWithMetaFiles string, useMetaFiles bool, dbName string) []string {
	if useMetaFiles {
		if dirWithMetaFiles == "" {
			panic("no directory given to find the needed meta files for offline processing, so no idea from where to get the data")
		}
		ret, err := importHelper.AllCollectionsForDb(dirWithMetaFiles, dbName)
		if err != nil {
			log.Fatalf("error while reading collections for db (%s) from 'dump_dir': %v", dbName, err)
		}
		return ret
	} else {
		if client == nil {
			panic("mongo client not initialized to query databases")
		}
		return mongoHelper.ReadCollectionsOrPanic(client, dbName)
	}
}

func queryCollection(client *mongo.Client, dbName string, collName string, callback mongoHelper.HandleDataCallback) error {
	if useDumps {
		if dumpDir == "" {
			panic(fmt.Sprintf("queryCollection - [%s:%s] no 'dump_dir' flag given, so no idea from where to get the data", dbName, collName))
		}
		importFile := utils.GetFileName(dumpDir, "bson", dbName, collName)
		return getCollectionFromLocalFile(importFile, callback)
	} else {
		if client == nil {
			panic("mongo client not initialized to query databases")
		}
		return mongoHelper.QueryCollection(client, dbName, collName, int(itemCount), useAggregation, mongoV44, callback)
	}
}

func getCollectionFromLocalFile(importFile string, callback mongoHelper.HandleDataCallback) error {
	file, err := os.Open(importFile)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	buf := make([]byte, 4)
	readCount := uint64(0)
	for {
		_, err := io.ReadFull(file, buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to read document size to buffer: %v, readCount: %d", err, readCount)
		}
		docLength := int32(binary.LittleEndian.Uint32(buf))
		docBuf := make([]byte, docLength)
		copy(docBuf, buf)
		_, err = io.ReadFull(file, docBuf[4:])
		if err != nil {
			return fmt.Errorf("failed to read document to buffer: %v, readCount: %d", err, readCount)
		}
		readCount++
		err = callback(docBuf)
		if err != nil {
			return fmt.Errorf("failed to call callback: %v, readCount: %d", err, readCount)
		}
	}
	return nil
}

func removeBlacklisted(collections []string, blacklist []string) []string {
	ret := make([]string, 0)
	for _, c := range collections {
		if slices.Contains(blacklist, c) {
			log.Printf("Skip blacklisted collection: %s\n", c)
			continue
		}
		ret = append(ret, c)
	}
	return ret
}
