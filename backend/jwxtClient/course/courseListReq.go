package course

import (
	"fmt"
	"server/backend/jwxtClient/request"
	"server/backend/jwxtClient/util"
)

type ReqOption struct {
	campusId         string
	courseName       string
	collectionStatus string
}

func NewReqOption(campusId, courseName string, collectionStatus bool) *ReqOption {
	var collected string
	if collectionStatus {
		collected = "1"
	} else {
		collected = "0"
	}
	return &ReqOption{
		campusId:         campusId,
		courseName:       courseName,
		collectionStatus: collected,
	}
}

func (o *ReqOption) GetCampusId() string         { return o.campusId }
func (o *ReqOption) GetCourseName() string       { return o.courseName }
func (o *ReqOption) GetCollectionStatus() string { return o.collectionStatus }

type CourseListReq struct {
	pageNo           int
	pageSize         int
	yearTerm         string
	campusId         string
	collectionStatus string
	courseName       string
	courseType       *CourseType
}

func NewCourseListReq(yearTerm string, courseType *CourseType, option *ReqOption) *CourseListReq {
	return &CourseListReq{
		pageNo:           1,
		pageSize:         10,
		yearTerm:         yearTerm,
		campusId:         option.campusId,
		collectionStatus: option.collectionStatus,
		courseName:       option.courseName,
		courseType:       courseType,
	}
}

func (r *CourseListReq) SetCampusId(campus string) *CourseListReq {
	r.campusId = campus
	return r
}

func (r *CourseListReq) SetCourseName(courseName string) *CourseListReq {
	r.courseName = courseName
	return r
}

func (r *CourseListReq) SetCollection(isJustShowCollected string) *CourseListReq {
	r.collectionStatus = isJustShowCollected
	return r
}

func (r *CourseListReq) IncrePage() *CourseListReq {
	r.pageNo += 1
	return r
}

func (r *CourseListReq) Marshall() string {
	tpl := `{"pageNo":%d,"pageSize":%d,"param":{"semesterYear":"%s","selectedType":"%s","selectedCate":"%s","hiddenConflictStatus":"0","hiddenSelectedStatus":"0","hiddenEmptyStatus":"0","vacancySortStatus":"0","collectionStatus":"%s","campusId":"%s","courseName":"%s"}}`
	return fmt.Sprintf(tpl, r.pageNo, r.pageSize, r.yearTerm, r.courseType.SelectedType, r.courseType.SelectedCate, r.collectionStatus, r.campusId, r.courseName)
}

func (reqJson *CourseListReq) Do(c request.Clienter) *CourseList {
	log.WithField("reqJson", reqJson.Marshall()).Debug(util.WhereAmI())

	url := "https://jwxt.sysu.edu.cn/jwxt/choose-course-front-server/classCourseInfo/course/list"
	ref := "https://jwxt.sysu.edu.cn/jwxt/mk/courseSelection/"
	respJson := request.PostJson(url, reqJson.Marshall()).Referer(ref).Do(c).Bytes()

	log.WithField("respJson", util.Truncate100(string(respJson))).Debug(util.WhereAmI())

	var resp CourseListResp
	request.JsonToStruct(respJson, &resp)
	if resp.Code != 200 {
		log.WithField("respJson.msg", resp.Message).Warn("get list error")
	}
	if resp.Code == 52021136 {
		log.Error("黑名单")
	}
	totalCourses := resp.Data.Rows

	n_course := resp.Data.Total
	total_pages := n_course / 10 // 10 is one page size
	for i := 0; i < total_pages; i++ {
		respJson := request.PostJson(url, reqJson.IncrePage().Marshall()).Referer(ref).Do(c).Bytes()
		var resp CourseListResp
		request.JsonToStruct(respJson, &resp)
		totalCourses = append(totalCourses, resp.Data.Rows...)
	}

	return NewCourseList(reqJson.yearTerm, reqJson.courseType, totalCourses)
}

func (reqJson *CourseListReq) DoPageN(c request.Clienter, n_pages int) *CourseList {
	log.WithField("reqJson", reqJson.Marshall()).Debug(util.WhereAmI())

	url := "https://jwxt.sysu.edu.cn/jwxt/choose-course-front-server/classCourseInfo/course/list"
	ref := "https://jwxt.sysu.edu.cn/jwxt/mk/courseSelection/"
	respJson := request.PostJson(url, reqJson.Marshall()).Referer(ref).Do(c).Bytes()

	log.WithField("respJson", util.Truncate100(string(respJson))).Debug(util.WhereAmI())

	var resp CourseListResp
	request.JsonToStruct(respJson, &resp)
	if resp.Code != 200 {
		log.WithField("respJson.msg", resp.Message).Warn("get list error")
	}
	if resp.Code == 52021136 {
		log.Error("黑名单")
	}
	totalCourses := resp.Data.Rows

	n_course := resp.Data.Total
	var total_pages = n_course / 10 // 10 is one page size
	n_loop := util.Min(total_pages-1, n_pages-1)

	for i := 0; i < n_loop; i++ {
		respJson := request.PostJson(url, reqJson.IncrePage().Marshall()).Referer(ref).Do(c).Bytes()
		var resp CourseListResp
		request.JsonToStruct(respJson, &resp)
		totalCourses = append(totalCourses, resp.Data.Rows...)
	}

	return NewCourseList(reqJson.yearTerm, reqJson.courseType, totalCourses)
}
