package logger

import (
	"FallGuys66/config"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

var debugLogger *log.Logger
var infoLogger *log.Logger
var warnLogger *log.Logger
var errorLogger *log.Logger

func init() {
	logFilePath := filepath.Join(config.UserConfigDir, "FallGuys66.log")
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatalln("Failed to open logger file:", err)
	}
	debugLogger = log.New(io.MultiWriter(file, os.Stdout), "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile)
	infoLogger = log.New(io.MultiWriter(file, os.Stdout), "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	warnLogger = log.New(io.MultiWriter(file, os.Stderr), "[WARN] ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(io.MultiWriter(file, os.Stderr), "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Errorf(format string, v ...interface{}) {
	errorLogger.Output(2, fmt.Sprintf(format, v...))
}
func Error(v interface{}) {
	errorLogger.Output(2, fmt.Sprintf("%v", v))
}
func Err(err error) {
	errorLogger.Output(2, err.Error())
}
func Warnf(format string, v ...interface{}) {
	warnLogger.Output(2, fmt.Sprintf(format, v...))
}
func Warn(v interface{}) {
	warnLogger.Output(2, fmt.Sprintf("%v", v))
}
func Infof(format string, v ...interface{}) {
	infoLogger.Output(2, fmt.Sprintf(format, v...))
}
func Info(v interface{}) {
	infoLogger.Output(2, fmt.Sprintf("%v", v))
}
func Debugf(format string, v ...interface{}) {
	debugLogger.Output(2, fmt.Sprintf(format, v...))
}
func Debug(v interface{}) {
	debugLogger.Output(2, fmt.Sprintf("%v", v))
}
func Deb(err error) {
	debugLogger.Output(2, err.Error())
}
func ShowJson(format string, v ...interface{}) {
	var vt []any
	for _, i := range v {
		b, err := json.Marshal(i)
		if err == nil {
			vt = append(vt, string(b))
		}
	}
	debugLogger.Output(2, fmt.Sprintf(format, vt...))
}
