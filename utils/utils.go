package utils

import (
	"github.com/snwfdhmp/errlog"
)

// CheckErr checks error and print it if it exists
func CheckErr(e error) {
	if e != nil {
		errlog.Debug(e)
		panic(e)
	}
}
