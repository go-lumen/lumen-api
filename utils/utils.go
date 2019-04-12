package utils

import (
	"github.com/sirupsen/logrus"
)

// CheckErr checks error and print it if it exists
func CheckErr(e error) {
	if e != nil {
		logrus.Errorln(e)
		panic(e)
	}
}
