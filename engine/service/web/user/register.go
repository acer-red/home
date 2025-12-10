package user

import (
	"time"

	"github.com/acer-red/official/engine/service/cache"
	"github.com/acer-red/official/engine/service/modb"
	"github.com/acer-red/official/engine/service/web/common"
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
func issueRegisterCookie(c *gin.Context, req modb.RequestUserRegister, id string) bool {
	expireAt := util.JWTDefaultExpireAt()

	token, err := util.CreateJWT(id, id, req.Username, req.Email, req.Category, time.Until(expireAt))
	if err != nil {
		log.Errorf("issueRegisterCookie CreateJWT error: %v", err)
		common.InternalServerError(c)
		return false
	}

	_, err = cache.SaveSession(id, id, token, req.Category, expireAt)
	if err != nil {
		log.Errorf("issueRegisterCookie SaveSession error: %v", err)
		common.InternalServerError(c)
		return false
	}

	common.SetJWTCookie(c, token, expireAt)
	return true
}

func userRegisterVisitor(c *gin.Context) {
	log.Info("游客注册")

	var req modb.RequestUserRegister

	type response struct {
		ID  string     `json:"id"`
		API []modb.API `json:"api"`
	}

	var err error
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		common.BadRequest(c)
		return
	}
	if !req.IsFromIndex() {
		common.BadRequest(c)
		return
	}

	if err := req.CheckAndSetCatetory(); err != nil {
		common.BadRequest(c)
		return
	}
	// 账号使用随机生成
	req.RandomAccount()

	id, api, err := req.Register(util.RoleVisitor)
	if err != nil {
		log.Errorf("userRegisterVisitor Register error: %v", err)
		common.InternalServerError(c)
		return
	}

	if err := req.BuildProfile(); err != nil {
		log.Errorf("userRegisterVisitor BuildProfile error: %v", err)
		req.CancelRegister()
		common.InternalServerError(c)
		return

	}

	if ok := issueRegisterCookie(c, req, id); !ok {
		return
	}

	log.Info("游客注册成功")

	// 不同的注册源，返回不同的验证方式
	common.OkData(c, response{ID: id, API: []modb.API{api}})

}

func userRegisterNormal(c *gin.Context) {
	log.Info("用户注册")

	var req modb.RequestUserRegister
	type response struct {
		ID  string     `json:"id"`
		API []modb.API `json:"api"`
	}
	var err error
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		common.BadRequest(c)
		return
	}

	if ok := req.Check(); !ok {
		common.BadRequest(c)
		return
	}

	user, found, err := req.FindUser()
	if err != nil {
		log.Errorf("userRegisterNormal FindUser error: %v", err)
		common.InternalServerError(c)
		return
	}

	var id string
	isNewUser := false

	if found {
		if user.IsDeleted {
			id, _, err = req.Reactivate(user)
			if err != nil {
				log.Errorf("userRegisterNormal Reactivate error: %v", err)
				common.InternalServerError(c)
				return
			}
		} else {
			common.Conflict(c)
			return
		}
	} else {
		isNewUser = true
		id, _, err = req.Register(util.RoleNormal)
		if err != nil {
			log.Errorf("userRegisterNormal Register error: %v", err)
			common.InternalServerError(c)
			return
		}
	}

	if err := req.BuildProfile(); err != nil {
		log.Errorf("userRegisterNormal BuildProfile error: %v", err)
		if isNewUser {
			req.CancelRegister()
		}
		common.InternalServerError(c)
		return

	}

	if ok := issueRegisterCookie(c, req, id); !ok {
		return
	}

	log.Info("用户注册成功")

	common.CreatedData(c, response{ID: id})
}
