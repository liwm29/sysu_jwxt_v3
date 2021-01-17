package jwxtClient

import (
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

func (c *JwxtClient) CasFirstGet(captchaSavePath string, loginForm *loginForm) {
	indexUrl := "https://cas.sysu.edu.cn/cas/login?service=https://jwxt.sysu.edu.cn/jwxt/api/sso/cas/login?pattern=student-login"
	resp := c.Get(indexUrl)

	doc, _ := goquery.NewDocumentFromReader(resp.Reader())
	execution, isexist := doc.Find("[name=execution]").Attr("value")
	loginForm.execution = execution
	if !isexist {
		log.Fatal("未找到页面参数\"execution\",可能登陆逻辑改变,检查 " + indexUrl)
	}
	captchaUrl := "https://cas.sysu.edu.cn/cas/captcha.jsp"
	respImg := c.Get(captchaUrl).Bytes()
	ioutil.WriteFile(captchaSavePath, respImg, 0666)
}

func (c *JwxtClient) Login(loginForm *url.Values) {
	log.Debug("登陆表单:", loginForm.Encode())
	log.Info("User:", loginForm.Get("username"), loginForm.Get("password"), loginForm.Get("captcha"))
	url := "https://cas.sysu.edu.cn/cas/login?service=https://jwxt.sysu.edu.cn/jwxt/api/sso/cas/login?pattern=student-login"
	html := c.PostForm(url, loginForm.Encode()).Bytes()
	// 检查是否登陆成功
	if c.CheckLogin() {
		log.WithField("user", c.username).Info("登陆成功")
	} else {
		log.WithField("user", c.username).Info("登陆失败")
		ioutil.WriteFile("./loginFailResp.html", html, 0666)
	}
}

func (c *JwxtClient) CheckLogin() bool {
	url := "https://jwxt.sysu.edu.cn/jwxt/api/login/status"
	respJson := c.Get(url).Bytes()
	log.WithField("respJson", string(respJson)).Debug("检查登陆状态")
	if request.JsonToMap(respJson)["data"].(float64) == 1 {
		c.isLogin = true // 也许不需要这个字段,考虑是否弃用
		return true
	} else {
		return false
	}
}
