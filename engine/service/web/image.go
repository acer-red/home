package web

import (
	"net/http"
	"strings"

	"github.com/acer-red/official/engine/service/db"
	Err "github.com/acer-red/official/engine/service/web/error"
	"github.com/acer-red/official/engine/util"
	"github.com/gin-gonic/gin"
	log "github.com/tengfei-xy/go-log"
)

func RouterImageGet(g *gin.Engine) {
	a := g.Group("/images/:file")
	{
		a.GET("", ImageGet)
	}
}

//	func RouteImage(g *gin.Engine) {
//		a := g.Group("/image/:file")
//		{
//			a.DELETE("", ImageDelete)
//		}
//		b := g.Group("/image")
//		{
//			b.POST("", ImagePost)
//		}
//	}
func ImageGet(c *gin.Context) {
	log.Infof("获取图片")
	name := c.Param("file")
	if !strings.Contains(name, ".") {
		c.AbortWithStatusJSON(http.StatusBadRequest, Err.FormatError.JSON())
		return
	}

	data, err := db.ImageGet(name)

	if err == util.ErrNoFound {
		c.AbortWithStatusJSON(http.StatusNotFound, Err.NoFound.JSON())
		return
	}
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, Err.InternalServer.JSON())
		return
	}

	fotmat := strings.ToLower(strings.Split(name, ".")[1])
	switch fotmat {
	case "png":
		c.Data(http.StatusOK, "image/png", data.Bytes())
	case "jpg":
	case "jpeg":
		c.Data(http.StatusOK, "image/jpeg", data.Bytes())
	default:
		c.AbortWithStatusJSON(http.StatusBadRequest, Err.UnknownType.JSON())
	}
}
