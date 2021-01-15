package main

import (
	"fmt"
	jwxt "server/backend/jwxtClient"
	// "testing"
)

// ok
func TestLogin() {
	c := jwxt.NewClient("primaryID")
	loginForm := jwxt.NewLoginForm()
	c.CasFirstGet("./captcha.jpg", loginForm)
	fmt.Println("输入用户名")
	fmt.Scanf("%s\n", &loginForm.Username)
	fmt.Println("输入密码")
	fmt.Scanf("%s\n", &loginForm.Password)
	fmt.Println("输入验证码:")
	fmt.Scanf("%s\n", &loginForm.Captcha)
	c.Login(loginForm.ConvertToUrlVal())
	// show cookies
	for k, v := range c.CookieJar.DB {
		fmt.Println(k, " : ", v)
	}
	c.StoreCookies("./cookie")
}

// ok
func TestLoadCookie() {
	c := jwxt.NewClient("liwm29")
	c.LoadCookies("./cookie")
	c.CheckLogin()
}

func main() {
	// TestLogin()
	TestLoadCookie()
}
