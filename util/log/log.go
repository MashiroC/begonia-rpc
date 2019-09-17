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

func Info(i ...interface{}) {
	log(info, i)
}

func Warn(i ...interface{}) {
	log(warn, i)
}

func Error(i ...interface{}) {
	log(error, i)
}

func Fatal(i ...interface{}) {
	log(error, i)
	os.Exit(-2)
}

func log(level int, in ...interface{}) {
	now := time.Now().String()
	now = now[:19]
	message := fmt.Sprintf("[magicloop] %s |", now)
	for i := 0; i < len(in); i++ {
		message = fmt.Sprintln(message, in[i])
	}
	logPrint(level, message)
}

func logPrint(color int, message string) {
	fmt.Printf("\x1b[0;%dm%s\x1b[0m", color, message)
}
