package web

import (
	"modb"
	"sys"

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

	if err == sys.ErrNoFound {
		notFound(g)
		return
	}
	if err != nil {
		internalServerError(g)
		return
	}

	okImage(g, res)
}
