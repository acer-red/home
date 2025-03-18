package sys

import (
	"encoding/json"
)

func JsonPrettyPrint(data any) string {
	ret, _ := json.MarshalIndent(data, "", "  ")
	return string(ret)
}
