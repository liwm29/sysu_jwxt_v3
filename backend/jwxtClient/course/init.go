package course

import "github.com/sirupsen/logrus"

var (
	_campus_id = map[string]string{
		"东校园":  "5063559",
		"北校园":  "5062202",
		"南校园":  "5062201",
		"深圳校区": "333291143",
		"珠海校区": "5062203",
		"所有校区": "",
	}
	_selectedType = map[string]string{
		"本专业":  "1",
		"校级公选": "4",
		"跨专业":  "2",
	}
	_selectedCate = map[string]string{
		"专必":     "11",
		"专选":     "21",
		"院内公选":   "30",
		"公必(体育)": "10",
		"公必(大英)": "10",
		"公必(其他)": "10",
	}
)

const (
	CAMPUS_EAST  = "5063559"
	CAMPUS_NORTH = "5062202"
	CAMPUS_SOUTH = "5062201"
	CAMPUS_SZ    = "333291143"
	CAMPUS_ZH    = "5062203"
	CAMPUS_ALL   = ""

	// TYPE_MAJOR      = "1"
	// TYPE_PUBLIC     = "4"
	// TYPE_CROSSMAJOR = "2"

	// CATE_COMPULSORY      = "11"
	// CATE_ELECTIVE        = "21"
	// CATE_SCHOOL_ELECTIVE = "30"
	// CATE_GYM             = "10"
	// CATE_ENGLISH         = "10"
	// CATE_OTHERS           = "10"

	TYPE_PUB_ELECTIVE   = "公选"
	TYPE_MAJ_ELECTIVE   = "专选"
	TYPE_MAJ_COMPULSORY = "专必"
	TYPE_GYM            = "体育"
	TYPE_ENGLISH        = "英语"
	TYPE_OTHERS         = "其他"

	NAME_ALL = ""
)

var YEAR_TERM = ""

func SetYearTerm(semester string) {
	YEAR_TERM = semester
}

var log *logrus.Logger

func SetLogger(logger *logrus.Logger) {
	log = logger
}

func init() {
}
