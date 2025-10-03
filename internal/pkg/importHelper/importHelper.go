package importHelper

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"okieoth/schemaguesser/internal/pkg/meta"
)

func AllDatabases(inputDir string) ([]string, error) {
	var ret []string

	err := filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(info.Name()) == ".meta" {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			var meta meta.MetaInfo
			if err := json.Unmarshal(data, &meta); err != nil {
				return err
			}
			if !slices.Contains(ret, meta.Db) {
				ret = append(ret, meta.Db)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return ret, nil
}

func AllCollectionsForDb(inputDir string, dbName string) ([]string, error) {
	var ret []string

	err := filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(info.Name()) == ".meta" {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			var meta meta.MetaInfo
			if err := json.Unmarshal(data, &meta); err != nil {
				return err
			}
			if meta.Db == dbName {
				if !slices.Contains(ret, meta.Collection) {
					ret = append(ret, meta.Collection)
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return ret, nil
}

func ImportData(client *mongo.Client, importFile string, dbName string, collName string, chunkSize int64, ctx *context.Context) (uint64, error) {
	file, err := os.Open(importFile)
	if err != nil {
		return 0, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	var docs []interface{}
	buf := make([]byte, 4)
	readCount := uint64(0)
	chunkElemCount := uint64(0)

	collection := client.Database(dbName).Collection(collName)
	for {
		_, err := io.ReadFull(file, buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return readCount, fmt.Errorf("failed to read document size to buffer: %v", err)
		}
		docLength := int32(binary.LittleEndian.Uint32(buf))
		docBuf := make([]byte, docLength)
		copy(docBuf, buf)
		_, err = io.ReadFull(file, docBuf[4:])
		if err != nil {
			return readCount, fmt.Errorf("failed to read document to buffer: %v", err)
		}
		readCount++
		if chunkElemCount < uint64(chunkSize) {
			docs = append(docs, docBuf)
		} else {
			err := insertChunk(docs, collection, ctx)
			if err != nil {
				return readCount, fmt.Errorf("failed to insert chunk into db: %v", err)
			}
			chunkElemCount = 0
			docs = docs[:0]
		}
	}
	if len(docs) > 0 {
		err = insertChunk(docs, collection, ctx)
		if err != nil {
			return readCount, fmt.Errorf("failed to insert final chunk into db: %v", err)
		}

	}
	log.Printf("[%s:%s] Imported %d documents", dbName, collName, readCount)
	return readCount, nil
}

func insertChunk(documents []interface{}, collection *mongo.Collection, ctx *context.Context) error {
	// Bulk insert the documents
	//wc := writeconcern.New(writeconcern.WMajority())
	//opts := options.InsertMany().SetWriteConcern(wc)

	opts := options.InsertMany()

	_, err := collection.InsertMany(*ctx, documents, opts)
	if err != nil {
		return err
	}
	return nil
}
