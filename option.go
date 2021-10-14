package ginx

import "github.com/gin-gonic/gin"

var (
	defaultOpt = option{
		cookieAge:   604800,
		cookieKey:   CookieKey,
		recoverable: true,
		recoverSkip: 0,
		logRequest:  false,
	}
)

type option struct {
	cookieAge   int
	cookieKey   string
	recoverable bool
	recoverSkip int
	logRequest  bool
	middlewares []gin.HandlerFunc
}

type OptionFn func(*option)

func WithCookieAge(age int) OptionFn {
	return func(o *option) {
		o.cookieAge = age
	}
}

func WithRecovery(enable bool, recoverSkip int) OptionFn {
	return func(o *option) {
		o.recoverable = enable
		o.recoverSkip = recoverSkip
	}
}

func WithLogRequest(enable bool) OptionFn {
	return func(o *option) {
		o.logRequest = enable
	}
}

func WithMiddleware(middlewares ...gin.HandlerFunc) OptionFn {
	return func(o *option) {
		if o.middlewares == nil {
			o.middlewares = make([]gin.HandlerFunc, 0)
		}
		o.middlewares = append(o.middlewares, middlewares...)
	}
}

func WithCookieKey(key string) OptionFn {
	return func(o *option) {
		o.cookieKey = key
	}
}
