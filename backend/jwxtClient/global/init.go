/*
 * 相当于一个全局变量容器,redux/vuex
 */

package global

import (
	"github.com/sirupsen/logrus"
	"os"
)

var (
	Log       *logrus.Logger
	HOST      string
	YEAR_TERM string
	MODE      int
)

func init() {
	Log = logrus.New()
	f, err := os.OpenFile("log", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		Log.Panic("创建global.log文件失败:", err.Error())
	}
	Log.Out = f
	SetLogLevel_DEBUG()
}

func SetLogLevel_INFO() {
	Log.SetLevel(logrus.InfoLevel)
}

func SetLogLevel_DEBUG() {
	Log.SetLevel(logrus.DebugLevel)
}
