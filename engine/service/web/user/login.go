package user

import (
	"net/http"
	"time"

	"github.com/acer-red/official/engine/service/cache"
	"github.com/acer-red/official/engine/service/modb"
	"github.com/acer-red/official/engine/service/web/common"
	"github.com/acer-red/official/engine/util"
	"github.com/gin-gonic/gin"
	log "github.com/tengfei-xy/go-log"
)

func userAutoLogin(c *gin.Context) {
	log.Info("用户自动登陆")

	user := c.MustGet("user").(modb.User)
	common.OkData(c, user)
}
func userLogin(c *gin.Context) {
	var req modb.RequestUserLogin
	log.Info("用户登陆")

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		log.Warnf("login bind body failed: %v", err)
		common.BadRequest(c)
		return
	}

	if ok := req.Check(); !ok {
		log.Warnf("login param check failed account=%s category=%s", req.Account, req.Category)
		common.BadRequest(c)
		return
	}

	if ok, err := req.Find(); err != nil {
		log.Errorf("login find user error account=%s err=%v", req.Account, err)
		common.InternalServerError(c)
		return

	} else if !ok {
		log.Warnf("login user not found account=%s", req.Account)
		common.BadRequest(c)
		return
	}

	err := req.ComparePassword()
	if err != nil {
		log.Warnf("login password mismatch account=%s err=%v", req.Account, err)
		common.BadRequest(c)
		return
	}
	if err = req.CheckNewDevice(); err != nil {
		log.Errorf("userLogin CheckNewDevice error: %v", err)
		common.InternalServerError(c)
		return
	}

	jwtExpireAt := util.JWTDefaultExpireAt()

	if _, err := c.Cookie("jwt"); err == http.ErrNoCookie {
		if token, claims, found, err := cache.GetSessionCookieByUID(req.CategoryValue(), req.GetAccountID()); err != nil {
			log.Warnf("login lookup session failed account=%s err=%v", req.Account, err)
		} else if found && claims != nil && claims.ExpiresAt != nil {
			res := req.BuildLoginResponse()
			common.SetJWTCookie(c, token, claims.ExpiresAt.Time)
			common.OkData(c, res)
			return
		}
	}

	res, token, err := req.Login(jwtExpireAt)

	if err != nil {
		log.Errorf("userLogin Login error: %v", err)
		common.InternalServerError(c)
		return
	}

	// 为请求头设置set-cookie
	common.SetJWTCookie(c, token, jwtExpireAt)
	common.OkData(c, res)

}
func userLogout(c *gin.Context) {
	log.Info("用户注销")

	// user := c.MustGet("user").(modb.User)
	claims, exist := c.Get("claims")
	if exist {
		if err := cache.DeleteSession(claims.(util.JWTClaims)); err != nil {
			log.Warnf("delete redis login session failed: %v", err)
		}
	}

	// 向客户端删除cookie
	common.SetJWTCookie(c, "", time.Now())
	common.OkData(c, nil)
}
