package jwxtClient

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

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

func init() {
	log.Out = os.Stdout
	log.SetLevel(logrus.DebugLevel)
}
