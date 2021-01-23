package jwxtClient

import (
	"os"

	"github.com/sirupsen/logrus"
	"server/backend/jwxtClient/course"
)

var log = logrus.New()

func SetLogLevel_INFO() {
	log.SetLevel(logrus.InfoLevel)
}

func SetLogLevel_DEBUG() {
	log.SetLevel(logrus.DebugLevel)
}

/* log level
// log.Trace("Something very low level.")
// log.Debug("Useful debugging information.")
// log.Info("Something noteworthy happened!")
// log.Warn("You should probably take a look at this.")
// log.Error("Something failed but I'm not quitting.")
// // Calls os.Exit(1) after logging
// log.Fatal("Bye.")
// // Calls panic() after logging
// log.Panic("I'm bailing.")
*/

var DEFAULT_COOKIE_PATH string = "./cookie"
var DEFAULT_CAPTCHA_PATH string = "./captcha.jpg"

func init() {
	log.Out = os.Stdout
	log.SetLevel(logrus.DebugLevel)

	course.SetLogger(log)
}
