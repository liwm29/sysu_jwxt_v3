package course

import ()

var (
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

	TYPE_PUB_ELECTIVE   = "公选"
	TYPE_MAJ_ELECTIVE   = "专选"
	TYPE_MAJ_COMPULSORY = "专必"
	TYPE_GYM            = "体育"
	TYPE_ENGLISH        = "英语"
	TYPE_OTHERS         = "其他"

	NAME_ALL = ""
)

func CampusEasyMap(c string) string {
	switch c {
	case "东", "east", "e":
		return CAMPUS_EAST
	case "北", "north", "n":
		return CAMPUS_NORTH
	case "南", "south", "s":
		return CAMPUS_SOUTH
	case "深圳", "shenzhen", "sz":
		return CAMPUS_SZ
	case "珠海", "zhuhai", "zh":
		return CAMPUS_ZH
	default:
		return c
	}
}
