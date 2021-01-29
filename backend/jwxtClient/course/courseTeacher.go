package course

import (
	"errors"
	"server/backend/jwxtClient/global"
	"server/backend/jwxtClient/request"
	"server/backend/jwxtClient/util"
)

type courseOutlineResp struct {
	Code int `json:"code"`
	Data struct {
		OutlineInfo struct {
			Creator                       string  `json:"creator"`
			CreatorName                   string  `json:"creatorName"`
			CreateTime                    string  `json:"createTime"`
			Editor                        string  `json:"editor"`
			EditorName                    string  `json:"editorName"`
			EditeTime                     string  `json:"editeTime"`
			OutlineCourseInfoID           string  `json:"outlineCourseInfoId"`
			CourseID                      string  `json:"courseId"`
			CourseNum                     string  `json:"courseNum"`
			OutlineType                   string  `json:"outlineType"`
			CourseType                    string  `json:"courseType"`
			CourseTypeName                string  `json:"courseTypeName"`
			LanguageNum                   string  `json:"languageNum"`
			PlanClassSize                 string  `json:"planClassSize"`
			ReferenceBook                 string  `json:"referenceBook"`
			CourseContentInChinese        string  `json:"courseContentInChinese"`
			CourseObjectiveAndRequirement string  `json:"courseObjectiveAndRequirement"`
			TeachMethod                   string  `json:"teachMethod"`
			EvaluationMethod              string  `json:"evaluationMethod"`
			CurrentFlowNum                string  `json:"currentFlowNum"`
			AuditStatus                   string  `json:"auditStatus"`
			GiveLessSemtster              string  `json:"giveLessSemtster"`
			LecturesCreHours              float64 `json:"lecturesCreHours"`
			TutorialsCreHours             float64 `json:"tutorialsCreHours"`
			LabCreHours                   float64 `json:"labCreHours"`
			OtherLectureCreHours          float64 `json:"otherLectureCreHours"`
			CourseResource                string  `json:"courseResource"`
			LastAuditStatus               string  `json:"lastAuditStatus"`
			SubCourseType                 string  `json:"subCourseType"`
			SubCourseTypeName             string  `json:"subCourseTypeName"`
			CourseName                    string  `json:"courseName"`
			CourseEngName                 string  `json:"courseEngName"`
			Credit                        string  `json:"credit"`
			TotalHours                    string  `json:"totalHours"`
			EstablishUnitNumberName       string  `json:"establishUnitNumberName"`
		} `json:"outlineInfo"`
		ScheduleList []struct {
			OutlineTeachingScheduleID string `json:"outlineTeachingScheduleId"`
			OutlineCourseInfoID       string `json:"outlineCourseInfoId"`
			WeekNum                   int    `json:"weekNum"`
			TeachingMainContent       string `json:"teachingMainContent"`
			TeachingHours             string `json:"teachingHours"`
			Sort                      string `json:"sort"`
		} `json:"scheduleList"`
		TeacherList []teacherInfo `json:"teacherList"`
	} `json:"data"`
}

type teacherInfo struct {
	Editor               string `json:"editor"`
	EditeTime            string `json:"editeTime"`
	ID                   string `json:"id"`
	TeacherNum           string `json:"teacherNum"`
	DepartmentNum        string `json:"departmentNum"`
	Name                 string `json:"name"`
	NameSpell            string `json:"nameSpell"`
	IDCardTypeNum        string `json:"idCardTypeNum"`
	IDCardNum            string `json:"idCardNum"`
	GenderNum            string `json:"genderNum"`
	BirthDate            string `json:"birthDate"`
	NationNum            string `json:"nationNum"`
	NationalityNum       string `json:"nationalityNum"`
	PoliticsNum          string `json:"politicsNum"`
	TeacherTypeNum       string `json:"teacherTypeNum"`
	ProfessionNum        string `json:"professionNum"`
	ProfessionNumName    string `json:"professionNumName"`
	BestEducationNum     string `json:"bestEducationNum"`
	BestEducationNumName string `json:"bestEducationNumName"`
	ThisStateNum         string `json:"thisStateNum"`
	InEstablishment      string `json:"inEstablishment"`
	OnDuty               string `json:"onDuty"`
	CampusID             string `json:"campusId"`
	Email                string `json:"email"`
	MobilePhone          string `json:"mobilePhone"`
	BestDegree           string `json:"bestDegree"`
	BestDegreeName       string `json:"bestDegreeName"`
	StationType          string `json:"stationType"`
	EmployTime           string `json:"employTime"`
}

func (c *Course) GetTeachers(cl request.Clienter) ([]teacherInfo, error) {
	url := global.HOST + "jwxt/training-programe/courseoutline/getalloutlineinfo?courseNum=" + c.BaseInfo.CourseNum
	ref := global.HOST + "jwxt/mk/courseSelection/"
	respJson := request.Get(url).Referer(ref).Do(cl).Bytes()
	global.Log.WithField("respJson", util.Truncate100(string(respJson))).Debug(util.WhereAmI())
	// 首先检查是否有大纲信息 {"code":52000000,"message":"未查询到大纲信息"}
	m := request.JsonToMap(respJson)
	if m["code"].(float64) == 52000000 {
		return nil, errors.New(m["message"].(string))
	}

	var outline courseOutlineResp
	request.JsonToStruct(respJson, &outline)

	return outline.Data.TeacherList, nil
}
