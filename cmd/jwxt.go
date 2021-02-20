package main

import (
	"errors"
	"fmt"
	"github.com/liwm29/sysu_jwxt_v3/pkg/jwxtClient/client"
	"github.com/liwm29/sysu_jwxt_v3/pkg/jwxtClient/course"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func newClientAndCheck() (*client.JwxtClient, error) {
	cl := client.NewClient("1")
	// 尝试使用cookie登陆
	if LoginOk, _ := cl.LoginWithCookies(); LoginOk {
		return cl, nil
	}
	return nil, errors.New("not logined")
}

func main() {
	app := cli.App{Commands: []*cli.Command{
		{
			Name:  "login",
			Usage: "sign in the sysu jwxt, get the cookie",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "username",
					Aliases:  []string{"u"},
					Required: true,
				},
				&cli.StringFlag{
					Name:     "password",
					Aliases:  []string{"p"},
					Required: true,
				},
			},
			Action: func(c *cli.Context) error {
				cl := client.NewClient("1")
				// 尝试使用cookie登陆
				if LoginOk, _ := cl.LoginWithCookies(); LoginOk {
					return nil
				}
				fmt.Println("使用cookie登陆失败,尝试使用form表单登陆")

				LoginForm := client.NewLoginForm()
				cl.CasFirstGet(LoginForm)
				LoginForm.Username = c.String("username")
				LoginForm.Password = c.String("password")
				var captcha string
				fmt.Println("输入验证码(已下载到./captcha/jpg):")
				fmt.Scanf("%s", &captcha)
				LoginForm.Captcha = captcha
				isLogin := cl.LoginWithForm(LoginForm)
				if !isLogin {
					return errors.New("登陆失败,检查log")
				}
				return nil
			},
		},
		{
			Name:  "list",
			Usage: "list course",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "type",
					Aliases:  []string{"t"},
					Usage:    "课程类型:公选 专选 专必 英语 体育 其他",
					Required: true,
				},
				&cli.StringFlag{
					Name:    "campus",
					Aliases: []string{"c"},
					Usage:   "校园代码:e n s sz zh (east/north/south/shenzhen/zhuhai)",
					Value:   "",
				},
				&cli.StringFlag{
					Name:    "name",
					Aliases: []string{"n"},
					Usage:   "课程名,如:高等数学",
					Value:   "",
				},
				&cli.IntFlag{
					Name:    "page",
					Aliases: []string{"p"},
					Usage:   "显示多少课程",
					Value:   -1,
				},
			},
			Action: func(c *cli.Context) error {
				cl, err := newClientAndCheck()
				if err != nil {
					return errors.New("cmd:list |" + err.Error())
				}

				// todo: optimize: use ListCourseWithPage instead
				list := cl.ListCourse(c.String("type"),
					course.WithCampus(c.String("campus")),
					course.WithCourseName(c.String("name")))

				page := c.Int("page")
				if page == -1 {
					page = len(list.Courses)
				}
				fmt.Println(list.CourseNames()[:page])
				return nil
			},
		},
		{
			Name:  "course",
			Usage: "courseTeacherInfo",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "name",
					Aliases:  []string{"n"},
					Usage:    "课程名,如:高等数学",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "type",
					Aliases:  []string{"t"},
					Usage:    "课程类型:公选 专选 专必 英语 体育 其他",
					Required: true,
				},
			},
			Action: func(c *cli.Context) error {
				cl, err := newClientAndCheck()
				if err != nil {
					return fmt.Errorf("cmd:cancel |%s", err.Error())
				}

				list := cl.ListCourse(c.String("type"), course.WithCourseName(c.String("name")))
				course, err := list.First()
				if err != nil {
					return errors.New("未找到课程:" + c.String("name") + " | " + err.Error())
				}
				fmt.Println("找到课程:" + course.CourseName())

				tInfos, err := course.GetTeachers(cl)
				if err != nil {
					return err
				}
				for _, v := range tInfos {
					fmt.Printf("%#v", v)
				}
				return nil
			},
		},
		{
			Name:  "choose",
			Usage: "choose course",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "name",
					Aliases:  []string{"n"},
					Usage:    "课程名,如:高等数学",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "type",
					Aliases:  []string{"t"},
					Usage:    "课程类型:公选 专选 专必 英语 体育 其他",
					Required: true,
				},
				&cli.BoolFlag{
					Name:    "autoChoose",
					Usage:   "是否自动选课",
					Aliases: []string{"auto"},
					Value:   false,
				},
				&cli.IntFlag{
					Name:    "interval",
					Usage:   "选课间隔时长:秒",
					Aliases: []string{"i"},
					Value:   -1,
				},
			},
			Action: func(c *cli.Context) error {
				cl, err := newClientAndCheck()
				if err != nil {
					return errors.New("cmd:select |" + err.Error())
				}
				list := cl.ListCourse(c.String("type"), course.WithCourseName(c.String("name")))
				course, err := list.First()
				if err != nil {
					return errors.New("未找到课程:" + c.String("name") + " | " + err.Error())
				}
				fmt.Println("找到课程:" + course.CourseName())

				fmt.Println("检查课程是否已选...")
				if course.IsSelected() {
					fmt.Println("课程已选")
					return nil
				} else {
					fmt.Println("课程未选")
				}

				fmt.Println("检查课程余量...")
				if course.VacancyNum() == 0 {
					fmt.Println("课程已满", course.VacancyInfo())
				} else {
					fmt.Println("课程可选", course.VacancyInfo())
				}

				if !c.Bool("autochoose") {
					fmt.Println("选课中...")
					if course.Choose(cl) {
						fmt.Println("选课成功")
					} else {
						fmt.Println("选课失败")
					}
					fmt.Println("课程余量", course.VacancyInfo())
					return nil
				} else {
					if c.Int("interval") == -1 {
						return fmt.Errorf("not set autoChoose sleep interval")
					}
					ch := course.AutoChoose(cl, time.Second*time.Duration(c.Int("interval")))
					for msg := range ch {
						fmt.Println(msg)
					}
				}
				return nil
			},
		},
		{
			Name:  "cancel",
			Usage: "cancel course",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "name",
					Aliases:  []string{"n"},
					Usage:    "课程名,如:高等数学",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "type",
					Aliases:  []string{"t"},
					Usage:    "课程类型:公选 专选 专必 英语 体育 其他",
					Required: true,
				},
			},
			Action: func(c *cli.Context) error {
				cl, err := newClientAndCheck()
				if err != nil {
					return fmt.Errorf("cmd:cancel |%s", err.Error())
				}

				list := cl.ListCourse(c.String("type"), course.WithCourseName(c.String("name")))
				course, err := list.First()
				if err != nil {
					return errors.New("未找到课程:" + c.String("name") + " | " + err.Error())
				}
				fmt.Println("找到课程:" + course.CourseName())

				fmt.Println("检查课程是否已选...")
				if course.IsSelected() {
					fmt.Println("课程已选")
					return nil
				} else {
					fmt.Println("课程未选")
				}

				fmt.Println("退课中...")
				if course.Cancel(cl) {
					fmt.Println("退课成功")
				} else {
					fmt.Println("退课失败")
				}
				fmt.Println("课程余量", course.VacancyInfo())
				return nil
			},
		},
		{
			Name:  "img",
			Usage: "teacher's/student's image",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "teacher",
					Aliases: []string{"t"},
				},
				&cli.BoolFlag{
					Name:    "student",
					Aliases: []string{"s"},
				},
				&cli.StringFlag{
					Name:    "id",
					Aliases: []string{"i"},
				},
			},
			Action: func(c *cli.Context) error {
				cl, err := newClientAndCheck()
				if err != nil {
					return fmt.Errorf("cmd:img |%s", err.Error())
				}
				id := c.String("id")
				if c.Bool("teacher") {
					buf := cl.GetTeacherImg(id)
					file := "teacher" + id + ".jpg"
					ioutil.WriteFile(file, buf, 0666)
					fmt.Printf("图片已下载到./%s\n", file)
				}
				if c.Bool("student") {
					buf := cl.GetStudentImg(id)
					file := "student" + id + ".jpg"
					ioutil.WriteFile(file, buf, 0666)
					fmt.Printf("图片已下载到./%s\n", file)
				}
				return nil
			},
		},
	}}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
