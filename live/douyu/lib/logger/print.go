package logger

import (
	"encoding/json"
	"log"
	"os"
)

var debugLogger = log.New(os.Stdout, "[DEBUG]", log.LstdFlags)
var infoLogger = log.New(os.Stdout, "[INFO]", log.LstdFlags)
var warnLogger = log.New(os.Stderr, "[WARN]", log.LstdFlags)
var errorLogger = log.New(os.Stderr, "[ERROR]", log.LstdFlags)

func Errorf(format string, v ...interface{}) {
	errorLogger.Printf(format, v...)
}
func Error(v ...interface{}) {
	errorLogger.Println(v)
}
func Err(err error) {
	errorLogger.Println(err.Error())
}
func Warnf(format string, v ...interface{}) {
	warnLogger.Printf(format, v...)
}
func Warn(v ...interface{}) {
	warnLogger.Println(v)
}
func Infof(format string, v ...interface{}) {
	infoLogger.Printf(format, v...)
}
func Info(v ...interface{}) {
	infoLogger.Println(v...)
}
func Debugf(format string, v ...interface{}) {
	debugLogger.Printf(format, v...)
}
func Debug(v ...interface{}) {
	debugLogger.Println(v)
}
func Deb(err error) {
	debugLogger.Println(err.Error())
}
func ShowJson(format string, v ...interface{}) {
	var vt []any
	for _, i := range v {
		b, err := json.Marshal(i)
		if err == nil {
			vt = append(vt, string(b))
		}
	}
	debugLogger.Printf(format, vt...)
}
