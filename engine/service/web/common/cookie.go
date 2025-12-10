package common

import (
	"time"

	"github.com/gin-gonic/gin"
)

func SetCookie(c *gin.Context, key, value string, ex int) {
	c.SetCookie(
		key,                               // Cookie 的名称
		value,                             // Cookie 的值
		ex,                                // Cookie 的过期时间 (Unix 时间戳)
		"/",                               // Cookie 的路径 (通常设置为 "/")
		c.MustGet("cors_origin").(string), // Cookie 的域名 (留空表示当前域名)
		false,                             // 是否只允许 HTTPS 访问
		false,                             // 是否禁止 JavaScript 访问 (HttpOnly)
	)
}

func SetJWTCookie(c *gin.Context, token string, expireAt time.Time) {
	maxAge := int(time.Until(expireAt).Seconds())
	if token == "" {
		maxAge = -1 // instruct browser to delete the cookie
		expireAt = time.Unix(0, 0)
	} else if maxAge <= 0 {
		maxAge = int((30 * 24 * time.Hour).Seconds())
	}

	c.SetCookie(
		"jwt",
		token,
		maxAge,
		"/",
		c.MustGet("cors_origin").(string),
		false,
		true,
	)
}
