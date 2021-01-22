package jwxtClient

import (
	"server/backend/jwxtClient/student"
	"server/backend/jwxtClient/teacher"
)

func (c *JwxtClient) GetStudentImg(id string) []byte {
	s := student.NewStudent(id)
	return s.GetImg(c)
}

func (c *JwxtClient) GetTeacherImg(id string) []byte {
	t := teacher.NewTeacher(id)
	return t.GetImg(c)
}

