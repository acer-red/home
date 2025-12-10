package common

import (
	"bytes"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/tengfei-xy/go-log"
)

type msgErr int

type message struct {
	Err  msgErr      `json:"err"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

const (
	mseqOK msgErr = iota
	mseqCreated
)
const mstrOK string = "ok"
const mstrCreated string = "创建成功"

func (msg message) setData(data interface{}) message {
	msg.Data = data
	return msg
}

func msgOK(msg ...string) message {
	m := message{Err: mseqOK}
	if len(msg) > 0 {
		m.Msg = msg[0]
	} else {
		m.Msg = mstrOK
	}
	return m
}

func msgCreated(msg ...string) message {
	m := message{Err: mseqOK}
	if len(msg) > 0 {
		m.Msg = msg[0]
	} else {
		m.Msg = mstrCreated
	}
	return m
}

func Ok(g *gin.Context) {
	d := msgOK()
	log.Debug3j(d)
	g.JSON(http.StatusOK, d)
}

func OkData(g *gin.Context, obj any) {
	d := msgOK().setData(obj)
	log.Debug3j(d)
	g.JSON(http.StatusOK, d)
}

func OkImage(g *gin.Context, data bytes.Buffer) {
	name := g.Param("file")
	if !strings.Contains(name, ".") {
		BadRequest(g)
		return
	}
	fotmat := strings.ToLower(strings.Split(name, ".")[1])
	switch fotmat {
	case "png":
		g.Data(http.StatusOK, "image/png", data.Bytes())
	case "jpg":
	case "jpeg":
		g.Data(http.StatusOK, "image/jpeg", data.Bytes())
	default:
		log.Errorf("未知图片类型,%s", fotmat)
	}
}

func BadRequest(g *gin.Context) {
	g.AbortWithStatus(http.StatusBadRequest)
}
func InternalServerError(g *gin.Context) {
	g.AbortWithStatus(http.StatusInternalServerError)
}

func NotFound(g *gin.Context) {
	g.AbortWithStatus(http.StatusNotFound)
}
func Conflict(g *gin.Context) {
	g.AbortWithStatus(http.StatusConflict)
}
func CreatedData(g *gin.Context, obj any) {
	d := msgCreated().setData(obj)
	log.Debug3j(d)
	g.JSON(http.StatusCreated, d)
}

func Unauthorized(g *gin.Context) {
	g.AbortWithStatus(http.StatusUnauthorized)
}
