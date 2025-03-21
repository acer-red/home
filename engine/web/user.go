package web

import (
	"modb"
	"sys"

	"github.com/gin-gonic/gin"
	log "github.com/tengfei-xy/go-log"
)

func RouteUser(c *gin.Engine) {
	v1 := c.Group("/api/v1")
	{
		// v1.Use(outputRequestHeader())
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

	var req modb.RequestUserRegister
	type response struct {
		ID  string     `json:"id"`
		API []modb.API `json:"api"`
	}
	var err error
	log.Info("用户注册")
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

	id, api, err := req.Register()
	if err != nil {
		internalServerError(c)
		return
	}

	if err := req.BuildProfile(); err != nil {
		req.CancelRegister()
		internalServerError(c)
		return

	}

	log.Info("用户注册成功")

	// 不同的注册源，返回不同的验证方式
	if !req.IsFromIndex() {
		log.Debug3f("%s", api.APIKey)
		okData(c, response{ID: id, API: []modb.API{api}})
		return
	}

	req.GetCookie()
	setCookie(c, req.Cookie.Key, req.Cookie.Value, int(req.Cookie.EXTime.Unix()))
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
	} else if !ok {
		badRequest(c)
		return
	}

	err := req.ComparePassword()
	if err != nil {
		badRequest(c)
		return
	}

	// 官网注册方式，返回 cookie
	if req.IsFromIndex() {
		req.GetCookie()
		setCookie(c, req.Cookie.Key, req.Cookie.Value, int(req.Cookie.EXTime.Unix()))
	}
	okData(c, req.Login())

}
func userLogout(c *gin.Context) {
	log.Info("用户注销")

	user := c.MustGet("user").(modb.User)
	if err := user.DeleteCookie(); err != nil {
		internalServerError(c)
		return
	}

	setCookie(c, "login", "", 0)
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
