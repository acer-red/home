package web

import (
	"net/http"
	"time"

	"github.com/acer-red/home/engine/modb"
	"github.com/acer-red/home/engine/sys"

	"github.com/acer-red/home/engine/storage"
	"github.com/gin-gonic/gin"
	log "github.com/tengfei-xy/go-log"
)

func RouteUser(c *gin.Engine) {
	v1 := c.Group("/api/v1")
	{
		v1.Use(outputRequestHeader())
		v1User := v1.Group("/user")
		{
			v1User.POST("/register", userRegister)
			v1User.POST("/login", userLogin)
			v1User.GET("/randomonfo", userRandomInfo)
			v1User.Use(auth())
			v1User.POST("/autologin", userAutoLogin)
			v1User.POST("/logout", userLogout)
			v1User.DELETE("/info", userDelete)
			v1User.GET("/info", getUserInfo)
			v1User.PUT("/info", putUserInfo)
			v1User.PUT("/profile", putUserProfile)
		}
	}

}

func userRegister(c *gin.Context) {
	if c.Query("visitor") == "1" {
		userRegisterVisitor(c)
		return
	}
	userRegisterNormal(c)

}

// issueRegisterCookie builds a JWT for the newly registered user and attaches it via Set-Cookie.
func issueRegisterCookie(c *gin.Context, req modb.RequestUserRegister, id string) bool {
	expireAt := sys.JWTDefaultExpireAt()
	jti := sys.CreateUUID()

	token, err := sys.CreateJWT(jti, id, req.Username, req.Email, req.Category, time.Until(expireAt))
	if err != nil {
		internalServerError(c)
		return false
	}

	_, err = storage.SaveSession(jti, id, token, req.Category, expireAt)
	if err != nil {
		internalServerError(c)
		return false
	}

	setJWTCookie(c, token, expireAt)
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
		badRequest(c)
		return
	}
	if !req.IsFromIndex() {
		badRequest(c)
		return
	}

	if err := req.CheckAndSetCatetory(); err != nil {
		badRequest(c)
		return
	}
	// 账号使用随机生成
	req.RandomAccount()

	id, api, err := req.Register(sys.RoleVisitor)
	if err != nil {
		internalServerError(c)
		return
	}

	if err := req.BuildProfile(); err != nil {
		req.CancelRegister()
		internalServerError(c)
		return

	}

	if ok := issueRegisterCookie(c, req, id); !ok {
		return
	}

	log.Info("游客注册成功")

	// 不同的注册源，返回不同的验证方式
	okData(c, response{ID: id, API: []modb.API{api}})

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
		badRequest(c)
		return
	}

	if ok := req.Check(); !ok {
		badRequest(c)
		return
	}

	if ok, err := req.Find(); err != nil {
		internalServerError(c)
		return
	} else if ok {
		conflict(c)
		return
	}

	id, _, err := req.Register(sys.RoleNormal)
	if err != nil {
		internalServerError(c)
		return
	}

	if err := req.BuildProfile(); err != nil {
		req.CancelRegister()
		internalServerError(c)
		return

	}

	if ok := issueRegisterCookie(c, req, id); !ok {
		return
	}

	log.Info("用户注册成功")

	createdData(c, response{ID: id})
}
func userAutoLogin(c *gin.Context) {
	log.Info("用户自动登陆")

	user := c.MustGet("user").(modb.User)
	okData(c, user)
}
func userLogin(c *gin.Context) {
	var req modb.RequestUserLogin
	log.Info("用户登陆")

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		log.Warnf("login bind body failed: %v", err)
		badRequest(c)
		return
	}

	if ok := req.Check(); !ok {
		log.Warnf("login param check failed account=%s category=%s", req.Account, req.Category)
		badRequest(c)
		return
	}

	if ok, err := req.Find(); err != nil {
		log.Errorf("login find user error account=%s err=%v", req.Account, err)
		internalServerError(c)
		return

	} else if !ok {
		log.Warnf("login user not found account=%s", req.Account)
		badRequest(c)
		return
	}

	err := req.ComparePassword()
	if err != nil {
		log.Warnf("login password mismatch account=%s err=%v", req.Account, err)
		badRequest(c)
		return
	}
	if err = req.CheckNewDevice(); err != nil {
		log.Error(err)
		internalServerError(c)
		return
	}

	jwtExpireAt := sys.JWTDefaultExpireAt()

	if _, err := c.Cookie("jwt"); err == http.ErrNoCookie {
		if token, claims, found, err := storage.GetSessionCookieByUID(req.CategoryValue(), req.UserID()); err != nil {
			log.Warnf("login lookup session failed account=%s err=%v", req.Account, err)
		} else if found && claims != nil && claims.ExpiresAt != nil {
			res := req.BuildLoginResponse()
			setJWTCookie(c, token, claims.ExpiresAt.Time)
			okData(c, res)
			return
		}
	}

	res, token, err := req.Login(jwtExpireAt)

	if err != nil {
		internalServerError(c)
		return
	}

	// 为请求头设置set-cookie
	setJWTCookie(c, token, jwtExpireAt)
	okData(c, res)

}
func userLogout(c *gin.Context) {
	log.Info("用户注销")

	// user := c.MustGet("user").(modb.User)
	claims, exist := c.Get("claims")
	if exist {
		if err := storage.DeleteSession(claims.(sys.JWTClaims)); err != nil {
			log.Warnf("delete redis login session failed: %v", err)
		}
	}

	// 向客户端删除cookie
	setJWTCookie(c, "", time.Now())
	okData(c, nil)
}
func getUserInfo(c *gin.Context) {
	log.Info("用户获取信息")

	user := c.MustGet("user").(modb.User)
	okData(c, user)
}
func putUserInfo(c *gin.Context) {
	log.Info("用户修改信息")
	var req modb.RequestPutUserInfo
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		badRequest(c)
		return
	}
	req.UOID = c.MustGet("user").(modb.User).UOID

	if err := req.Update(); err != nil {
		internalServerError(c)
		return
	}
	ok(c)
}
func userRandomInfo(c *gin.Context) {
	log.Info("用户随机信息")

	type response struct {
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
	}

	okData(c, response{
		Nickname: sys.RandomNickname(),
		Avatar:   sys.RandomAvatarBase64(sys.CreateUUID()),
	})
}
func putUserProfile(c *gin.Context) {

	type response struct {
		URL string `json:"url"`
	}
	hasNickname := false
	hasAvatar := false
	user := c.MustGet("user").(modb.User)
	if user.IsNoID() {
		internalServerError(c)
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		badRequest(c)
		return
	}

	nicknames := form.Value["nickname"]

	if len(nicknames) != 0 {
		nickname := nicknames[0]
		log.Info("更新用户昵称")

		hasNickname = true
		if err := user.UpdateNickname(nickname); err != nil {
			internalServerError(c)
			return
		}
	}
	avatars := form.File["avatar"]
	if len(avatars) != 0 {
		hasAvatar = true
	}

	// 没有头像和昵称
	if hasNickname {
		if !hasAvatar {
			log.Info("更新用户昵称完成")
			ok(c)
			return
		}
	} else {
		if !hasAvatar {
			badRequest(c)
			return
		}
	}
	log.Info("更新用户头像")

	ext := form.Value["ext"][0]
	if ext == "" {
		badRequest(c)
		return
	}

	file, err := avatars[0].Open()
	if err != nil {
		badRequest(c)
		return
	}
	defer file.Close()

	if err := user.UpdateAvatar(file, ext); err != nil {
		internalServerError(c)
		return
	}
	log.Info("更新用户头像完成")

	okData(c, response{URL: user.Profile.Avatar.URL})
}
func userDelete(c *gin.Context) {
	log.Info("用户删除")

	user := c.MustGet("user").(modb.User)
	if err := user.Delete(); err != nil {
		internalServerError(c)
		return
	}
	setCookie(c, "login", "", 0)
	ok(c)
}
