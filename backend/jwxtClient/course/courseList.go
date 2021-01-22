package course

import ()

type CourseList struct {
	Courses []*Course
}

func NewCourseList(yearTerm string, courseType *CourseType, courseInfos []CourseInfo) *CourseList {
	courses := make([]*Course, 0, len(courseInfos))
	for i := range courseInfos {
		courses = append(courses, newCourse(yearTerm, courseType, courseInfos[i]))
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

func (l *CourseList) First() *Course {
	return l.Courses[0]
}
