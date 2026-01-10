package common

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/acer-red/official/engine/service/cache"
	"github.com/acer-red/official/engine/service/db"
	Err "github.com/acer-red/official/engine/service/web/error"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func Cors(origin string) gin.HandlerFunc {
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
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {

		u, pid, exist, err := authCookie(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, Err.Unauthorized.JSON())

			return
		}
		if exist {
			if u.IsDeleted {
				c.AbortWithStatusJSON(http.StatusGone, Err.AlreadyDeleted.JSON())
				return
			}
			log.Debug("cookie通过")
			c.Set("user", u)
			if pid != nil {
				c.Set("product_id", *pid)
			}
			c.Next()
			return
		}

		u, pid, exist, err = authAPI(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, Err.Unauthorized.JSON())

			return
		}
		if exist {
			if u.IsDeleted {
				c.AbortWithStatusJSON(http.StatusGone, Err.AlreadyDeleted.JSON())
				return
			}
			log.Debug("API通过")
			c.Set("user", u)
			if pid != nil {
				c.Set("product_id", *pid)
			}
			c.Next()
			return
		}

		log.Debug("无任何认证方式")
		c.AbortWithStatusJSON(http.StatusUnauthorized, Err.Unauthorized.JSON())

	}
}
func authCookie(c *gin.Context) (db.User, *uuid.UUID, bool, error) {
	log.Debug3("尝试Cookie认证")
	cookie, err := c.Cookie("jwt")

	if err != nil {
		if err == http.ErrNoCookie {
			log.Debug3("无cookie")
			return db.User{}, nil, false, nil
		}
		return db.User{}, nil, false, err
	}

	uid, category, claims, err := cache.GetUserFromCookie(cookie)
	if err != nil {
		return db.User{}, nil, false, err
	}
	if claims == nil {
		return db.User{}, nil, false, nil
	}

	c.Set("claims", *claims)
	return db.GetUserByIDAndCategory(uid, category)
}
func authAPI(c *gin.Context) (db.User, *uuid.UUID, bool, error) {
	log.Debug3f("尝试API认证")

	api := c.Request.Header.Get("Authorization")
	if api == "" {
		return db.User{}, nil, false, nil
	}

	return db.GetUserFromAPI(api)
}
