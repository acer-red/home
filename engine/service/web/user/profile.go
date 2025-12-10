package user

import (
	"github.com/acer-red/official/engine/service/modb"
	"github.com/acer-red/official/engine/service/web/common"
	"github.com/gin-gonic/gin"
	log "github.com/tengfei-xy/go-log"
)

func putUserProfile(c *gin.Context) {

	type response struct {
		URL string `json:"url"`
	}
	hasNickname := false
	hasAvatar := false
	user := c.MustGet("user").(modb.User)
	if user.IsNoID() {
		common.InternalServerError(c)
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		common.BadRequest(c)
		return
	}

	nicknames := form.Value["nickname"]

	if len(nicknames) != 0 {
		nickname := nicknames[0]
		log.Info("更新用户昵称")

		hasNickname = true
		if err := user.UpdateNickname(nickname); err != nil {
			common.InternalServerError(c)
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
			common.Ok(c)
			return
		}
	} else {
		if !hasAvatar {
			common.BadRequest(c)
			return
		}
	}
	log.Info("更新用户头像")

	ext := form.Value["ext"][0]
	if ext == "" {
		common.BadRequest(c)
		return
	}

	file, err := avatars[0].Open()
	if err != nil {
		common.BadRequest(c)
		return
	}
	defer file.Close()

	if err := user.UpdateAvatar(file, ext); err != nil {
		common.InternalServerError(c)
		return
	}
	log.Info("更新用户头像完成")

	common.OkData(c, response{URL: user.Profile.Avatar.URL})
}
