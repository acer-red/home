package modb

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"regexp"
	"strings"
	"sys"
	"time"

	"github.com/tengfei-xy/go-log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type API struct {
	APIKey   string    `bson:"apikey" json:"apikey"`
	EXTime   time.Time `bson:"extime" json:"extime"`
	LUTime   time.Time `bson:"lutime" json:"lutime"`
	UsedTims int32     `bson:"used_times" json:"used_times"`
}
type Cookie struct {
	Key    string    `bson:"key"`
	Value  string    `bson:"value"`
	CRTime time.Time `bson:"crtime"`
	EXTime time.Time `bson:"extime"`
}
type Avatar struct {
	Name string `json:"name" bson:"name"`
	URL  string `json:"url" bson:"url"`
}
type Profile struct {
	Avatar   Avatar `json:"avatar"`
	Nickname string `json:"nickname"`
}
type ResponseGetUserInfo struct {
	ID       string    `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Profile  Profile   `json:"profile"`
	CRTime   time.Time `json:"crtime"`
	API      []API     `json:"api"`
}
type User struct {
	UOID      primitive.ObjectID `bson:"_id" json:"-"`
	auatarOID primitive.ObjectID `bson:"" json:"-"`
	ID        string             `bson:"id" json:"-"`
	Username  string             `bson:"username" json:"username"`
	Email     string             `bson:"email" json:"email"`
	CRTime    time.Time          `bson:"crtime" json:"crtime"`
	Profile   Profile            `bson:"profile" json:"profile"`
	// UTime    time.Time       `bson:"uptime"`
	Cookies []Cookie `bson:"cookies" json:"-"`
	API     []API    `json:"api"`
}
type RequestUserRegister struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	CategoryStr string `json:"category"`
	Category    sys.CAtegory
	Cookie      Cookie
	profile     Profile
	uoid        primitive.ObjectID
}
type RequestUserLogin struct {
	Account     string `json:"account"`
	Password    string `json:"password"`
	CategoryStr string `json:"category"`
	Category    sys.CAtegory
	m           bson.M
	Cookie      Cookie `json:"-"`
}
type RequestPutUserInfo struct {
	Nickname string             `json:"nickname"`
	Password string             `json:"password"`
	UOID     primitive.ObjectID `json:"-"`
}

// cookie
func (c *Cookie) setLoginCookie() {
	c.Key = "login"
	c.Value = sys.CreateUUID()
	c.CRTime = time.Now()
	c.EXTime = time.Now().AddDate(0, 1, 0)
}

// 用户注册
func (req *RequestUserRegister) outputSrc() {
	switch req.Category {
	case sys.CAtegoryIndex:
		log.Info("注册源:官网")
	case sys.CAtegoryWT:
		log.Info("注册源:枫迹")
	default:
		log.Warn("注册源:未知")
	}
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
	// if !regexp.MustCompile(`[!@#\$%^&*(),.?":{}|<>]`).MatchString(password) {
	// 	return false
	// }
	return true
}
func (req *RequestUserRegister) CheckAndSetCatetory() error {
	c, err := sys.GetCategory(req.CategoryStr)
	if err != nil {
		return err
	}
	req.Category = c
	return nil
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
	if err := req.CheckAndSetCatetory(); err != nil {
		log.Error(err)
		return false
	}
	req.outputSrc()

	return true
}
func (req *RequestUserRegister) IsFromIndex() bool {
	return req.Category == sys.CAtegoryIndex
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
func (req *RequestUserRegister) GetCookie() {
	req.Cookie.setLoginCookie()
	filter := bson.M{"$or": []bson.M{
		{"username": req.Username},
		{"email": req.Email},
	}}
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "cookies", Value: req.Cookie}}}}
	_, err := db.Collection("user").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Error(err)
		return
	}
}
func (req *RequestUserRegister) RandomAccount() {
	req.Username = sys.CreateUUID()
	req.Password = sys.CreateUUID()
	req.Email = fmt.Sprintf("%s@%s.com", sys.CreateUUID(), sys.CreateUUID())
}
func (req *RequestUserRegister) BuildProfile() error {

	data := bytes.NewBuffer(sys.RandomAvatar())
	filename := fmt.Sprintf("%s.png", sys.CreateUUID())
	if err := ImageAvatarCreate(filename, data, req.uoid); err != nil {
		log.Error(err)
		return err
	}
	req.profile = Profile{
		// 用户名，暂时只支持中文
		Nickname: sys.RandomNickname(),

		Avatar: Avatar{
			Name: filename,
			URL:  setAvatarUrl(filename),
		},
	}
	m := bson.D{{Key: "profile", Value: bson.M{
		"nickname": req.profile.Nickname,
		"avatar": bson.M{
			"name": req.profile.Avatar.Name,
			"url":  req.profile.Avatar.URL,
		},
	}}}
	filter := bson.M{"_id": req.uoid}
	update := bson.D{{Key: "$set", Value: m}}
	_, err := db.Collection("user").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
func (req *RequestUserRegister) CancelRegister() {
	filter := bson.M{"_id": req.uoid}
	_, err := db.Collection("user").DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Error(err)
		return
	}
}
func (req *RequestUserRegister) Register(role sys.Role) (string, API, error) {
	var err error

	req.Password, err = sys.HashPassword(req.Password)
	if err != nil {
		log.Error(err)
		return "", API{}, err
	}
	id := sys.CreateUUID()

	m := bson.D{
		{Key: "id", Value: id},
		{Key: "role", Value: int(role)},
		{Key: "username", Value: req.Username},
		{Key: "password", Value: req.Password},
		{Key: "email", Value: req.Email},
		{Key: "crtime", Value: time.Now()},
		{Key: "uptime", Value: time.Now()},
		{Key: "cookies", Value: []Cookie{
			req.Cookie,
		}},
	}

	// 添加产品API
	api := newAPI()

	if req.Category != sys.CAtegoryIndex {
		m = append(m, bson.E{Key: "products", Value: bson.M{string(req.Category): bson.M{
			string("api"): []bson.M{
				{
					"apikey":     api.APIKey,
					"extime":     api.EXTime,
					"lutime":     api.LUTime,
					"used_times": api.UsedTims,
				},
			},
		}}})
	}

	result, err := db.Collection("user").InsertOne(context.TODO(), m)
	if err != nil {
		log.Error(err)
		return "", API{}, err
	}
	req.uoid = result.InsertedID.(primitive.ObjectID)
	if req.Category == sys.CAtegoryIndex {
		req.Cookie.setLoginCookie()
	}

	return id, api, nil
}

