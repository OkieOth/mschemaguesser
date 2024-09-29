package optionaltypes

import (
	"encoding/json"
)

type EnumType interface {
	String() string
	ValueFromStr(str string) error
}

type Optional[C any] struct {
	Value C
	IsSet bool
}

func (m *Optional[C]) Set(v C) {
	m.Value = v
	m.IsSet = true
}

func (m *Optional[C]) UnSet() {
	m.IsSet = false
}

func (v *Optional[C]) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		v.IsSet = false
		return nil
	}

	err := json.Unmarshal(data, &v.Value)
	if err != nil {
		return err
	}

	v.IsSet = true
	return nil
}

func (v *Optional[C]) MarshalJSON() ([]byte, error) {
	if !v.IsSet {
		return []byte("null"), nil
	} else {
		return json.Marshal(v.Value)
	}
}

type OptionalEnum[C EnumType] struct {
	Value C
	IsSet bool
}

func (m *OptionalEnum[C]) Set(v C) {
	m.Value = v
	m.IsSet = true
}

func (m *OptionalEnum[C]) UnSet() {
	m.IsSet = false
}

func (v *OptionalEnum[C]) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	var s string
	err := json.Unmarshal(data, &s)

	if err != nil {
		return err
	}
	err = v.Value.ValueFromStr(s)
	if err != nil {
		return err
	}
	v.IsSet = true
	return nil
}

func (v *OptionalEnum[C]) MarshalJSON() ([]byte, error) {
	if !v.IsSet {
		return []byte("null"), nil
	} else {
		return json.Marshal(v.Value.String())
	}
}
