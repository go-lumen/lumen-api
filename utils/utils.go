package utils

import (
	"context"
	"github.com/go-lumen/lumen-api/config"
	"github.com/sirupsen/logrus"
	"github.com/snwfdhmp/errlog"
)

// CheckErr checks error and print it if it exists
func CheckErr(e error) {
	if e != nil {
		errlog.Debug(e)
	}
}

// Log logs if debug env var is set at true
func Log(ctxt context.Context, level string, msg ...interface{}) {
	if ctxt == nil || config.GetBool(ctxt, "debug") {
		switch level {
		case "info":
			logrus.Infoln(msg)
		case "warn":
			logrus.Warnln(msg)
		case "error":
			logrus.Errorln(msg)
		default:
			logrus.Infoln(msg)
		}
	}
}
