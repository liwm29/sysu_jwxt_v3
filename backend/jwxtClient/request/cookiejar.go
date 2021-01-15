/*
actually u dont need to simulate a browser to maintain the cookies,
simply we can just store it,never care others like expire
*/

package request

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

type CookieJar struct {
	DB map[string][]*http.Cookie
}

func NewSimpleJar() *CookieJar {
	return &CookieJar{DB: make(map[string][]*http.Cookie)}

}

func (j *CookieJar) SetCookies(url *url.URL, cookies []*http.Cookie) {
	j.DB[url.Host] = append(j.DB[url.Host], cookies...)
}

func (j *CookieJar) Cookies(url *url.URL) []*http.Cookie {
	return j.DB[url.Path]
}

func (j *CookieJar) StoreCookies(filepath string) {
	cookieJson, err := json.MarshalIndent(j.DB, "", "\t")
	PanicIf(err)
	ioutil.WriteFile(filepath, cookieJson, 0666)
}

func (j *CookieJar) LoadCookies(filepath string) {
	cookieJson, err := ioutil.ReadFile(filepath)
	PanicIf(err)
	json.Unmarshal(cookieJson, &j.DB)
}
