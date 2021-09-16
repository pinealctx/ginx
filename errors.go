package ginx

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/pinealctx/neptune/tex"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ResCode 响应码
type ResCode uint32

const (
	// RCodeSuccess 请求成功。
	RCodeSuccess ResCode = 2000
	// RCodeBadRequest 一般情况下这类错误是用户行为产生的，这类错误需要将错误描述根据语言进行翻译，客户端可直接显示出来。
	RCodeBadRequest ResCode = 4000
	// RCodeNeedLogin 鉴权失败，用户需重新登录。
	RCodeNeedLogin ResCode = 4001
	// RCodeNoPermission 没有权限，用户已登录但没有访问该功能的权限。
	RCodeNoPermission ResCode = 4002
	// RCodeInvalidRequest 无效的请求，一般存在客户端对协议的实现有问题（参数错误）或非法客户端等。
	RCodeInvalidRequest ResCode = 4003
	// RCodeInternal 服务端内部错误，需提示用户稍后重试等。
	RCodeInternal ResCode = 5000

	// DefaultLang 默认语言
	DefaultLang = "zh"
	// LangKey 语言字段在cookie/url/header的名称
	LangKey = "lang"
)

var (
	errMsgHash map[string]map[string]string
)

// SetupErrI18nFile 设置i18n文件
func SetupErrI18nFile(filepath string) error {
	if errMsgHash != nil {
		return nil
	}
	return tex.LoadJSONFile2Obj(filepath, &errMsgHash)
}

// ParseResCode 根据错误解析出ResCode和描述。
// 所有返回给客户端的错误，都应该是可控制的，因此错误码均使用rpc的错误码来限制错误类型，
// 当通过rpc的status解析不出来均认为是服务端内部错误 -> RCodeInternal
// 当错误码为 msg == RCodeBadRequest 时会将错误描述翻译出来。
// 当错误码 >= 5000 时不显示具体错误
func ParseResCode(err error, lang string) (code ResCode, errmsg string) {
	if err == nil {
		code = RCodeSuccess
		return
	}
	var s, ok = status.FromError(err)
	if !ok || s.Code() < codes.Code(RCodeSuccess) {
		code = RCodeInternal
		return
	}
	code = ResCode(s.Code())
	if code >= RCodeInternal {
		return
	}
	errmsg = s.Message()
	if code != RCodeBadRequest {
		return
	}
	var hash map[string]string
	hash, ok = errMsgHash[lang]
	if !ok {
		hash, ok = errMsgHash[DefaultLang]
		if !ok {
			return
		}
	}
	var dist string
	dist, ok = hash[errmsg]
	if ok {
		errmsg = dist
	}
	return
}

// MakeErr 构造错误
func MakeErr(code ResCode, msg string) error {
	return status.Error(codes.Code(code), msg)
}

// MakeBadRequestErr 构造BadRequest错误
func MakeBadRequestErr(msg string) error {
	return MakeErr(RCodeBadRequest, msg)
}

// MakeNeedLoginErr 构造NeedLogin错误
func MakeNeedLoginErr() error {
	return MakeErr(RCodeNeedLogin, "need.login")
}

// MakeNoPermissionErr 构造NoPermission错误
func MakeNoPermissionErr() error {
	return MakeErr(RCodeNoPermission, "no.permission")
}

// MakeInvalidRequestErr 构造InvalidRequest错误
func MakeInvalidRequestErr(msg string) error {
	return MakeErr(RCodeInvalidRequest, msg)
}

// MakeInternalErr 构造Internal错误
func MakeInternalErr() error {
	return MakeErr(RCodeInternal, "internal.error")
}

// IsError 判断错误是否相等，在原生错误判断的基础上加上了grpc status判断
func IsError(err, target error) bool {
	if target == nil {
		return err == target
	}
	var errStatus, ok = status.FromError(err)
	if !ok {
		return errors.Is(err, target)
	}
	var targetStatus *status.Status
	targetStatus, ok = status.FromError(target)
	if !ok {
		return false
	}
	return errStatus.Code() == targetStatus.Code() && errStatus.Message() == targetStatus.Message()
}

// GetLangFromContext 从gin的context里获取语言类型
func GetLangFromContext(c *gin.Context) string {
	var lang = c.Query(LangKey)
	if IsValidLang(lang) {
		return lang
	}
	lang = c.GetHeader(LangKey)
	if IsValidLang(lang) {
		return lang
	}
	lang, _ = c.Cookie(LangKey)
	if IsValidLang(lang) {
		return lang
	}
	lang = c.PostForm(LangKey)
	if IsValidLang(lang) {
		return lang
	}
	return DefaultLang
}

// IsValidLang 判断是否为合法的语言类型
func IsValidLang(lang string) bool {
	return lang == "zh" || lang == "en"
}
