package web

import (
	"modb"

	"github.com/gin-gonic/gin"
	log "github.com/tengfei-xy/go-log"
)

func RouteUser(c *gin.Engine) {
	v1 := c.Group("/api/v1")
	{
		v1User := v1.Group("/user")
		{
			v1User.POST("/register", userRegister)
			v1User.POST("/login", userLogin)
			v1User.Use(authMiddleware())
			v1User.POST("/autologin", userAutoLogin)
			v1User.POST("/logout", userLogout)
			v1User.GET("/info", getUserInfo)
			v1User.PUT("/info", putUserInfo)
		}
	}

}
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("login")
		if err != nil {
			unauthorized(c)
			return
		}
		u, exist, err := modb.GetUser(cookie)
		if err != nil {
			internalServerError(c)
			return
		}
		if !exist {
			unauthorized(c)
			return
		}

		c.Set("user", u)
		c.Next()
	}
}

func userRegister(c *gin.Context) {

	var req modb.RequestUserRegister

	type responseOK struct {
		ID string `json:"id"`
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
	createdData(c, responseOK{ID: id})

}
func userAutoLogin(c *gin.Context) {
	user := c.MustGet("user").(modb.User)
	okData(c, user)
}
func userLogin(c *gin.Context) {
	var req modb.RequestUserLogin
	type response struct {
		ID string `json:"id"`
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
	} else if !ok {
		badRequest(c)
		return
	}

	err := req.ComparePassword()
	if err != nil {
		log.Warn("compare password")
		badRequest(c)
		return
	}

	req.GetCookie()

	setCookie(c, req.Cookie.Key, req.Cookie.Value, int(req.Cookie.EXTime.Unix()))

	id := req.GetInfo().ID
	okData(c, response{ID: id})
}
func userLogout(c *gin.Context) {

	user := c.MustGet("user").(modb.User)
	if err := user.DeleteCookie(); err != nil {
		internalServerError(c)
		return
	}

	setCookie(c, "login", "", 0)
	okData(c, nil)
}
func getUserInfo(c *gin.Context) {
	user := c.MustGet("user").(modb.User)
	okData(c, user)
}
func putUserInfo(c *gin.Context) {

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
