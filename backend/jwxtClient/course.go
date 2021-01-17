package jwxtClient

import (
	"encoding/json"
	"fmt"
	"server/backend/jwxtClient/request"
)

type Course struct {
	SelectedType string
	SelectedCate string
	YearTerm     string
	BaseInfo     courseInfo
}

func newCourse(courseType, yearTerm string, baseInfo courseInfo) *Course {
	selectedType, selectedCate := getCourseType(courseType)

	return &Course{
		SelectedType: selectedType,
		SelectedCate: selectedCate,
		YearTerm:     yearTerm,
		BaseInfo:     baseInfo,
	}
}

func newCourses(courseType, yearTerm string, courseBases []courseInfo) []*Course {
	courses := make([]*Course, 0, len(courseBases))
	for i := range courseBases {
		courses = append(courses, newCourse(courseType, yearTerm, courseBases[i]))
	}
	return courses
}

func (c *Course) clazzId() string         { return c.BaseInfo.TeachingClassID }
func (c *Course) courseId() string        { return c.BaseInfo.CourseID }
func (c *Course) teachingClassId() string { return c.BaseInfo.TeachingClassID }
func (c *Course) selectedType() string    { return c.SelectedType }
func (c *Course) selectedCate() string    { return c.SelectedCate }

func (c *Course) MarshallIndent(prefix, indent string) string {
	b, err := json.MarshalIndent(c, "", "\t")
	PanicIf(err)
	return string(b)
}

func (c *Course) Choose(cl *JwxtClient) bool {
	url := "https://jwxt.sysu.edu.cn/jwxt/choose-course-front-server/classCourseInfo/course/choose"
	tpl := `{"clazzId":"%s","selectedType":"%s","selectedCate":"%s","check":true}`
	body := fmt.Sprintf(tpl, c.clazzId(), c.selectedType(), c.selectedCate())
	respJson := request.PostJson(url, body).Referer("https://jwxt.sysu.edu.cn/jwxt/mk/courseSelection/").Do(cl).Bytes()
	resp := request.JsonToMap(respJson)

	if resp["code"].(float64) != 200 {
		log.WithField("req", body).Debug("course choose fail:", resp["message"])
		return false
	} else {
		log.WithField("req", body).Debug("course choose success:", resp["data"])
		return true
	}
}

func (c *Course) Cancel(cl *JwxtClient) bool {
	url := "https://jwxt.sysu.edu.cn/jwxt/choose-course-front-server/classCourseInfo/course/back"
	tpl := `{"courseId":"%s","clazzId":"%s","selectedType":"%s"}`
	body := fmt.Sprintf(tpl, c.courseId(), c.clazzId(), c.selectedType())
	respJson := request.PostJson(url, body).Referer("https://jwxt.sysu.edu.cn/jwxt/mk/courseSelection/").Do(cl).Bytes()
	var resp normalResp
	request.JsonToStruct(respJson, &resp)

	if resp.Code != 200 {
		log.WithField("req", body).Debug("course cancel fail:", resp.Message)
		return false
	} else {
		log.WithField("req", body).Debug("course cancel success:", resp.Data)
		return true
	}
}

func (c *Course) doCollection(cl *JwxtClient) []byte {
	url := "https://jwxt.sysu.edu.cn/jwxt/choose-course-front-server/stuCollectedCourse/create"
	tpl := `{"classesID":"%s","selectedType":"%s"}`
	body := fmt.Sprintf(tpl, c.teachingClassId(), c.selectedType())
	return request.PostJson(url, body).Referer("https://jwxt.sysu.edu.cn/jwxt/mk/courseSelection/?resourceName=%25E9%2580%2589%25E8%25AF%25BE").Do(cl).Bytes()

}

func (c *Course) PushCollection(cl *JwxtClient) bool {
	respJson := c.doCollection(cl)
	var resp normalResp
	request.JsonToStruct(respJson, &resp)
	if resp.Code != 200 {
		log.WithField("req", c.teachingClassId()).Debug("courseCollection add fail:", resp.Message)
		return false
	} else {
		log.WithField("req", c.teachingClassId()).Debug("courseCollection add success:", resp.Data)
		return true
	}
}

func (c *Course) PopCollection(cl *JwxtClient) bool {
	respJson := c.doCollection(cl)
	var resp normalResp
	request.JsonToStruct(respJson, &resp)
	if resp.Code != 200 {
		log.WithField("req", c.teachingClassId()).Debug("courseCollection remove fail:", resp.Message)
		return false
	} else {
		log.WithField("req", c.teachingClassId()).Debug("courseCollection remove success:", resp.Data)
		return true
	}
}

// 这个函数假设course是empty的,取响应的第一个course,replace这个course
func (c *Course) SearchFirst(courseName string) {

}
