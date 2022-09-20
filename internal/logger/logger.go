package logger

import (
	"fmt"
	"log"
)

type LogLevel int

const (
	LevelError LogLevel = iota // 0
	LevelWarning
	LevelInfo
	LevelDebug
)

var (
	logLevel LogLevel = LevelInfo
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorGray   = "\033[0:37m"
	ColorWhite  = "\033[1:37m"
)

const (
	prefixDebug   = "[DEBUG]"
	prefixError   = "[ERROR]"
	prefixWarning = "[WARN]"
	prefixInfo    = "[INFO]"
	spaceElem     = " "
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	logLevel = LevelInfo
}

func SetLogLevel(level LogLevel) {
	logLevel = level
}

func Printf(format string, args ...interface{}) {
	log.SetPrefix(ColorWhite)
	log.Print(prefixInfo+spaceElem, fmt.Sprintf(format, args...))
}

func Println(v ...any) {
	log.SetPrefix(ColorWhite)
	log.Print(prefixInfo, fmt.Sprintln(v...))
}

func Debugf(format string, args ...interface{}) {
	if logLevel >= LevelDebug {
		log.SetPrefix(ColorGray)
		log.Print(prefixDebug+spaceElem, fmt.Sprintf(format, args...))
	}
}

func Debugln(msg string) {
	if logLevel >= LevelDebug {
		log.SetPrefix(ColorGray)
		log.Println(prefixDebug, msg)
	}
}

func Warnf(format string, args ...interface{}) {
	if logLevel >= LevelWarning {
		log.SetPrefix(ColorYellow)
		log.Print(prefixWarning+spaceElem, fmt.Sprintf(format, args...))
	}
}

func Warnln(msg string) {
	if logLevel >= LevelWarning {
		log.SetPrefix(ColorYellow)
		log.Println(prefixWarning, msg)
	}
}

func Errorln(v ...interface{}) {
	if logLevel >= LevelError {
		log.SetPrefix(ColorRed)
		v = append([]interface{}{prefixError}, v...)
		log.Println(v...)
	}
}

func Errorf(format string, args ...interface{}) {
	if logLevel >= LevelError {
		log.SetPrefix(ColorRed)
		log.Print(prefixWarning+spaceElem, fmt.Sprintf(format, args...))
	}
}

func Fatalln(v ...interface{}) {
	log.SetPrefix(ColorRed)
	v = append([]interface{}{prefixError}, v...)
	log.Fatalln(v...)
}
