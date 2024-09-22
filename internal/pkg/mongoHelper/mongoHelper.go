package mongoHelper

import (
	"context"
	"fmt"
	"okieoth/schemaguesser/internal/pkg/utils"

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

func QueryCollection(client *mongo.Client, databaseName string, collectionName string, itemCount int) ([]bson.Raw, error) {
	var ret []bson.Raw

	db := client.Database(databaseName)
	collection := db.Collection(collectionName)
	cursor, err := collection.Find(context.Background(), bson.M{}, options.Find().SetLimit(int64(itemCount)))

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
