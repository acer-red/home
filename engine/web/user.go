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
			v1User.POST("/logout", userLogout)
		}
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
	id, err := req.Register()
	if err != nil {
		internalServerError(c)
		return
	}

	c.SetCookie(
		req.Cookie.Key,                // Cookie 的名称
		req.Cookie.Value,              // Cookie 的值
		int(req.Cookie.EXTime.Unix()), // Cookie 的过期时间 (Unix 时间戳)
		"/",                           // Cookie 的路径 (通常设置为 "/")
		"",                            // Cookie 的域名 (留空表示当前域名)
		false,                         // 是否只允许 HTTPS 访问
		false,                         // 是否禁止 JavaScript 访问 (HttpOnly)
	)

	createdData(c, responseOK{ID: id})

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
		log.Warn("check failed")
		badRequest(c)
		return
	}

	if ok, err := req.Find(); err != nil {
		internalServerError(c)
		return
	} else if !ok {
		log.Warn("find failed")
		badRequest(c)
		return
	}

	err := req.ComparePassword()
	if err != nil {
		log.Warn("compare password")
		badRequest(c)
		return
	}

	id, err := req.GetID()
	if err != nil {
		internalServerError(c)
		return
	}

	req.GetCookie()
	c.SetCookie(
		req.Cookie.Key,                // Cookie 的名称
		req.Cookie.Value,              // Cookie 的值
		int(req.Cookie.EXTime.Unix()), // Cookie 的过期时间 (Unix 时间戳)
		"/",                           // Cookie 的路径 (通常设置为 "/")
		"",                            // Cookie 的域名 (留空表示当前域名)
		false,                         // 是否只允许 HTTPS 访问
		false,                         // 是否禁止 JavaScript 访问 (HttpOnly)
	)

	okData(c, response{ID: id})
}

func userLogout(c *gin.Context) {
	var req modb.RequestUserLogout

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
		notFound(c)
		return
	}

	if err := req.DeleteCookie(); err != nil {
		internalServerError(c)
		return
	}

	c.SetCookie(
		"login", // Cookie 的名称
		"",      // Cookie 的值
		0,       // Cookie 的过期时间 (Unix 时间戳)
		"/",     // Cookie 的路径 (通常设置为 "/")
		"",      // Cookie 的域名 (留空表示当前域名)
		false,   // 是否只允许 HTTPS 访问
		false,   // 是否禁止 JavaScript 访问 (HttpOnly)
	)

	okData(c, nil)
}
