package jwxtClient

import (
	// "net/http"
	"server/backend/jwxtClient/request"
)

type JwxtClient struct {
	*request.HttpClient
	username string
	isLogin  bool
	yearTerm string
}

func NewClient(username string) *JwxtClient {
	c := &JwxtClient{
		HttpClient: request.NewClient(),
		username:   username,
		isLogin:    false,
		yearTerm:   "2020-2",
	}
	return c
}

func (c *JwxtClient) GetYearTerm() string {
	url := "https://jwxt.sysu.edu.cn/jwxt/choose-course-front-server/stuCollectedCourse/getYearTerm"
	ref := "https://jwxt.sysu.edu.cn/jwxt/mk/courseSelection/"
	var resp normalResp
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

func (c *JwxtClient) ListCourseWithReq(reqJson *CourseListReq) []courseInfo {
	log.WithField("reqJson", reqJson.marshall()).Debug("ListCourseWithReq")

	url := "https://jwxt.sysu.edu.cn/jwxt/choose-course-front-server/classCourseInfo/course/list"
	ref := "https://jwxt.sysu.edu.cn/jwxt/mk/courseSelection/"
	respJson := request.PostJson(url, reqJson.marshall()).Referer(ref).Do(c).Bytes()

	log.WithField("respJson", string(respJson)).Debug("ListCourseWithReq")

	var resp courseListResp
	request.JsonToStruct(respJson, &resp)
	if resp.Code != 200 {
		log.WithField("msg", resp.Message).Warn("get list error")
	}
	if resp.Code == 52021136 {
		log.Error("黑名单")
	}
	totalCourses := resp.Data.Rows

	n_course := resp.Data.Total
	var times = n_course / 10 // 10 is one page size
	for i := 0; i < times; i++ {
		respJson := request.PostJson(url, reqJson.setNextPage().marshall()).Referer(ref).Do(c).Bytes()
		var resp courseListResp
		request.JsonToStruct(respJson, &resp)
		totalCourses = append(totalCourses, resp.Data.Rows...)
	}
	return totalCourses
}

func (c *JwxtClient) ListCourse(courseType, courseName, campusId string) []*Course {
	req := NewCourseListReq(c.yearTerm, courseType)
	req.SetCourseName(courseName)
	req.SetCampusId(campusId)
	courses := c.ListCourseWithReq(req)
	return newCourses(courseType, c.yearTerm, courses)
}

func (c *JwxtClient) ListPubOpCourse(courseName string) []*Course {
	return c.ListCourse("公选", courseName, "")
}

func (c *JwxtClient) ListPubOpCourseEast(courseName string) []*Course {
	return c.ListCourse("公选", courseName, _campus_id["东校园"])
}
func (c *JwxtClient) ListPubOpCourseSZ(courseName string) []*Course {
	return c.ListCourse("公选", courseName, _campus_id["深圳校区"])
}
func (c *JwxtClient) ListPubOpCourseNorth(courseName string) []*Course {
	return c.ListCourse("公选", courseName, _campus_id["北校园"])
}
func (c *JwxtClient) ListPubOpCourseSouth(courseName string) []*Course {
	return c.ListCourse("公选", courseName, _campus_id["南校园"])
}
func (c *JwxtClient) ListPubOpCourseZH(courseName string) []*Course {
	return c.ListCourse("公选", courseName, _campus_id["珠海校区"])
}

func (c *JwxtClient) ListMajOpCourse(courseName string) []*Course {
	return c.ListCourse("专选", courseName, "")
}

func (c *JwxtClient) GetCoursePhase() coursePhase {
	url := "https://jwxt.sysu.edu.cn/jwxt/choose-course-front-server/classCourseInfo/selectCourseInfo"
	ref := "https://jwxt.sysu.edu.cn/jwxt/mk/courseSelection/"
	respJson := request.Get(url).Referer(ref).Do(c).Bytes()
	var resp coursePhaseResp
	request.JsonToStruct(respJson, &resp)
	if resp.Code != 200 {
		log.WithField("msg", resp.Message).Error("无法获取选课阶段")
	}
	return resp.Data
	// {"code":200,"message":null,"data":{"electiveCourseStageName":"改补选","retreatCourseStatus":"1","code":200,"semesterYear":"2020-2","courseSelectType":"0","chooseCourseStatus":"1","electiveCourseStageCode":"3","startTime":"2021-01-08 13:00:00","endTime":"2021-03-04 23:00:00","crossMajor":"1"}}
}

func (c *JwxtClient) GetFavicon() []byte {
	url := "https://jwxt.sysu.edu.cn/jwxt/mk/courseSelection/favicon.ico"
	return c.Get(url).Bytes()
}
