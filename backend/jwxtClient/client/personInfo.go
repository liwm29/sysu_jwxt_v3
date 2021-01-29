package client

import (
	"server/backend/jwxtClient/request"
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

type user struct {
	Name   string
	Id     string
	School string
	Major  string
	Grade  string
}

func (c *JwxtClient) GetUserInfo() *user {
	url := "https://jwxt-443.webvpn.sysu.edu.cn/jwxt/student-status/countrystu/studentRollView"
	ref := "https://jwxt-443.webvpn.sysu.edu.cn/jwxt/mk/studentWeb/"
	respJson := request.Get(url).Referer(ref).Do(c).Bytes()
	resp := request.JsonToMap(respJson)
	data := resp["data"].(map[string]interface{})
	return &user{
		Name:   data["basicName"].(string),
		Id:     data["studentNumber"].(string),
		School: data["rollCollegeNumNAME"].(string),
		Major:  data["rollStandardNAME"].(string),
		Grade:  data["rollGrade"].(string),
	}
}
