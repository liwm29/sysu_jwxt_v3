package course

import "errors"

type CourseList struct {
	Courses []*Course
}

func NewCourseList(courseType *CourseType, baseInfos []CourseInfo, reqOption ReqOptions) *CourseList {
	courses := make([]*Course, 0, len(baseInfos))
	for i := range baseInfos {
		courses = append(courses, newCourse(courseType, baseInfos[i], reqOption))
	}
	return &CourseList{courses}
}

func (l *CourseList) CourseNames() []string {
	courseNames := make([]string, 0, len(l.Courses))
	for _, v := range l.Courses {
		courseNames = append(courseNames, v.BaseInfo.CourseName)
	}
	return courseNames
}

func (l *CourseList) First() (*Course, error) {
	if len(l.Courses) > 0 {
		return l.Courses[0], nil
	}
	return nil, errors.New("len(courseList.Courses) == 0,can't call First()")
}
