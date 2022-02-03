package util

import (
	"encoding/json"
	"strings"

	"github.com/rs/zerolog/log"
)

func GetSubstringAfter(value string, a string) string {
	pos := strings.LastIndex(value, a)
	if pos == -1 {
		return ""
	}
	adjustedPos := pos + len(a)
	if adjustedPos >= len(value) {
		return ""
	}
	return value[adjustedPos:]
}

func ObjectToJson(object interface{}) []byte {
	mobj, err := json.Marshal(object)
	if err != nil {
		print()
		log.Error().Stack().Err(err).Msg("Iterating headers...")
		return []byte{}
	}
	return mobj
}
