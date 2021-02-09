package course

import (
	"fmt"
	"github.com/liwm29/sysu_jwxt_v3/backend/jwxtClient/global"
	"github.com/liwm29/sysu_jwxt_v3/backend/jwxtClient/request"
	"github.com/liwm29/sysu_jwxt_v3/backend/jwxtClient/util"
)

// 之所以这样设计,一方面是往New函数里面传option是一种风格,另一方面,也是为了可拓展性,因为请求还可以加入其他字段
type ReqOptions struct {
	campusId         string
	courseName       string
	collectionStatus string
}

func defaultReqOptions() ReqOptions {
	return ReqOptions{
		campusId:         "",
		courseName:       "",
		collectionStatus: "",
	}
}

type reqOptionSetFunc func(*ReqOptions) ReqOptionSetter

// for further expansion ,use struct to wrap
type ReqOptionSetter struct {
	// value interface{}
	f reqOptionSetFunc
}

func (r *ReqOptionSetter) apply(ropts *ReqOptions) {
	r.f(ropts)
}

func WithCampus(campusId string) ReqOptionSetter {
	return ReqOptionSetter{
		func(ro *ReqOptions) ReqOptionSetter {
			prev := ro.campusId
			ro.campusId = campusId
			return WithCampus(prev)
		},
	}
}

func WithCourseName(courseName string) ReqOptionSetter {
	return ReqOptionSetter{
		func(ro *ReqOptions) ReqOptionSetter {
			prev := ro.courseName
			ro.courseName = courseName
			return WithCourseName(prev)
		},
	}
}

func WithShowCollected(isOnlyShowCollected bool) ReqOptionSetter {
	return ReqOptionSetter{
		func(ro *ReqOptions) ReqOptionSetter {
			prev := ro.collectionStatus
			ro.collectionStatus = util.Bool2Str(isOnlyShowCollected)
			return WithShowCollected(util.Str2Bool(prev))
		},
	}
}

func (o *ReqOptions) GetCampusId() string         { return o.campusId }
func (o *ReqOptions) GetCourseName() string       { return o.courseName }
func (o *ReqOptions) GetCollectionStatus() string { return o.collectionStatus }

type CourseListReq struct {
	pageNo     int
	pageSize   int
	options    ReqOptions
	courseType *CourseType
}

func NewCourseListReq(courseType *CourseType, opts ...ReqOptionSetter) *CourseListReq {
	if global.YEAR_TERM == "" {
		global.Log.WithField("semester:", global.YEAR_TERM).Info("未设置学期 ", util.WhereAmI())
		return nil
	}

	req := &CourseListReq{
		pageNo:     1,
		pageSize:   10,
		options:    defaultReqOptions(),
		courseType: courseType,
	}
	for _, o := range opts {
		o.apply(&req.options)
	}

	return req
}

// func (r *CourseListReq) SetCampusId(campus string) *CourseListReq {
// 	r.options.campusId = campus
// 	return r
// }

// func (r *CourseListReq) SetCourseName(courseName string) *CourseListReq {
// 	r.options.courseName = courseName
// 	return r
// }

// func (r *CourseListReq) SetCollection(isJustShowCollected string) *CourseListReq {
// 	r.options.collectionStatus = isJustShowCollected
// 	return r
// }

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
	return fmt.Sprintf(tpl, r.pageNo, r.pageSize, global.YEAR_TERM, r.courseType.SelectedType, r.courseType.SelectedCate, r.options.collectionStatus, r.options.campusId, r.options.courseName)
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
	courseList = NewCourseList(reqJson.courseType, resp.Data.Rows, reqJson.options)
	return courseList, n_page
}
