package utils

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type JsonArray []interface{}

func (a JsonArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *JsonArray) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}

func (a *JsonArray) String() string {
	bytes, _ := json.Marshal(a)
	return string(bytes)
}

type JsonStringArray []string

func (a JsonStringArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *JsonStringArray) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}

func (a *JsonStringArray) String() string {
	bytes, _ := json.Marshal(a)
	return string(bytes)
}

type JsonObject map[string]interface{}

func (a JsonObject) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *JsonObject) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}

func (a *JsonObject) String() string {
	bytes, _ := json.Marshal(a)
	return string(bytes)
}

func NewJsonObjectFromString(text string) JsonObject {
	obj := JsonObject{}
	MustUnmarshal([]byte(text), &obj)
	return obj
}

func NewJsonObjectFromObject(v interface{}) JsonObject {
	obj := JsonObject{}
	bytes := MustMarshal(v)
	MustUnmarshal(bytes, &obj)
	return obj
}