// 用户登陆
func (req *RequestUserLogin) checkCatetory() error {
	c, err := sys.GetCategory(req.CategoryStr)
	if err != nil {
		return err
	}
	req.Category = c
	return nil
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
	if err := req.checkCatetory(); err != nil {
		log.Error(err)
		return false
	}
	return true
}
func (req *RequestUserLogin) IsFromIndex() bool {
	return req.Category == sys.CAtegoryIndex
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
	req.Cookie.setLoginCookie()
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
func (req *RequestUserLogin) Login() ResponseGetUserInfo {

	avatar := req.m["profile"].(bson.M)["avatar"].(bson.M)

	res := ResponseGetUserInfo{
		ID:       req.m["id"].(string),
		Username: req.m["username"].(string),
		Email:    req.m["email"].(string),
		CRTime:   req.m["crtime"].(primitive.DateTime).Time(),
		Profile: Profile{
			Nickname: req.m["profile"].(bson.M)["nickname"].(string),
			Avatar: Avatar{
				Name: avatar["name"].(string),
				URL:  avatar["url"].(string),
			},
		},
	}

	if req.IsFromIndex() {
		return res
	}

	// 客户端用户返回API
	if l, ok := req.m["products"].(bson.M)[string(req.Category)].(bson.M)["api"]; ok {
		for _, g := range l.(primitive.A) {
			res.API = append(res.API, API{
				APIKey:   g.(bson.M)["apikey"].(string),
				EXTime:   g.(bson.M)["extime"].(primitive.DateTime).Time(),
				LUTime:   g.(bson.M)["lutime"].(primitive.DateTime).Time(),
				UsedTims: g.(bson.M)["used_times"].(int32),
			})
		}
	}
	return res
}

// 用户信息
func (req *RequestPutUserInfo) Update() error {

	var err error

	filter := bson.M{"_id": req.UOID}

	update := bson.D{{Key: "$set", Value: bson.D{}}}
	if req.Nickname != "" {
		update[0].Value = append(update[0].Value.(bson.D), bson.E{Key: "profile.nickname", Value: req.Nickname})
	}

	if req.Password != "" {
		req.Password, err = sys.HashPassword(req.Password)
		if err != nil {
			log.Error(err)
			return err
		}
		update[0].Value = append(update[0].Value.(bson.D), bson.E{Key: "password", Value: req.Password})
	}

	if len(update[0].Value.(bson.D)) == 0 {
		return nil
	}
	update[0].Value = append(update[0].Value.(bson.D), bson.E{Key: "uptime", Value: time.Now()})

	if _, err = db.Collection("user").UpdateOne(context.TODO(), filter, update); err != nil {
		log.Error(err)
		return err
	}
	return nil
}
func (u *User) IsNoID() bool {
	return u.UOID == primitive.NilObjectID
}
func (u *User) DeleteCookie() error {
	filter := bson.M{"_id": u.UOID}
	update := bson.D{{Key: "$pull", Value: bson.D{{Key: "cookies", Value: bson.M{"key": "login"}}}}}
	_, err := db.Collection("user").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
func (u *User) UpdateNickname(nickname string) error {

	filter := bson.M{"_id": u.UOID}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "profile.nickname", Value: nickname}}}}
	if _, err := db.Collection("user").UpdateOne(context.TODO(), filter, update); err != nil {
		log.Error(err)
		return err
	}
	return nil
}
func (u *User) UpdateAvatar(data io.Reader, ext string) error {
	err := u.getAuatarOID()
	if err != nil {
		log.Error(err)
		return err
	}
	if err := u.deleteAuatar(); err != nil {
		log.Error(err)
		return err
	}
	filename := fmt.Sprintf("%s%s", sys.CreateUUID(), ext)

	if err := u.createAuatar(filename, data); err != nil {
		log.Error(err)
		return err
	}

	if err := u.updateProfileAvatar(filename); err != nil {
		log.Error(err)
		return err
	}
	return nil

}
func (u *User) getAuatarOID() error {
	filter := bson.D{{Key: "metadata.uoid", Value: u.UOID}, {Key: "metadata.setup", Value: "avatar"}, {Key: "metadata.type", Value: "image"}}
	projection := bson.D{{Key: "_id", Value: 1}}
	findOptions := options.FindOne().SetProjection(projection)
	var resultDoc bson.M
	err := db.Collection("fs.files").FindOne(context.TODO(), filter, findOptions).Decode(&resultDoc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return sys.ErrNoFound
		} else {
			log.Error(err)
			return sys.ErrInternalServer
		}
	}
	objectID, ok := resultDoc["_id"].(primitive.ObjectID)
	if !ok {
		return sys.ErrInternalServer
	}
	u.auatarOID = objectID
	return nil
}
func (u *User) updateProfileAvatar(filename string) error {

	u.Profile.Avatar.Name = filename
	u.Profile.Avatar.URL = setAvatarUrl(filename)
	m := bson.D{{Key: "profile.avatar", Value: bson.M{
		"name": filename,
		"url":  setAvatarUrl(filename),
	}}}
	filter := bson.M{"_id": u.UOID}
	update := bson.D{{Key: "$set", Value: m}}
	_, err := db.Collection("user").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
func (u *User) deleteAuatar() error {
	if err := ImageDelete(u.auatarOID); err != nil {
		log.Error(err)
		return err
	}
	return nil
}
func (u *User) createAuatar(filename string, data io.Reader) error {

	if err := ImageAvatarCreate(filename, data, u.UOID); err != nil {
		log.Error(err)
		return err
	}
	return nil
}
func (u *User) Delete() error {
	filter := bson.D{{Key: "_id", Value: u.UOID}}
	_, err := db.Collection("user").DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil

}

// 根据cookie获取用户信息，用在auth中间件
func GetUserFromCookie(cookie string) (User, bool, error) {

	filter := bson.M{"cookies": bson.M{"$elemMatch": bson.M{"key": "login"}}}
	var m bson.M

	err := db.Collection("user").FindOne(context.TODO(), filter).Decode(&m)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return User{}, false, nil
		}
		log.Error(err)

		return User{}, false, err
	}
	profile := m["profile"].(bson.M)
	avatar := profile["avatar"].(bson.M)

	return User{
		ID:       m["id"].(string),
		Username: m["username"].(string),
		Email:    m["email"].(string),
		CRTime:   m["crtime"].(primitive.DateTime).Time(),
		Profile: Profile{
			Nickname: profile["nickname"].(string),
			Avatar: Avatar{
				Name: avatar["name"].(string),
				URL:  avatar["url"].(string),
			},
		},
	}, true, nil

}

