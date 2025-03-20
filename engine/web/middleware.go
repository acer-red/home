package web

import (
	"github.com/gin-gonic/gin"
	log "github.com/tengfei-xy/go-log"
	"modb"
	"net/http"
	"net/url"
	"strings"
)

func getHost(origin string) string {
	uri, err := url.Parse(origin)
	if err != nil {
		log.Errorf("Failed to parse origin: %s, error: %v", origin, err)
		return ""
	}
	if strings.Contains(uri.Host, ":") {
		return strings.Split(uri.Host, ":")[0]
	}
	return uri.Host
}

func cors(origin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("cors_origin", getHost(origin))

		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, Cookie")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
func auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("login")
		if err != nil {
			unauthorized(c)
			return
		}
		u, exist, err := modb.GetUser(cookie)
		if err != nil {
			internalServerError(c)
			return
		}
		if !exist {
			unauthorized(c)
			return
		}

		c.Set("user", u)
		c.Next()
	}
}
