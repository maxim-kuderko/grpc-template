package initializers

import (
	"github.com/sirupsen/logrus"
	"os"
)

const defaultLogLevel = logrus.WarnLevel

func init() {
	initLogger()
}

func initLogger() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	lvl := defaultLogLevel
	if v := os.Getenv(`LOG_LEVEL`); v != `` {
		if tmp, err := logrus.ParseLevel(v); err != nil {
			lvl = tmp
		}
	}
	logrus.SetLevel(lvl)
}
