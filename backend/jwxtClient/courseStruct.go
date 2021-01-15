package jwxtClient

import "fmt"

type courseListReq struct {
	pageNo           int
	pageSize         int
	yearTerm         string
	selectedType     string
	selectedCate     string
	campusId         string
	collectionStatus string
}

func (r *courseListReq) setNextPage() *courseListReq {
	r.pageNo += 1
	return r
}

func (r *courseListReq) marshall() string {
	tpl := `{"pageNo":%d,"pageSize":%d,"param":{"semesterYear":"%s","selectedType":"%s","selectedCate":"%s","hiddenConflictStatus":"0","hiddenSelectedStatus":"0","hiddenEmptyStatus":"0","vacancySortStatus":"0","collectionStatus":"%s","campusId":"%s"}}`
	return fmt.Sprintf(tpl, r.pageNo, r.pageSize, r.yearTerm, r.selectedType, r.selectedCate, r.collectionStatus, r.campusId)
}

type courseListResp struct {
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
	Data    struct {
		Total int   `json:"total"`
		Rows  []Row `json:"rows"`
	} `json:"data"`
}

type Row struct {
	MainClassesID        string      `json:"mainClassesID"`
	TeachingClassID      string      `json:"teachingClassId"`
	TeachingClassNum     string      `json:"teachingClassNum"`
	TeachingClassName    interface{} `json:"teachingClassName"`
	CourseNum            string      `json:"courseNum"`
	CourseName           string      `json:"courseName"`
	Credit               float64     `json:"credit"`
	ExamFormName         string      `json:"examFormName"`
	CourseUnitNum        string      `json:"courseUnitNum"`
	CourseUnitName       string      `json:"courseUnitName"`
	TeachingTeacherNum   string      `json:"teachingTeacherNum"`
	TeachingTeacherName  string      `json:"teachingTeacherName"`
	BaseReceiveNum       int         `json:"baseReceiveNum"`
	AddReceiveNum        int         `json:"addReceiveNum"`
	TeachingTimePlace    string      `json:"teachingTimePlace"`
	StudyCampusID        string      `json:"studyCampusId"`
	Week                 string      `json:"week"`
	ClassTimes           string      `json:"classTimes"`
	CourseSelectedNum    string      `json:"courseSelectedNum"`
	FilterSelectedNum    string      `json:"filterSelectedNum"`
	SelectedStatus       string      `json:"selectedStatus"`
	CollectionStatus     string      `json:"collectionStatus"`
	TeachingLanguageCode string      `json:"teachingLanguageCode"`
	PubCourseTypeCode    interface{} `json:"pubCourseTypeCode"`
	CourseCateCode       string      `json:"courseCateCode"`
	SpecialClassCode     interface{} `json:"specialClassCode"`
	SportItemID          interface{} `json:"sportItemId"`
	RecordMode           string      `json:"recordMode"`
	ClazzNum             string      `json:"clazzNum"`
	ExamFormCode         string      `json:"examFormCode"`
	CourseID             string      `json:"courseId"`
	ScheduleExamTime     interface{} `json:"scheduleExamTime"`
}

type courseInfo struct {
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
		TeacherList []struct {
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
		} `json:"teacherList"`
	} `json:"data"`
}

type coursePhaseResp struct {
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
	Data    coursePhase `json:"data"`
}

type coursePhase struct {
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
