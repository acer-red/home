package modb

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/acer-red/official/engine/service/cache"
	"github.com/acer-red/official/engine/util"

	"github.com/tengfei-xy/go-log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const passwordMinLen = 8
const loginJWTDuration = 30 * 24 * time.Hour

type API struct {
	APIKey     string    `bson:"apikey" json:"apikey"`
	ExpiresAt  time.Time `bson:"expiresAt" json:"expiresAt"`
	LastUsedAt time.Time `bson:"lastusedAt" json:"lastusedAt"`
	UsedTims   int32     `bson:"used_times" json:"used_times"`
}
type Cookie struct {
	Key       string    `bson:"key"`
	Value     string    `bson:"value"`
	CeateAt   time.Time `bson:"createAt"`
	ExpiresAt time.Time `bson:"expiresAt"`
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
	CeateAt  time.Time `json:"createAt"`
	API      []API     `json:"api"`
}
type User struct {
	UOID      primitive.ObjectID `bson:"_id" json:"-"`
	auatarOID primitive.ObjectID `bson:"" json:"-"`
	ID        string             `bson:"id" json:"-"`
	Username  string             `bson:"username" json:"username"`
	Email     string             `bson:"email" json:"email"`
	CeateAt   time.Time          `bson:"createAt" json:"createAt"`
	Profile   Profile            `bson:"profile" json:"profile"`
	// UTime    time.Time       `bson:"updateAt"`
	API       []API `json:"api"`
	IsDeleted bool  `bson:"-" json:"-"`
}
type RequestUserRegister struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	CategoryStr string `json:"category"`
	UID         string `json:"uid"`
	PublicKey   string `json:"public_key"`
	Category    util.CAtegory
	Cookie      Cookie
	profile     Profile
	uoid        primitive.ObjectID
}
type RequestUserLogin struct {
	Account  string `json:"account"`
	Password string `json:"password"`
	Category string `json:"category"`
	UID      string `json:"uid"`
	category util.CAtegory
	m        bson.M
}
type RequestPutUserInfo struct {
	Nickname string             `json:"nickname"`
	Password string             `json:"password"`
	UOID     primitive.ObjectID `json:"-"`
}

// cookie
func (c *Cookie) setLoginCookie() {
	c.Key = "login"
	c.Value = util.CreateUUID()
	c.CeateAt = time.Now()
	c.ExpiresAt = time.Now().AddDate(0, 1, 0)
}

