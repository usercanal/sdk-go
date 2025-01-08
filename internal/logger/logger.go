// logger/logger.go
package logger

import (
	"fmt"
	"log"
	"os"
)

type Logger struct {
	debug bool
	log   *log.Logger
}

var defaultLogger = &Logger{
	debug: false,
	log:   log.New(os.Stderr, "[usercanal] ", log.LstdFlags),
}

func SetDebug(debug bool) {
	defaultLogger.debug = debug
}

func Debug(format string, v ...interface{}) {
	if defaultLogger.debug {
		defaultLogger.log.Output(2, fmt.Sprintf("DEBUG: "+format, v...))
	}
}

func Info(format string, v ...interface{}) {
	defaultLogger.log.Output(2, fmt.Sprintf("INFO: "+format, v...))
}

func Warn(format string, v ...interface{}) {
	defaultLogger.log.Output(2, fmt.Sprintf("WARN: "+format, v...))
}

func Error(format string, v ...interface{}) {
	defaultLogger.log.Output(2, fmt.Sprintf("ERROR: "+format, v...))
}
