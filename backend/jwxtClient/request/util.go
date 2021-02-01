package request

import (
	"bytes"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	"github.com/mattn/go-runewidth"
	"github.com/rodaine/table"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
)

var JsonParseErrFilePath = "jsonParseFailSrcData.log"

func JsonErr(err error, data []byte) {
	if err != nil {
		ioutil.WriteFile(JsonParseErrFilePath, data, 0666)
		printParseErrDetailMsg(data)
		PanicIf(err)
	}
}

func ReactIf(err error, f func()) {
	if err != nil {
		f()
		PanicIf(err)
	}
}

func JsonToMap(data []byte) map[string]interface{} {
	var d map[string]interface{}
	JsonErr(json.Unmarshal(data, &d), data)
	return d
}

func JsonToStruct(data []byte, v interface{}) {
	JsonErr(json.Unmarshal(data, v), data)
}

func JsonConvert(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// 403
func is403Forbidden(data []byte) bool {
	dom, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	PanicIf(err)
	text, err := dom.Find("head > title:nth-child(1)").Html()
	if err != nil {
		return false
	}
	if text == "403 Forbidden" {
		return true
	}
	return false
}

// 不对外网开放
func isAccessForbidden(data []byte) bool {
	dom, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	PanicIf(err)
	text, err := dom.Find("title").Html()
	if err != nil {
		return false
	}
	if text == "资源或业务被限制访问  Access Forbidden" {
		return true
	}
	return false
}

func IsHtml(data []byte) bool {
	return bytes.Contains(data, []byte("<!DOCTYPE html"))
}

func IsJson(data []byte) bool {
	return data[0] == byte('{')
}

func printParseErrDetailMsg(data []byte) {
	color.Red("Error detected: json parse")
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tab := table.New("Data type", "Error details", "Suggest").WithWidthFunc(runewidth.StringWidth)
	tab.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	dataType := "unknown"
	details := "null"
	suggest := "检查" + JsonParseErrFilePath
	if IsHtml(data) {
		dataType = "html"
		if is403Forbidden(data) {
			details = "403 forbidden"
			suggest = "检查请求的referer等头部"
		} else if isAccessForbidden(data) {
			details = "资源或业务被限制访问"
			suggest = "检查是否处于校园网环境"
		}
	} else if IsJson(data) {
		dataType = "json"
	}
	tab.AddRow(dataType, details, suggest)
	tab.Print()
}

func LogRequest(req *HttpReq, resp *HttpResp) {
	logger.Println("reqUrl=", req.Request.URL.String(), "reqCookie=", cookie2names(req.Cookies()),
		"respSetCookie=", resp.Header.Get("Set-Cookie"), "resp=", resp.String())
	// color.Blue("request.url=", req.Request.URL.String(), "cookie=", cookie2names(req.Cookies()))
}

func PanicIf(err error) {
	if err == nil {
		return
	}

	color.Red("Error catched: %s", err.Error())
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tab := table.New("# Traceback", "Func", "# Line")
	tab.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	for i := 0; i < 10; i++ {
		pc, _, line, _ := runtime.Caller(i)
		f := runtime.FuncForPC(pc)
		tab.AddRow("Error stack"+itoa(i), f.Name(), line)
		if f.Name() == "main.main" {
			break
		}
	}
	tab.Print()
	os.Exit(0)
}

func itoa(number int) string {
	return strconv.Itoa(number)
}

func atoi(number string) int {
	i, err := strconv.Atoi(number)
	PanicIf(err)
	return i
}
