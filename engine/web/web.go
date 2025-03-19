package web

import (
	"fmt"
	"sys"

	"github.com/gin-gonic/gin"
	log "github.com/tengfei-xy/go-log"
)

func Init(env sys.Web) {
	gin.SetMode(gin.ReleaseMode)
	g := gin.Default()

	g.Use(setEnv(env))

	if env.CORS.Enable {
		log.Infof("Enable CORS, Origin:%s", env.CORS.AllowOrigin)
		g.Use(cors(env.CORS.AllowOrigin))
	} else {

	}

	// 设定路由
	RouteFeedback(g)
	RouteUser(g)
	RouterImageGet(g)

	if env.Server.SslEnable {
		err := g.RunTLS(fmt.Sprintf(":%d", env.Server.Port), env.Server.CrtFile, env.Server.KeyFile)
		if err != nil {
			log.Fatal(err)
		}
		return
	} else {
		err := g.Run(fmt.Sprintf(":%d", env.Server.Port))
		if err != nil {
			log.Fatal(err)
		}
	}

}
func setEnv(env sys.Web) gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Set("env", env)
		c.Set("cors_origin", "")
	}
}

func setCookie(c *gin.Context, key, value string, ex int) {
	c.SetCookie(
		key,                               // Cookie 的名称
		value,                             // Cookie 的值
		ex,                                // Cookie 的过期时间 (Unix 时间戳)
		"/",                               // Cookie 的路径 (通常设置为 "/")
		c.MustGet("cors_origin").(string), // Cookie 的域名 (留空表示当前域名)
		false,                             // 是否只允许 HTTPS 访问
		false,                             // 是否禁止 JavaScript 访问 (HttpOnly)
	)
}
