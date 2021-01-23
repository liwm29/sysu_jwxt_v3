package main

import (
	"fmt"
	"io/ioutil"
	jwxt "server/backend/jwxtClient"
	"server/backend/jwxtClient/course"
	"server/backend/jwxtClient/util"
	// "testing"
)

// ok
func TestLogin() *jwxt.JwxtClient {
	c := jwxt.NewClient("primaryID")
	// 尝试使用cookie登陆
	if loginOk, _ := c.LoginWithCookies(); loginOk {
		return c
	}

	loginForm := jwxt.NewLoginForm()
	c.CasFirstGet(loginForm)
	loginForm = jwxt.LoginFormCli(loginForm)

	isLogin, err := c.LoginWithForm(loginForm)
	if !isLogin {
		util.PanicIf(err)
	}
	c.StoreCookies(jwxt.DEFAULT_COOKIE_PATH)
	return c
}

// ok
func TestLoadCookie() *jwxt.JwxtClient {
	c := jwxt.NewClient("primaryID")
	isLogin, err := c.LoginWithCookies()
	if !isLogin {
		fmt.Println("login false:", err)
	}
	return c
}

func TestGetBaseInfo(c *jwxt.JwxtClient) {
	fmt.Println(c.GetCoursePhase())
	fmt.Println(c.GetYearTerm())
	ioutil.WriteFile("favicon.ico", c.GetFavicon(), 0666)
}

func TestCourseList(c *jwxt.JwxtClient) {
	// 查课
	courses := c.ListMajElecCourse(course.NAME_ALL, course.CAMPUS_ALL)
	fmt.Println("专选", courses.CourseNames())

	courses = c.ListMajCompCourse(course.NAME_ALL, course.CAMPUS_ALL)
	fmt.Println("专必", courses.CourseNames())

	courses = c.ListPubElecCourse("3211", course.CAMPUS_ALL)
	fmt.Println("某个公选", courses.CourseNames())

	courses = c.GetCourseList(course.NAME_ALL, course.CAMPUS_ALL, course.TYPE_GYM)
	fmt.Println("体育", courses.CourseNames())

	courses = c.GetCourseList(course.NAME_ALL, course.CAMPUS_ALL, course.TYPE_ENGLISH)
	fmt.Println("英语", courses.CourseNames())

	courses = c.GetCourseList(course.NAME_ALL, course.CAMPUS_ALL, course.TYPE_OTHERS)
	fmt.Println("其他", courses.CourseNames())
}

func TestCourseOutline(c *jwxt.JwxtClient) {
	courses := c.GetCourseListPage(course.NAME_ALL, course.CAMPUS_ALL, course.TYPE_PUB_ELECTIVE, 1)
	for _, v := range courses.Courses {
		fmt.Println("某个公选:", v.CourseName(), "剩余容量:", v.VacancyNum())
		teachers, err := v.GetTeachers(c)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%#v\n", teachers)
	}
}

func main() {
	// jwxt.SetLogLevel_INFO()
	c := TestLogin()

	TestGetBaseInfo(c)
	TestCourseList(c)
	TestCourseOutline(c)
}
