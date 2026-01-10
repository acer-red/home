package user

import (
	"net/http"

	"github.com/acer-red/official/engine/service/db"
	"github.com/acer-red/official/engine/service/web/common"
	"github.com/acer-red/official/engine/service/web/error"
	Err "github.com/acer-red/official/engine/service/web/error"
	"github.com/acer-red/official/engine/util"
	"github.com/gin-gonic/gin"
	log "github.com/tengfei-xy/go-log"
)

func userDelete(c *gin.Context) {
	log.Info("用户删除")

	user := c.MustGet("user").(db.User)
	if err := user.Delete(); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	common.SetCookie(c, "login", "", 0)
	c.JSON(http.StatusOK, error.OK.JSON())

}

func userUnbindApp(c *gin.Context) {
	log.Info("用户解绑应用")

	var req struct {
		Category string `json:"category"`
	}
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, Err.DataParseFailed)
		return
	}

	category, valid := util.GetCategory(req.Category)
	if !valid {
		c.AbortWithStatusJSON(http.StatusInternalServerError, Err.InternalServer.JSON())
		return
	}

	user := c.MustGet("user").(db.User)
	if err := user.UnbindApp(category); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, error.OK.JSON())
}
