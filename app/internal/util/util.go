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

func ObjectToJson(object interface{}) []byte {
	mobj, err := json.Marshal(object)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Error to marshal object")
		return []byte{}
	}
	return mobj
}
