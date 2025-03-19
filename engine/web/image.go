package web

import (
	"modb"
	"sys"

	"github.com/gin-gonic/gin"
)

func RouterImageGet(g *gin.Engine) {
	a := g.Group("/image/:file")
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
	name := g.Param("file")
	res, err := modb.ImageGet(name)

	if err == sys.ERR_NO_FOUND {
		notFound(g)
		return
	}
	if err != nil {
		internalServerError(g)
		return
	}

	okImage(g, res)
}
