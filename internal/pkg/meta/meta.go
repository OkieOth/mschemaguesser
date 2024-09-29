package meta

// This package helps to write meta information about the specific export
// Since collection and database names can not really be guessed from the specific file
// name, the meta file can be used to import persisted data to the proper database and collection

import (
	"encoding/json"
	"errors"
	"log"
	"okieoth/schemaguesser/internal/pkg/utils"
	"time"
)

var version = "1.0.0"

type TimeoutInfo struct {
	Reached bool   `json:"reached,omitempty"`
	Seconds int64  `json:"seconds,omitempty"`
	Error   string `json:"error,omitempty"`
}

type MetaInfo struct {
	Version    string       `json:"version,omitempty"`
	Comment    string       `json:"comment,omitempty"`
	Db         string       `json:"db,omitempty"`
	Collection string       `json:"collection,omitempty"`
	ExportTime time.Time    `json:"exportTime,omitempty"`
	ItemCount  uint64       `json:"itemCount,omitempty"`
	Timeout    *TimeoutInfo `json:"timeout,omitempty"`
}

func WriteMetaInfo(outputDir string, dbName string, collName string, itemCount uint64, comment string, timeout *TimeoutInfo) error {
	outputFile, err := utils.CreateOutputFile(outputDir, "meta", dbName, collName)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	var metaInfo MetaInfo
	metaInfo.Version = version
	metaInfo.Collection = collName
	metaInfo.Db = dbName
	metaInfo.Comment = comment
	metaInfo.ExportTime = time.Now()
	metaInfo.ItemCount = itemCount
	if timeout != nil {
		metaInfo.Timeout = timeout
	}

	jsonData, err := json.MarshalIndent(metaInfo, "", "  ")

	if err != nil {
		log.Printf("Error while marshalling meta info to JSON: %v", err)
		return errors.Join(err)
	}
	_, err = outputFile.Write(jsonData)
	if err != nil {
		log.Printf("Failed to write document length to file: %v", err)
		return err
	}

	return nil
}
