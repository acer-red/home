package user

import (
	"net/http"

	"github.com/acer-red/official/engine/service/db"
	"github.com/acer-red/official/engine/service/web/error"
	Err "github.com/acer-red/official/engine/service/web/error"
	"github.com/gin-gonic/gin"
	log "github.com/tengfei-xy/go-log"
)

func putUserProfile(c *gin.Context) {

	type response struct {
		URL string `json:"url"`
	}
	hasNickname := false
	hasAvatar := false
	user := c.MustGet("user").(db.User)
	if user.IsNoID() {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, Err.InternalServer.JSON())
		return
	}

	nicknames := form.Value["nickname"]

	if len(nicknames) != 0 {
		nickname := nicknames[0]
		log.Info("更新用户昵称")

		hasNickname = true
		if err := user.UpdateNickname(nickname); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
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
			c.JSON(http.StatusOK, error.OK.JSON())
			return
		}
	} else {
		if !hasAvatar {
			c.AbortWithStatusJSON(http.StatusInternalServerError, Err.InternalServer.JSON())
			return
		}
	}
	log.Info("更新用户头像")

	ext := form.Value["ext"][0]
	if ext == "" {
		c.AbortWithStatusJSON(http.StatusInternalServerError, Err.InternalServer.JSON())
		return
	}

	file, err := avatars[0].Open()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, Err.InternalServer.JSON())
		return
	}
	defer file.Close()

	if err := user.UpdateAvatar(file, ext); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	log.Info("更新用户头像完成")

	c.JSON(http.StatusOK, error.OK.Data(response{URL: user.AvatarURL}))
}
