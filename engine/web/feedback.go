package web

import (
	"fmt"
	"io"
	"modb"
	"path/filepath"
	"strconv"
	"sys"

	"github.com/gin-gonic/gin"
	log "github.com/tengfei-xy/go-log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func RouteFeedback(g *gin.Engine) {
	a := g.Group("/feedback")
	{
		a.POST("", fbPost)
	}

	b := g.Group("/feedbacks")
	{
		b.GET("", fbsGet)
	}
}

func fbPost(g *gin.Context) {
	log.Info("创建反馈")

	type response struct {
		ID string `json:"id"`
	}
	res := response{
		ID: sys.CreateUUID(),
	}

	var req modb.RequestFeedbackPost
	req.UOID = g.MustGet("uoid").(primitive.ObjectID)
	req.FBID = res.ID

	form, err := g.MultipartForm()
	if err != nil {
		log.Error(err)
		badRequest(g)
		return
	}

	req.FbType = atoi(form.Value["fb_type"][0])
	req.Title = form.Value["title"][0]
	req.Content = form.Value["content"][0]
	req.IsPublic = atob(form.Value["is_public"][0])
	deviceFiles := form.File["device_file"]
	imageFiles := form.File["images"]

	// 如果有设备文件信息，则上传
	if len(deviceFiles) > 0 {
		file, err := deviceFiles[0].Open()
		if err != nil {
			log.Error(err)
			badRequest(g)
			return
		}
		defer file.Close()
		req.DeviceFileName = fmt.Sprintf("%s_device.txt", res.ID)
		req.DeviceFile = file
	}

	// 如果有图片信息，则上传
	if len(imageFiles) > 0 {
		req.Images = make([]io.Reader, len(imageFiles))
		req.ImagesName = make([]string, len(imageFiles))
		for i, fileHeader := range imageFiles {
			file, err := fileHeader.Open()
			if err != nil {
				log.Error(err)
				badRequest(g)
				return
			}
			defer file.Close()
			req.ImagesName[i] = fmt.Sprintf("%s_fb_%d%s", res.ID, i, filepath.Ext(fileHeader.Filename))
			log.Info("image file")
			req.Images[i] = file
		}
	}

	if err := modb.FeedbackPost(&req); err != nil {
		internalServerError(g)
		return
	}

	log.Infof("创建反馈成功 %s", res.ID)
	okData(g, res)
}
func atoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
func atob(s string) bool {
	b, _ := strconv.ParseBool(s)
	return b
}

func fbsGet(g *gin.Context) {
	log.Infof("获取反馈列表(已公开)")
	feedbacks, err := modb.FeedbacksGet(modb.FBFilter{
		Text: g.Query("text"),
	})
	if err != nil {
		log.Error(err)
		internalServerError(g)
		return
	}
	okData(g, feedbacks)
}

// 	g.JSON(sys.StatusOK, response{Feedbacks: feedbacks})
// }

// func fbGet(g *gin.Context) {
// 	type response struct {
// 		Feedback modb.Feedback `json:"feedback"`
// 	}

// 	goid := g.MustGet("goid").(primitive.ObjectID)
// 	fbid := g.Param("fbid")

// 	feedback, err := modb.FeedbackGet(goid, fbid)
// 	if err != nil {
// 		log.Error(err)
// 		badRequest(g)
// 		return
// 	}

// 	g.JSON(sys.StatusOK, response{Feedback: feedback})
// }
// func fbPut(g *gin.Context) {
// 	type response struct {
// 		Feedback modb.Feedback `json:"feedback"`
// 	}

// 	goid := g.MustGet("goid").(primitive.ObjectID)
// 	fbid := g.Param("fbid")

// 	var req modb.RequestFeedbackPut
// 	if err := g.ShouldBindBodyWithJSON(&req); err != nil {
// 		log.Error(err)
// 		badRequest(g)
// 		return
// 	}

// 	feedback, err := modb.FeedbackPut(goid, fbid, &req)
// 	if err != nil {
// 		log.Error(err)
// 		badRequest(g)
// 		return
// 	}

// 	g.JSON(sys.StatusOK, response{Feedback: feedback})
// }
