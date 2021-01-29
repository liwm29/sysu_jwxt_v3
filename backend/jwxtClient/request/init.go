package request

import (
	"log"
	"os"
)

var DefaultClient *HttpClient
var logger *log.Logger

func init() {
	DefaultClient = NewClient()
	f, err := os.OpenFile("requestLog", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Panic("创建requsetLog失败", err)
	}
	logger = log.New(f, "[request.Log]", 0)
}
