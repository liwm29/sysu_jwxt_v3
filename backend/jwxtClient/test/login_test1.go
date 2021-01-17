package main

import (
	"fmt"
	"io/ioutil"
	"os"
	jwxt "server/backend/jwxtClient"
	// "testing"
)

// ok
func TestLogin() *jwxt.JwxtClient {
	c := jwxt.NewClient("primaryID")
	loginForm := jwxt.NewLoginForm()
	c.CasFirstGet("./captcha.jpg", loginForm)

	username := os.Getenv("jwxtUsername")
	fmt.Println("输入用户名:", username)
	if username == "" {
		fmt.Scanf("%s\n", &username)
	}
	loginForm.Username = username

	password := os.Getenv("jwxtPassword")
	fmt.Println("输入密码:", password)
	if password == "" {
		fmt.Scanf("%s\n", &password)
	}
	loginForm.Password = password

	fmt.Println("输入验证码:")
	fmt.Scanf("%s\n", &loginForm.Captcha)

	c.Login(loginForm.ConvertToUrlVal())
	c.StoreCookies("./cookie")
	return c
}

// ok
func TestLoadCookie() *jwxt.JwxtClient {
	c := jwxt.NewClient("primaryID")
	c.LoadCookies("./cookie")
	c.CheckLogin()
	return c
}

func TestGetBaseInfo(c *jwxt.JwxtClient) {
	fmt.Println(c.GetCoursePhase())
	fmt.Println(c.GetYearTerm())
	ioutil.WriteFile("favicon.ico", c.GetFavicon(), 0666)
}

func TestCourse(c *jwxt.JwxtClient) {
	// 查课
	courses := c.ListMajOpCourse("")
	for _, v := range courses {
		fmt.Println(v.BaseInfo.CourseName)
	}

	courses = c.ListMajOpCourse("123")
	for _, v := range courses {
		fmt.Println(v.BaseInfo.CourseName)
	}

	courses = c.ListPubOpCourseZH("3211")
	for _, v := range courses {
		fmt.Println(v.BaseInfo.CourseName)
	}
}

func main() {
	// 以下二选一
	// c := TestLogin()
	c := TestLoadCookie()

	TestGetBaseInfo(c)
	TestCourse(c)
}
