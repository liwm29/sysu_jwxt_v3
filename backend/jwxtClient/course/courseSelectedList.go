package course

import "github.com/liwm29/sysu_jwxt_v3/backend/jwxtClient/request"

func GetSelectedCourseNames(cl request.Clienter) (courseNames []string) {
	url := "https://jwxt-443.webvpn.sysu.edu.cn/jwxt/choose-course-front-server/selectedCourse/list"
	ref := "https://jwxt-443.webvpn.sysu.edu.cn/jwxt/mk/courseSelection/"
	body := `{"pageNo":1,"pageSize":15,"total":true,"param":{"successStatus":"1","failureStatus":"0","retiredClass":"0","waitingScreen":"0"}}`
	respJson := request.PostJson(url, body).Referer(ref).Do(cl).Bytes()
	data := request.JsonToMap(respJson)
	courses := data["data"].(map[string]interface{})["rows"].([]interface{})
	for _, v := range courses {
		courseNames = append(courseNames, v.(map[string]interface{})["courseName"].(string))
	}
	return courseNames
}
