package jwxtClient

import (
	"errors"
	"fmt"
	"os"

	// "fmt"
	"io/ioutil"
	"net/url"
	"server/backend/jwxtClient/course"
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
func (c *JwxtClient) CasFirstGet(loginForm *loginForm) {
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
	ioutil.WriteFile(DEFAULT_CAPTCHA_PATH, respImg, 0666)
}

func (c *JwxtClient) LoginWithForm(loginForm *loginForm) (isLogin bool, err error) {
	log.WithFields(logrus.Fields{"Username": loginForm.Username, "Passwrod": loginForm.Password, "captcha": loginForm.Captcha}).Info("User")
	log.WithField("loginForm", loginForm.ConvertToUrlVal().Encode()).Debug("loginReq")
	log.Debug(util.WhereAmI())
	url := "https://cas.sysu.edu.cn/cas/login?service=https://jwxt.sysu.edu.cn/jwxt/api/sso/cas/login?pattern=student-login"
	html := c.PostForm(url, loginForm.ConvertToUrlVal().Encode()).Bytes()
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

func (c *JwxtClient) LoginWithCookies() (isLogin bool, err error) {
	err = c.LoadCookies(DEFAULT_COOKIE_PATH)
	if err != nil {
		log.WithField("error", err.Error()).Debug(util.WhereAmI())
		return false, err
	}
	c.activeCookie()
	isLogin = c.CheckLogin()
	log.WithField("isLogin", isLogin).Debug(util.WhereAmI())
	return isLogin, err
}

func (c *JwxtClient) Login() (bool, error) {
	// 尝试使用cookie登陆
	if loginOk, _ := c.LoginWithCookies(); loginOk {
		return true, nil
	}

	loginForm := NewLoginForm()
	c.CasFirstGet(loginForm)
	loginForm = LoginFormCli(loginForm)

	isLogin, err := c.LoginWithForm(loginForm)
	if err != nil {
		return isLogin, err
	}
	c.StoreCookies(DEFAULT_COOKIE_PATH)
	return isLogin, nil
}

func (c *JwxtClient) CheckLogin() bool {
	url := "https://jwxt.sysu.edu.cn/jwxt/api/login/status"
	respJson := c.Get(url).Bytes()
	log.WithField("respJson", util.Truncate100(string(respJson))).Debug(util.WhereAmI())
	if request.JsonToMap(respJson)["data"].(float64) == 1 {
		// 获取学期信息
		course.SetYearTerm(c.GetYearTerm())
		return true
	} else {
		return false
	}
}

// 当多端登陆时,比如使用多个浏览器登陆,每个浏览器都会维护各自的cookie
// 在教务系统的后端,虽然同一时间只能有一个cookie有效,但是只要访问登陆页,便会将自己的cookie激活
func (c *JwxtClient) activeCookie() {
	url := "https://jwxt.sysu.edu.cn/jwxt/api/sso/cas/login?pattern=student-login"
	ref := "https://jwxt.sysu.edu.cn/jwxt/"
	request.Get(url).Referer(ref).Do(c)
}

// cli,填充loginform表单
func LoginFormCli(form *loginForm) *loginForm {
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
