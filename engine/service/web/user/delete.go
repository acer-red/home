package user

import (
	"github.com/acer-red/official/engine/service/modb"
	"github.com/acer-red/official/engine/service/web/common"
	"github.com/acer-red/official/engine/util"
	"github.com/gin-gonic/gin"
	log "github.com/tengfei-xy/go-log"
)

func userDelete(c *gin.Context) {
	log.Info("用户删除")

	user := c.MustGet("user").(modb.User)
	if err := user.Delete(); err != nil {
		common.InternalServerError(c)
		return
	}
	common.SetCookie(c, "login", "", 0)
	common.Ok(c)
}

func userUnbindApp(c *gin.Context) {
	log.Info("用户解绑应用")

	var req struct {
		Category string `json:"category"`
	}
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		common.BadRequest(c)
		return
	}

	category, valid := util.GetCategory(req.Category)
	if !valid {
		common.BadRequest(c)
		return
	}

	user := c.MustGet("user").(modb.User)
	if err := user.UnbindApp(category); err != nil {
		common.InternalServerError(c)
		return
	}
	common.Ok(c)
}
