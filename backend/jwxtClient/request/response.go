package request

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"server/backend/jwxtClient/util"
)

type HttpResp struct {
	*http.Response
	copy   []byte
	isRead bool
}

func NewResponse(resp *http.Response) *HttpResp {
	return &HttpResp{resp, nil, false}
}

func (resp *HttpResp) Bytes() []byte {
	if resp.isRead {
		return resp.copy
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	util.PanicIf(err)
	resp.Body.Close()
	resp.copy = bytes
	resp.isRead = true
	return bytes
}

func (resp *HttpResp) String() string {
	if !resp.isRead {
		resp.isRead = true
		return string(resp.Bytes())
	}

	return string(resp.copy)
}

func (resp *HttpResp) Reader() *bytes.Reader {
	if !resp.isRead {
		resp.Bytes()
	}

	return bytes.NewReader(resp.copy)
}
