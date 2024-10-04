package cmd

import (
	"fmt"
	"log"

	"okieoth/schemaguesser/internal/pkg/importHelper"
	"okieoth/schemaguesser/internal/pkg/mongoHelper"

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

func getAllDatabasesOrPanic(client *mongo.Client) []string {
	if useDumps {
		if dumpDir == "" {
			panic("no 'dump_dir' flag given, so no idea from where to get the data")
		}
		ret, err := importHelper.AllDatabases(dumpDir)
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

func getAllCollectionsOrPanic(client *mongo.Client, dbName string) []string {
	if useDumps {
		if dumpDir == "" {
			panic("no 'dump_dir' flag given, so no idea from where to get the data")
		}
		ret, err := importHelper.AllCollectionsForDb(dumpDir, dbName)
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
