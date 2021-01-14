package jwxtClient

import (
	"server/backend/jwxtClient/request"
)

func (c *jwxtClient) getStudentImg(id string) []byte {
	url := "https://jwxt.sysu.edu.cn/jwxt/student-status/stu-photo/photo?photoType=1&stuNumber=" + id
	referer := "https://jwxt.sysu.edu.cn/jwxt/mk/studentWeb/"
	return request.Get(url).Referer(referer).Do(c.HttpClient)
}
