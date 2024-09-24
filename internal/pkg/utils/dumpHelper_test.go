package utils

import (
	"fmt"
	"os"
	"testing"
)

func Test_dumpToFile(t *testing.T) {
	dumpFile, err := os.Create("../../../tmp" + string(os.PathSeparator) + "Test_dumpToFile.txt")
	if err != nil {
		t.Errorf("Fail to open file: %v", err)
		return
	}
	defer dumpFile.Close()
	for i := 0; i < 10; i++ {
		s := fmt.Sprintf("Das ist ein Test: %d\n", i)
		bytes := []byte(s)
		dumpToFile(bytes, dumpFile)
	}
}
