package jwxtClient

import (
	"bytes"
	"io/ioutil"
	"net/url"
	"server/backend/jwxtClient/request"

	"github.com/PuerkitoBio/goquery"
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

func (c *jwxtClient) CasFirstGet(captchaSavePath string, loginForm *loginForm) {
	indexUrl := "https://cas.sysu.edu.cn/cas/login?service=https://jwxt.sysu.edu.cn/jwxt/api/sso/cas/login?pattern=student-login"
	respHtml := c.Get(indexUrl)

	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(respHtml))
	execution, isexist := doc.Find("[name=execution]").Attr("value")
	loginForm.execution = execution
	if !isexist {
		log.Fatal("未找到页面参数\"execution\",可能登陆逻辑改变,检查 " + indexUrl)
	}
	captchaUrl := "https://cas.sysu.edu.cn/cas/captcha.jsp"
	respImg := c.Get(captchaUrl)
	ioutil.WriteFile(captchaSavePath, respImg, 0666)
}

func (c *jwxtClient) Login(loginForm *url.Values) {
	url := "https://cas.sysu.edu.cn/cas/login?service=https://jwxt.sysu.edu.cn/jwxt/api/sso/cas/login?pattern=student-login"
	c.PostForm(url, loginForm.Encode())
	// 检查是否登陆成功
	if c.CheckLogin() {
		log.WithField("user", c.username).Info("登陆成功")
	} else {
		log.WithField("user", c.username).Info("登陆失败")
	}
}

func (c *jwxtClient) CheckLogin() bool {
	url := "https://jwxt.sysu.edu.cn/jwxt/api/login/status"
	respJson := c.Get(url)
	log.WithField("respJson", string(respJson)).Info("checkLoginStatus")
	if request.JsonToMap(respJson)["code"].(float64) == 200 {
		c.isLogin = true // 也许不需要这个字段,考虑是否弃用
		return true
	} else {
		return false
	}
}
