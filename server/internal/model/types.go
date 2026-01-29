package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// JSONArray JSON数组类型
type JSONArray []string

// Scan 实现sql.Scanner接口
func (j *JSONArray) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("类型断言为[]byte失败")
	}
	if len(bytes) == 0 || string(bytes) == "null" {
		*j = nil
		return nil
	}
	return json.Unmarshal(bytes, j)
}

// Value 实现driver.Valuer接口
func (j JSONArray) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// JSONMap JSON对象类型
type JSONMap map[string]interface{}

// Scan 实现sql.Scanner接口
func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("类型断言为[]byte失败")
	}
	if len(bytes) == 0 || string(bytes) == "null" {
		*j = nil
		return nil
	}
	return json.Unmarshal(bytes, j)
}

// Value 实现driver.Valuer接口
func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// JSONIntArray JSON整数数组类型
type JSONIntArray []int

// Scan 实现sql.Scanner接口
func (j *JSONIntArray) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("类型断言为[]byte失败")
	}
	if len(bytes) == 0 || string(bytes) == "null" {
		*j = nil
		return nil
	}
	return json.Unmarshal(bytes, j)
}

// Value 实现driver.Valuer接口
func (j JSONIntArray) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}
