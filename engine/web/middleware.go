package web

import (
	"fmt"
	"modb"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/tengfei-xy/go-log"
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

		u, exist, err := authCookie(c)
		if err != nil {
			internalServerError(c)
			return
		}
		if exist {
			c.Set("user", u)
			c.Next()
			return
		}
		u, exist, err = authAPI(c)
		if err != nil {
			internalServerError(c)
			return
		}
		if exist {
			c.Set("user", u)
			c.Next()
			return
		} else {
			unauthorized(c)
			return
		}
	}
}
func authCookie(c *gin.Context) (modb.User, bool, error) {
	cookie, err := c.Cookie("login")

	if err != nil {
		if err == http.ErrNoCookie {
			return modb.User{}, false, nil
		}
		return modb.User{}, false, err
	}

	return modb.GetUserFromCookie(cookie)
}
func authAPI(c *gin.Context) (modb.User, bool, error) {
	api := c.Request.Header.Get("Authorization")
	if api == "" {
		return modb.User{}, false, nil
	}
	return modb.GetUserFromAPI(api)

}

func outputRequestHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		for key, values := range c.Request.Header {
			for _, value := range values {
				fmt.Printf("%s: %s\n", key, value)
			}
		}
		c.Next()
	}
}
