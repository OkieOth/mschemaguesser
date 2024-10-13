package meta

// This package helps to write meta information about the specific export
// Since collection and database names can not really be guessed from the specific file
// name, the meta file can be used to import persisted data to the proper database and collection

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"okieoth/schemaguesser/internal/pkg/utils"
	"os"
	"path/filepath"
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
	FileName   string       `json:"fileName,omitempty"`
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
	metaInfo.FileName = filepath.Base(outputFile.Name())
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

func GetAllMetaInfos(metaDir string) ([]MetaInfo, error) {

	ret := make([]MetaInfo, 0)

	// checks if metaDir exists
	// if so find all files with the extension '.meta'
	// unmarshal them from json to the MetaInfo type
	// append the MetaInfo to the 'ret' array

	if _, err := os.Stat(metaDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("directory does not exist: %s", metaDir)
	}

	err := filepath.Walk(metaDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(info.Name()) == ".meta" {
			data, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("error reading file %s: %w", path, err)
			}

			var meta MetaInfo
			if err := json.Unmarshal(data, &meta); err != nil {
				return fmt.Errorf("error while unmarshalling file %s: %w", path, err)
			}
			ret = append(ret, meta)
		}
		return nil
	})
	return ret, err
}
