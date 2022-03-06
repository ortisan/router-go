package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	objAsString = "{\"id\":\"123456\"}"
)

type TestObject struct {
	Id string `json:"id,omitempty"`
}

func newTestObject() TestObject {
	return TestObject{Id: "123456"}
}

func TestGetSubstringAfter(t *testing.T) {
	str := "http://testing.com/api/xpto"
	strRes := GetSubstringAfter(str, "api/")
	assert.Equal(t, "xpto", strRes)
}

func TestObjectToJson(t *testing.T) {
	obj := newTestObject()
	objBytes, _ := ObjectToJson(obj)
	assert.Equal(t, objAsString, string(objBytes))
}

func TestStringToObject(t *testing.T) {
	obj := TestObject{}
	StringToObject(objAsString, &obj)
	assert.Equal(t, newTestObject(), obj)
}