// 根据API获取用户信息，用在auth中间件
func GetUserFromAPI(api string) (User, bool, error) {
	log.Debug3f("API验证:%s", api)
	filter := bson.M{fmt.Sprintf("%s.%s.%s", "products", string(sys.CAtegoryWT), "api"): bson.M{"$elemMatch": bson.M{"apikey": api}}}
	var m bson.M
	var apis []API

	err := db.Collection("user").FindOne(context.TODO(), filter).Decode(&m)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return User{}, false, nil
		}
		log.Error(err)
		return User{}, false, err
	}

	for _, api := range m["products"].(bson.M)[string(sys.CAtegoryWT)].(bson.M)["api"].(primitive.A) {
		apis = append(apis, API{
			APIKey:   api.(bson.M)["apikey"].(string),
			EXTime:   api.(bson.M)["extime"].(primitive.DateTime).Time(),
			LUTime:   api.(bson.M)["lutime"].(primitive.DateTime).Time(),
			UsedTims: api.(bson.M)["used_times"].(int32),
		})
	}
	profile := m["profile"].(bson.M)
	avatar := profile["avatar"].(bson.M)
	return User{
		UOID:     m["_id"].(primitive.ObjectID),
		ID:       m["id"].(string),
		Username: m["username"].(string),
		Email:    m["email"].(string),
		CRTime:   m["crtime"].(primitive.DateTime).Time(),
		Profile: Profile{
			Nickname: profile["nickname"].(string),
			Avatar: Avatar{
				Name: avatar["name"].(string),
				URL:  avatar["url"].(string),
			},
		},
		API: apis,
	}, true, nil
}

// 只能用在设置路径的地方，不能用不在获取路径的地方
func setAvatarUrl(f string) string {
	return "/images/" + f
}
func newAPI() API {
	return API{
		APIKey:   sys.CreateAPIKey(),
		EXTime:   time.Now().AddDate(0, 3, 0),
		LUTime:   time.Now(),
		UsedTims: 0,
	}
}
