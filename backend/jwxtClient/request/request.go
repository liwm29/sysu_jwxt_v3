package request

import (
	"io"
	"net/http"
	"strings"
)

type HttpReq struct {
	*http.Request
}

func setDefaultHdr(req *http.Request) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36")
}

func NewRequest(method, url string, body io.Reader) *HttpReq {
	req, _ := http.NewRequest(method, url, body)
	setDefaultHdr(req)
	return &HttpReq{req}
}

func Get(url string) *HttpReq {
	return NewRequest("GET", url, nil)
}

func Post(url string, body io.Reader) *HttpReq {
	return NewRequest("POST", url, body)
}

func PostJson(url string, body string) *HttpReq {
	return NewRequest("POST", url, strings.NewReader(body)).Json()
}

func PostForm(url string, body string) *HttpReq {
	return NewRequest("POST", url, strings.NewReader(body)).Form()
}

func (r *HttpReq) Referer(referer string) *HttpReq {
	r.Header.Set("Referer", referer)
	return r
}

func (r *HttpReq) Json() *HttpReq {
	r.Header.Set("Content-Type", "application/json;charset=UTF-8")
	return r
}

func (r *HttpReq) Form() *HttpReq {
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return r
}

type Clienter interface {
	Do(req *HttpReq) *HttpResp
}

func (r *HttpReq) Do(c Clienter) *HttpResp {
	return c.Do(r)
}
