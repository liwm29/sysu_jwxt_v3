package jwxtClient

import (
	// "net/http"
	"server/backend/jwxtClient/request"
)

type jwxtClient struct {
	*request.HttpClient
	username string
	isLogin  bool
	yearTerm string
}

func NewClient(username string) *jwxtClient {
	c := &jwxtClient{
		HttpClient: request.NewClient(),
		username:   username,
		isLogin:    false,
		yearTerm:   "2020-2",
	}
	return c
}

func (c *jwxtClient) GetYearTerm() string {
	url := "https://jwxt.sysu.edu.cn/jwxt/choose-course-front-server/stuCollectedCourse/getYearTerm"
	var resp normalResp
	request.JsonToStruct(c.Get(url), &resp)
	if resp.Code != 200 {
		log.WithField("url", url).Error("can't get year term")
	} else {
		log.WithField("yearterm", resp.Data).Info("get year term ok")
		c.yearTerm = resp.Data
		return resp.Data
	}
	return ""
}

func (c *jwxtClient) ListCourseWithReq(reqJson *courseListReq) []Row {
	url := "https://jwxt.sysu.edu.cn/jwxt/choose-course-front-server/schoolCourse/pageList"
	ref := "https://jwxt.sysu.edu.cn/jwxt/mk/courseSelection/"
	respJson := request.PostJson(url, reqJson.marshall()).Referer(ref).Do(c)
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
		respJson := request.PostJson(url, reqJson.setNextPage().marshall()).Referer(ref).Do(c)
		var resp courseListResp
		request.JsonToStruct(respJson, &resp)
		totalCourses = append(totalCourses, resp.Data.Rows...)
	}
	return totalCourses
}

const (
	CTYPE_PUBOP int32 = iota // COURSE TYPE PUBLIC OPTIONAL 公选
	CTYPE_MAJOP
)

func (c *jwxtClient) ListCourse(courseType int32) []Row {
	var reqJson courseListReq
	switch courseType {
	case CTYPE_PUBOP:
		reqJson = courseListReq{
			pageNo:           1,
			pageSize:         10,
			yearTerm:         c.yearTerm,
			selectedType:     _selectedType["校级公选"],
			selectedCate:     "", // 公选时,不需要
			campusId:         _campus_id["东校区"],
			collectionStatus: "0",
		}
	case CTYPE_MAJOP:
		reqJson = courseListReq{
			pageNo:           1,
			pageSize:         10,
			yearTerm:         c.yearTerm,
			selectedType:     _selectedType["本专业"],
			selectedCate:     _selectedCate["专选"],
			campusId:         _campus_id["东校区"],
			collectionStatus: "0",
		}
	}
	courses := c.ListCourseWithReq(&reqJson)
	return courses
}

func (c *jwxtClient) ListPubOpCourse() []Row {
	return c.ListCourse(CTYPE_PUBOP)
}

func (c *jwxtClient) ListMajOpCourse() []Row {
	return c.ListCourse(CTYPE_MAJOP)
}

func (c *jwxtClient) GetCoursePhase() coursePhase {
	url := "https://jwxt.sysu.edu.cn/jwxt/choose-course-front-server/classCourseInfo/selectCourseInfo"
	respJson := c.Get(url)
	var resp coursePhaseResp
	request.JsonToStruct(respJson, &resp)
	if resp.Code != 200 {
		log.WithField("msg", resp.Message).Warn("无法获取选课阶段")
	}
	return resp.Data
	// {"code":200,"message":null,"data":{"electiveCourseStageName":"改补选","retreatCourseStatus":"1","code":200,"semesterYear":"2020-2","courseSelectType":"0","chooseCourseStatus":"1","electiveCourseStageCode":"3","startTime":"2021-01-08 13:00:00","endTime":"2021-03-04 23:00:00","crossMajor":"1"}}
}

func (c *jwxtClient) GetFavicon() []byte {
	url := "https://jwxt.sysu.edu.cn/jwxt/mk/courseSelection/favicon.ico"
	return c.Get(url)
}
