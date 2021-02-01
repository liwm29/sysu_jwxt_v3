package client

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"github.com/liwm29/sysu_jwxt_v3/backend/jwxtClient/global"
	"github.com/liwm29/sysu_jwxt_v3/backend/jwxtClient/request"
	"github.com/liwm29/sysu_jwxt_v3/backend/jwxtClient/util"
	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
)

type LoginForm struct {
	Username,
	Password,
	Captcha,
	_eventId,
	execution,
	geolocation string
}

func NewLoginForm() *LoginForm {
	return &LoginForm{
		_eventId: "submit",
	}
}

func (f *LoginForm) ConvertToUrlVal() *url.Values {
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
func (c *JwxtClient) CasFirstGet(LoginForm *LoginForm) {
	global.Log.Debug(util.WhereAmI())
	LoginUrl, captchaUrl := getLoginCaptchaUrl()
	resp := request.Get(LoginUrl).Do(c)

	doc, _ := goquery.NewDocumentFromReader(resp.Reader())
	execution, isexist := doc.Find("[name=execution]").Attr("value")
	LoginForm.execution = execution
	if !isexist {
		ioutil.WriteFile("./casFirstGetResp.html", resp.Bytes(), 0666)
		util.PanicIf(errors.New("未找到页面参数'execution'| 可能是带着失效cookie访问了登录页,检查'./casFirstGetResp.html'或" + LoginUrl))
	}
	respImg := c.Get(captchaUrl).Bytes()
	ioutil.WriteFile(DEFAULT_CAPTCHA_PATH, respImg, 0666)

	// 重写postform的Url
	jwxt443LoginUrl = resp.Request.URL.String()
}

func (c *JwxtClient) LoginWithForm(LoginForm *LoginForm) bool {
	global.Log.WithField("LoginForm", LoginForm.ConvertToUrlVal().Encode()).Debug("LoginReq ", util.WhereAmI())

	LoginUrl, _ := getLoginCaptchaUrl()
	request.PostForm(LoginUrl, LoginForm.ConvertToUrlVal().Encode()).Do(c)

	return c.CheckLogin()
}

func (c *JwxtClient) LoginWithCookies() (isLogin bool, err error) {
	err = c.LoadCookies(DEFAULT_COOKIE_PATH)
	if err != nil {
		return false, err
	}
	c.activeCookie()
	isLogin = c.CheckLogin()
	return isLogin, err
}

func (c *JwxtClient) Login() bool {
	// 尝试使用cookie登陆
	if LoginOk, _ := c.LoginWithCookies(); LoginOk {
		return true
	}
	color.Red("使用cookie登陆失败,尝试使用form表单登陆")

	LoginForm := NewLoginForm()
	c.CasFirstGet(LoginForm)
	LoginForm = LoginFormCli(LoginForm)
	isLogin := c.LoginWithForm(LoginForm)
	return isLogin
}

func (c *JwxtClient) CheckLogin() bool {
	url := global.HOST + "jwxt/api/login/status"
	respJson := c.Get(url).Bytes()
	// {"code":53000000,"message":"系统错误"}
	var respMap map[string]interface{}
	if err := request.JsonConvert(respJson, &respMap); err != nil {
		return false
	}

	color.Yellow("正在检查登陆状态 url=jwxt/api/login/status resp=%#v", respMap)

	if respMap["data"].(float64) == 1 {
		global.YEAR_TERM = c.GetYearTerm()
		initClient(c)
		loginOkPrintInfo(c)
		c.StoreCookies(DEFAULT_COOKIE_PATH)
		return true
	}
	return false
}

// 当多端登陆时,比如使用多个浏览器登陆,每个浏览器都会维护各自的cookie
// 在教务系统的后端,虽然同一时间只能有一个cookie有效,但是只要再访问登陆页,便会将自己的cookie激活,而不需要再次post表单
func (c *JwxtClient) activeCookie() {
	if MODE_JWXT443 == global.MODE {
		request.Get(jwxt443LoginUrl).Do(c)
		return
	}
	color.Yellow("正在刷新cookie, url=jwxt/api/sso/cas/Login?pattern=student-Login")
	url := global.HOST + "jwxt/api/sso/cas/Login?pattern=student-Login"
	ref := global.HOST + "jwxt/"
	request.Get(url).Referer(ref).Do(c)
}

// cli,填充Loginform表单
func LoginFormCli(form *LoginForm) *LoginForm {
	username := os.Getenv("jwxtUsername")
	fmt.Println("输入用户名:", username)
	if username == "" {
		fmt.Scanf("%s\n", &username)
	}
	form.Username = username

	password := os.Getenv("jwxtPassword")
	fmt.Println("输入密码:", password)
	if password == "" {
		fmt.Scanf("%s\n", &password)
	}
	form.Password = password

	fmt.Println("输入验证码:")
	fmt.Scanf("%s\n", &form.Captcha)

	return form
}

func loginOkPrintInfo(c *JwxtClient) {
	color.HiMagenta("登陆成功,欢迎:%#v", c.User)
}
