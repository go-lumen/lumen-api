package utils

import (
	"github.com/sirupsen/logrus"
)

func CheckErr(e error) {
	if e != nil {
		logrus.Errorln(e)
		panic(e)
	}
}
