package sys

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func CreateUUID() string {
	u, err := uuid.NewV7()
	if err != nil {
		fmt.Println(err)
	}
	return strings.ReplaceAll(u.String(), "-", "")
}
