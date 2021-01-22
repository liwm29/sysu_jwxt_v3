package jwxtClient

import (
	"errors"
	// "fmt"
	"io/ioutil"
	"net/url"
	"server/backend/jwxtClient/request"
	"server/backend/jwxtClient/util"

	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
)

type loginForm struct {
	Username,
	Password,
	Captcha,
	_eventId,
	execution,
	geolocation string
}

func NewLoginForm() *loginForm {
	return &loginForm{
		_eventId: "submit",
	}
}

func (f *loginForm) ConvertToUrlVal() *url.Values {
	return &url.Values{
		"username":    {f.Username},
		"password":    {f.Password},
		"captcha":     {f.Captcha},
		"_eventId":    {f._eventId},
		"execution":   {f.execution},
		"geolocation": {f.geolocation},
	}
}

// 向cas中央认证系统获取验证码
func (c *JwxtClient) CasFirstGet(captchaSavePath string, loginForm *loginForm) {
	log.Debug(util.WhereAmI())
	indexUrl := "https://cas.sysu.edu.cn/cas/login?service=https://jwxt.sysu.edu.cn/jwxt/api/sso/cas/login?pattern=student-login"
	ref := "https://jwxt.sysu.edu.cn/jwxt/"
	resp := request.Get(indexUrl).Referer(ref).Do(c)

	doc, _ := goquery.NewDocumentFromReader(resp.Reader())
	execution, isexist := doc.Find("[name=execution]").Attr("value")
	loginForm.execution = execution
	if !isexist {

		ioutil.WriteFile("./casFirstGetResp.html", resp.Bytes(), 0666)
		util.PanicIf(errors.New("未找到页面参数'execution',可能登陆逻辑改变;可能是带着失效cookie访问了登录页,检查'./casFirstGetResp.html'或" + indexUrl))
	}
	captchaUrl := "https://cas.sysu.edu.cn/cas/captcha.jsp"
	respImg := c.Get(captchaUrl).Bytes()
	ioutil.WriteFile(captchaSavePath, respImg, 0666)
}

func (c *JwxtClient) Login(loginForm *url.Values) (isLogin bool, err error) {
	log.WithFields(logrus.Fields{"Username": loginForm.Get("username"), "Passwrod": loginForm.Get("password"), "captcha": loginForm.Get("captcha")}).Info("User")
	log.WithField("loginForm", loginForm.Encode()).Debug("loginReq")
	log.Debug(util.WhereAmI())
	url := "https://cas.sysu.edu.cn/cas/login?service=https://jwxt.sysu.edu.cn/jwxt/api/sso/cas/login?pattern=student-login"
	html := c.PostForm(url, loginForm.Encode()).Bytes()
	// 检查是否登陆成功
	if c.CheckLogin() {
		log.WithField("user", c.username).Info("登陆成功")
		return true, nil
	} else {
		log.WithField("user", c.username).Info("登陆失败")
		ioutil.WriteFile("./loginFailResp.html", html, 0666)
		return false, errors.New("详情请看./loginFailResp.html")
	}
}

func (c *JwxtClient) CheckLogin() bool {
	url := "https://jwxt.sysu.edu.cn/jwxt/api/login/status"
	respJson := c.Get(url).Bytes()
	log.WithField("respJson", util.Truncate100(string(respJson))).Debug(util.WhereAmI())
	if request.JsonToMap(respJson)["data"].(float64) == 1 {
		return true
	} else {
		return false
	}
}

func (c *JwxtClient) LoginWithCookies(cookiePath string) (isLogin bool, err error) {
	err = c.LoadCookies(cookiePath)
	if err != nil {
		log.WithField("error", err.Error()).Debug(util.WhereAmI())
		return false, err
	}
	c.activeCookie()
	isLogin = c.CheckLogin()
	log.WithField("isLogin", isLogin).Debug(util.WhereAmI())
	return isLogin, err
}

// 当多端登陆时,比如使用多个浏览器登陆,每个浏览器都会维护各自的cookie
// 在教务系统的后端,虽然同一时间只能有一个cookie有效,但是只要访问登陆页,便会将自己的cookie激活
func (c *JwxtClient) activeCookie() {
	url := "https://jwxt.sysu.edu.cn/jwxt/api/sso/cas/login?pattern=student-login"
	ref := "https://jwxt.sysu.edu.cn/jwxt/"
	request.Get(url).Referer(ref).Do(c)
}
