package ginx

import (
	"github.com/gin-gonic/gin"
	"net/http/pprof"
)

type ProfX struct {
	*GinX
}

func (s *ProfX) handleIndex(c *gin.Context) {
	pprof.Index(c.Writer, c.Request)
}

func (s *ProfX) handleCmdLine(c *gin.Context) {
	pprof.Cmdline(c.Writer, c.Request)
}

func (s *ProfX) handleProfile(c *gin.Context) {
	pprof.Profile(c.Writer, c.Request)
}

func (s *ProfX) handleSymbol(c *gin.Context) {
	pprof.Symbol(c.Writer, c.Request)
}

func (s *ProfX) handleTrace(c *gin.Context) {
	pprof.Trace(c.Writer, c.Request)
}

func (s *ProfX) genHandle(name string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var f = pprof.Handler(name)
		f.ServeHTTP(c.Writer, c.Request)
	}
}

func (s *ProfX) setupRoutine() {
	s.GET("/", s.handleIndex)
	s.GET("/cmdline", s.handleCmdLine)
	s.GET("/profile", s.handleProfile)
	s.GET("/symbol", s.handleSymbol)
	s.GET("/trace", s.handleTrace)

	s.GET("/goroutine", s.genHandle(`goroutine`))
	s.GET("/heap", s.genHandle(`heap`))
	s.GET("/threadcreate", s.genHandle(`threadcreate`))
	s.GET("/block", s.genHandle(`block`))
	s.GET("/mutex", s.genHandle(`mutex`))
	s.GET("/allocs", s.genHandle(`allocs`))
}

func NewProfX(addr string, opts ...OptionFn) *ProfX {
	var p = &ProfX{}
	p.GinX = New(addr, opts...)
	p.setupRoutine()
	return p
}
