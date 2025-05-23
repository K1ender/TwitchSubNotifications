package logger

import (
	log "github.com/sirupsen/logrus"
)

var Log = log.New()

func init() {
	Log.SetLevel(log.DebugLevel)
	Log.SetFormatter(&log.TextFormatter{})
	Log.SetReportCaller(true)
}
