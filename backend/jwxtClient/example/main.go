package main

import (
	"errors"
	"fmt"
	jwxtClient "github.com/liwm29/sysu_jwxt_v3/backend/jwxtClient/client"
	"github.com/liwm29/sysu_jwxt_v3/backend/jwxtClient/course"
	"github.com/liwm29/sysu_jwxt_v3/backend/jwxtClient/global"
	"github.com/liwm29/sysu_jwxt_v3/backend/jwxtClient/util"
	"io/ioutil"
	"os"
	"time"
)

func main() {
	// 设置log级别
	global.SetLogLevel_DEBUG()
	// 构造客户端
	c := jwxtClient.NewClient("")
	// 构造登陆表单, 获取验证码,默认将验证码图片下载到"./captcha.jpg", 登陆cli ,都被集成到了jwxtClient.Login()
	isLogin := c.Login()
	if isLogin {
		fmt.Println("登陆成功")
	} else {
		fmt.Println("登陆失败")
		os.Exit(0)
	}

	fmt.Println("已选课程:", course.GetSelectedCourseNames(c))

	// 获取选课阶段,不在选课阶段时,不能使用
	selectPhase := c.GetCoursePhase()
	fmt.Printf("选课阶段:%s %s 学期:%s 时间:%s=>%s\n",
		selectPhase.ElectiveCourseStageCode, selectPhase.ElectiveCourseStageName, selectPhase.SemesterYear,
		selectPhase.StartTime, selectPhase.EndTime)

	// 获取课程列表,专选
	courseList1 := c.ListCourse(course.TYPE_MAJ_ELECTIVE)
	// equal to
	// courseList1 := c.GetCourseList(course.TYPE_MAJ_ELECTIVE, course.WithCampus(course.CAMPUS_ALL), course.WithCourseName(course.TYPE_MAJ_ELECTIVE))
	fmt.Printf("%#v\n", courseList1.CourseNames())

	// 东校区,公选
	// courseList2 := c.ListPubElecCourse(course.WithCampus(course.CAMPUS_EAST))
	// fmt.Printf("%#v", courseList2.CourseNames()[:5])

	// 单个课程,比如热门课程photoshop
	course, err := c.ListPubElecCourse(course.WithCourseName("photoshop")).First()
	// equal to
	// course, err := c.ListCourse(course.TYPE_PUB_ELECTIVE, course.WithCourseName("photoshop")).First()

	if err != nil {
		fmt.Println("未找到课程信息:", err)
	} else {
		fmt.Println("找到:", course.VacancyInfo())
	}

	// 课程教师信息
	teachers, err := course.GetTeachers(c)
	if err != nil {
		fmt.Println(err)
	} else {
		if len(teachers) > 0 {
			fmt.Printf("课程:%s 教师信息: 姓名:%s Email:%s\n", course.CourseName(), teachers[0].Name, teachers[0].Email)
		} else {
			fmt.Println("无教师信息")
		}
	}

	// 如果在选课第三阶段
	if !selectPhase.CanSelect() {
		util.PanicIf(errors.New("不在选课阶段" + fmt.Sprintf("%#v", selectPhase)))
	}
	if course.VacancyNum() > 0 {
		isOk := course.Choose(c)
		fmt.Println(course.VacancyInfo(), "选课", isOk)
	} else {
		fmt.Println(course.VacancyInfo(), "课程已满")
	}

	// 退课
	if course.IsSelected() {
		isOk := course.Cancel(c)
		fmt.Println(course.VacancyInfo(), "退课", isOk)
	}

	// 课程剩余名额
	fmt.Println("课程容量:", course.VacancyInfo())

	// 刷新课程剩余名额
	if (course.Refresh(c)) != nil {
		fmt.Println("课程刷新失败")
	} else {
		fmt.Println("课程刷新成功: ", course.VacancyInfo())
	}

	// 教师照片
	teacherId := "123456"
	ioutil.WriteFile("teacher"+teacherId+".jpg", c.GetTeacherImg(teacherId), 0666)

	// 学生照片
	studentId := "123456"
	ioutil.WriteFile("student"+studentId+".jpg", c.GetStudentImg(studentId), 0666)

	// 自动选课,5s查询一次,异步
	isOkChan := course.AutoChoose(c, time.Second*5)
	for err := range isOkChan {
		if err == nil {
			fmt.Println(course.VacancyInfo(), "选课成功")
			break
		}
		fmt.Println(err)
	}
}
