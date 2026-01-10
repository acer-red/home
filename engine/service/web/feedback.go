package web

import (
	"io"
	"net/http"
	"strconv"

	"github.com/acer-red/official/engine/service/db"

	"github.com/acer-red/official/engine/service/web/error"
	Err "github.com/acer-red/official/engine/service/web/error"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/tengfei-xy/go-log"
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

func fbPost(c *gin.Context) {
	log.Info("创建反馈")

	type response struct {
		ID string `json:"id"`
	}

	var req db.RequestFeedbackPost

	if pid, ok := c.Get("product_id"); ok {
		req.ProductID = pid.(uuid.UUID)
	} else {
		log.Error("product_id not found in context")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		log.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, Err.InternalServer.JSON())
		return
	}

	if val, ok := form.Value["fb_type"]; ok && len(val) > 0 {
		req.FbType = atoi(val[0])
	}
	if val, ok := form.Value["title"]; ok && len(val) > 0 {
		req.Title = val[0]
	}
	if val, ok := form.Value["content"]; ok && len(val) > 0 {
		req.Content = val[0]
	}
	if val, ok := form.Value["is_public"]; ok && len(val) > 0 {
		req.IsPublic = atob(val[0])
	}

	deviceFiles := form.File["device_file"]
	imageFiles := form.File["images"]

	if len(deviceFiles) > 0 {
		file, err := deviceFiles[0].Open()
		if err != nil {
			log.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, Err.InternalServer.JSON())
			return
		}
		defer file.Close()
		req.DeviceFileName = deviceFiles[0].Filename
		req.DeviceFile = file
	}

	if len(imageFiles) > 0 {
		req.Images = make([]io.Reader, len(imageFiles))
		req.ImagesName = make([]string, len(imageFiles))
		for i, fileHeader := range imageFiles {
			file, err := fileHeader.Open()
			if err != nil {
				log.Error(err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, Err.InternalServer.JSON())
				return
			}
			defer file.Close()
			req.ImagesName[i] = fileHeader.Filename
			log.Info("image file")
			req.Images[i] = file
		}
	}

	fbID, err := db.FeedbackPost(&req)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	res := response{ID: fbID.String()}

	log.Infof("创建反馈成功 %s", res.ID)
	c.JSON(http.StatusOK, error.OK.Data(res))
}
func atoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func atob(s string) bool {
	b, _ := strconv.ParseBool(s)
	return b
}

func fbsGet(c *gin.Context) {
	log.Infof("获取反馈列表(已公开)")
	feedbacks, err := db.FeedbacksGet(db.FBFilter{
		Text: c.Query("text"),
	})
	if err != nil {
		log.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, error.OK.Data(feedbacks))
}
