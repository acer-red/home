package user

import (
	"net/http"
	"time"

	"github.com/acer-red/official/engine/service/cache"
	"github.com/acer-red/official/engine/service/db"
	"github.com/acer-red/official/engine/service/web/common"
	Err "github.com/acer-red/official/engine/service/web/error"
	"github.com/acer-red/official/engine/util"
	"github.com/gin-gonic/gin"
	log "github.com/tengfei-xy/go-log"
)

func userAutoLogin(c *gin.Context) {
	log.Info("用户自动登陆")

	user := c.MustGet("user").(db.User)
	apis, _ := db.GetAPIsForUser(user.ID)
	c.JSON(http.StatusOK, Err.OK.Data(user.ToResponse(apis)))

}
func userLogin(c *gin.Context) {
	var req db.RequestUserLogin
	log.Info("用户登陆")

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		log.Warnf("login bind body failed: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, Err.InternalServer.JSON())
		return
	}

	if ok := req.Check(); !ok {
		log.Warnf("login param check failed account=%s category=%s", req.Account, req.Category)
		c.AbortWithStatusJSON(http.StatusInternalServerError, Err.InternalServer.JSON())
		return
	}

	if ok, err := req.Find(); err != nil {
		log.Errorf("login find user error account=%s err=%v", req.Account, err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return

	} else if !ok {
		log.Warnf("login user not found account=%s", req.Account)
		c.AbortWithStatusJSON(http.StatusInternalServerError, Err.InternalServer.JSON())
		return
	}

	err := req.ComparePassword()
	if err != nil {
		log.Warnf("login password mismatch account=%s err=%v", req.Account, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, Err.InternalServer.JSON())
		return
	}
	if err = req.CheckNewDevice(); err != nil {
		log.Errorf("userLogin CheckNewDevice error: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	jwtExpireAt := util.JWTDefaultExpireAt()

	if _, err := c.Cookie("jwt"); err == http.ErrNoCookie {
		if token, claims, found, err := cache.GetSessionCookieByUID(req.CategoryValue(), req.GetAccountID()); err != nil {
			log.Warnf("login lookup session failed account=%s err=%v", req.Account, err)
		} else if found && claims != nil && claims.ExpiresAt != nil {
			res := req.BuildLoginResponse()
			common.SetJWTCookie(c, token, claims.ExpiresAt.Time)
			c.JSON(http.StatusOK, Err.OK.Data(res))
			return
		}
	}

	res, token, err := req.Login(jwtExpireAt)

	if err != nil {
		log.Errorf("userLogin Login error: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// 为请求头设置set-cookie
	common.SetJWTCookie(c, token, jwtExpireAt)
	c.JSON(http.StatusOK, Err.OK.Data(res))

}
func userLogout(c *gin.Context) {
	log.Info("用户注销")

	// user := c.MustGet("user").(db.User)
	claims, exist := c.Get("claims")
	if exist {
		if err := cache.DeleteSession(claims.(util.JWTClaims)); err != nil {
			log.Warnf("delete redis login session failed: %v", err)
		}
	}

	// 向客户端删除cookie
	common.SetJWTCookie(c, "", time.Now())
	c.JSON(http.StatusOK, Err.OK)
}
