package course

import (
	"fmt"
	"server/backend/jwxtClient/request"
	"server/backend/jwxtClient/util"
)

// 之所以这样设计,一方面是往New函数里面传option是一种风格,另一方面,也是为了可拓展性,因为请求还可以加入其他字段
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
	pageNo     int
	pageSize   int
	yearTerm   string
	option     *ReqOption
	courseType *CourseType
}

func NewCourseListReq(courseType *CourseType, option *ReqOption) *CourseListReq {
	if option == nil {
		option = new(ReqOption)
	}
	if YEAR_TERM == "" {
		log.WithField("semester:", YEAR_TERM).Info("未设置学期 ", util.WhereAmI())
		return nil
	}
	return &CourseListReq{
		pageNo:     1,
		pageSize:   10,
		yearTerm:   YEAR_TERM,
		option:     option,
		courseType: courseType,
	}
}

func (r *CourseListReq) SetCampusId(campus string) *CourseListReq {
	r.option.campusId = campus
	return r
}

func (r *CourseListReq) SetCourseName(courseName string) *CourseListReq {
	r.option.courseName = courseName
	return r
}

func (r *CourseListReq) SetCollection(isJustShowCollected string) *CourseListReq {
	r.option.collectionStatus = isJustShowCollected
	return r
}

func (r *CourseListReq) IncrePage() *CourseListReq {
	r.pageNo += 1
	return r
}

func (r *CourseListReq) SetPage(pageNo int) *CourseListReq {
	r.pageNo = pageNo
	return r
}

func (r *CourseListReq) Marshall() string {
	tpl := `{"pageNo":%d,"pageSize":%d,"param":{"semesterYear":"%s","selectedType":"%s","selectedCate":"%s","hiddenConflictStatus":"0","hiddenSelectedStatus":"0","hiddenEmptyStatus":"0","vacancySortStatus":"0","collectionStatus":"%s","campusId":"%s","courseName":"%s"}}`
	return fmt.Sprintf(tpl, r.pageNo, r.pageSize, r.yearTerm, r.courseType.SelectedType, r.courseType.SelectedCate, r.option.collectionStatus, r.option.campusId, r.option.courseName)
}

// 返回所有课程列表,从第一页开始
func (reqJson *CourseListReq) Do(c request.Clienter) *CourseList {
	log.WithField("reqJson", reqJson.Marshall()).Debug(util.WhereAmI())
	courses, n_page := reqJson.SetPage(1).DoPage(c)

	for i := 0; i < n_page; i++ {
		courseListTmp, _ := reqJson.SetPage(i + 2).DoPage(c)
		courses.Courses = append(courses.Courses, courseListTmp.Courses...)
	}

	return courses
}

// 返回一页课程列表
func (reqJson *CourseListReq) DoPage(c request.Clienter) (*CourseList, int) {
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

	n_page := (resp.Data.Total + reqJson.pageSize - 1) / reqJson.pageSize
	courseList := NewCourseList(reqJson.courseType, resp.Data.Rows, reqJson.option)
	return courseList, n_page
}
