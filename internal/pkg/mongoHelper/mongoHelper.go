package mongoHelper

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"okieoth/schemaguesser/internal/pkg/utils"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	//"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type HandleDataCallback func(bson.Raw) error

var ConStr string

func Connect(conStr string) (*mongo.Client, error) {
	conStr2Use := utils.ReplaceWithEnvVar(conStr, "MONGO_USER", "admin")
	conStr2Use = utils.ReplaceWithEnvVar(conStr2Use, "MONGO_PASSWORD", "secretpassword")
	conStr2Use = utils.ReplaceWithEnvVar(conStr2Use, "MONGO_HOST", "localhost")
	conStr2Use = utils.ReplaceWithEnvVar(conStr2Use, "MONGO_PORT", "27017")

	log.Printf("Connect to mongo: conStr2Use=%s", conStr2Use)

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

func queryCollectionWithAggregation(client *mongo.Client, databaseName string, collectionName string, itemCount int, handleDataCallback HandleDataCallback) error {
	db := client.Database(databaseName)
	collection := db.Collection(collectionName)

	// Define a simple aggregation pipeline that acts like a find
	pipeline := mongo.Pipeline{
		{{"$match", bson.M{}}}, // Add any match conditions if needed
		{{"$limit", itemCount}},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Set allowDiskUse to true in aggregation options
	aggregationOptions := options.Aggregate().SetAllowDiskUse(true)

	startTime := time.Now()
	cursor, err := collection.Aggregate(ctx, pipeline, aggregationOptions)
	if err != nil {
		log.Printf("[%s:%s] Collection query error: %v\n", databaseName, collectionName, err)
		return err
	}
	log.Printf("[%s:%s] Collection query executed in %v\n", databaseName, collectionName, time.Since(startTime))

	for cursor.Next(ctx) {
		bsonRaw := cursor.Current
		if err := handleDataCallback(bsonRaw); err != nil {
			log.Printf("[%s:%s] error while processing the data: %v\n", databaseName, collectionName, err)
			return err
		}
	}

	return nil
}

func queryCollection(client *mongo.Client, databaseName string, collectionName string, itemCount int, mongo44 bool, handleDataCallback HandleDataCallback) error {
	db := client.Database(databaseName)
	collection := db.Collection(collectionName)
	// setAllowDiskUse requires mongodb 4.4 at minimum
	startTime := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	findOptions := options.Find().SetLimit(int64(itemCount))
	if mongo44 {
		findOptions = findOptions.SetAllowDiskUse(true)
	}
	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		//panic(err)
		log.Printf("[%s:%s] Collection query error: %v\n", databaseName, collectionName, err)
		return err
	}
	log.Printf("Query executed in %v\n", time.Since(startTime))

	for cursor.Next(ctx) {
		bsonRaw := cursor.Current
		if err := handleDataCallback(bsonRaw); err != nil {
			log.Printf("[%s:%s] error while processing the data: %v\n", databaseName, collectionName, err)
			return err
		}
	}

	return nil
}

func DumpCollectionToFile(ctx context.Context, outputFile *os.File, client *mongo.Client, databaseName string, collectionName string, itemCount int64, useAggregation bool, mongo44 bool) (uint64, error) {
	if useAggregation {
		return dumpCollectionWithAggregationToFile(ctx, outputFile, client, databaseName, collectionName, itemCount)
	} else {
		return dumpCollectionToFile(ctx, outputFile, client, databaseName, collectionName, itemCount, mongo44)
	}
}

func dumpCollectionWithAggregationToFile(ctx context.Context, outputFile *os.File, client *mongo.Client, databaseName string, collectionName string, itemCount int64) (uint64, error) {
	db := client.Database(databaseName)
	collection := db.Collection(collectionName)
	// setAllowDiskUse requires mongodb 4.4 at minimum
	startTime := time.Now()

	// Define a simple aggregation pipeline that acts like a find

	var pipeline []bson.D
	if itemCount > 0 {
		pipeline = mongo.Pipeline{
			{{"$match", bson.M{}}}, // Add any match conditions if needed
			{{"$limit", itemCount}},
		}
	} else {
		pipeline = mongo.Pipeline{
			{{"$match", bson.M{}}}, // Add any match conditions if needed
		}
	}

	// Set allowDiskUse to true in aggregation options
	aggregationOptions := options.Aggregate().SetAllowDiskUse(true)

	cursor, err := collection.Aggregate(ctx, pipeline, aggregationOptions)
	if err != nil {
		log.Printf("[%s:%s] Collection query error: %v\n", databaseName, collectionName, err)
		return 0, err
	}

	log.Printf("[%s:%s] dumpCollectionToFile - Query executed in %v\n", databaseName, collectionName, time.Since(startTime))
	var dumpCount uint64
	for cursor.Next(ctx) {
		bsonRaw := cursor.Current
		_, err = outputFile.Write(bsonRaw)
		if err != nil {
			log.Printf("Failed to write BSON to file: %v", err)
			return dumpCount, err
		}
		dumpCount++
	}
	return dumpCount, nil
}

func dumpCollectionToFile(ctx context.Context, outputFile *os.File, client *mongo.Client, databaseName string, collectionName string, itemCount int64, mongo44 bool) (uint64, error) {
	db := client.Database(databaseName)
	collection := db.Collection(collectionName)
	// setAllowDiskUse requires mongodb 4.4 at minimum
	startTime := time.Now()

	findOptions := options.Find()
	if itemCount > 0 {
		findOptions.SetLimit(int64(itemCount))
		if mongo44 {
			findOptions = findOptions.SetAllowDiskUse(true)
		}
	}
	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		//panic(err)
		log.Printf("[%s:%s] dumpCollectionToFile - Collection query error: %v\n", databaseName, collectionName, err)
		return 0, err
	}
	log.Printf("[%s:%s] dumpCollectionToFile - Query executed in %v\n", databaseName, collectionName, time.Since(startTime))
	var dumpCount uint64
	for cursor.Next(ctx) {
		bsonRaw := cursor.Current
		_, err = outputFile.Write(bsonRaw)
		if err != nil {
			log.Printf("Failed to write BSON to file: %v", err)
			return dumpCount, err
		}
		dumpCount++
	}
	return dumpCount, nil
}

