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

		v1.Use(outputRequestHeader())
		v1User := v1.Group("/user")
		{
			v1User.POST("/register", userRegister)
			v1User.POST("/login", userLogin)
			v1User.GET("/randomonfo", userRandomInfo)
			v1User.Use(auth())
			v1User.POST("/autologin", userAutoLogin)
			v1User.POST("/logout", userLogout)
			v1User.GET("/info", getUserInfo)
			v1User.PUT("/info", putUserInfo)
		}
	}

}

func userRegister(c *gin.Context) {

	var req modb.RequestUserRegister

	type response struct {
		ID string `json:"id"`
	}
	var err error
	log.Info("用户注册")

	req.Category, err = sys.GetCategory(req.CategoryStr)
	if err != nil {
		badRequest(c)
		return
	}

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

	if err := req.BuildProfile(); err != nil {
		internalServerError(c)
		return
	}

	id, err := req.Register()
	if err != nil {
		internalServerError(c)
		return
	}

	req.GetCookie()
	setCookie(c, req.Cookie.Key, req.Cookie.Value, int(req.Cookie.EXTime.Unix()))

	log.Infof("用户注册成功 产品:%s", string(req.Category))
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

	req.GetCookie()

	setCookie(c, req.Cookie.Key, req.Cookie.Value, int(req.Cookie.EXTime.Unix()))

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
