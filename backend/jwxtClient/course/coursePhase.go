package course

import "server/backend/jwxtClient/request"

type CoursePhaseResp struct {
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
	Data    CoursePhase `json:"data"`
}

type CoursePhase struct {
	ElectiveCourseStageName string `json:"electiveCourseStageName"`
	RetreatCourseStatus     string `json:"retreatCourseStatus"`
	Code                    int    `json:"code"`
	SemesterYear            string `json:"semesterYear"`
	CourseSelectType        string `json:"courseSelectType"`
	ChooseCourseStatus      string `json:"chooseCourseStatus"`
	ElectiveCourseStageCode string `json:"electiveCourseStageCode"`
	StartTime               string `json:"startTime"`
	EndTime                 string `json:"endTime"`
	CrossMajor              string `json:"crossMajor"`
}

func GetCoursePhase(c request.Clienter) *CoursePhase {
	url := "https://jwxt.sysu.edu.cn/jwxt/choose-course-front-server/classCourseInfo/selectCourseInfo"
	ref := "https://jwxt.sysu.edu.cn/jwxt/mk/courseSelection/"
	respJson := request.Get(url).Referer(ref).Do(c).Bytes()
	var resp CoursePhaseResp
	request.JsonToStruct(respJson, &resp)
	if resp.Code != 200 {
		log.WithField("msg", resp.Message).Error("无法获取选课阶段")
	}
	return &resp.Data
	// {"code":200,"message":null,"data":{"electiveCourseStageName":"改补选","retreatCourseStatus":"1","code":200,"semesterYear":"2020-2","courseSelectType":"0","chooseCourseStatus":"1","electiveCourseStageCode":"3","startTime":"2021-01-08 13:00:00","endTime":"2021-03-04 23:00:00","crossMajor":"1"}}
}

func (p *CoursePhase) CanSelect() bool {
	if p.ElectiveCourseStageCode == "3" {
		return true
	}
	return false
}
