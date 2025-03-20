package sys

import (
	"fmt"
	"strings"

	rs "github.com/acer-red/randResources"
	"github.com/google/uuid"
)

func CreateUUID() string {
	u, err := uuid.NewV7()
	if err != nil {
		fmt.Println(err)
	}
	return strings.ReplaceAll(u.String(), "-", "")
}
func CreateAPIKey() string {
	return strings.ToUpper(CreateUUID())

}
func RandomNickname() string {
	return rs.Text()
}
func RandomAvatar(random string) []byte {
	i, err := rs.BuildImage(random)
	if err != nil {
		fmt.Println(err)
	}
	return i.Bytes()
}
