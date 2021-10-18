package ginx

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/pinealctx/neptune/ulog"
	"go.uber.org/zap"
	"net/http"
	"time"
)

const (
	CookieKey = "api_token"
)

type GinGroup struct {
	*gin.RouterGroup
	*GinX
}

type GinX struct {
	*gin.Engine
	srv       *http.Server
	cookieAge int
	cookieKey string
}

func New(addr string, optFns ...OptionFn) *GinX {
	var opt = &option{}
	*opt = defaultOpt
	for _, optFn := range optFns {
		optFn(opt)
	}
	var s = &GinX{}
	s.cookieAge = opt.cookieAge
	s.cookieKey = opt.cookieKey
	s.srv = &http.Server{Addr: addr}
	s.Engine = gin.New()
	if opt.recoverable {
		s.Engine.Use(recovery)
	}
	if opt.logRequest {
		s.Engine.Use(logRequest)
	}
	for _, middleware := range opt.middlewares {
		s.Engine.Use(middleware)
	}
	return s
}

func (s *GinX) Serve() error {
	s.srv.Handler = s.Engine
	return s.srv.ListenAndServe()
}

func (s *GinX) Run(errChan chan error) {
	go func() {
		errChan <- s.Serve()
	}()
}

func (s *GinX) Stop() error {
	return s.srv.Close()
}

func (s *GinX) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func (s *GinX) GinGroup(relativePath string, handlers ...gin.HandlerFunc) *GinGroup {
	return &GinGroup{
		RouterGroup: s.Group(relativePath, handlers...),
		GinX:        s,
	}
}

func (s *GinX) SetCookie(c *gin.Context, val string) {
	var ck = &http.Cookie{
		Name:     s.cookieKey,
		Value:    val,
		Path:     "/",
		MaxAge:   s.cookieAge,
		HttpOnly: true,
	}
	c.Writer.Header().Set("Set-Cookie", ck.String())
}

func (s *GinX) GetCookie(c *gin.Context) (string, error) {
	var ck, err = c.Request.Cookie(s.cookieKey)
	if err != nil {
		return "", err
	}
	return ck.Value, nil
}

func (s *GinX) GetToken(c *gin.Context) string {
	var tk = c.Request.Header.Get(s.cookieKey)
	if tk != "" {
		return tk
	}
	var ck, err = s.GetCookie(c)
	if err != nil {
		return ""
	}
	return ck
}

func (s *GinX) ClearCookie(c *gin.Context) {
	var ck = &http.Cookie{
		Name:     s.cookieKey,
		MaxAge:   -1,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	}
	c.Writer.Header().Set("Set-Cookie", ck.String())
}

func recovery(c *gin.Context) {
	defer func() {
		var capture = recover()
		if capture != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			ulog.Error("http.panic", ctxFields(c, zap.Reflect("capture", capture), zap.Stack("stack"))...)
		}
	}()
	c.Next()
}

func logRequest(c *gin.Context) {
	var start = time.Now()
	c.Next()
	var end = time.Now()
	ulog.Debug("http.request", ctxFields(c, zap.Time("end", end), zap.Duration("latency", end.Sub(start)))...)
}

func ctxFields(c *gin.Context, fields ...zap.Field) []zap.Field {
	var dist = []zap.Field{
		zap.String("path", c.Request.URL.Path),
		zap.String("method", c.Request.Method),
		zap.Int("status", c.Writer.Status()),
		zap.String("ip", c.ClientIP()),
		zap.String("ua", c.Request.UserAgent()),
	}
	dist = append(dist, fields...)
	return dist
}

func init() {
	gin.SetMode(gin.ReleaseMode)
}
