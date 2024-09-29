package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var databaseName string

var collectionName string

var outputDir string

var itemCount int64

var blacklist []string

var useAggregation bool

var mongoV44 bool

var timeout int64

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
}
