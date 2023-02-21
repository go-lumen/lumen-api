package utils

import (
	"context"
	"github.com/go-lumen/lumen-api/config"
	"github.com/sirupsen/logrus"
	"github.com/snwfdhmp/errlog"
	"reflect"
	"strconv"
	"time"
)

// CheckErr checks error and print it if it exists
func CheckErr(e error) bool {
	if e != nil {
		errlog.Debug(e)
		return true
	}
	return false
}

// Log logs if debug env var is set at true
func Log(ctxt context.Context, level string, msg ...interface{}) {
	if ctxt == nil || config.GetBool(ctxt, "debug") {
		switch level {
		case "info":
			logrus.Infoln(msg...)
		case "warn":
			logrus.Warnln(msg...)
		case "error":
			logrus.Errorln(msg...)
		default:
			logrus.Infoln(msg...)
		}
	}
}

// RemoveStringFromSlice allows to remove a string in a slice
func RemoveStringFromSlice(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

// FindStringInSlice allows to find a string in a slice
func FindStringInSlice(val string, slice []string) (isFound bool) {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

// FindIntInSlice allows to find an int in a slice
func FindIntInSlice(val int64, slice []int64) (isFound bool) {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

// MongoIDToTimestamp allows extracting timestamp from Mongo Object ID
func MongoIDToTimestamp(mongoID string) (ret time.Time) {
	ts, _ := strconv.ParseInt(mongoID[0:8], 16, 32)
	return time.Unix(ts, 0)
}

// FormatDate allows to format a date
func FormatDate(timestamp int64, format string) (ret string) {
	return time.Unix(timestamp, 0).Format(format)
}

// GenerateTimestampArray allows to generate a Timestamp array
func GenerateTimestampArray(startTS, endTS int64) (tsArray []int64) {
	daysNbr := (endTS - startTS) / (24 * 3600)
	var i int64
	for i = 0; i <= daysNbr; i++ {
		tsArray = append(tsArray, startTS+(i*24*3600))
	}
	return tsArray
}

// EnsurePointer ensures that an interface{} is passed as a pointer. Panic if it is not a pointer.
func EnsurePointer(obj interface{}) {
	if reflect.ValueOf(obj).Kind() != reflect.Ptr {
		panic("the argument must be passed as a pointer")
	}
}

// EnforcePointer takes an object and enforce a pointer. Returning a pointer if it is not the case.
func EnforcePointer(obj interface{}) interface{} {
	if reflect.ValueOf(obj).Kind() != reflect.Ptr {
		ptr := reflect.New(reflect.TypeOf(obj))
		ptr.Elem().Set(reflect.ValueOf(obj))
		return ptr.Interface()
	}
	return obj
}

// Typ returns object type.
func Typ(any interface{}) reflect.Type {
	return reflect.TypeOf(EnforcePointer(any)).Elem()
}

// ValidateInputField allows returning correct value from an input
func ValidateInputField(oldValue, newValue string) string {
	if newValue == "undefined" || newValue == "" || newValue == " " {
		return oldValue
	}
	return newValue
}

// CalcBusinessDays allows computing number of business days between 2 dates
func CalcBusinessDays(from, to time.Time) int {
	totalDays := float32(to.Sub(from) / (24 * time.Hour))
	weekDays := float32(from.Weekday()) - float32(to.Weekday())
	businessDays := int(1 + (totalDays*5-weekDays*2)/7)
	if from.Weekday() == time.Saturday || from.Weekday() == time.Sunday || to.Weekday() == time.Saturday || to.Weekday() == time.Sunday {
		businessDays--
	}

	return businessDays
}
