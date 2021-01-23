package main

import (
	"fmt"
	"io/ioutil"
	jwxt "server/backend/jwxtClient"
	"server/backend/jwxtClient/course"
	"server/backend/jwxtClient/util"
	"time"
)

func main() {
	// 设置log级别
	jwxt.SetLogLevel_INFO()
	// 构造客户端
	c := jwxt.NewClient("")
	// 构造登陆表单
	form := jwxt.NewLoginForm()
	// 获取验证码
	c.CasFirstGet(form)
	// 登陆cli
	_, err := c.Login()
	util.PanicIf(err)

	// 获取选课阶段,不在选课阶段时,不能使用
	selectPhase := c.GetCoursePhase()
	fmt.Printf("%#v", selectPhase)

	// 获取课程列表,专选
	courseList1 := c.GetCourseList(course.NAME_ALL, course.CAMPUS_ALL, course.TYPE_MAJ_ELECTIVE)
	fmt.Printf("%#v", courseList1.CourseNames())
	// 东校区,公选
	courseList2 := c.GetCourseList(course.NAME_ALL, course.CAMPUS_EAST, course.TYPE_PUB_ELECTIVE)
	fmt.Printf("%#v", courseList2.CourseNames())

	// 单个课程,比如热门课程photoshop
	course, err := c.GetCourseList("photoshop", course.CAMPUS_ALL, course.TYPE_PUB_ELECTIVE).First()
	if err != nil {
		fmt.Println("找到: ", course.VacancyInfo())
	} else {
		fmt.Println("未找到课程信息:", err)
	}

	// 课程教师信息
	teachers, err := course.GetTeachers(c)
	if err != nil {
		fmt.Println(err)
	} else {
		if len(teachers) > 0 {
			fmt.Printf("%#v", teachers)
		} else {
			fmt.Println("无教师信息")
		}
	}

	// 如果在选课第三阶段
	if selectPhase.CanSelect() {
		if course.VacancyNum() > 0 {
			isOk := course.Choose(c)
			fmt.Println(course.VacancyInfo(), "选课", isOk)
		} else {
			fmt.Println(course.VacancyInfo(), "课程已满")
		}
	}

	// 退课
	isOk := course.Cancel(c)
	fmt.Println(course.VacancyInfo(), "退课", isOk)

	// 课程剩余名额
	fmt.Println(course.VacancyInfo())

	// 刷新课程剩余名额
	if (course.Refresh(c)) != nil {
		fmt.Println("课程刷新失败")
	} else {
		fmt.Println("课程刷新成功: ", course.VacancyInfo())
	}

	// 教师照片
	teacherId := "123456"
	ioutil.WriteFile(teacherId+".jpg", c.GetTeacherImg(teacherId), 0666)

	// 学生照片
	studentId := "123456"
	ioutil.WriteFile(studentId+".jpg", c.GetTeacherImg(studentId), 0666)

	// 自动选课,5s查询一次,异步
	isOkChan := course.AutoChoose(c, time.Second*5)
	for err := range isOkChan {
		if err == nil {
			fmt.Println(course.VacancyInfo(), "选课成功")
		}
		fmt.Println(err)
	}
}
