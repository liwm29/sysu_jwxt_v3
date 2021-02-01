package request

import (
	"io"
	"net/http"
)

type HttpClient struct {
	Cl *http.Client
}

func NewClient() *HttpClient {
	c := &HttpClient{
		Cl: &http.Client{
			Jar: NewSimpJar(),
		},
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
				if len(via) == 1 {
					logger.Println("[  Request  ] via ", via[0].URL.String(), "\t\t", "Cookie: ", cookie2names(via[0].Cookies()))
				}
				logger.Println("[Redirecting] via ", req.URL.String(), "\t\t", "Cookie: ", cookie2names(req.Cookies()))
			}
			return nil
		}
	}

	c.Cl.CheckRedirect = f
}

func (c *HttpClient) Do(req *HttpReq) *HttpResp {
	resp, err := c.Cl.Do(req.Request)
	respWrapper := NewResponse(resp)
	PanicIf(err)
	LogRequest(req, respWrapper)
	return respWrapper
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
	return c.Cl.Jar.(*simpJar).StoreCookies(filepath)
}
func (c *HttpClient) LoadCookies(filepath string) error {
	return c.Cl.Jar.(*simpJar).LoadCookies(filepath)
}
