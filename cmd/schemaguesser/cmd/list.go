package cmd

import (
	"fmt"

	"okieoth/schemaguesser/internal/pkg/mongoHelper"

	"github.com/spf13/cobra"
)

var databaseName string

var collectionName string

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
		dbs := mongoHelper.ReadDatabasesOrPanic()
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
		if databaseName == "all" {
			printAllCollections()
		} else {
			printOneCollection(databaseName, false)
		}
	},
}

var indexesCmd = &cobra.Command{
	Use:   "indexes",
	Short: "Indexes to a given collection",
	Long:  "Provides information about indexes of a given collection",
	Run: func(cmd *cobra.Command, args []string) {

		if databaseName == "all" {
			printIndexesForAllDatabases()
		} else {
			if collectionName == "all" {
				printIndexesForAllCollections(databaseName)
			} else {
				printIndexesForOneCollection(databaseName, collectionName, false)
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

func printOneCollection(dbName string, verbose bool) {
	collections := mongoHelper.ReadCollectionsOrPanic(dbName)
	for _, s := range *collections {
		if verbose {
			fmt.Printf("Database: %s, Collection: %s\n", dbName, s)
		} else {
			fmt.Println(s)
		}
	}
}

func printAllCollections() {
	dbs := mongoHelper.ReadDatabasesOrPanic()
	for _, db := range *dbs {
		printOneCollection(db, true)
	}
}

func printIndexesForOneCollection(dbName string, collName string, verbose bool) {
	indexes, err := mongoHelper.ListIndexes(mongoHelper.ConStr, dbName, collName)
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

func printIndexesForAllCollections(dbName string) {
	collections := mongoHelper.ReadCollectionsOrPanic(dbName)
	for _, coll := range *collections {
		printIndexesForOneCollection(dbName, coll, true)
	}
}

func printIndexesForAllDatabases() {
	dbs := mongoHelper.ReadDatabasesOrPanic()
	for _, db := range *dbs {
		printIndexesForAllCollections(db)
	}
}
