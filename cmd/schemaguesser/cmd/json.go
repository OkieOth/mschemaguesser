package cmd

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"okieoth/schemaguesser/internal/pkg/mongoHelper"
	"okieoth/schemaguesser/internal/pkg/progressbar"

	"okieoth/schemaguesser/internal/pkg/utils"

	"github.com/spf13/cobra"
)

var jsonCmd = &cobra.Command{
	Use:   "json",
	Short: "dump bson content converted to JSON",
	Long:  "With this command you can dump raw content as converted JSON of one or more mongodb collections. The usecase is comparing collection content in an editor for instance.",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := mongoHelper.Connect(mongoHelper.ConStr)
		if err != nil {
			msg := fmt.Sprintf("Failed to connect to db: %v", err)
			panic(msg)
		}
		defer mongoHelper.CloseConnection(client)

		if databaseName == "all" {
			jsonForAllDatabases(client, true)
		} else {
			if collectionName == "all" {
				jsonForAllCollections(client, databaseName, true)
			} else {
				jsonForOneCollection(client, databaseName, collectionName, false, true)
			}
		}

	},
}

func replaceUuidValues(version byte, jsonStr string) (string, error) {
	re := regexp.MustCompile(fmt.Sprintf(`"Subtype":\s*%d,\s*"Data":"([A-Za-z0-9+/=]+)"`, version))

	matches := re.FindAllStringSubmatch(jsonStr, -1)

	ret := jsonStr

	for _, match := range matches {
		if len(match) > 1 {
			base64Str := match[1]
			decoded, err := base64.StdEncoding.DecodeString(base64Str)
			if err != nil {
				log.Printf("Error decoding Base64 data: %v", err)
				return "", err
			}
			if len(decoded) == 16 { // UUIDs are 16 bytes long
				uuidVal, err := uuid.FromBytes(decoded)
				if err == nil {
					origStr := fmt.Sprintf(`{%s}`, match[0])
					newStr := fmt.Sprintf(`"%s"`, &uuidVal)
					ret = strings.ReplaceAll(ret, origStr, newStr)
				}
			}
		}
	}
	return ret, nil
}

func getJsonBytes(b *bson.Raw) ([]byte, error) {
	var doc bson.M
	err := bson.Unmarshal(*b, &doc)
	if err != nil {
		log.Printf("Error while unmarshalling BSON: %v", err)
		return make([]byte, 0), errors.Join(err)
	}

	// Convert BSON document to JSON
	jsonData, err := json.Marshal(doc)

	if err != nil {
		log.Printf("Error while marshalling to JSON: %v", err)
		return make([]byte, 0), errors.Join(err)
	}

	jsonStr := string(jsonData)
	jsonStr, err = replaceUuidValues(3, jsonStr)
	if err != nil {
		log.Printf("Error while replacing UUID v3: %v", err)
		return make([]byte, 0), errors.Join(err)
	}
	jsonStr, err = replaceUuidValues(4, jsonStr)
	if err != nil {
		log.Printf("Error while replacing UUID v4: %v", err)
		return make([]byte, 0), errors.Join(err)
	}

	return []byte(jsonStr), nil
}

func jsonForOneCollection(client *mongo.Client, dbName string, collName string, doRecover bool, initProgressBar bool) {
	defer func() {
		if doRecover {
			if r := recover(); r != nil {
				log.Printf("Recovered while handling collection (db: %s, collection: %s): %v", dbName, collName, r)
			}
		}
	}()
	if initProgressBar {
		descr := fmt.Sprintf("JSON export of %s:%s", dbName, collName)
		progressbar.Init(1, descr)
	}

	outputFile, err := utils.CreateOutputFile(outputDir, "json", dbName, collName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	startTime := time.Now()
	i := 0

	utils.DumpBytesToFile([]byte("["), outputFile)
	err = mongoHelper.QueryCollection(client, dbName, collName, int(itemCount), useAggregation, mongoV44, func(data bson.Raw) error {
		bytes, err := getJsonBytes(&data)
		if err != nil {
			log.Printf("Error while converting to JSON: %v", err)
			return err
		}
		if i > 0 {
			utils.DumpBytesToFile([]byte(","), outputFile)
		}
		utils.DumpBytesToFile(bytes, outputFile)
		utils.DumpBytesToFile([]byte("\n"), outputFile)
		i++
		return nil // TODO
	})
	utils.DumpBytesToFile([]byte("]"), outputFile)

	if err != nil {
		msg := fmt.Sprintf("Error while reading data for collection (%s.%s): \n%v\n", dbName, collName, err)
		panic(msg)
	}
	log.Printf("[%s:%s] JSON exported for collection in %v\n", dbName, collName, time.Since(startTime))
	if initProgressBar {
		progressbar.ProgressOne()
	}
}

func jsonForAllCollections(client *mongo.Client, dbName string, initProgressBar bool) {
	collections := mongoHelper.ReadCollectionsOrPanic(client, dbName)
	var wg sync.WaitGroup
	if initProgressBar {
		progressbar.Init(int64(len(collections)), "JSON export for all collections")
	}

	for _, coll := range collections {
		if slices.Contains(blacklist, coll) {
			log.Printf("[%s:%s] skip blacklisted collection\n", dbName, coll)
			continue
		}
		wg.Add(1)
		go func(s string) {
			startTime := time.Now()
			defer func() {
				log.Printf("[%s:%s] JSON export of collection in %v\n", dbName, s, time.Since(startTime))
				wg.Done()
				if initProgressBar {
					progressbar.ProgressOne()
				}
			}()
			jsonForOneCollection(client, dbName, s, true, false)
		}(coll)
	}
	wg.Wait()
}

func jsonForAllDatabases(client *mongo.Client, initProgressBar bool) {
	dbs := mongoHelper.ReadDatabasesOrPanic(client)
	var wg sync.WaitGroup
	if initProgressBar {
		progressbar.Init(int64(len(dbs)), "JSON export for all databases")
	}
	for _, db := range dbs {
		if slices.Contains(blacklist, db) {
			log.Printf("[%s] skip blacklisted DB\n", db)
			continue
		}
		wg.Add(1)
		go func(s string) {
			startTime := time.Now()
			defer func() {
				log.Printf("[%s] JSON exported from DB in %v\n", s, time.Since(startTime))
				wg.Done()
				if initProgressBar {
					progressbar.ProgressOne()
				}
			}()
			jsonForAllCollections(client, s, false)
		}(db)
	}
	wg.Wait()
}
