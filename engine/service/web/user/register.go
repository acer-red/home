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

// / /api/v1/user/register
func userRegister(c *gin.Context) {
	if c.Query("visitor") == "1" {
		userRegisterVisitor(c)
		return
	}
	userRegisterNormal(c)

}

// issueRegisterCookie builds a JWT for the newly registered user and attaches it via Set-Cookie.
func issueRegisterCookie(c *gin.Context, req db.RequestUserRegister, id string) bool {
	expireAt := util.JWTDefaultExpireAt()

	token, err := util.CreateJWT(id, id, req.Username, req.Email, req.Category, time.Until(expireAt))
	if err != nil {
		log.Errorf("issueRegisterCookie CreateJWT error: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return false
	}

	_, err = cache.SaveSession(id, id, token, req.Category, expireAt)
	if err != nil {
		log.Errorf("issueRegisterCookie SaveSession error: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return false
	}

	common.SetJWTCookie(c, token, expireAt)
	return true
}

func userRegisterVisitor(c *gin.Context) {
	log.Info("游客注册")

	var req db.RequestUserRegister

	type response struct {
		ID        string   `json:"id"`
		API       []db.API `json:"api"`
		ProductID string   `json:"product_id,omitempty"`
		DeviceID  string   `json:"device_id,omitempty"`
	}

	var err error
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, Err.DataParseFailed)
		return
	}
	if !req.IsFromIndex() {
		c.AbortWithStatusJSON(http.StatusInternalServerError, Err.InternalServer.JSON())
		return
	}

	if err := req.CheckAndSetCatetory(); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, Err.InternalServer.JSON())
		return
	}
	// 账号使用随机生成
	req.RandomAccount()

	id, api, deviceID, err := req.Register(util.RoleVisitor)
	if err != nil {
		log.Errorf("userRegisterVisitor Register error: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if err := req.BuildProfile(id); err != nil {
		log.Errorf("userRegisterVisitor BuildProfile error: %v", err)
		req.CancelRegister(id)
		c.AbortWithStatus(http.StatusInternalServerError)
		return

	}

	if ok := issueRegisterCookie(c, req, id); !ok {
		return
	}

	log.Info("游客注册成功")

	// 不同的注册源，返回不同的验证方式
	c.JSON(http.StatusCreated, response{
		ID:        id,
		API:       []db.API{api},
		ProductID: api.ProductID.String(),
		DeviceID:  deviceID,
	})

}

func userRegisterNormal(c *gin.Context) {
	log.Info("用户注册")

	var req db.RequestUserRegister
	type response struct {
		ID        string `json:"user_id"`
		ProductID string `json:"product_id"`
		DeviceID  string `json:"device_id"`
	}
	var err error
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		log.Debug2(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, Err.InternalServer.JSON())
		return
	}

	if ok := req.Check(); !ok {
		log.Debug2("Check failed")
		c.AbortWithStatusJSON(http.StatusInternalServerError, Err.InternalServer.JSON())
		return
	}

	user, found, err := req.FindUser()
	if err != nil {
		log.Errorf("userRegisterNormal FindUser error: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var id string
	var api db.API
	var deviceID string
	isNewUser := false

	if found {
		if user.IsDeleted {
			id, api, deviceID, err = req.Reactivate(user)
			if err != nil {
				log.Errorf("userRegisterNormal Reactivate error: %v", err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		} else {
			log.Debug2("用户已存在")
			c.AbortWithStatusJSON(http.StatusConflict, Err.UserAlreadyExist.JSON())
			return
		}
	} else {
		isNewUser = true
		id, api, deviceID, err = req.Register(util.RoleNormal)
		if err != nil {
			log.Errorf("userRegisterNormal Register error: %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	if err := req.BuildProfile(id); err != nil {
		log.Errorf("userRegisterNormal BuildProfile error: %v", err)
		if isNewUser {
			log.Debug2("BuildProfile failed for new user")
			req.CancelRegister(id)
		}
		log.Debug2(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return

	}

	if ok := issueRegisterCookie(c, req, id); !ok {
		return
	}

	log.Info("用户注册成功")

	c.JSON(http.StatusOK, Err.OK.Data(response{
		ID:        id,
		ProductID: api.ProductID.String(),
		DeviceID:  deviceID,
	}))
}
