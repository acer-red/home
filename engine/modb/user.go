package modb

import (
	"context"
	"regexp"
	"strings"
	"sys"
	"time"

	"github.com/tengfei-xy/go-log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Cookie struct {
	Key    string    `bson:"key"`
	Value  string    `bson:"value"`
	CRTime time.Time `bson:"crtime"`
	EXTime time.Time `bson:"extime"`
}
type RequestUserRegister struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Cookie   Cookie `json:"-"`
}

func (req *RequestUserRegister) checkUser() bool {
	username := req.Username
	if username == "" || len(username) > 20 {
		return false
	}
	if len(username) < 3 {
		return false
	}

	if len(username) > 20 {
		return false
	}

	// 允许字母、数字、下划线和点
	allowedChars := regexp.MustCompile(`^[a-zA-Z0-9_.]+$`)
	if !allowedChars.MatchString(username) {
		return false
	}

	// 避免使用过于简单的数字组合
	onlyNumbers := regexp.MustCompile(`^+$`)
	if onlyNumbers.MatchString(username) && len(username) < 6 {
		return false
	}

	// 更全面的禁用词列表 (可以从配置文件或数据库加载)
	forbiddenWords := []string{"admin", "test", "guest", "root", "administrator", "administrators", "superuser"}
	for _, word := range forbiddenWords {
		if strings.ToLower(username) == word {
			return false
		}
	}
	return true
}
func (req *RequestUserRegister) checkPasswd() bool {
	password := req.Password
	if len(password) < 8 {
		return false
	}

	// 检查是否包含大写字母
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return false
	}

	// 检查是否包含小写字母
	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return false
	}

	// 检查是否包含数字
	if !regexp.MustCompile(``).MatchString(password) {
		return false
	}

	// 检查是否包含特殊字符
	if !regexp.MustCompile(`[!@#\$%^&*(),.?":{}|<>]`).MatchString(password) {
		return false
	}
	return true
}
func (req *RequestUserRegister) Check() bool {
	if !req.checkUser() {
		return false
	}
	if !req.checkPasswd() {
		return false
	}

	// 检查邮箱
	if req.Email == "" || !strings.Contains(req.Email, "@") {
		return false
	}
	return true
}
func (req *RequestUserRegister) Find() (bool, error) {
	// 检查用户名是否存在
	filter := bson.M{"username": req.Username}
	coll := db.Collection("user")
	if count, err := coll.CountDocuments(context.TODO(), filter); err != nil {
		log.Error(err)
		return false, err
	} else if count > 0 {
		return true, nil
	}

	// 检查邮箱是否存在
	filter = bson.M{"email": req.Email}
	if count, err := coll.CountDocuments(context.TODO(), filter); err != nil {
		log.Error(err)
		return false, err
	} else if count > 0 {
		return true, nil
	}
	return false, nil
}
func (req *RequestUserRegister) Register() (string, error) {

	var err error
	req.Password, err = sys.HashPassword(req.Password)
	if err != nil {
		log.Error(err)
		return "", err
	}
	id := sys.CreateUUID()

	m := bson.D{
		{Key: "username", Value: req.Username},
		{Key: "password", Value: req.Password},
		{Key: "email", Value: req.Email},
		{Key: "id", Value: id},
		{Key: "cookies", Value: []Cookie{
			req.Cookie,
		}},
	}

	if _, err := db.Collection("user").InsertOne(context.TODO(), m); err != nil {
		log.Error(err)
		return "", err
	}
	req.Cookie.new()
	return id, nil
}
func (c *Cookie) new() {
	c.Key = "login"
	c.Value = sys.CreateUUID()
	c.CRTime = time.Now()
	c.EXTime = time.Now().AddDate(0, 1, 0)
}

type RequestUserLogin struct {
	Account  string `json:"account"`
	Password string `json:"password"`
	m        bson.M
	Cookie   Cookie `json:"-"`
}

func (req *RequestUserLogin) Check() bool {
	// 检查用户名
	if req.Account == "" || len(req.Account) > 20 {
		return false
	}
	// 检查密码
	if len(req.Password) < 6 {
		return false
	}
	return true
}
func (req *RequestUserLogin) Find() (bool, error) {

	filter := bson.M{"$or": []bson.M{
		{"username": req.Account},
		{"email": req.Account},
	}}
	err := db.Collection("user").FindOne(context.TODO(), filter).Decode(&req.m)

	if err == mongo.ErrNoDocuments {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
func (req *RequestUserLogin) ComparePassword() error {
	return sys.ComparePassword(req.m["password"].(string), req.Password)
}
func (req *RequestUserLogin) GetID() (string, error) {
	id := req.m["id"].(string)
	return id, nil
}

func (req *RequestUserLogin) GetCookie() {
	cookies := req.m["cookies"].(primitive.A)

	for _, cookie := range cookies {
		c := cookie.(bson.M)
		if c["key"] != "login" {
			continue
		}
		req.Cookie.CRTime = c["crtime"].(primitive.DateTime).Time()
		req.Cookie.EXTime = c["extime"].(primitive.DateTime).Time()
		req.Cookie.Value = c["value"].(string)
		req.Cookie.Key = c["key"].(string)
		return
	}
	req.Cookie.new()
	filter := bson.M{"$or": []bson.M{
		{"username": req.Account},
		{"email": req.Account},
	}}
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "cookies", Value: req.Cookie}}}}
	_, err := db.Collection("user").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Error(err)
		return
	}
}

type RequestUserLogout struct {
	ID   string             `json:"id"`
	UOID primitive.ObjectID `json:"-" bson:"_id"`
}

func (req *RequestUserLogout) Check() bool {
	return req.ID != ""
}
func (req *RequestUserLogout) Find() (bool, error) {
	filter := bson.M{"id": req.ID}
	err := db.Collection("user").FindOne(context.TODO(), filter).Decode(&req)
	if err == mongo.ErrNoDocuments {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}
func (req *RequestUserLogout) DeleteCookie() error {
	filter := bson.M{"_id": req.UOID}
	update := bson.D{{Key: "$pull", Value: bson.D{{Key: "cookies", Value: bson.M{"key": "login"}}}}}
	_, err := db.Collection("user").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
