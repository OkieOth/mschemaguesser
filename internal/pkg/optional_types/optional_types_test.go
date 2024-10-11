package optionaltypes

import (
	"errors"
	"fmt"
	"testing"
)

func TestOptionalString(t *testing.T) {
	var s Optional[string]
	if s.IsSet {
		t.Errorf("uninitialized value is set after creation")
	}
}

type DummyEnum int64

const (
	value1 DummyEnum = iota
	value2
	value3
	value4
	value5
	value6
)

func (s DummyEnum) String() string {
	switch s {
	case value1:
		return "value1"
	case value2:
		return "value2"
	case value3:
		return "value3"
	case value4:
		return "value4"
	case value5:
		return "value5"
	default:
		return "value6"
	}
}

func (s DummyEnum) ValueFromStr(str string) error {
	switch str {
	case "value1":
		s = value1
	case "value2":
		s = value2
	case "value3":
		s = value3
	case "value4":
		s = value4
	case "value5":
		s = value5
	case "value6":
		s = value6
	default:
		return errors.New(fmt.Sprintf("input not part of the enum: %v", str))
	}
	return nil
}

func TestOptionalEnum(t *testing.T) {
	var s OptionalEnum[DummyEnum]
	if s.IsSet {
		t.Errorf("uninitialized value is set after creation")
	}
}
