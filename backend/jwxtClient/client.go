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

func NewClient(username string) *JwxtClient {
	c := &JwxtClient{
		HttpClient: request.NewClient(),
		username:   username,
		yearTerm:   "2020-2",
	}
	return c
}

func (c *JwxtClient) GetYearTerm() string {
	url := "https://jwxt.sysu.edu.cn/jwxt/choose-course-front-server/stuCollectedCourse/getYearTerm"
	ref := "https://jwxt.sysu.edu.cn/jwxt/mk/courseSelection/"
	var resp util.NormalResp
	request.JsonToStruct(request.Get(url).Referer(ref).Do(c).Bytes(), &resp)
	if resp.Code != 200 {
		log.WithField("url", url).Error("can't get year term")
	} else {
		log.WithField("yearterm", resp.Data).Info("get year term ok")
		c.yearTerm = resp.Data
		return resp.Data
	}
	return ""
}

func (c *JwxtClient) ListCourse(courseName, campusId, courseType string) *course.CourseList {
	option := course.NewReqOption(campusId, courseName, false)
	req := course.NewCourseListReq(c.yearTerm, course.NewCourseType(courseType), option)
	return req.Do(c)
}

func (c *JwxtClient) ListCoursePageN(courseName, campusId, courseType string, n_pages int) *course.CourseList {
	option := course.NewReqOption(campusId, courseName, false)
	req := course.NewCourseListReq(c.yearTerm, course.NewCourseType(courseType), option)
	return req.DoPageN(c, n_pages)
}

// 公选
func (c *JwxtClient) ListPubElecCourse(courseName, campusId string) *course.CourseList {
	return c.ListCourse(courseName, campusId, course.TYPE_PUB_ELECTIVE)
}

// func (c *JwxtClient) ListPubOpCourseEast(courseName string) *course.CourseList {
// 	return c.ListCourse(courseName, course.CAMPUS_EAST, course.TYPE_PUB_ELECTIVE)
// }
// func (c *JwxtClient) ListPubOpCourseSZ(courseName string) *course.CourseList {
// 	return c.ListCourse(courseName, course.CAMPUS_SZ, course.TYPE_PUB_ELECTIVE)
// }
// func (c *JwxtClient) ListPubOpCourseNorth(courseName string) *course.CourseList {
// 	return c.ListCourse(courseName, course.CAMPUS_NORTH, course.TYPE_PUB_ELECTIVE)
// }
// func (c *JwxtClient) ListPubOpCourseSouth(courseName string) *course.CourseList {
// 	return c.ListCourse(courseName, course.CAMPUS_SOUTH, course.TYPE_PUB_ELECTIVE)
// }
// func (c *JwxtClient) ListPubOpCourseZH(courseName string) *course.CourseList {
// 	return c.ListCourse(courseName, course.CAMPUS_ZH, course.TYPE_PUB_ELECTIVE)
// }

// 专选
func (c *JwxtClient) ListMajElecCourse(courseName, campusId string) *course.CourseList {
	return c.ListCourse(courseName, campusId, course.TYPE_MAJ_ELECTIVE)
}

// 专必
func (c *JwxtClient) ListMajCompCourse(courseName, campusId string) *course.CourseList {
	return c.ListCourse(courseName, campusId, course.TYPE_MAJ_COMPULSORY)
}

func (c *JwxtClient) GetCoursePhase() *course.CoursePhase {
	return course.GetCoursePhase(c)
}

func (c *JwxtClient) GetFavicon() []byte {
	url := "https://jwxt.sysu.edu.cn/jwxt/mk/courseSelection/favicon.ico"
	return c.Get(url).Bytes()
}
