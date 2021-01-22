package student

import (
	"server/backend/jwxtClient/request"
)

func (s *student) GetImg(c request.Clienter) []byte {
	url := "https://jwxt.sysu.edu.cn/jwxt/student-status/stu-photo/photo?photoType=1&stuNumber=" + s.id
	referer := "https://jwxt.sysu.edu.cn/jwxt/mk/studentWeb/"
	return request.Get(url).Referer(referer).Do(c).Bytes()
}
