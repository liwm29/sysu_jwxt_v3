package request

import (
	"fmt"
	"io"
	"net/http"
	// "net/http/cookiejar"
)

type HttpClient struct {
	Cl        *http.Client
	CookieJar *CookieJar
}

func NewClient() *HttpClient {
	// jar, _ := cookiejar.New(nil)
	jar := NewSimpleJar()
	return &HttpClient{
		Cl: &http.Client{
			Jar: jar,
		},
		CookieJar: jar,
	}
}

func (c *HttpClient) SetRedirectCallback(f func(req *http.Request, via []*http.Request) error) {
	if f == nil {
		f = func(req *http.Request, via []*http.Request) error {
			if len(via) > 15 {
				return http.ErrUseLastResponse
			} else {
				fmt.Println("[Redirecting] via ", req.URL.Path, "Method: ", req.Method, "Cookie: ", req.Cookies())
			}
			return nil
		}
	}

	c.Cl.CheckRedirect = f
}

func (c *HttpClient) Do(req *HttpReq) []byte {
	resp, err := c.Cl.Do(req.Request)
	PanicIf(err)
	return NewResponse(resp).ReadAll()
}

func (c *HttpClient) Get(url string) []byte {
	return Get(url).Do(c)
}

func (c *HttpClient) Post(url string, body io.Reader) []byte {
	return Post(url, body).Do(c)
}

func (c *HttpClient) PostJson(url, body string) []byte {
	return PostJson(url, body).Do(c)
}

func (c *HttpClient) PostForm(url, body string) []byte {
	return PostForm(url, body).Do(c)
}

func (c *HttpClient) StoreCookies(filepath string) {
	c.CookieJar.StoreCookies(filepath)
}
func (c *HttpClient) LoadCookies(filepath string) {
	c.CookieJar.LoadCookies(filepath)
}
