package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Log *logrus.Logger

func Init(debug bool) {
	Log = &logrus.Logger{
		Out:          os.Stdout,
		ExitFunc:     os.Exit,
		ReportCaller: true,
		Formatter:    &logrus.JSONFormatter{},
	}

	if debug {
		Log.SetLevel(logrus.DebugLevel)
		//Log.SetFormatter(&logrus.TextFormatter{
		//	ForceColors:   true,
		//	FullTimestamp: true,
		//})
		Log.SetOutput(os.Stdout)
	} else {
		Log.SetLevel(logrus.WarnLevel)
		// logger := logrus.New()
		// Log.SetOutput(logger.Writer())
	}

}
