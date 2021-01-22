package course

import (
	"errors"
	"server/backend/jwxtClient/util"
)

type CourseType struct {
	SelectedType string
	SelectedCate string
}

// param: 公选/专选
func NewCourseType(courseTypeStr string) *CourseType {
	var selectedType, selectedCate string
	switch courseTypeStr {
	case TYPE_PUB_ELECTIVE:
		selectedType = _selectedType["校级公选"]
		selectedCate = ""
	case TYPE_MAJ_ELECTIVE:
		selectedType = _selectedType["本专业"]
		selectedCate = _selectedCate["专选"]
	case TYPE_MAJ_COMPULSORY:
		selectedType = _selectedType["本专业"]
		selectedCate = _selectedCate["专必"]
	case TYPE_GYM:
		selectedType = _selectedType["本专业"]
		selectedCate = _selectedCate["公必(体育)"]
	case TYPE_ENGLISH:
		selectedType = _selectedType["本专业"]
		selectedCate = _selectedCate["公必(大英)"]
	case TYPE_OTHERS:
		selectedType = _selectedType["本专业"]
		selectedCate = _selectedCate["公必(其他)"]
	default:
		util.PanicIf(errors.New("not registered courseType"))
	}
	return &CourseType{
		SelectedType: selectedType,
		SelectedCate: selectedCate,
	}
}
