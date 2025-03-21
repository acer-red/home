package modb

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sys"

	log "github.com/tengfei-xy/go-log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ImageGet(name string) (bytes.Buffer, error) {

	// 获取 GridFS Bucket 对象
	bucket, err := gridfs.NewBucket(db)
	if err != nil {
		log.Error(err)
		return bytes.Buffer{}, err
	}

	// 查询条件：根据 filename 查找
	filter := bson.D{{Key: "filename", Value: name}}

	// 投影：只返回 _id 字段
	projection := bson.D{{Key: "_id", Value: 1}}

	// FindOneOptions 用于设置投影
	findOptions := options.FindOne().SetProjection(projection)

	// 执行 FindOne 查询
	var resultDoc bson.M // 用 bson.M 存储结果，也可以定义 struct
	err = db.Collection("fs.files").FindOne(context.TODO(), filter, findOptions).Decode(&resultDoc)

	if err != nil {
		log.Error(err)
		if err == mongo.ErrNoDocuments {
			return bytes.Buffer{}, sys.ErrNoFound
		} else {
			return bytes.Buffer{}, sys.ErrInternalServer
		}
	}

	// 从结果印迹中获取 _id
	objectID, ok := resultDoc["_id"].(primitive.ObjectID)
	if !ok {
		return bytes.Buffer{}, sys.ErrInternalServer
	}

	var downloadBuffer bytes.Buffer

	downloadStreamByID, err := bucket.OpenDownloadStream(objectID) // 使用文件 ID 下载
	if err != nil {
		log.Error(err)
		return downloadBuffer, err
	}
	defer downloadStreamByID.Close()

	if _, err := io.Copy(&downloadBuffer, downloadStreamByID); err != nil {
		log.Error(err)
		return downloadBuffer, err
	}

	return downloadBuffer, nil
}
func ImageCreate(filename string, category string, data []byte) error {

	bucket, err := gridfs.NewBucket(db)
	if err != nil {
		log.Error(err)
		return err
	}

	uploadStream, err := bucket.OpenUploadStream(
		filename,
		options.GridFSUpload().SetMetadata(map[string]interface{}{"type": "image", "category": category}),
	)
	if err != nil {
		log.Error(err)
		return err
	}
	defer uploadStream.Close()

	fileSize, err := io.Copy(uploadStream, bytes.NewReader(data))
	if err != nil {
		log.Error(err)
		return err
	}
	log.Infof("创建图片: %s(%s)", filename, ByteCountSI(fileSize))
	return nil
}
func ImageAvatarCreate(filename string, category string, data io.Reader, uoid primitive.ObjectID) error {

	bucket, err := gridfs.NewBucket(db)
	if err != nil {
		log.Error(err)
		return err
	}

	uploadStream, err := bucket.OpenUploadStream(
		filename,
		options.GridFSUpload().SetMetadata(map[string]interface{}{"type": "image", "uoid": uoid}),
	)
	if err != nil {
		log.Error(err)
		return err
	}
	defer uploadStream.Close()

	fileSize, err := io.Copy(uploadStream, data)
	if err != nil {
		log.Error(err)
		return err
	}
	log.Infof("创建图片: %s(%s)", filename, ByteCountSI(fileSize))
	return nil
}

// 不一定有用
func getAvatarFileIDFromUOID(uoid primitive.ObjectID) (primitive.ObjectID, error) {
	filter := bson.D{{Key: "metadata.uoid", Value: uoid}, {Key: "metadata.setup", Value: "avatar"}, {Key: "type", Value: "image"}}
	projection := bson.D{{Key: "_id", Value: 1}}
	findOptions := options.FindOne().SetProjection(projection)
	var resultDoc bson.M
	err := db.Collection("fs.files").FindOne(context.TODO(), filter, findOptions).Decode(&resultDoc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return primitive.NilObjectID, sys.ErrNoFound
		} else {
			log.Error(err)
			return primitive.NilObjectID, sys.ErrInternalServer
		}
	}
	objectID, ok := resultDoc["_id"].(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, sys.ErrInternalServer
	}
	return objectID, nil
}

func ImageCreateRandomAvatar(random string) (string, error) {
	d := sys.RandomAvatar(random)
	f := fmt.Sprintf("%s.png", sys.CreateUUID())
	return f, ImageCreate(f, "avatar", d)
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
func ImageDelete(fileID primitive.ObjectID) error {
	bucket, err := gridfs.NewBucket(db)
	if err != nil {
		log.Error(err)
		return err
	}
	err = bucket.Delete(fileID)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
