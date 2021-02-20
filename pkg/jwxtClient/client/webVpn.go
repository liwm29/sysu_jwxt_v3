package client

import (
	"errors"

	"github.com/liwm29/sysu_jwxt_v3/pkg/jwxtClient/global"
	"github.com/liwm29/sysu_jwxt_v3/pkg/jwxtClient/internal/util"
	"github.com/liwm29/sysu_jwxt_v3/pkg/request"
)

const (
	jwxt443Host string = "https://jwxt-443.webvpn.sysu.edu.cn/"
	jwxtHost    string = "https://jwxt.sysu.edu.cn/"
)

const (
	jwxt443Url = "https://jwxt-443.webvpn.sysu.edu.cn/jwxt/#/login"
	jwxtUrl    = "https://jwxt.sysu.edu.cn/jwxt/#/login"
)

var (
	// portalLoginUrl   = "https://cas-443.webvpn.sysu.edu.cn/cas/login?service=https://portal.sysu.edu.cn/management/shiro-cas"

	// 注意这个loginUrl不是最终的提交表单的url,之所以要请求这个url,是为了获取_astraeus_session这个cookie;
	// 这个url将会被重定向cas-443,见casFirstGet(),loginUrl将会被重写
	jwxt443LoginUrl   = "https://jwxt-443.webvpn.sysu.edu.cn/jwxt/api/toCasUrl?toCasUrl=/"
	jwxt443CaptchaUrl = "https://cas-443.webvpn.sysu.edu.cn/cas/captcha.jsp"

	jwxtLoginUrl   = "https://cas.sysu.edu.cn/cas/login?service=https://jwxt.sysu.edu.cn/jwxt/api/sso/cas/login?pattern=student-login"
	jwxtCaptchaUrl = "https://cas.sysu.edu.cn/cas/captcha.jsp"
)

const (
	MODE_JWXT int = iota
	MODE_JWXT443
)

// 判断在校园网内部
func canAccessJwxt(c request.Clienter) bool {
	resp := request.Get(jwxtHost + "jwxt/").Do(c).Bytes()
	return util.IsAccessForbidden(resp)
}

func init() {
	switch canAccessJwxt(request.DefaultClient) {
	case true:
		global.MODE = MODE_JWXT443
		global.HOST = jwxt443Host
	case false:
		global.MODE = MODE_JWXT
		global.HOST = jwxtHost
	}
	printVersion()
}

func getLoginCaptchaUrl() (string, string) {
	switch global.MODE {
	case MODE_JWXT:
		return jwxtLoginUrl, jwxtCaptchaUrl
	case MODE_JWXT443:
		return jwxt443LoginUrl, jwxt443CaptchaUrl
	default:
		util.PanicIf(errors.New("unexpected MODE,want MODE_JWXT/MODE_PORTAL"))
		return "", ""
	}
}
