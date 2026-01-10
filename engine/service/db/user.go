package db

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"errors"

	"github.com/acer-red/official/engine/service/cache"
	"github.com/acer-red/official/engine/util"

	"github.com/google/uuid"
	"github.com/tengfei-xy/go-log"
	"gorm.io/gorm"
)

const passwordMinLen = 8
const loginJWTDuration = 30 * 24 * time.Hour

type Cookie struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}
type Avatar struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type Profile struct {
	Avatar   Avatar `json:"avatar" gorm:"embedded"`
	Nickname string `json:"nickname"`
}
type ResponseGetUserInfo struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Profile   Profile   `json:"profile"`
	CreatedAt time.Time `json:"created_at"`
	API       []API     `json:"api"`
	ProductID string    `json:"product_id,omitempty"`
	DeviceID  string    `json:"device_id,omitempty"`
}

type User struct {
	Base
	Username   string    `gorm:"unique;size:20" json:"username"`
	Email      string    `gorm:"unique" json:"email"`
	Password   string    `json:"-"`
	PublicKey  []byte    `json:"-"`
	Nickname   string    `json:"-"`
	AvatarName string    `json:"-"`
	AvatarURL  string    `json:"-"`
	Products   []Product `gorm:"foreignKey:UserID" json:"-"`
	IsDeleted  bool      `gorm:"-" json:"-"`
}

// Helpers
func (u *User) GetProfile() Profile {
	return Profile{
		Nickname: u.Nickname,
		Avatar: Avatar{
			Name: u.AvatarName,
			URL:  u.AvatarURL,
		},
	}
}

// 请求用户注册的解析结构体
type RequestUserRegister struct {
	Username    string        `json:"username"`
	Password    string        `json:"password"`
	Email       string        `json:"email"`
	CategoryStr string        `json:"category"`
	PublicKey   string        `json:"public_key"`
	Category    util.CAtegory `json:"-"`
	Cookie      Cookie        `json:"-"`
	profile     Profile       `json:"-"`
}

type RequestUserLogin struct {
	Account   string        `json:"account"`
	Password  string        `json:"password"`
	Category  string        `json:"category"`
	UID       string        `json:"uid"`
	category  util.CAtegory `json:"-"`
	user      User          `json:"-"`
	productID string        `json:"-"`
	deviceID  string        `json:"-"`
}

type RequestPutUserInfo struct {
	Nickname string `json:"nickname"`
	Password string `json:"password"`
	UserID   string `json:"-"`
}

// cookie
func (c *Cookie) setLoginCookie() {
	c.Key = "login"
	c.Value = util.RandomString32()
	c.CreatedAt = time.Now()
	c.ExpiresAt = time.Now().AddDate(0, 1, 0)
}

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
	allowedChars := regexp.MustCompile(`^[a-zA-Z0-9_.]+$`)
	if !allowedChars.MatchString(username) {
		return false
	}
	onlyNumbers := regexp.MustCompile(`^\d+$`)
	if onlyNumbers.MatchString(username) && len(username) < 6 {
		return false
	}
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
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return false
	}
	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return false
	}
	if !regexp.MustCompile(`\d`).MatchString(password) {
		return false
	}
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
	var count int64
	if err := db.Model(&User{}).Where("username = ?", req.Username).Count(&count).Error; err != nil {
		log.Error(err)
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	if err := db.Model(&User{}).Where("email = ?", req.Email).Count(&count).Error; err != nil {
		log.Error(err)
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func (req *RequestUserRegister) FindUser() (*User, bool, error) {
	var user User
	err := db.Where("username = ? OR email = ?", req.Username, req.Email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}
	// Logic regarding deleted product removed as Product table in DDL has no deleted_at
	return &user, true, nil
}

func (req *RequestUserRegister) Reactivate(u *User) (string, API, string, error) {
	var err error
	req.Password, err = util.HashPassword(req.Password)
	if err != nil {
		return "", API{}, "", err
	}

	u.Password = req.Password
	u.PublicKey = []byte(req.PublicKey)

	if err := db.Save(u).Error; err != nil {
		return "", API{}, "", err
	}

	// Ensure Product exists
	var p Product
	if err := db.Where("user_id = ? AND category = ?", u.ID, req.Category).First(&p).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			p = Product{
				UserID:   u.ID,
				Category: string(req.Category),
			}
			if err := db.Create(&p).Error; err != nil {
				return "", API{}, "", err
			}
		} else {
			return "", API{}, "", err
		}
	}

	// Create Device
	// Since we cannot store client UID in device table (DDL restriction), we create a new Device.
	// If req.UID is a valid UUID, maybe we could try to use it as ID?
	// But let's just generate new one or rely on default.
	d := Device{
		ProductID: p.ID,
	}
	if err := db.Create(&d).Error; err != nil {
		return "", API{}, "", err
	}

	api := newAPI()
	api.ProductID = p.ID
	// api.Category not in struct matches DDL? No "category" in API DDL.

	if err := db.Create(&api).Error; err != nil {
		return "", API{}, "", err
	}
	return u.ID.String(), api, d.ID.String(), nil
}

