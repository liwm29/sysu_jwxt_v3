/*
 * 相当于一个全局变量容器,redux/vuex
 */

package global

import (
	"fmt"
	"os"
	"path"

	"github.com/liwm29/sysu_jwxt_v3/pkg/request"
	"github.com/sirupsen/logrus"
)

var (
	Log          *logrus.Logger
	HOST         string
	YEAR_TERM    string
	MODE         int
	kLogFilePath string = "jwxtLog"
)

func setLogFolder(folderPath string) error {
	err := os.MkdirAll(folderPath, 0666)
	if err != nil {
		return wrapError("setlogfolder:", err)
	}
	kLogFilePath = path.Join(folderPath, kLogFilePath)
	return nil
}

func wrapError(newErr string, wrappedErr error) error {
	return fmt.Errorf("%s%s", newErr, wrappedErr.Error())
}

func init() {
	Log = logrus.New()
	err := setLogFolder("./log")
	if err != nil {
		Log.Panic(err)
	}
	f, err := os.OpenFile(kLogFilePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		Log.Panic(wrapError("创建global.log文件失败:", err))
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

type JwxtClienter interface {
	// Do(*request.HttpReq) *request.HttpResp
	request.Clienter
}
