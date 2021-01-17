package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"runtime"

	"github.com/PuerkitoBio/goquery"
)

func PanicIf(err error) {
	if err != nil {
		pc0, _, line0, _ := runtime.Caller(0)
		pc1, _, line1, _ := runtime.Caller(1)
		pc2, _, line2, _ := runtime.Caller(2)
		pc3, _, line3, _ := runtime.Caller(3)
		fmt.Printf("---------\nError happans:\n")
		fmt.Printf("Error stack0: %s \t \t Line#%d\n", runtime.FuncForPC(pc0).Name(), line0)
		fmt.Printf("Error stack1: %s \t \t Line#%d\n", runtime.FuncForPC(pc1).Name(), line1)
		fmt.Printf("Error stack2: %s \t \t Line#%d\n", runtime.FuncForPC(pc2).Name(), line2)
		fmt.Printf("Error stack3: %s \t \t Line#%d\n", runtime.FuncForPC(pc3).Name(), line3)

		if err.Error() == "invalid character '<' looking for beginning of value" {
			fmt.Println("猜测原因:响应返回了html,被当作了json进行解析")
		}
		panic(err)
	}
}

func ReactIf(err error, f func(), args ...interface{}) {
	if err != nil {
		f()
		checkParseData(args[0].([]byte))
		PanicIf(err)
	}
}

func JsonToMap(data []byte) map[string]interface{} {
	var d map[string]interface{}
	ReactIf(json.Unmarshal(data, &d), func() {
		ioutil.WriteFile("jsonParseFailSrcData", data, 0666)
	}, data)
	return d
}

func JsonToStruct(data []byte, v interface{}) {
	ReactIf(json.Unmarshal(data, v), func() {
		ioutil.WriteFile("jsonParseFailSrcData", data, 0666)
	}, data)
}

func JsonConvert(data []byte, v interface{}) {
	ReactIf(json.Unmarshal(data, v), func() {
		ioutil.WriteFile("jsonParseFailSrcData", data, 0666)
	}, data)
}

func find403Forbidden(data []byte) bool {
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

func checkParseData(data []byte) {
	if find403Forbidden(data) {
		fmt.Println("在待解析数据中发现了403 forbidden,检查请求的referer等头部")
	} else {
		fmt.Println("未在待解析数据中发现了403 forbidden")
	}
}
