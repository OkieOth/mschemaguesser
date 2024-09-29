package importHelper

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"

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
