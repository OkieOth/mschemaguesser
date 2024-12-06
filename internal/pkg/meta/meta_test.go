package meta_test

import (
	"okieoth/schemaguesser/internal/pkg/meta"
	"testing"
)

func TestGetAllMetaInfos(t *testing.T) {
	metaDir := "../../../resources/key_values"
	metaInfos, err := meta.GetAllMetaInfos(metaDir)
	if err != nil {
		t.Errorf("Error while reading all meta files: %v", err)
		return
	}
	expectedLen := 3
	gotLen := len(metaInfos)
	if gotLen != expectedLen {
		t.Errorf("Retrieved wrong number of meta files: got=%d, expected=%d", gotLen, expectedLen)
	}
}
