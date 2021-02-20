package course

import (
	"errors"
	"fmt"
	"strings"

	"github.com/liwm29/sysu_jwxt_v3/pkg/jwxtClient/global"
	"github.com/liwm29/sysu_jwxt_v3/pkg/jwxtClient/internal/util"
	"github.com/liwm29/sysu_jwxt_v3/pkg/request"
)

type Course struct {
	CourseType *CourseType
	Option     ReqOptions
	BaseInfo   CourseInfo
}

func newCourse(courseType *CourseType, baseInfo CourseInfo, reqOption ReqOptions) *Course {
	return &Course{
		CourseType: courseType,
		Option:     reqOption,
		BaseInfo:   baseInfo,
	}
}

func (c *Course) clazzId() string         { return c.BaseInfo.TeachingClassID }
func (c *Course) courseId() string        { return c.BaseInfo.CourseID }
func (c *Course) teachingClassId() string { return c.BaseInfo.TeachingClassID }
func (c *Course) selectedType() string    { return c.CourseType.SelectedType }
func (c *Course) selectedCate() string    { return c.CourseType.SelectedCate }

func (c *Course) Choose(cl global.JwxtClienter) bool {
	url := global.HOST + "jwxt/choose-course-front-server/classCourseInfo/course/choose"
	ref := global.HOST + "jwxt/mk/courseSelection/"
	tpl := `{"clazzId":"%s","selectedType":"%s","selectedCate":"%s","check":true}`
	body := fmt.Sprintf(tpl, c.clazzId(), c.selectedType(), c.selectedCate())
	respJson := request.PostJson(url, body).Referer(ref).Do(cl).Bytes()
	resp := request.JsonToMap(respJson)

	if resp["code"].(float64) != 200 {
		global.Log.WithField("req", body).Debug("course choose fail:", resp["message"])
		return false
	} else {
		global.Log.WithField("req", body).Debug("course choose success:", resp["data"])
		return true
	}
}

func (c *Course) Cancel(cl global.JwxtClienter) bool {
	url := global.HOST + "jwxt/choose-course-front-server/classCourseInfo/course/back"
	ref := global.HOST + "jwxt/mk/courseSelection/"
	tpl := `{"courseId":"%s","clazzId":"%s","selectedType":"%s"}`
	body := fmt.Sprintf(tpl, c.courseId(), c.clazzId(), c.selectedType())
	respJson := request.PostJson(url, body).Referer(ref).Do(cl).Bytes()
	var resp util.NormalResp
	request.JsonToStruct(respJson, &resp)

	if resp.Code != 200 {
		global.Log.WithField("req", body).Debug("course cancel fail:", resp.Message)
		return false
	} else {
		global.Log.WithField("req", body).Debug("course cancel success:", resp.Data)
		return true
	}
}

func (c *Course) doCollection(cl global.JwxtClienter) []byte {
	url := global.HOST + "jwxt/choose-course-front-server/stuCollectedCourse/create"
	ref := global.HOST + "jwxt/mk/courseSelection/?resourceName=%25E9%2580%2589%25E8%25AF%25BE"
	tpl := `{"classesID":"%s","selectedType":"%s"}`
	body := fmt.Sprintf(tpl, c.teachingClassId(), c.selectedType())
	return request.PostJson(url, body).Referer(ref).Do(cl).Bytes()
}

func (c *Course) PushCollection(cl global.JwxtClienter) bool {
	respJson := c.doCollection(cl)
	var resp util.NormalResp
	request.JsonToStruct(respJson, &resp)
	if resp.Code != 200 {
		global.Log.WithField("req", c.teachingClassId()).Debug("courseCollection add fail:", resp.Message)
		return false
	} else {
		global.Log.WithField("req", c.teachingClassId()).Debug("courseCollection add success:", resp.Data)
		return true
	}
}

func (c *Course) PopCollection(cl global.JwxtClienter) bool {
	respJson := c.doCollection(cl)
	var resp util.NormalResp
	request.JsonToStruct(respJson, &resp)
	if resp.Code != 200 {
		global.Log.WithField("req", c.teachingClassId()).Debug("courseCollection remove fail:", resp.Message)
		return false
	} else {
		global.Log.WithField("req", c.teachingClassId()).Debug("courseCollection remove success:", resp.Data)
		return true
	}
}

func (c *Course) VacancyNum() int {
	selecedNum := util.AtoI(c.BaseInfo.CourseSelectedNum)
	receiveNum := c.BaseInfo.AddReceiveNum + c.BaseInfo.BaseReceiveNum
	return receiveNum - selecedNum
}

func (c *Course) CourseName() string {
	return c.BaseInfo.CourseName
}

func (c *Course) VacancyInfo() string {
	return fmt.Sprintf("Name:%s Vacancy:%d", c.CourseName(), c.VacancyNum())
}

func (c *Course) Refresh(cl global.JwxtClienter) error {
	req := NewCourseListReq(c.CourseType, WithCourseName(c.BaseInfo.CourseNum))
	courseList, n_page := req.DoPage(cl)
	if n_page != 1 || len(courseList.Courses) != 1 {
		global.Log.WithField("n_page", n_page).Info("注意:查询到不止一个页面或不止一个课程\t", util.WhereAmI())
	}
	course, err := courseList.First()
	if err != nil {
		return err
	}
	if course.courseId() != c.courseId() {
		return errors.New("课程不存在,仅找到课程:" + strings.Join(courseList.CourseNames(), ","))
	}
	c.BaseInfo = course.BaseInfo
	return nil
}

// note: u should call course.refresh to update it
func (c *Course) IsSelected() bool {
	return c.BaseInfo.SelectedStatus == "4"
}
