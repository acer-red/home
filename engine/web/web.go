package web

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/acer-red/home/engine/sys"

	"github.com/gin-gonic/gin"
	log "github.com/tengfei-xy/go-log"
)

func Init(env sys.Web) {
	gin.SetMode(gin.ReleaseMode)
	g := gin.Default()

	g.Use(setEnv(env))
	g.Use(loggerMiddleware())

	if env.CORS.Enable {
		log.Infof("Enable CORS, Origin:%s", env.CORS.AllowOrigin)
		g.Use(cors(env.CORS.AllowOrigin))
	} else {

	}

	// 设定路由
	RouteFeedback(g)
	RouteUser(g)
	RouterImageGet(g)

	if env.Server.SslEnable {
		err := g.RunTLS(fmt.Sprintf(":%d", env.Server.Port), env.Server.CrtFile, env.Server.KeyFile)
		if err != nil {
			log.Fatal(err)
		}
		return
	} else {
		err := g.Run(fmt.Sprintf(":%d", env.Server.Port))
		if err != nil {
			log.Fatal(err)
		}
	}

}
func setEnv(env sys.Web) gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Set("env", env)
		c.Set("cors_origin", "")
	}
}

func setCookie(c *gin.Context, key, value string, ex int) {
	c.SetCookie(
		key,                               // Cookie 的名称
		value,                             // Cookie 的值
		ex,                                // Cookie 的过期时间 (Unix 时间戳)
		"/",                               // Cookie 的路径 (通常设置为 "/")
		c.MustGet("cors_origin").(string), // Cookie 的域名 (留空表示当前域名)
		false,                             // 是否只允许 HTTPS 访问
		false,                             // 是否禁止 JavaScript 访问 (HttpOnly)
	)
}

func setJWTCookie(c *gin.Context, token string, expireAt time.Time) {
	maxAge := int(time.Until(expireAt).Seconds())
	if maxAge <= 0 {
		maxAge = int((30 * 24 * time.Hour).Seconds())
	}

	c.SetCookie(
		"jwt",
		token,
		maxAge,
		"/",
		c.MustGet("cors_origin").(string),
		false,
		true,
	)
}

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
func loggerMiddleware() gin.HandlerFunc {
	return func(g *gin.Context) {
		var requestInfo strings.Builder
		requestInfo.WriteString("========== REQUEST ==========\n")
		requestInfo.WriteString(fmt.Sprintf("%s %s\n", g.Request.Method, g.Request.URL.Path))

		if len(g.Request.URL.RawQuery) > 0 {
			requestInfo.WriteString(fmt.Sprintf("Query: %s\n", g.Request.URL.RawQuery))
		}

		requestInfo.WriteString("Headers:\n")
		for key, values := range g.Request.Header {
			for _, value := range values {
				requestInfo.WriteString(fmt.Sprintf("  %s: %s\n", key, value))
			}
		}

		var bodyBytes []byte
		if g.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(g.Request.Body)
			g.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		if len(bodyBytes) > 0 {
			contentType := g.Request.Header.Get("Content-Type")
			if strings.Contains(contentType, "text") ||
				strings.Contains(contentType, "json") ||
				strings.Contains(contentType, "xml") ||
				contentType == "" {
				requestInfo.WriteString(fmt.Sprintf("Body:\n%s\n", string(bodyBytes)))
			} else {
				requestInfo.WriteString(fmt.Sprintf("Body: [%s, %d bytes]\n", contentType, len(bodyBytes)))
			}
		}
		requestInfo.WriteString("=============================")

		log.Debug3f("\n%s", requestInfo.String())

		// 使用自定义ResponseWriter捕获响应体
		writer := &responseWriter{
			ResponseWriter: g.Writer,
			body:           bytes.NewBufferString(""),
		}
		g.Writer = writer

		g.Next()

		var responseInfo strings.Builder
		responseInfo.WriteString("========== RESPONSE ==========\n")
		responseInfo.WriteString(fmt.Sprintf("%s %s | %d\n", g.Request.Method, g.Request.URL.Path, g.Writer.Status()))

		responseInfo.WriteString("Headers:\n")
		for key, values := range g.Writer.Header() {
			for _, value := range values {
				responseInfo.WriteString(fmt.Sprintf("  %s: %s\n", key, value))
			}
		}

		if writer.body.Len() > 0 {
			contentType := g.Writer.Header().Get("Content-Type")
			if strings.Contains(contentType, "text") ||
				strings.Contains(contentType, "json") ||
				strings.Contains(contentType, "xml") ||
				contentType == "" {
				responseInfo.WriteString(fmt.Sprintf("Body:\n%s\n", writer.body.String()))
			} else {
				responseInfo.WriteString(fmt.Sprintf("Body: [%s, %d bytes]\n", contentType, writer.body.Len()))
			}
		}
		responseInfo.WriteString("==============================")

		log.Debug2f("\n%s", responseInfo.String())
	}
}
