package jwxtClient

import (
	"fmt"
	"server/backend/jwxtClient/request"
)

type course struct {
	courseId        string
	clazzId         string
	teachingClassId string // clazzId is the same as teachingClassId
	selectedType    string
	selectedCate    string
	yearTerm        string
}

type normalResp struct {
	Code    float64
	Message string
	Data    string
}

func NewCourse(clazzId, selectedType, selectedCate string) {
	return
}

func (c *course) Choose(cl *jwxtClient) bool {
	url := "https://jwxt.sysu.edu.cn/jwxt/choose-course-front-server/classCourseInfo/course/choose"
	tpl := `{"clazzId":"%s","selectedType":"%s","selectedCate":"%s","check":true}`
	body := fmt.Sprintf(tpl, c.clazzId, c.selectedType, c.selectedCate)
	respJson := request.PostJson(url, body).Referer("https://jwxt.sysu.edu.cn/jwxt/mk/courseSelection/").Do(cl)
	resp := request.JsonToMap(respJson)

	if resp["code"].(float64) != 200 {
		log.WithField("req", body).Debug("course choose fail:", resp["message"])
		return false
	} else {
		log.WithField("req", body).Debug("course choose success:", resp["data"])
		return true
	}
}

func (c *course) Cancel(cl *jwxtClient) bool {
	url := "https://jwxt.sysu.edu.cn/jwxt/choose-course-front-server/classCourseInfo/course/back"
	tpl := `{"courseId":"%s","clazzId":"%s","selectedType":"%s"}`
	body := fmt.Sprintf(tpl, c.courseId, c.clazzId, c.selectedType)
	respJson := request.PostJson(url, body).Referer("https://jwxt.sysu.edu.cn/jwxt/mk/courseSelection/").Do(cl)
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

func (c *course) doCollection(cl *jwxtClient) []byte {
	url := "https://jwxt.sysu.edu.cn/jwxt/choose-course-front-server/stuCollectedCourse/create"
	tpl := `{"classesID":"%s","selectedType":"%s"}`
	body := fmt.Sprintf(tpl, c.teachingClassId, c.selectedType)
	return request.PostJson(url, body).Referer("https://jwxt.sysu.edu.cn/jwxt/mk/courseSelection/?resourceName=%25E9%2580%2589%25E8%25AF%25BE").Do(cl)

}

func (c *course) PushCollection(cl *jwxtClient) bool {
	respJson := c.doCollection(cl)
	var resp normalResp
	request.JsonToStruct(respJson, &resp)
	if resp.Code != 200 {
		log.WithField("req", c.teachingClassId).Debug("courseCollection add fail:", resp.Message)
		return false
	} else {
		log.WithField("req", c.teachingClassId).Debug("courseCollection add success:", resp.Data)
		return true
	}
}

func (c *course) PopCollection(cl *jwxtClient) bool {
	respJson := c.doCollection(cl)
	var resp normalResp
	request.JsonToStruct(respJson, &resp)
	if resp.Code != 200 {
		log.WithField("req", c.teachingClassId).Debug("courseCollection remove fail:", resp.Message)
		return false
	} else {
		log.WithField("req", c.teachingClassId).Debug("courseCollection remove success:", resp.Data)
		return true
	}
}
