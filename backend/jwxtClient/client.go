package jwxtClient

import (
	// "net/http"
	"server/backend/jwxtClient/course"
	"server/backend/jwxtClient/request"
	"server/backend/jwxtClient/util"
)

type JwxtClient struct {
	*request.HttpClient
	username string
	yearTerm string
}

func NewClient(UID string) *JwxtClient {
	c := &JwxtClient{
		HttpClient: request.NewClient(),
		username:   UID,
		yearTerm:   "",
	}
	return c
}

// 设置学期信息,并返回
func (c *JwxtClient) GetYearTerm() string {
	url := "https://jwxt.sysu.edu.cn/jwxt/choose-course-front-server/stuCollectedCourse/getYearTerm"
	ref := "https://jwxt.sysu.edu.cn/jwxt/mk/courseSelection/"
	var resp util.NormalResp
	request.JsonToStruct(request.Get(url).Referer(ref).Do(c).Bytes(), &resp)
	if resp.Code != 200 {
		log.WithField("url", url).Error("无法获取学期信息,使用默认:2020-2")
		c.yearTerm = "2020-2"
	} else {
		log.WithField("semester:", resp.Data).Info("学期信息")
		c.yearTerm = resp.Data
	}
	return c.yearTerm
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
	courseList, _ := req.SetPage(pageNo).DoPage(c)
	return courseList
}

// 公选
func (c *JwxtClient) ListPubElecCourse(courseName, campusId string) *course.CourseList {
	return c.GetCourseList(courseName, campusId, course.TYPE_PUB_ELECTIVE)
}

func (c *JwxtClient) ListPubElecCourseEast(courseName string) *course.CourseList {
	return c.GetCourseList(courseName, course.CAMPUS_EAST, course.TYPE_PUB_ELECTIVE)
}
func (c *JwxtClient) ListPubElecCourseSZ(courseName string) *course.CourseList {
	return c.GetCourseList(courseName, course.CAMPUS_SZ, course.TYPE_PUB_ELECTIVE)
}
func (c *JwxtClient) ListPubElecCourseNorth(courseName string) *course.CourseList {
	return c.GetCourseList(courseName, course.CAMPUS_NORTH, course.TYPE_PUB_ELECTIVE)
}
func (c *JwxtClient) ListPubElecCourseSouth(courseName string) *course.CourseList {
	return c.GetCourseList(courseName, course.CAMPUS_SOUTH, course.TYPE_PUB_ELECTIVE)
}
func (c *JwxtClient) ListPubElecCourseZH(courseName string) *course.CourseList {
	return c.GetCourseList(courseName, course.CAMPUS_ZH, course.TYPE_PUB_ELECTIVE)
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
	url := "https://jwxt.sysu.edu.cn/jwxt/mk/courseSelection/favicon.ico"
	return c.Get(url).Bytes()
}
