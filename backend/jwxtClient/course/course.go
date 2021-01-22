package course

import (
	"fmt"
	"server/backend/jwxtClient/request"
	"server/backend/jwxtClient/util"
)

type Course struct {
	CourseType *CourseType
	YearTerm   string
	BaseInfo   CourseInfo
}

func newCourse(yearTerm string, courseType *CourseType, baseInfo CourseInfo) *Course {
	return &Course{
		CourseType: courseType,
		YearTerm:   yearTerm,
		BaseInfo:   baseInfo,
	}
}

func (c *Course) clazzId() string         { return c.BaseInfo.TeachingClassID }
func (c *Course) courseId() string        { return c.BaseInfo.CourseID }
func (c *Course) teachingClassId() string { return c.BaseInfo.TeachingClassID }
func (c *Course) selectedType() string    { return c.CourseType.SelectedType }
func (c *Course) selectedCate() string    { return c.CourseType.SelectedCate }

// func (c *Course) MarshallIndent(prefix, indent string) string {
// 	b, err := json.MarshalIndent(c, "", "\t")
// 	util.PanicIf(err)
// 	return string(b)
// }

func (c *Course) Choose(cl request.Clienter) bool {
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

func (c *Course) Cancel(cl request.Clienter) bool {
	url := "https://jwxt.sysu.edu.cn/jwxt/choose-course-front-server/classCourseInfo/course/back"
	tpl := `{"courseId":"%s","clazzId":"%s","selectedType":"%s"}`
	body := fmt.Sprintf(tpl, c.courseId(), c.clazzId(), c.selectedType())
	respJson := request.PostJson(url, body).Referer("https://jwxt.sysu.edu.cn/jwxt/mk/courseSelection/").Do(cl).Bytes()
	var resp util.NormalResp
	request.JsonToStruct(respJson, &resp)

	if resp.Code != 200 {
		log.WithField("req", body).Debug("course cancel fail:", resp.Message)
		return false
	} else {
		log.WithField("req", body).Debug("course cancel success:", resp.Data)
		return true
	}
}

func (c *Course) doCollection(cl request.Clienter) []byte {
	url := "https://jwxt.sysu.edu.cn/jwxt/choose-course-front-server/stuCollectedCourse/create"
	tpl := `{"classesID":"%s","selectedType":"%s"}`
	body := fmt.Sprintf(tpl, c.teachingClassId(), c.selectedType())
	return request.PostJson(url, body).Referer("https://jwxt.sysu.edu.cn/jwxt/mk/courseSelection/?resourceName=%25E9%2580%2589%25E8%25AF%25BE").Do(cl).Bytes()

}

func (c *Course) PushCollection(cl request.Clienter) bool {
	respJson := c.doCollection(cl)
	var resp util.NormalResp
	request.JsonToStruct(respJson, &resp)
	if resp.Code != 200 {
		log.WithField("req", c.teachingClassId()).Debug("courseCollection add fail:", resp.Message)
		return false
	} else {
		log.WithField("req", c.teachingClassId()).Debug("courseCollection add success:", resp.Data)
		return true
	}
}

func (c *Course) PopCollection(cl request.Clienter) bool {
	respJson := c.doCollection(cl)
	var resp util.NormalResp
	request.JsonToStruct(respJson, &resp)
	if resp.Code != 200 {
		log.WithField("req", c.teachingClassId()).Debug("courseCollection remove fail:", resp.Message)
		return false
	} else {
		log.WithField("req", c.teachingClassId()).Debug("courseCollection remove success:", resp.Data)
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
