package util

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	"github.com/rodaine/table"
)

type NormalResp struct {
	Code    float64
	Message string
	Data    string
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
		tab.AddRow("Error stack"+ItoA(i), f.Name(), line)
		if f.Name() == "main.main" {
			break
		}
	}
	tab.Print()
	os.Exit(0)
}

func Truncate100(data string) string {
	return TruncateN(data, 100)
}

func TruncateN(data string, N int) string {
	if len(data) < N {
		return data
	}
	return data[:N]
}

func WhereAmI() string {
	pc, _, line, _ := runtime.Caller(1)
	return fmt.Sprintf("%s#Line%d", runtime.FuncForPC(pc).Name(), line)
}

func ItoA(number int) string {
	return strconv.Itoa(number)
}

func AtoI(number string) int {
	i, err := strconv.Atoi(number)
	PanicIf(err)
	return i
}

func Bool2Str(is bool) string {
	if is {
		return "1"
	} else {
		return "0"
	}
}

func Str2Bool(is string) bool {
	if is == "" || is == "0" || strings.ToLower(is) == "null" {
		return false
	} else {
		return true
	}
}

func Min(a int, b int) int {
	x := math.Min(float64(a), float64(b))
	return int(x)
}

func IsAccessForbidden(data []byte) bool {
	dom, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	PanicIf(err)
	text, _ := dom.Find("title").Html()
	return text == "资源或业务被限制访问  Access Forbidden"
}