func (req *RequestUserRegister) RandomAccount() {
	req.Username = util.RandomString32()
	req.Password = util.RandomString32()
	req.Email = fmt.Sprintf("%s@%s.com", util.RandomString32(), util.RandomString32())
}

func (req *RequestUserRegister) BuildProfile(userID string) error {
	var user User
	if err := db.Select("id").Where("id = ?", userID).First(&user).Error; err != nil {
		return err
	}

	data := bytes.NewBuffer(util.RandomAvatarBytes())
	filename := fmt.Sprintf("%s.png", util.RandomString32())

	if err := ImageAvatarCreate(filename, data, user.ID); err != nil {
		log.Error(err)
		return err
	}
	req.profile = Profile{
		Nickname: util.RandomNickname(),
		Avatar:   Avatar{Name: filename, URL: setAvatarUrl(filename)},
	}
	return db.Model(&User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"nickname":    req.profile.Nickname,
		"avatar_name": req.profile.Avatar.Name,
		"avatar_url":  req.profile.Avatar.URL,
	}).Error
}

func (req *RequestUserRegister) CancelRegister(userID string) {
	db.Where("id = ?", userID).Delete(&User{})
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

func (req *RequestUserRegister) Register(role util.Role) (string, API, string, error) {
	var err error
	req.Password, err = util.HashPassword(req.Password)
	if err != nil {
		log.Error(err)
		return "", API{}, "", err
	}

	ok := req.checkPublicKey()
	if !ok {
		return "", API{}, "", fmt.Errorf("不是有效的公钥")
	}

	user := User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password,
		PublicKey: []byte(req.PublicKey),
	}

	api := newAPI()
	// api.Category = string(req.Category)

	if err := db.Create(&user).Error; err != nil {
		log.Error(err)
		return "", API{}, "", err
	}

	var deviceID string
	if req.Category != util.CAtegoryOfficial {
		// api.UserID = user.ID

		p := Product{
			UserID:   user.ID,
			Category: string(req.Category),
		}
		if err := db.Create(&p).Error; err != nil {
			return "", API{}, "", err
		}

		// Create associated device
		d := Device{
			ProductID: p.ID,
		}
		if err := db.Create(&d).Error; err != nil {
			return "", API{}, "", err
		}
		deviceID = d.ID.String()

		api.ProductID = p.ID // Link API to Product

		if err := db.Create(&api).Error; err != nil {
			return "", API{}, "", err
		}
	} else {
		req.Cookie.setLoginCookie()
	}

	return user.ID.String(), api, deviceID, nil
}

func (req *RequestUserLogin) checkCatetory() error {
	c, ok := util.GetCategory(req.Category)
	if !ok {
		return fmt.Errorf("无效的产品类别")
	}
	req.category = c
	return nil
}
func (req *RequestUserLogin) Check() bool {
	if req.Account == "" {
		log.Warn("用户名为空")
		return false
	}
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
	var user User
	err := db.Where("username = ? OR email = ?", req.Account, req.Account).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	var activeCount int64
	db.Model(&Product{}).Where("user_id = ? AND category = ?", user.ID, req.category).Count(&activeCount)
	if activeCount == 0 {
		return false, nil
	}

	req.user = user
	return true, nil
}

func (req *RequestUserLogin) ComparePassword() error {
	return util.ComparePassword(req.user.Password, req.Password)
}
func (req *RequestUserLogin) CategoryValue() util.CAtegory {
	return req.category
}
func (req *RequestUserLogin) GetAccountID() string {
	return req.user.ID.String()
}
func (req *RequestUserLogin) GetUserID() string {
	return req.UID
}

func (req *RequestUserLogin) CheckNewDevice() error {
	log.Debugf("检查是否新设备UID=%s", req.UID)
	// With DDL restrictions, we cannot look up device by client UID.
	// We will create a new Device record to represent this login session on a device.

	var p Product
	if err := db.Where("user_id = ? AND category = ?", req.user.ID, req.category).First(&p).Error; err != nil {
		return err
	}
	req.productID = p.ID.String()

	// Create new Device
	d := Device{
		ProductID: p.ID,
	}
	if err := db.Create(&d).Error; err != nil {
		return err
	}
	req.deviceID = d.ID.String()
	return nil
}

