package teacher

import "server/backend/jwxtClient/request"

func (t *teacher) GetImg(c request.Clienter) []byte {
	url := "https://jwxt.sysu.edu.cn/jwxt/evaluation-manage/evaluationMission/profile?no=" + t.id
	ref := "https://jwxt.sysu.edu.cn/jwxt/mk/evaluation/"
	return request.Get(url).Referer(ref).Do(c).Bytes() 
}

