package mongoHelper

import (
	"context"
	"fmt"
	"log"
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

func Dummy() {
	fmt.Println("'Dummy' is called")
}

func ListDatabases(conStr string) ([]string, error) {
	var ret []string
	client, err := Connect(conStr)
	if err != nil {
		return ret, err
	}

	defer func() {
		if client == nil {
			return
		}
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	cursor, err := client.ListDatabases(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
		return ret, err
	}

	for _, db := range cursor.Databases {
		ret = append(ret, db.Name)
	}
	return ret, nil
}

func ListCollections(conStr string, databaseName string) ([]string, error) {
	var ret []string
	client, err := Connect(conStr)
	if err != nil {
		return ret, err
	}

	defer func() {
		if client == nil {
			return
		}
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	db := client.Database(databaseName)
	cursor, err := db.ListCollectionNames(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
		return ret, err
	}

	for _, collName := range cursor {
		ret = append(ret, collName)
	}
	return ret, nil
}

func ListIndexes(conStr string, databaseName string, collectionName string) ([]string, error) {
	var ret []string
	client, err := Connect(conStr)
	if err != nil {
		return ret, err
	}

	defer func() {
		if client == nil {
			return
		}
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	db := client.Database(databaseName)
	collection := db.Collection(collectionName)
	indexView := collection.Indexes()
	cursor, err := indexView.List(context.Background())

	if err != nil {
		log.Fatal(err)
		return ret, err
	}

	for cursor.Next(context.Background()) {
		bsonRaw := cursor.Current
		ret = append(ret, bsonRaw.String())
	}

	return ret, nil
}

func QueryCollection(conStr string, databaseName string, collectionName string, itemCount int) ([]bson.Raw, error) {
	var ret []bson.Raw
	client, err := Connect(conStr)
	if err != nil {
		return ret, err
	}

	defer func() {
		if client == nil {
			return
		}
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	db := client.Database(databaseName)
	collection := db.Collection(collectionName)
	cursor, err := collection.Find(context.Background(), bson.M{})

	if err != nil {
		log.Fatal(err)
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

func ReadCollectionsOrPanic(dbName string) *[]string {
	collections, err := ListCollections(ConStr, dbName)
	if err != nil {
		msg := fmt.Sprintf("Error while reading collections for database (%s): \n%v\n", dbName, err)
		panic(msg)
	}
	return &collections
}

func ReadDatabasesOrPanic() *[]string {
	dbs, err := ListDatabases(ConStr)
	if err != nil {
		msg := fmt.Sprintf("Error while reading existing databases: \n%v\n", err)
		panic(msg)
	}
	return &dbs
}
