package client

import (
	// "net/http"
	"server/backend/jwxtClient/course"
	"server/backend/jwxtClient/global"
	"server/backend/jwxtClient/request"
	"server/backend/jwxtClient/util"
)

type JwxtClient struct {
	*request.HttpClient
	UId  string
	User *user
}

func NewClient(UId string) *JwxtClient {
	c := &JwxtClient{
		HttpClient: request.NewClient(),
		UId:        UId,
		User:       nil,
	}
	return c
}

// 设置学期信息,并返回
func (c *JwxtClient) GetYearTerm() string {
	url := global.HOST + "jwxt/choose-course-front-server/stuCollectedCourse/getYearTerm"
	ref := global.HOST + "jwxt/mk/courseSelection/"
	var resp util.NormalResp
	request.JsonToStruct(request.Get(url).Referer(ref).Do(c).Bytes(), &resp)
	if resp.Code != 200 {
		global.Log.WithField("url", url).Error("无法获取学期信息,使用默认:2020-2")
		global.YEAR_TERM = "2020-2"
	} else {
		global.Log.WithField("semester:", resp.Data).Info("学期信息")
		global.YEAR_TERM = resp.Data
	}
	return global.YEAR_TERM
}

// 获取所有页的课程
func (c *JwxtClient) GetCourseList(courseName, campusId, courseType string) *course.CourseList {
	option := course.NewReqOption(campusId, courseName, false)
	req := course.NewCourseListReq(course.NewCourseType(courseType), option)
	return req.Do(c)
}

// 获取特定页
func (c *JwxtClient) GetCourseListPage(courseName, campusId, courseType string, pageNo int) *course.CourseList {
	option := course.NewReqOption(campusId, courseName, false)
	req := course.NewCourseListReq(course.NewCourseType(courseType), option)
	courseList, _ := req.SetPageNo(pageNo).DoPage(c)
	return courseList
}

// 公选
func (c *JwxtClient) ListPubElecCourse(courseName, campusId string) *course.CourseList {
	return c.GetCourseList(courseName, campusId, course.TYPE_PUB_ELECTIVE)
}

// 专选
func (c *JwxtClient) ListMajElecCourse(courseName, campusId string) *course.CourseList {
	return c.GetCourseList(courseName, campusId, course.TYPE_MAJ_ELECTIVE)
}

// 专必
func (c *JwxtClient) ListMajCompCourse(courseName, campusId string) *course.CourseList {
	return c.GetCourseList(courseName, campusId, course.TYPE_MAJ_COMPULSORY)
}

func (c *JwxtClient) GetCoursePhase() *course.CoursePhase {
	return course.GetCoursePhase(c)
}

func (c *JwxtClient) GetFavicon() []byte {
	url := global.HOST + "jwxt/mk/courseSelection/favicon.ico"
	return c.Get(url).Bytes()
}
