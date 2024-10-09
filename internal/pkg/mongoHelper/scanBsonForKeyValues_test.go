package mongoHelper

import (
	"path/filepath"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

func TestCheckIfBinaryIsUUID(t *testing.T) {
	// Happy case: binary with subtype 3 (Legacy UUID)
	legacyUUID := bson.RawValue{
		Type:  bson.TypeBinary,
		Value: bsoncore.AppendBinary(nil, 0x03, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}),
	}

	// Happy case: binary with subtype 4 (RFC 4122 UUID)
	rfcUUID := bson.RawValue{
		Type:  bson.TypeBinary,
		Value: bsoncore.AppendBinary(nil, 0x04, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}),
	}

	// Failing case: binary with a different subtype (not 3 or 4)
	randomBinary := bson.RawValue{
		Type:  bson.TypeBinary,
		Value: bsoncore.AppendBinary(nil, 0x05, []byte{1, 2, 3, 4}),
	}

	// Failing case: value is not binary (e.g., a string)
	nonBinaryValue := bson.RawValue{
		Type:  bson.TypeString,
		Value: bsoncore.AppendString(nil, "not a binary type"),
	}

	// Run the tests
	tests := []struct {
		name      string
		input     bson.RawValue
		expect    bool
		expectErr bool
	}{
		{"Valid UUID (Legacy, subtype 3)", legacyUUID, true, false},
		{"Valid UUID (RFC 4122, subtype 4)", rfcUUID, true, false},
		{"Invalid binary (random subtype)", randomBinary, false, false},
		{"Invalid value (non-binary type)", nonBinaryValue, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := checkIfBinaryIsUUID(tt.input)
			if (err != nil) != tt.expectErr {
				t.Errorf("expected error: %v, got error: %v", tt.expectErr, err)
			}
			if result != tt.expect {
				t.Errorf("expected result: %v, got: %v", tt.expect, result)
			}
		})
	}
}

func TestCheckIfStringIsUUIDString(t *testing.T) {
	tests := []struct {
		name    string
		value   bson.RawValue
		want    bool
		wantErr bool
	}{
		{
			name: "Valid UUID",
			value: bson.RawValue{
				Type:  bson.TypeString,
				Value: bsoncore.AppendString(nil, "123e4567-e89b-12d3-a456-426614174000"),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Invalid UUID format (no hyphens)",
			value: bson.RawValue{
				Type:  bson.TypeString,
				Value: bsoncore.AppendString(nil, "123e4567e89b12d3a456426614174000"),
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Invalid UUID format (incorrect length)",
			value: bson.RawValue{
				Type:  bson.TypeString,
				Value: bsoncore.AppendString(nil, "123e4567-e89b-12d3-a456-426614174"),
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Valid UUID with uppercase letters",
			value: bson.RawValue{
				Type:  bson.TypeString,
				Value: bsoncore.AppendString(nil, "123E4567-E89B-12D3-A456-426614174000"),
			},
			want:    true,
			wantErr: false,
		},
		{
			value: bson.RawValue{
				Type:  bson.TypeInt32,
				Value: []byte{42}, // Example of an int32 BSON value
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Empty string value",
			value: bson.RawValue{
				Type:  bson.TypeString,
				Value: bsoncore.AppendString(nil, ""),
			},
			want:    false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkIfStringIsUUIDString(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkIfStringIsUUIDString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("checkIfStringIsUUIDString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPersistenceFileName(t *testing.T) {
	outputDir := "/base/output"

	// Define test cases
	tests := []struct {
		dbName     string
		collName   string
		attribName string
		expected   string
	}{
		{
			dbName:     "myDB",
			collName:   "myCollection",
			attribName: "attribute",
			expected:   filepath.Join(outputDir, "myDB", "myCollection", "attribute.keyvalues.txt"),
		},
		{
			dbName:     "my DB",    // space in dbName
			collName:   "coll@123", // special chars in collName
			attribName: "attrib$1", // special chars in attribName
			expected:   filepath.Join(outputDir, "my_DB", "coll_123", "attrib_1.keyvalues.txt"),
		},
		{
			dbName:     "dbWithÜnicode", // non-ASCII char in dbName
			collName:   "normalColl",
			attribName: "üattrib", // non-ASCII char in attribName
			expected:   filepath.Join(outputDir, "dbWith_nicode", "normalColl", "_attrib.keyvalues.txt"),
		},
		{
			dbName:     "nameWithSpace",
			collName:   "collWithSpace",
			attribName: "attribute with space", // space in attribName
			expected:   filepath.Join(outputDir, "nameWithSpace", "collWithSpace", "attribute_with_space.keyvalues.txt"),
		},
		{
			dbName:     "123中文",      // non-ASCII and numeric characters in dbName
			collName:   "測試Coll",     // non-ASCII in collName
			attribName: "someAttr測試", // non-ASCII in attribName
			expected:   filepath.Join(outputDir, "123__", "__Coll", "someAttr__.keyvalues.txt"),
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.dbName+"_"+test.collName+"_"+test.attribName, func(t *testing.T) {
			result := GetPersistenceFileName(outputDir, test.dbName, test.collName, test.attribName)
			if result != test.expected {
				t.Errorf("Expected %s, but got %s", test.expected, result)
			}
		})
	}
}