func readBSONFileAndInsertToMongo(collection *mongo.Collection, filePath string) error {
	// Open the BSON file for reading
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("failed to open input file: %v", err)
		return err
	}
	defer file.Close()

	buf := make([]byte, 4)

	var documents []bson.M
	for {

		_, err := io.ReadFull(file, buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		docLength := int32(binary.LittleEndian.Uint32(buf))
		docBuf := make([]byte, docLength)
		_, err = io.ReadFull(file, docBuf)
		if err != nil {
			return err
		}
		var doc bson.M
		err = bson.Unmarshal(docBuf, &doc)
		if err != nil {
			log.Printf("failed to unmarshal bytes to bson doc: %v", err)
			return err
		}
		documents = append(documents, doc)
		if len(documents) == 100 {
			// TODO do bulk insert
			documents = documents[:0]
		}
	}
	if len(documents) > 0 {
		// TODO do bulk insert
	}

	return nil
}

func insertChunk(documents []interface{}, collection *mongo.Collection) error {
	// Bulk insert the documents
	// TODO
	//wc := writeconcern.New(writeconcern.WMajority())
	opts := options.InsertMany()

	_, err := collection.InsertMany(context.TODO(), documents, opts)
	if err != nil {
		return err
	}

	fmt.Printf("Inserted %d documents in batch\n", len(documents))
	return nil
}

// This version only works from mongodb v4.4
func QueryCollection(client *mongo.Client, databaseName string, collectionName string, itemCount int, useAggregation bool, mongo44 bool, handleDataCallback HandleDataCallback) error {
	if useAggregation {
		return queryCollectionWithAggregation(client, databaseName, collectionName, itemCount, handleDataCallback)
	} else {
		return queryCollection(client, databaseName, collectionName, itemCount, mongo44, handleDataCallback)
	}
}

func CountCollection(client *mongo.Client, dbName string, collName string) (int64, error) {
	db := client.Database(dbName)
	collection := db.Collection(collName)
	startTime := time.Now()

	pipeline := mongo.Pipeline{
		{{"$count", "totalCount"}},
	}

	// Set aggregation options with AllowDiskUse
	aggOpts := options.Aggregate().SetAllowDiskUse(true)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Run the aggregation
	cursor, err := collection.Aggregate(ctx, pipeline, aggOpts)
	if err != nil {
		return -1, err
	}

	defer cursor.Close(ctx)

	var result []bson.M
	if err = cursor.All(ctx, &result); err != nil {
		return -1, err
	}
	var c int64 = 0
	if len(result) > 0 {
		if count, ok := result[0]["totalCount"].(int32); ok {
			c = int64(count)
		} else if count, ok := result[0]["totalCount"].(int64); ok {
			c = count
		} else {
			return -1, errors.New("Failed to cast the count to int64")
		}
	}
	log.Printf("Query executed in %v, count=%d\n", time.Since(startTime), c)
	return c, nil
}

func ReadCollectionsOrPanic(client *mongo.Client, dbName string) []string {
	collections, err := ListCollections(client, dbName)
	if err != nil {
		msg := fmt.Sprintf("Error while reading collections for database (%s): \n%v\n", dbName, err)
		panic(msg)
	}
	return collections
}

func ReadDatabasesOrPanic(client *mongo.Client) []string {
	dbs, err := ListDatabases(client)
	if err != nil {
		msg := fmt.Sprintf("Error while reading existing databases: \n%v\n", err)
		panic(msg)
	}
	return dbs
}
