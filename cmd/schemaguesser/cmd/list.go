package cmd

import (
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"

	"okieoth/schemaguesser/internal/pkg/mongoHelper"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Discovery the database you are connected to",
	Long:  "Provides information about existing databases, collections and their indexes",
	Run: func(cmd *cobra.Command, args []string) {
		// Logic for the greet command
		fmt.Println("To learn about the possible options type:\n`schemaguesser list --help`")
	},
}

var databasesCmd = &cobra.Command{
	Use:   "databases",
	Short: "Names of databases available in the connected server",
	Long:  "Prints a list of available databases i the connected server",
	Run: func(cmd *cobra.Command, args1 []string) {
		client, err := mongoHelper.Connect(mongoHelper.ConStr)
		if err != nil {
			msg := fmt.Sprintf("Failed to connect to db: %v", err)
			panic(msg)
		}
		defer mongoHelper.CloseConnection(client)

		dbs := mongoHelper.ReadDatabasesOrPanic(client)
		for _, s := range *dbs {
			fmt.Println(s)
		}
	},
}

var collectionsCmd = &cobra.Command{
	Use:   "collections",
	Short: "Collections available in the specified database",
	Long:  "Provides information about existing collections in the specified database",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := mongoHelper.Connect(mongoHelper.ConStr)
		if err != nil {
			msg := fmt.Sprintf("Failed to connect to db: %v", err)
			panic(msg)
		}
		defer mongoHelper.CloseConnection(client)

		if databaseName == "all" {
			printAllCollections(client)
		} else {
			printOneCollection(client, databaseName, false)
		}
	},
}

var indexesCmd = &cobra.Command{
	Use:   "indexes",
	Short: "Indexes to a given collection",
	Long:  "Provides information about indexes of a given collection",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := mongoHelper.Connect(mongoHelper.ConStr)
		if err != nil {
			msg := fmt.Sprintf("Failed to connect to db: %v", err)
			panic(msg)
		}
		defer mongoHelper.CloseConnection(client)

		if databaseName == "all" {
			printIndexesForAllDatabases(client)
		} else {
			if collectionName == "all" {
				printIndexesForAllCollections(client, databaseName)
			} else {
				printIndexesForOneCollection(client, databaseName, collectionName, false)
			}
		}
	},
}

func init() {
	listCmd.AddCommand(databasesCmd)
	listCmd.AddCommand(collectionsCmd)
	listCmd.AddCommand(indexesCmd)

	collectionsCmd.Flags().StringVar(&databaseName, "database", "all", "Database to query existing collections. If 'all', then the collections of all databases are printed.")

	indexesCmd.Flags().StringVar(&databaseName, "database", "all", "Database to query existing collections. If 'all', then the collections of all databases are printed.")
	indexesCmd.Flags().StringVar(&collectionName, "collection", "all", "Name of the collection to show the indexes.  If 'all', then the collections of all databases are printed.")
}

func printOneCollection(client *mongo.Client, dbName string, verbose bool) {
	collections := mongoHelper.ReadCollectionsOrPanic(client, dbName)
	for _, s := range *collections {
		if verbose {
			fmt.Printf("Database: %s, Collection: %s\n", dbName, s)
		} else {
			fmt.Println(s)
		}
	}
}

func printAllCollections(client *mongo.Client) {
	dbs := mongoHelper.ReadDatabasesOrPanic(client)
	for _, db := range *dbs {
		printOneCollection(client, db, true)
	}
}

func printIndexesForOneCollection(client *mongo.Client, dbName string, collName string, verbose bool) {
	indexes, err := mongoHelper.ListIndexes(client, dbName, collName)
	if err != nil {
		msg := fmt.Sprintf("Error while reading indexes for collection (%s.%s): \n%v\n", dbName, collName, err)
		panic(msg)
	}
	for _, s := range indexes {
		if verbose {
			fmt.Printf("Database: %s, Collection: %s, Index: %s\n", dbName, collName, s)
		} else {
			fmt.Println(s)
		}
	}
}

func printIndexesForAllCollections(client *mongo.Client, dbName string) {
	collections := mongoHelper.ReadCollectionsOrPanic(client, dbName)
	for _, coll := range *collections {
		printIndexesForOneCollection(client, dbName, coll, true)
	}
}

func printIndexesForAllDatabases(client *mongo.Client) {
	dbs := mongoHelper.ReadDatabasesOrPanic(client)
	for _, db := range *dbs {
		printIndexesForAllCollections(client, db)
	}
}
