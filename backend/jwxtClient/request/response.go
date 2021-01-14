package request

import (
	"io/ioutil"
	"net/http"
)

type HttpResp struct {
	*http.Response
}

func NewResponse(resp *http.Response) *HttpResp {
	return &HttpResp{resp}
}

func (resp *HttpResp) ReadAll() []byte {
	bytes, err := ioutil.ReadAll(resp.Body)
	PanicIf(err)
	resp.Body.Close()
	return bytes
}
