package util

import (
	"regexp"
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
	objStr, _ := ObjectToJson(obj)
	matched, _ := regexp.MatchString("\"id\":\"123456\"", string(objStr))
	assert.Equal(t, objAsString, matched)
}

func TestStringToObject(t *testing.T) {
	obj = make(TestObject)

	StringToObject(objAsString, TestObject{})
}
