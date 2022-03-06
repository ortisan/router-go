package util

import (
	"encoding/json"
	"strings"

	"github.com/rs/zerolog/log"
)

func GetSubstringAfter(value string, matchString string) string {
	pos := strings.LastIndex(value, matchString)
	if pos == -1 {
		return ""
	}
	adjustedPos := pos + len(matchString)
	if adjustedPos >= len(value) {
		return ""
	}
	return value[adjustedPos:]
}

func ObjectToJsonBytes(object interface{}) ([]byte, error) {
	objRes, err := json.Marshal(object)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Error to marshal object")
		return nil, err
	}
	return objRes, nil
}

func ObjectToJsonStr(object interface{}) (string, error) {
	jsonBytes, err := ObjectToJsonBytes(object)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

func StringToObject(strObj string, object interface{}) (interface{}, error) {
	err := json.Unmarshal([]byte(strObj), object)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Error to marshal object")
		return nil, err
	}
	return object, nil
}
