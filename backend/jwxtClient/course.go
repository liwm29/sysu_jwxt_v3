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
}

func NewCourse(clazzId, selectedType, selectedCate string) {
	return
}

// type chooseResp struct {
// 	code ,message, data string,
// }

func (c *course) Choose(cl *jwxtClient) bool {
	url := "https://jwxt.sysu.edu.cn/jwxt/choose-course-front-server/classCourseInfo/course/choose"
	tpl := `{"clazzId":"%s","selectedType":"%s","selectedCate":"%s","check":true}`
	body := fmt.Sprintf(tpl, c.clazzId, c.selectedType, c.selectedCate)
	respJson := request.PostJson(url, body).Referer("https://jwxt.sysu.edu.cn/jwxt/mk/courseSelection/").Do(cl.HttpClient)
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
	respJson := request.PostJson(url, body).Referer("https://jwxt.sysu.edu.cn/jwxt/mk/courseSelection/").Do(cl.HttpClient)
	resp := request.JsonToMap(respJson)

	if resp["code"].(float64) != 200 {
		log.WithField("req", body).Debug("course cancel fail:", resp["message"])
		return false
	} else {
		log.WithField("req", body).Debug("course cancel success:", resp["data"])
		return true
	}
}

func (c *course) doCollection(cl *jwxtClient) []byte {
	url := "https://jwxt.sysu.edu.cn/jwxt/choose-course-front-server/stuCollectedCourse/create"
	tpl := `{"classesID":"%s","selectedType":"%s"}`
	body := fmt.Sprintf(tpl, c.teachingClassId, c.selectedType)
	return request.PostJson(url, body).Referer("https://jwxt.sysu.edu.cn/jwxt/mk/courseSelection/?resourceName=%25E9%2580%2589%25E8%25AF%25BE").Do(cl.HttpClient)

}

func (c *course) PushCollection(cl *jwxtClient) bool {
	respJson := c.doCollection(cl)
	resp := request.JsonToMap(respJson)
	if resp["code"].(float64) != 200 {
		log.WithField("req", c.teachingClassId).Debug("courseCollection add fail:", resp["message"])
		return false
	} else {
		log.WithField("req", c.teachingClassId).Debug("courseCollection add success:", resp["data"])
		return true
	}
}

func (c *course) PopCollection(cl *jwxtClient) bool {
	respJson := c.doCollection(cl)
	resp := request.JsonToMap(respJson)
	if resp["code"].(float64) != 200 {
		log.WithField("req", c.teachingClassId).Debug("courseCollection remove fail:", resp["message"])
		return false
	} else {
		log.WithField("req", c.teachingClassId).Debug("courseCollection remove success:", resp["data"])
		return true
	}
}

func (c *course) getCollections(cl *jwxtClient) {
	// tpl := `{"pageNo":1,"pageSize":10,"param":{"semesterYear":"%s","selectedType":"%s","selectedCate":"%s","hiddenConflictStatus":"0","hiddenSelectedStatus":"0","hiddenEmptyStatus":"0","vacancySortStatus":"0","collectionStatus":"1","studyCampusId":"5063559"}}`
	// payload := fmt.Sprintf(tpl, "2020-2", c.selectedType, c.selectedCate)
	// rows, err := client.getCourseList(payload)
	// if err != nil {
	// 	log.Println(err)
	// 	return nil, err
	// }
	// return rows, nil
}
