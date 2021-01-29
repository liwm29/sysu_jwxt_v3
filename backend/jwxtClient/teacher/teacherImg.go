package teacher

import (
	"server/backend/jwxtClient/global"
	"server/backend/jwxtClient/request"
)

func (t *teacher) GetImg(c request.Clienter) []byte {
	url := global.HOST + "jwxt/evaluation-manage/evaluationMission/profile?no=" + t.id
	ref := global.HOST + "jwxt/mk/evaluation/"
	return request.Get(url).Referer(ref).Do(c).Bytes()
}
