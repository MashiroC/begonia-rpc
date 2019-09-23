package log

import (
	"fmt"
	"os"
	"time"
)

const (
	warn  = 33
	error = 31
	info  = 2
)

func Info(format string, i ...interface{}) {
	log(info, format, i...)
}

func Warn(format string, i ...interface{}) {
	log(warn, format, i...)
}

func Error(format string, i ...interface{}) {
	log(error, format, i...)
}

func Fatal(format string, i ...interface{}) {
	log(error, format, i...)
	os.Exit(-2)
}

func log(level int, format string, in ...interface{}) {
	now := time.Now().String()
	now = now[:19]
	header := fmt.Sprintf("[magicloop] %s | ", now)
	message := header + fmt.Sprintf(format, in...)
	logPrint(level, message)
}

func logPrint(color int, message string) {
	fmt.Printf("\x1b[0;%dm%s\x1b[0m\n", color, message)
}
