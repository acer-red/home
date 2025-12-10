package web

import (
	"github.com/acer-red/official/engine/service/modb"
	"github.com/acer-red/official/engine/service/web/common"
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
func ImageGet(g *gin.Context) {
	log.Infof("获取图片")
	name := g.Param("file")
	res, err := modb.ImageGet(name)

	if err == util.ErrNoFound {
		common.NotFound(g)
		return
	}
	if err != nil {
		common.InternalServerError(g)
		return
	}

	common.OkImage(g, res)
}
