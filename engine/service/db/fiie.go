package db

import (
	"bytes"
	"fmt"
	"io"

	"github.com/acer-red/official/engine/util"
	"github.com/google/uuid"
	log "github.com/tengfei-xy/go-log"
	"gorm.io/gorm"
)

type File struct {
	gorm.Model
	Name     string `gorm:"unique"`
	Data     []byte
	Category string
	MimeType string
	Size     int64
	UserID   uuid.UUID `gorm:"type:uuid;not null"`
	Metadata string
}

func ImageGet(name string) (bytes.Buffer, error) {
	var file File
	if err := db.Where("name = ?", name).First(&file).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return bytes.Buffer{}, util.ErrNoFound
		}
		return bytes.Buffer{}, util.ErrInternalServer
	}
	return *bytes.NewBuffer(file.Data), nil
}

func ImageCreate(filename string, category string, data []byte) error {
	file := File{
		Name:     filename,
		Data:     data,
		Category: category,
		Size:     int64(len(data)),
		MimeType: "image",
	}
	if err := db.Create(&file).Error; err != nil {
		log.Error(err)
		return err
	}
	log.Infof("创建图片: %s(%s)", filename, ByteCountSI(file.Size))
	return nil
}

func ImageAvatarCreate(filename string, data io.Reader, userID uuid.UUID) error {
	buf := new(bytes.Buffer)
	size, err := buf.ReadFrom(data)
	if err != nil {
		log.Error(err)
		return err
	}

	file := File{
		Name:     filename,
		Data:     buf.Bytes(),
		Category: "avatar",
		Size:     size,
		MimeType: "image",
		UserID:   userID,
		Metadata: "setup=avatar",
	}

	if err := db.Create(&file).Error; err != nil {
		log.Error(err)
		return err
	}
	log.Infof("创建头像: %s(%s)", filename, ByteCountSI(file.Size))
	return nil
}

func ByteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

func ImageDelete(filename string) error {
	return db.Where("name = ?", filename).Delete(&File{}).Error
}
