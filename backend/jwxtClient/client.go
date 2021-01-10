package jwxtClient

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
)

type jwxtClient struct {
	*http.Client
	username string
	isLogin  bool
}

func NewClient() *jwxtClient {
	jar, _ := cookiejar.New(nil)

	return &jwxtClient{
		Client: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) > 15 {
					return http.ErrUseLastResponse
				} else {
					fmt.Println("[Redirecting] via ", req.URL.Path, "Method: ", req.Method, "Cookie: ", req.Cookies())
				}
				return nil
			},
			Jar: jar,
		},
		username: "NULL",
		isLogin:  false,
	}
}
