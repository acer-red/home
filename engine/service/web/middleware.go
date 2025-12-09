package web

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/acer-red/official/engine/service/cache"
	"github.com/acer-red/official/engine/service/modb"
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
			log.Debug("cookie通过")
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
			log.Debug("API通过")
			c.Set("user", u)
			c.Next()
			return
		}

		log.Debug("无任何认证方式")
		unauthorized(c)

	}
}
func authCookie(c *gin.Context) (modb.User, bool, error) {
	log.Debug3("尝试Cookie认证")
	cookie, err := c.Cookie("jwt")

	if err != nil {
		if err == http.ErrNoCookie {
			log.Debug3("无cookie")
			return modb.User{}, false, nil
		}
		return modb.User{}, false, err
	}

	uid, category, claims, err := cache.GetUserFromCookie(cookie)
	if err != nil {
		return modb.User{}, false, err
	}
	if claims == nil {
		return modb.User{}, false, nil
	}

	c.Set("claims", *claims)
	return modb.GetUserByIDAndCategory(uid, category)
}
func authAPI(c *gin.Context) (modb.User, bool, error) {
	log.Debug3f("尝试API认证")

	api := c.Request.Header.Get("Authorization")
	if api == "" {
		return modb.User{}, false, nil
	}

	return modb.GetUserFromAPI(api)
}