func (req *RequestUserLogin) buildJWT(accountID, uid string, expireAt time.Time) (string, error) {
	ttl := time.Until(expireAt)
	if ttl <= 0 {
		ttl = loginJWTDuration
	}
	return util.CreateJWT(accountID, uid, req.user.Username, req.user.Email, req.category, ttl)
}

func (u *User) ToResponse(apis []API) ResponseGetUserInfo {
	res := ResponseGetUserInfo{
		ID:        u.ID.String(),
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		Profile:   u.GetProfile(),
		API:       apis,
	}
	return res
}

func (req *RequestUserLogin) BuildLoginResponse() ResponseGetUserInfo {
	var apis []API
	db.Table("api").Joins("JOIN product ON product.id = api.product_id").Where("product.user_id = ? AND product.category = ?", req.user.ID, req.category).Find(&apis)

	res := req.user.ToResponse(apis)
	res.ProductID = req.productID
	res.DeviceID = req.deviceID
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

func (req *RequestPutUserInfo) Update() error {
	updates := make(map[string]interface{})
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Password != "" {
		hashed, err := util.HashPassword(req.Password)
		if err != nil {
			return err
		}
		updates["password"] = hashed
	}
	if len(updates) == 0 {
		return nil
	}
	// Use ID (UUID)
	return db.Model(&User{}).Where("id = ?", req.UserID).Updates(updates).Error
}

func (u *User) IsNoID() bool {
	return u.ID == uuid.Nil
}
func (u *User) DeleteCookie() error {
	return nil
}
func (u *User) UpdateNickname(nickname string) error {
	return db.Model(&User{}).Where("id = ?", u.ID).Update("nickname", nickname).Error
}
func (u *User) UpdateAvatar(data io.Reader, ext string) error {
	if u.AvatarName != "" {
		db.Where("name = ?", u.AvatarName).Delete(&File{})
	}
	filename := fmt.Sprintf("%s%s", util.RandomString32(), ext)
	if err := ImageAvatarCreate(filename, data, u.ID); err != nil {
		log.Error(err)
		return err
	}
	u.AvatarName = filename
	u.AvatarURL = setAvatarUrl(filename)
	return db.Model(&User{}).Where("id = ?", u.ID).Updates(map[string]interface{}{
		"avatar_name": u.AvatarName,
		"avatar_url":  u.AvatarURL,
	}).Error
}

func (u *User) Delete() error {
	return db.Where("id = ?", u.ID).Delete(&User{}).Error
}

func GetUserFromAPI(apiKey string) (User, *uuid.UUID, bool, error) {
	var api API
	if err := db.Where("api_key = ?", apiKey).First(&api).Error; err != nil {
		return User{}, nil, false, nil
	}

	// API links to Product
	var p Product
	if err := db.Where("id = ?", api.ProductID).First(&p).Error; err != nil {
		return User{}, nil, false, err
	}

	var user User
	// removed Preload("APIs")
	if err := db.First(&user, "id = ?", p.UserID).Error; err != nil {
		return User{}, nil, false, err
	}
	return user, &p.ID, true, nil
}

func GetUserByIDAndCategory(uid string, category util.CAtegory) (User, *uuid.UUID, bool, error) {
	var user User
	err := db.Where("id = ?", uid).First(&user).Error
	if err != nil {
		return User{}, nil, false, nil
	}
	var p Product
	if err := db.Where("user_id = ? AND category = ?", user.ID, category).First(&p).Error; err != nil {
		return User{}, nil, false, nil
	}
	return user, &p.ID, true, nil
}

func setAvatarUrl(f string) string {
	return "/images/" + f
}
func newAPI() API {
	return API{
		APIKey:     util.CreateAPIKey(),
		ExpiresAt:  time.Now().AddDate(0, 3, 0),
		LastUsedAt: time.Now(),
		UsedTimes:  0,
	}
}

func (u *User) UnbindApp(category util.CAtegory) error {
	// Delete Product (UserProduct replaced by Product)
	return db.Where("user_id = ? AND category = ?", u.ID, category).Delete(&Product{}).Error
}

func GetAPIsForUser(userID uuid.UUID) ([]API, error) {
	var apis []API
	err := db.Table("api").Joins("JOIN product ON product.id = api.product_id").Where("product.user_id = ?", userID).Find(&apis).Error
	return apis, err
}
