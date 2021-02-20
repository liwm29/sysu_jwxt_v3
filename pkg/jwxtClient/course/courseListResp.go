package course

type CourseListResp struct {
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
	Data    struct {
		Total int          `json:"total"`
		Rows  []CourseInfo `json:"rows"`
	} `json:"data"`
}

type CourseInfo struct {
	MainClassesID        string      `json:"mainClassesID"`
	TeachingClassID      string      `json:"teachingClassId"` // note this is clazzId
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