// 用户注册
func (req *RequestUserRegister) outputSrc() {
	switch req.Category {
	case util.CAtegoryOfficial:
		log.Info("注册源:官网")
	case util.CAtegoryWT:
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
	if len(password) < passwordMinLen {
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
	c, ok := util.GetCategory(req.CategoryStr)
	if !ok {
		return fmt.Errorf("无效的产品类别")
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
	return req.Category == util.CAtegoryOfficial
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

// FindUser 查找用户，返回用户对象
func (req *RequestUserRegister) FindUser() (*User, bool, error) {
	filter := bson.M{"$or": []bson.M{
		{"username": req.Username},
		{"email": req.Email},
	}}
	var m bson.M
	err := db.Collection("user").FindOne(context.TODO(), filter).Decode(&m)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, false, nil
		}
		return nil, false, err
	}

	profile := m["profile"].(bson.M)
	avatar := profile["avatar"].(bson.M)

	u := &User{
		UOID:     m["_id"].(primitive.ObjectID),
		ID:       m["id"].(string),
		Username: m["username"].(string),
		Email:    m["email"].(string),
		CeateAt:  m["createAt"].(primitive.DateTime).Time(),
		Profile: Profile{
			Nickname: profile["nickname"].(string),
			Avatar: Avatar{
				Name: avatar["name"].(string),
				URL:  avatar["url"].(string),
			},
		},
	}

	if req.Category != util.CAtegoryOfficial {
		if products, ok := m["products"].(bson.M); ok {
			if data, ok := products[string(req.Category)].(bson.M); ok {
				if _, ok := data["deleted_at"]; ok {
					u.IsDeleted = true
				}
			}
		}
	}

	return u, true, nil
}

// Reactivate 重新激活用户
func (req *RequestUserRegister) Reactivate(u *User) (string, API, error) {
	var err error
	req.Password, err = util.HashPassword(req.Password)
	if err != nil {
		return "", API{}, err
	}

	req.uoid = u.UOID

	api := newAPI()

	keyDeleted := fmt.Sprintf("products.%s.deleted_at", string(req.Category))
	keyApi := fmt.Sprintf("products.%s.api", string(req.Category))
	keyClients := fmt.Sprintf("products.%s.clients", string(req.Category))

	filter := bson.M{"_id": u.UOID}
	update := bson.M{
		"$set": bson.M{
			"password":   req.Password,
			"public_key": []byte(req.PublicKey),
			"updateAt":   time.Now(),
		},
		"$unset": bson.M{
			keyDeleted: "",
		},
		"$push": bson.M{
			keyApi: bson.M{
				"apikey":     api.APIKey,
				"expiresAt":  api.ExpiresAt,
				"lastusedAt": api.LastUsedAt,
				"used_times": api.UsedTims,
			},
		},
		"$addToSet": bson.M{
			keyClients: bson.M{"uid": req.UID},
		},
	}

	_, err = db.Collection("user").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return "", API{}, err
	}

	return u.ID, api, nil
}
func (req *RequestUserRegister) RandomAccount() {
	req.Username = util.CreateUUID()
	req.Password = util.CreateUUID()
	req.Email = fmt.Sprintf("%s@%s.com", util.CreateUUID(), util.CreateUUID())
}
func (req *RequestUserRegister) BuildProfile() error {

	data := bytes.NewBuffer(util.RandomAvatar())
	filename := fmt.Sprintf("%s.png", util.CreateUUID())
	if err := ImageAvatarCreate(filename, data, req.uoid); err != nil {
		log.Error(err)
		return err
	}
	req.profile = Profile{
		// 用户名，暂时只支持中文
		Nickname: util.RandomNickname(),

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
func (req *RequestUserRegister) checkPublicKey() bool {
	if req.PublicKey == "" {
		return false
	}
	decoded, err := base64.StdEncoding.DecodeString(req.PublicKey)
	if err != nil {
		return false
	}

	zeroKey := make([]byte, 32)
	if bytes.Equal(decoded, zeroKey) {
		return false
	}
	return len(decoded) == 32
}
func (req *RequestUserRegister) Register(role util.Role) (string, API, error) {
	var err error

	req.Password, err = util.HashPassword(req.Password)
	if err != nil {
		log.Error(err)
		return "", API{}, err
	}

	ok := req.checkPublicKey()
	if !ok {
		return "", API{}, fmt.Errorf("不是有效的公钥")
	}

	id := util.CreateUUID()

	m := bson.D{
		{Key: "id", Value: id},
		{Key: "role", Value: int(role)},
		{Key: "username", Value: req.Username},
		{Key: "password", Value: req.Password},
		{Key: "email", Value: req.Email},
		{Key: "createAt", Value: time.Now()},
		{Key: "updateAt", Value: time.Now()},
		{Key: "public_key", Value: []byte(req.PublicKey)}}

	// 添加产品API
	api := newAPI()

	if req.Category != util.CAtegoryOfficial {
		m = append(m, bson.E{
			Key: "products", Value: bson.M{
				string(req.Category): bson.M{
					string("api"): []bson.M{
						{
							"apikey":     api.APIKey,
							"expiresAt":  api.ExpiresAt,
							"lastusedAt": api.LastUsedAt,
							"used_times": api.UsedTims,
						},
					},
					"clients": []bson.M{{"uid": req.UID}},
				},
			}})
	}

	result, err := db.Collection("user").InsertOne(context.TODO(), m)
	if err != nil {
		log.Error(err)
		return "", API{}, err
	}
	req.uoid = result.InsertedID.(primitive.ObjectID)
	if req.Category == util.CAtegoryOfficial {
		req.Cookie.setLoginCookie()
	}

	return id, api, nil
}

// 用户登陆
func (req *RequestUserLogin) checkCatetory() error {
	c, ok := util.GetCategory(req.Category)
	if !ok {
		return fmt.Errorf("无效的产品类别")
	}
	req.category = c
	return nil
}
func (req *RequestUserLogin) Check() bool {
	// 检查用户名
	if req.Account == "" {
		log.Warn("用户名为空")
		return false
	}
	// 检查密码
	if len(req.Password) < passwordMinLen {
		log.Warn("密码过短")
		return false
	}
	if err := req.checkCatetory(); err != nil {
		return false
	}
	if req.UID == "" {
		return false
	}
	return true
}
func (req *RequestUserLogin) Find() (bool, error) {

	filter := bson.M{"$or": []bson.M{
		{"username": req.Account},
		{"email": req.Account},
	}}

	// 检查是否已解绑
	key := fmt.Sprintf("products.%s.deleted_at", string(req.category))
	filter[key] = bson.M{"$exists": false}

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
	return util.ComparePassword(req.m["password"].(string), req.Password)
}
func (req *RequestUserLogin) CategoryValue() util.CAtegory {
	return req.category
}
func (req *RequestUserLogin) GetAccountID() string {
	if id, ok := req.m["id"].(string); ok {
		return id
	}
	return ""
}
func (req *RequestUserLogin) GetUserID() string {
	return req.UID
}
func (req *RequestUserLogin) CheckNewDevice() error {
	log.Debugf("检查是否新设备UID=%s", req.UID)

	products, ok := req.m["products"].(bson.M)
	if !ok || products == nil {
		log.Error(util.ErrStructure.Error())
		return util.ErrStructure
	}

	product, ok := products[string(req.category)].(bson.M)
	if !ok || product == nil {
		return util.ErrStructure
	}

	clientsVal, ok := product["clients"]
	if !ok || clientsVal == nil {
		log.Error(util.ErrStructure.Error())
		return util.ErrStructure

	}

	clients, ok := clientsVal.(primitive.A)
	if !ok {
		log.Error(util.ErrStructure.Error())
		return util.ErrStructure
	}

	for _, c := range clients {
		client, _ := c.(bson.M)
		if client != nil && client["uid"] == req.UID {
			return nil
		}
	}

	filter := bson.M{"_id": req.m["_id"].(primitive.ObjectID)}
	update := bson.D{{Key: "$push", Value: bson.D{{Key: fmt.Sprintf("products.%s.clients", string(req.category)), Value: bson.M{"uid": req.UID}}}}}
	_, err := db.Collection("user").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Error(err)
		return err
	}
	log.Infof("新设备登录，添加设备UID=%s", req.UID)
	return nil
}
func (req *RequestUserLogin) buildJWT(accountID, uid string, expireAt time.Time) (string, error) {
	ttl := time.Until(expireAt)
	if ttl <= 0 {
		ttl = loginJWTDuration
	}

	return util.CreateJWT(accountID,
		uid,
		req.m["username"].(string),
		req.m["email"].(string),
		req.category,
		ttl,
	)
}
func (req *RequestUserLogin) BuildLoginResponse() ResponseGetUserInfo {
	avatar := req.m["profile"].(bson.M)["avatar"].(bson.M)

	res := ResponseGetUserInfo{
		ID:       req.m["id"].(string),
		Username: req.m["username"].(string),
		Email:    req.m["email"].(string),
		CeateAt:  req.m["createAt"].(primitive.DateTime).Time(),
		Profile: Profile{
			Nickname: req.m["profile"].(bson.M)["nickname"].(string),
			Avatar: Avatar{
				Name: avatar["name"].(string),
				URL:  avatar["url"].(string),
			},
		},
	}

	// 返回API
	if products, ok := req.m["products"].(bson.M); ok {
		if product, ok := products[string(req.category)].(bson.M); ok {
			if l, ok := product["api"].(primitive.A); ok {
				for _, g := range l {
					apiItem, _ := g.(bson.M)
					if apiItem == nil {
						continue
					}
					res.API = append(res.API, API{
						APIKey:     apiItem["apikey"].(string),
						ExpiresAt:  apiItem["expiresAt"].(primitive.DateTime).Time(),
						LastUsedAt: apiItem["lastusedAt"].(primitive.DateTime).Time(),
						UsedTims:   apiItem["used_times"].(int32),
					})
				}
			}
		}
	}
	return res
}
func (req *RequestUserLogin) Login(expireAt time.Time) (ResponseGetUserInfo, string, error) {

	res := req.BuildLoginResponse()

	token, err := req.buildJWT(req.GetAccountID(), req.UID, expireAt)
	if err != nil {
		log.Error(err)
		return ResponseGetUserInfo{}, "", err
	}

	if _, err := cache.SaveSession(req.GetAccountID(), req.UID, token, req.category, expireAt); err != nil {
		log.Error(err)
		return ResponseGetUserInfo{}, "", err
	}
	return res, token, nil
}

type Clients struct {
	UID string `bson:"uid" json:"uid"`
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
		req.Password, err = util.HashPassword(req.Password)
		if err != nil {
			log.Error(err)
			return err
		}
		update[0].Value = append(update[0].Value.(bson.D), bson.E{Key: "password", Value: req.Password})
	}

	if len(update[0].Value.(bson.D)) == 0 {
		return nil
	}
	update[0].Value = append(update[0].Value.(bson.D), bson.E{Key: "updateAt", Value: time.Now()})

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
	update := bson.D{{Key: "$pull", Value: bson.D{{Key: "", Value: bson.M{"key": "login"}}}}}
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
	filename := fmt.Sprintf("%s%s", util.CreateUUID(), ext)

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
			return util.ErrNoFound
		} else {
			log.Error(err)
			return util.ErrInternalServer
		}
	}
	objectID, ok := resultDoc["_id"].(primitive.ObjectID)
	if !ok {
		return util.ErrInternalServer
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

// 根据API获取用户信息，用在auth中间件
func GetUserFromAPI(api string) (User, bool, error) {
	log.Debug3f("API验证:%s", api)
	filter := bson.M{fmt.Sprintf("%s.%s.%s", "products", string(util.CAtegoryWT), "api"): bson.M{"$elemMatch": bson.M{"apikey": api}}}
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

	productData := m["products"].(bson.M)[string(util.CAtegoryWT)].(bson.M)
	isDeleted := false
	if _, ok := productData["deleted_at"]; ok {
		isDeleted = true
	}

	for _, api := range productData["api"].(primitive.A) {
		apis = append(apis, API{
			APIKey:     api.(bson.M)["apikey"].(string),
			ExpiresAt:  api.(bson.M)["expiresAt"].(primitive.DateTime).Time(),
			LastUsedAt: api.(bson.M)["lastusedAt"].(primitive.DateTime).Time(),
			UsedTims:   api.(bson.M)["used_times"].(int32),
		})
	}
	profile := m["profile"].(bson.M)
	avatar := profile["avatar"].(bson.M)
	return User{
		UOID:     m["_id"].(primitive.ObjectID),
		ID:       m["id"].(string),
		Username: m["username"].(string),
		Email:    m["email"].(string),
		CeateAt:  m["createAt"].(primitive.DateTime).Time(),
		Profile: Profile{
			Nickname: profile["nickname"].(string),
			Avatar: Avatar{
				Name: avatar["name"].(string),
				URL:  avatar["url"].(string),
			},
		},
		API:       apis,
		IsDeleted: isDeleted,
	}, true, nil
}

// 根据ID获取用户信息
func GetUserByIDAndCategory(uid string, category util.CAtegory) (User, bool, error) {
	filter := bson.M{"id": uid}
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

	u := User{
		UOID:     m["_id"].(primitive.ObjectID),
		ID:       m["id"].(string),
		Username: m["username"].(string),
		Email:    m["email"].(string),
		CeateAt:  m["createAt"].(primitive.DateTime).Time(),
		Profile: Profile{
			Nickname: profile["nickname"].(string),
			Avatar: Avatar{
				Name: avatar["name"].(string),
				URL:  avatar["url"].(string),
			},
		},
	}

	if category != util.CAtegoryOfficial {
		if products, ok := m["products"].(bson.M); ok {
			if data, ok := products[string(category)].(bson.M); ok {
				if _, ok := data["deleted_at"]; ok {
					u.IsDeleted = true
				}
				if apis, ok := data["api"]; ok {
					for _, api := range apis.(primitive.A) {
						u.API = append(u.API, API{
							APIKey:     api.(bson.M)["apikey"].(string),
							ExpiresAt:  api.(bson.M)["expiresAt"].(primitive.DateTime).Time(),
							LastUsedAt: api.(bson.M)["lastusedAt"].(primitive.DateTime).Time(),
							UsedTims:   api.(bson.M)["used_times"].(int32),
						})
					}
				}
			}
		}
	}

	return u, true, nil
}

// 只能用在设置路径的地方，不能用不在获取路径的地方
func setAvatarUrl(f string) string {
	return "/images/" + f
}
func newAPI() API {
	return API{
		APIKey:     util.CreateAPIKey(),
		ExpiresAt:  time.Now().AddDate(0, 3, 0),
		LastUsedAt: time.Now(),
		UsedTims:   0,
	}
}

// UnbindApp 删除用户在某个应用下的所有数据
func (u *User) UnbindApp(category util.CAtegory) error {
	key := fmt.Sprintf("products.%s.deleted_at", string(category))
	filter := bson.M{"_id": u.UOID}
	update := bson.M{"$set": bson.M{key: time.Now()}}

	_, err := db.Collection("user").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
