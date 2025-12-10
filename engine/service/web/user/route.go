package user

import (
	"github.com/acer-red/official/engine/service/web/common"
	"github.com/gin-gonic/gin"
)

func RouteUser(c *gin.Engine) {
	v1 := c.Group("/api/v1")
	{
		v1User := v1.Group("/user")
		{
			v1User.POST("/register", userRegister)
			v1User.POST("/login", userLogin)
			v1User.GET("/randomonfo", userRandomInfo)
			v1User.Use(common.Auth())
			v1User.POST("/autologin", userAutoLogin)
			v1User.POST("/logout", userLogout)
			v1User.DELETE("/info", userDelete)
			v1User.DELETE("/app", userUnbindApp)
			v1User.GET("/info", getUserInfo)
			v1User.PUT("/info", putUserInfo)
			v1User.PUT("/profile", putUserProfile)
		}
	}
}
