package util

import (
	"encoding/json"
)

func JsonPrettyPrint(data any) string {
	ret, _ := json.MarshalIndent(data, "", "  ")
	return string(ret)
}
