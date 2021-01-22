package request

import (
	"fmt"
	"io"
	"net/http"
	"server/backend/jwxtClient/util"
)

type HttpClient struct {
	Cl        *http.Client
	CookieJar *CookieJar
}

func NewClient() *HttpClient {
	jar := NewSimpleJar()
	c := &HttpClient{
		Cl: &http.Client{
			Jar: jar,
		},
		CookieJar: jar,
	}
	c.SetRedirectCallback(nil)
	return c
}

func (c *HttpClient) SetRedirectCallback(f func(req *http.Request, via []*http.Request) error) {
	if f == nil {
		f = func(req *http.Request, via []*http.Request) error {
			if len(via) > 20 {
				return http.ErrUseLastResponse
			} else {
				fmt.Println("[Redirecting] via ", util.TruncateN(req.URL.Path, 40), "\t\t", "Method:", req.Method, "Cookie: ", req.Cookies())
			}
			return nil
		}
	}

	c.Cl.CheckRedirect = f
}

func (c *HttpClient) Do(req *HttpReq) *HttpResp {
	resp, err := c.Cl.Do(req.Request)
	util.PanicIf(err)
	return NewResponse(resp)
}

func (c *HttpClient) Get(url string) *HttpResp {
	return Get(url).Do(c)
}

func (c *HttpClient) Post(url string, body io.Reader) *HttpResp {
	return Post(url, body).Do(c)
}

func (c *HttpClient) PostJson(url, body string) *HttpResp {
	return PostJson(url, body).Do(c)
}

func (c *HttpClient) PostForm(url, body string) *HttpResp {
	return PostForm(url, body).Do(c)
}

func (c *HttpClient) StoreCookies(filepath string) error {
	return c.CookieJar.StoreCookies(filepath)
}
func (c *HttpClient) LoadCookies(filepath string) error {
	return c.CookieJar.LoadCookies(filepath)
}
