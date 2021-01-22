/*
由于标准库的http/cookiejar的jar结构体字段均未导出,只导出了两个函数用于jar的接口,
这使得我们无法dump cookie,因此只能自己实现一个简单的cookiejar作为代替.
这个cookiejar并没有模拟浏览器实现cookie的管理,而只是简单的存储.
如果想实现一个完整的cookie管理模块,将标准库的实现拷贝过来即可
*/

package request

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"server/backend/jwxtClient/util"
)

// todo: 并发加锁

type CookieJar struct {
	DB map[string][]*http.Cookie
}

func NewSimpleJar() *CookieJar {
	return &CookieJar{DB: make(map[string][]*http.Cookie)}
}

func (j *CookieJar) SetCookies(url *url.URL, cookies []*http.Cookie) {
	// attention: it is a stack
	j.DB[url.Host] = append(cookies, j.DB[url.Host]...)
}

func (j *CookieJar) Cookies(url *url.URL) []*http.Cookie {
	return j.DB[url.Host]
}

func (j *CookieJar) StoreCookies(filepath string) error {
	cookieJson, err := json.MarshalIndent(j.DB, "", "\t")
	util.PanicIf(err)
	return ioutil.WriteFile(filepath, cookieJson, 0666)
}

func (j *CookieJar) LoadCookies(filepath string) error {
	cookieJson, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	return json.Unmarshal(cookieJson, &j.DB)
}

func (j *CookieJar) Clear() {
	// j.DB = make(map[string][]*http.Cookie)
	for k, _ := range j.DB {
		delete(j.DB, k)
	}
}
