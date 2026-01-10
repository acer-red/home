package db

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/tengfei-xy/go-log"
	"gorm.io/datatypes"
)

type Feedback struct {
	Base
	ProductID  uuid.UUID `gorm:"index;type:uuid;not null"`
	FbType     int
	Title      string
	Content    string
	IsPublic   bool
	DeviceFile string
	Images     datatypes.JSON
}

func (Feedback) TableName() string {
	return "feedback"
}

type RequestFeedbackPost struct {
	ProductID      uuid.UUID `gorm:"type:uuid;not null"`
	FbType         int
	Title          string
	Content        string
	DeviceFile     io.Reader
	DeviceFileName string
	IsPublic       bool
	Images         []io.Reader
	ImagesName     []string
}
type FBFilter struct {
	Text string
}

func FeedbackPost(req *RequestFeedbackPost) (uuid.UUID, error) {
	fb := Feedback{
		ProductID: req.ProductID,
		FbType:    req.FbType,
		Title:     req.Title,
		Content:   req.Content,
		IsPublic:  req.IsPublic,
	}

	if err := db.Create(&fb).Error; err != nil {
		log.Error(err)
		return uuid.Nil, err
	}

	updates := make(map[string]interface{})

	if req.DeviceFile != nil {
		ext := filepath.Ext(req.DeviceFileName)
		if ext == "" {
			ext = ".txt"
		}
		fb.DeviceFile = fmt.Sprintf("%s_device%s", fb.ID.String(), ext)
		updates["device_file"] = fb.DeviceFile

		buf := new(bytes.Buffer)
		if _, err := io.Copy(buf, req.DeviceFile); err == nil {
			ImageCreate(fb.DeviceFile, "feedback", buf.Bytes())
		}
	}

	if len(req.Images) > 0 {
		names := make([]string, len(req.Images))
		for i, image := range req.Images {
			orig := ""
			if i < len(req.ImagesName) {
				orig = req.ImagesName[i]
			}
			ext := filepath.Ext(orig)
			names[i] = fmt.Sprintf("%s_fb_%d%s", fb.ID.String(), i, ext)

			buf := new(bytes.Buffer)
			if _, err := io.Copy(buf, image); err == nil {
				ImageCreate(names[i], "feedback", buf.Bytes())
			}
		}

		fb.Images = datatypes.JSON(func() []byte {
			b, _ := json.Marshal(names)
			return b
		}())
		updates["images"] = fb.Images
	}

	if len(updates) > 0 {
		if err := db.Model(&fb).Updates(updates).Error; err != nil {
			log.Error(err)
			return uuid.Nil, err
		}
	}

	return fb.ID, nil
}

func FeedbacksGet(f FBFilter) (any, error) {
	type response struct {
		ID         string   `json:"id"`
		FbType     int      `json:"fb_type"`
		Title      string   `json:"title"`
		Content    string   `json:"content"`
		IsPublic   bool     `json:"is_public"`
		DeviceFile string   `json:"device_file"`
		Images     []string `json:"images"`
		CreatedAt  string   `json:"created_at"`
		UpdatedAt  string   `json:"updated_at"`
	}

	var feedbacks []Feedback
	tx := db.Where("is_public = ?", true)
	if f.Text != "" {
		tx = tx.Where("title LIKE ? OR content LIKE ?", "%"+f.Text+"%", "%"+f.Text+"%")
	}

	if err := tx.Find(&feedbacks).Error; err != nil {
		return nil, err
	}

	var res []response
	for _, fb := range feedbacks {
		var images []string
		_ = json.Unmarshal(fb.Images, &images)
		res = append(res, response{
			ID:         fb.ID.String(),
			FbType:     fb.FbType,
			Title:      fb.Title,
			Content:    fb.Content,
			IsPublic:   fb.IsPublic,
			DeviceFile: fb.DeviceFile,
			Images:     images,
			CreatedAt:  fb.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:  fb.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return res, nil
}
