package util

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

	rs "github.com/acer-red/randResources"
)

func randString(length int) string {
	const pool = "qazwsxedcrfvtgbyhnujmikolpQAZWSXEDCRFVTGBYHNUJMIKOLP1234567890"
	bytes := make([]byte, length)
	poolLength := big.NewInt(int64(len(pool)))

	for i := 0; i < length; i++ {
		// 使用 crypto/rand 获取真随机索引
		// 它会返回一个 [0, poolLength) 范围内的安全随机数
		num, err := rand.Int(rand.Reader, poolLength)
		if err != nil {
			return "" // 处理系统熵源不足的罕见错误
		}
		bytes[i] = pool[num.Int64()]
	}

	return string(bytes)
}
func RandomString32() string {
	return randString(32)
}
func CreateAPIKey() string {
	return strings.ToUpper(RandomString32())
}

func RandomNickname() string {
	return rs.Text()
}

func RandomAvatarBytes() []byte {
	i, err := rs.BuildImage(randString(32))
	if err != nil {
		fmt.Println(err)
	}
	return i.Bytes()
}
func RandomAvatarBase64() string {
	i, err := rs.BuildImage(randString(32))
	if err != nil {
		fmt.Println(err)
	}
	return i.Base64()
}
