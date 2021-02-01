package course

import (
	"fmt"
	"github.com/liwm29/sysu_jwxt_v3/backend/jwxtClient/global"
	"github.com/liwm29/sysu_jwxt_v3/backend/jwxtClient/request"
	"github.com/liwm29/sysu_jwxt_v3/backend/jwxtClient/util"
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
	option     *ReqOption
	courseType *CourseType
}

func NewCourseListReq(courseType *CourseType, option *ReqOption) *CourseListReq {
	if option == nil {
		option = new(ReqOption)
	}
	if global.YEAR_TERM == "" {
		global.Log.WithField("semester:", global.YEAR_TERM).Info("未设置学期 ", util.WhereAmI())
		return nil
	}
	return &CourseListReq{
		pageNo:     1,
		pageSize:   10,
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

func (r *CourseListReq) IncrePageNo() *CourseListReq {
	r.pageNo += 1
	return r
}

func (r *CourseListReq) SetPageNo(pageNo int) *CourseListReq {
	r.pageNo = pageNo
	return r
}

func (r *CourseListReq) Marshall() string {
	tpl := `{"pageNo":%d,"pageSize":%d,"param":{"semesterYear":"%s","selectedType":"%s","selectedCate":"%s","hiddenConflictStatus":"0","hiddenSelectedStatus":"0","hiddenEmptyStatus":"0","vacancySortStatus":"0","collectionStatus":"%s","campusId":"%s","courseName":"%s"}}`
	return fmt.Sprintf(tpl, r.pageNo, r.pageSize, global.YEAR_TERM, r.courseType.SelectedType, r.courseType.SelectedCate, r.option.collectionStatus, r.option.campusId, r.option.courseName)
}

// 返回所有课程列表,从第一页开始
func (reqJson *CourseListReq) Do(c request.Clienter) *CourseList {
	global.Log.WithField("reqJson", reqJson.Marshall()).Debug(util.WhereAmI())
	courses, n_page := reqJson.SetPageNo(1).DoPage(c)

	for i := 2; i <= n_page; i++ {
		courseListTmp, _ := reqJson.SetPageNo(i).DoPage(c)
		courses.Courses = append(courses.Courses, courseListTmp.Courses...)
	}

	return courses
}

// 返回一页课程列表,设置courseListReq.SetPage()
// @return 返回本页内所有课程列表,和所有页的课程数
func (reqJson *CourseListReq) DoPage(c request.Clienter) (courseList *CourseList, n_page int) {
	global.Log.WithField("reqJson", reqJson.Marshall()).Debug(util.WhereAmI())

	url := global.HOST + "jwxt/choose-course-front-server/classCourseInfo/course/list"
	ref := global.HOST + "jwxt/mk/courseSelection/"
	respJson := request.PostJson(url, reqJson.Marshall()).Referer(ref).Do(c).Bytes()

	global.Log.WithField("respJson", string(respJson)).Debug(util.WhereAmI())
	var resp CourseListResp
	request.JsonToStruct(respJson, &resp)
	if resp.Code != 200 {
		global.Log.WithField("respJson.msg", resp.Message).Warn("get list error")
	}
	if resp.Code == 52021136 {
		global.Log.Error("黑名单")
	}

	n_page = (resp.Data.Total + reqJson.pageSize - 1) / reqJson.pageSize
	courseList = NewCourseList(reqJson.courseType, resp.Data.Rows, reqJson.option)
	return courseList, n_page
}
