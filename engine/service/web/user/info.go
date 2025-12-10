package user

import (
	"github.com/acer-red/official/engine/service/modb"
	"github.com/acer-red/official/engine/service/web/common"
	"github.com/acer-red/official/engine/util"
	"github.com/gin-gonic/gin"
	log "github.com/tengfei-xy/go-log"
)

func getUserInfo(c *gin.Context) {
	log.Info("用户获取信息")

	user := c.MustGet("user").(modb.User)
	common.OkData(c, user)
}
func putUserInfo(c *gin.Context) {
	log.Info("用户修改信息")
	var req modb.RequestPutUserInfo
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		common.BadRequest(c)
		return
	}
	req.UOID = c.MustGet("user").(modb.User).UOID

	if err := req.Update(); err != nil {
		common.InternalServerError(c)
		return
	}
	common.Ok(c)
}
func userRandomInfo(c *gin.Context) {
	log.Info("用户随机信息")

	type response struct {
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
	}

	common.OkData(c, response{
		Nickname: util.RandomNickname(),
		Avatar:   util.RandomAvatarBase64(util.CreateUUID()),
	})
}
