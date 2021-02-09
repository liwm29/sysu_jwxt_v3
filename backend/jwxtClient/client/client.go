package client

import (
	"github.com/liwm29/sysu_jwxt_v3/backend/jwxtClient/course"
	"github.com/liwm29/sysu_jwxt_v3/backend/jwxtClient/global"
	"github.com/liwm29/sysu_jwxt_v3/backend/jwxtClient/request"
	"github.com/liwm29/sysu_jwxt_v3/backend/jwxtClient/util"
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
		global.Log.WithField("url", url).Error("无法获取学期信息,使用默认:NULL")
		global.YEAR_TERM = "NULL"
	} else {
		global.Log.WithField("semester:", resp.Data).Info("学期信息")
		global.YEAR_TERM = resp.Data
	}
	return global.YEAR_TERM
}

// 获取所有页的课程
func (c *JwxtClient) ListCourse(courseType string, opts ...course.ReqOptionSetter) *course.CourseList {
	req := course.NewCourseListReq(course.NewCourseType(courseType), opts...)
	return req.Do(c)
}

// 获取特定页
func (c *JwxtClient) ListCourseWithPage(courseType string, pageNo int, opts ...course.ReqOptionSetter) *course.CourseList {
	req := course.NewCourseListReq(course.NewCourseType(courseType), opts...)
	courseList, _ := req.SetPageNo(pageNo).DoPage(c)
	return courseList
}

// 公选
func (c *JwxtClient) ListPubElecCourse(opts ...course.ReqOptionSetter) *course.CourseList {
	return c.ListCourse(course.TYPE_PUB_ELECTIVE, opts...)
}

// 专选
func (c *JwxtClient) ListMajElecCourse(opts ...course.ReqOptionSetter) *course.CourseList {
	return c.ListCourse(course.TYPE_MAJ_ELECTIVE, opts...)
}

// 专必
func (c *JwxtClient) ListMajCompCourse(opts ...course.ReqOptionSetter) *course.CourseList {
	return c.ListCourse(course.TYPE_MAJ_COMPULSORY, opts...)
}

func (c *JwxtClient) GetCoursePhase() *course.CoursePhase {
	return course.GetCoursePhase(c)
}

func (c *JwxtClient) GetFavicon() []byte {
	url := global.HOST + "jwxt/mk/courseSelection/favicon.ico"
	return c.Get(url).Bytes()
}
