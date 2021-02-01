package student

import (
	"github.com/liwm29/sysu_jwxt_v3/backend/jwxtClient/global"
	"github.com/liwm29/sysu_jwxt_v3/backend/jwxtClient/request"
)

func (s *student) GetImg(c request.Clienter) []byte {
	url := global.HOST + "jwxt/student-status/stu-photo/photo?photoType=1&stuNumber=" + s.id
	referer := global.HOST + "jwxt/mk/studentWeb/"
	return request.Get(url).Referer(referer).Do(c).Bytes()
}
