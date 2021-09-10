package ginx

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	NameCode = "code"
	NameMsg  = "errMsg"
	NameData = "data"

	NotFound    = `not.found`
	InternalErr = `internal.error`
	_NeedLogin  = `need.login`

	ResCodeOK         = 2000
	ResCodeBadRequest = 4000
	_ResCodeNeedLogin = 4001
	ResCodeInternal   = 5000
)

//RspOK response 200 - ok
func RspOK(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{NameCode: ResCodeOK})
}

//RspData response 200 -  with data
func RspData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		NameCode: ResCodeOK,
		NameData: data,
	})
}

//RspBadRequest response 400 - bad request
func RspBadRequest(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, gin.H{
		NameCode: ResCodeBadRequest,
		NameMsg:  msg,
	})
	c.Abort()
}

//RspNeedLogin response 4001 - need login
func RspNeedLogin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		NameCode: _ResCodeNeedLogin,
		NameMsg:  _NeedLogin,
	})
	c.Abort()
}

//RspInternalError response 500 - internal error
func RspInternalError(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		NameCode: ResCodeInternal,
		NameMsg:  InternalErr,
	})
	c.Abort()
}

//RspErrMsg response customer response code and msg
func RspErrMsg(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, gin.H{
		NameCode: code,
		NameMsg:  msg,
	})
	c.Abort()
}
