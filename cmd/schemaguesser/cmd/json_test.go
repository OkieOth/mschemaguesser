package cmd

import (
	"fmt"
	"testing"
)

func Test_replaceUuidValues(t *testing.T) {
	jsonStr := `{"Category":7,"CharacterDetails":null,"DefaultOffProgramId":{"Subtype":4,"Data":"AAAAAAAAAAAAAAAAAAAAAA=="},"Description":"Vehicle Activated Sign","FullMatrixDetails":null,"Name":"Vehicle Activated Sign","PrismDetails":null,"TenantId":{"Subtype":4,"Data":"BWvPWOF+QrqBhvJf+96LNQ=="},"TenantType":0,"_id":{"Subtype":4,"Data":"AN0UjlWcSeKLPZa8F59Xog=="}}`

	convertedStr, err := replaceUuidValues(4, jsonStr)

	if err != nil {
		t.Errorf("Fail to replace uuids in replaceUuidValues: %v", err)
		return
	}

	fmt.Println(convertedStr)
}
