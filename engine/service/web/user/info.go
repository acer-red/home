package user

import (
	"net/http"

	"github.com/acer-red/official/engine/service/db"
	Err "github.com/acer-red/official/engine/service/web/error"
	"github.com/acer-red/official/engine/util"
	"github.com/gin-gonic/gin"
	log "github.com/tengfei-xy/go-log"
)

func getUserInfo(c *gin.Context) {
	log.Info("用户获取信息")

	user := c.MustGet("user").(db.User)
	apis, _ := db.GetAPIsForUser(user.ID)
	c.JSON(http.StatusOK, Err.OK.Data(user.ToResponse(apis)))
}
func putUserInfo(c *gin.Context) {
	log.Info("用户修改信息")
	var req db.RequestPutUserInfo
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, Err.DataParseFailed)
		return
	}
	req.UserID = c.MustGet("user").(db.User).ID.String()

	if err := req.Update(); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, Err.OK.JSON())
}
func userRandomInfo(c *gin.Context) {
	log.Info("用户随机信息")

	type response struct {
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
	}
	c.JSON(http.StatusOK, Err.OK.Data(response{
		Nickname: util.RandomNickname(),
		Avatar:   util.RandomAvatarBase64(),
	}))
}
