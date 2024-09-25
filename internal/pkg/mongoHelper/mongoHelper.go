package mongoHelper

import (
	"context"
	"fmt"
	"log"
	"okieoth/schemaguesser/internal/pkg/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ConStr string

func Connect(conStr string) (*mongo.Client, error) {
	conStr2Use := utils.ReplaceWithEnvVar(conStr, "MONGO_USER", "admin")
	conStr2Use = utils.ReplaceWithEnvVar(conStr2Use, "MONGO_PASSWORD", "secretpassword")
	conStr2Use = utils.ReplaceWithEnvVar(conStr2Use, "MONGO_HOST", "localhost")
	conStr2Use = utils.ReplaceWithEnvVar(conStr2Use, "MONGO_PORT", "27017")

	return mongo.Connect(context.Background(), options.Client().ApplyURI(conStr2Use))
}

func CloseConnection(client *mongo.Client) {
	if client == nil {
		return
	}
	if err := client.Disconnect(context.Background()); err != nil {
		println("Error while disconnect: %v", err)
	}
}

func Dummy() {
	fmt.Println("'Dummy' is called")
}

func ListDatabases(client *mongo.Client) ([]string, error) {
	var ret []string

	cursor, err := client.ListDatabases(context.Background(), bson.M{})
	if err != nil {
		panic(err)
	}

	for _, db := range cursor.Databases {
		ret = append(ret, db.Name)
	}
	return ret, nil
}

func ListCollections(client *mongo.Client, databaseName string) ([]string, error) {
	var ret []string

	db := client.Database(databaseName)
	cursor, err := db.ListCollectionNames(context.Background(), bson.M{})
	if err != nil {
		panic(err)
	}

	for _, collName := range cursor {
		ret = append(ret, collName)
	}
	return ret, nil
}

func ListIndexes(client *mongo.Client, databaseName string, collectionName string) ([]string, error) {
	var ret []string

	db := client.Database(databaseName)
	collection := db.Collection(collectionName)
	indexView := collection.Indexes()
	cursor, err := indexView.List(context.Background())

	if err != nil {
		panic(err)
	}

	for cursor.Next(context.Background()) {
		bsonRaw := cursor.Current
		ret = append(ret, bsonRaw.String())
	}

	return ret, nil
}

func queryCollectionWithAggregation(client *mongo.Client, databaseName string, collectionName string, itemCount int) ([]bson.Raw, error) {
	var ret []bson.Raw

	db := client.Database(databaseName)
	collection := db.Collection(collectionName)

	// Define a simple aggregation pipeline that acts like a find
	pipeline := mongo.Pipeline{
		{{"$match", bson.M{}}}, // Add any match conditions if needed
		{{"$limit", itemCount}},
	}

	// Set allowDiskUse to true in aggregation options
	aggregationOptions := options.Aggregate().SetAllowDiskUse(true)

	startTime := time.Now()
	cursor, err := collection.Aggregate(context.Background(), pipeline, aggregationOptions)
	log.Printf("[%s:%s] Collection query executed in %v\n", databaseName, collectionName, time.Since(startTime))
	if err != nil {
		return ret, err
	}

	for cursor.Next(context.Background()) {
		bsonRaw := cursor.Current
		ret = append(ret, bsonRaw)
	}

	return ret, nil
}

func queryCollection(client *mongo.Client, databaseName string, collectionName string, itemCount int, mongo44 bool) ([]bson.Raw, error) {
	var ret []bson.Raw

	db := client.Database(databaseName)
	collection := db.Collection(collectionName)
	// setAllowDiskUse requires mongodb 4.4 at minimum
	startTime := time.Now()

	findOptions := options.Find().SetLimit(int64(itemCount))
	if mongo44 {
		findOptions = findOptions.SetAllowDiskUse(true)
	}
	cursor, err := collection.Find(context.Background(), bson.M{}, findOptions)
	log.Printf("Query executed in %v\n", time.Since(startTime))
	if err != nil {
		//panic(err)
		return ret, err
	}

	i := 0
	for cursor.Next(context.Background()) && i < itemCount {
		bsonRaw := cursor.Current
		ret = append(ret, bsonRaw)
		i++
	}

	return ret, nil
}

// This version only works from mongodb v4.4
func QueryCollection(client *mongo.Client, databaseName string, collectionName string, itemCount int, useAggregation bool, mongo44 bool) ([]bson.Raw, error) {
	if useAggregation {
		return queryCollectionWithAggregation(client, databaseName, collectionName, itemCount)
	} else {
		return queryCollection(client, databaseName, collectionName, itemCount, mongo44)
	}
}

func ReadCollectionsOrPanic(client *mongo.Client, dbName string) *[]string {
	collections, err := ListCollections(client, dbName)
	if err != nil {
		msg := fmt.Sprintf("Error while reading collections for database (%s): \n%v\n", dbName, err)
		panic(msg)
	}
	return &collections
}

func ReadDatabasesOrPanic(client *mongo.Client) *[]string {
	dbs, err := ListDatabases(client)
	if err != nil {
		msg := fmt.Sprintf("Error while reading existing databases: \n%v\n", err)
		panic(msg)
	}
	return &dbs
}
