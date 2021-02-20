package teacher

import (
	"github.com/liwm29/sysu_jwxt_v3/pkg/jwxtClient/global"
	"github.com/liwm29/sysu_jwxt_v3/pkg/request"
)

func (t *teacher) GetImg(c global.JwxtClienter) []byte {
	url := global.HOST + "jwxt/evaluation-manage/evaluationMission/profile?no=" + t.id
	ref := global.HOST + "jwxt/mk/evaluation/"
	return request.Get(url).Referer(ref).Do(c).Bytes()
}
