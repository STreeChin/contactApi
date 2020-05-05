package log

import (
	"os"

	"github.com/STreeChin/contactapi/pkg/config"
	"github.com/sirupsen/logrus"
)

//NewLogger instance
func NewLogger(c config.Config) *logrus.Logger {
	var log = logrus.New()
	/*	f, _ := os.Create("./logrus.log")
		log.Out = f*/
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(os.Stdout)
	logrus.SetReportCaller(false)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	return log
}
