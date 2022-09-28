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

func Printf(format string, v ...any) {
	log.SetPrefix(ColorWhite)
	log.Print(prefixInfo+spaceElem, fmt.Sprintf(format, v...))
}

func Println(v ...any) {
	log.SetPrefix(ColorWhite)
	log.Println(prefixInfo+spaceElem, fmt.Sprint(v...))
}

func Debugf(format string, v ...any) {
	if logLevel >= LevelDebug {
		log.SetPrefix(ColorGray)
		log.Print(prefixDebug+spaceElem, fmt.Sprintf(format, v...))
	}
}

func Debugln(v ...any) {
	if logLevel >= LevelDebug {
		log.SetPrefix(ColorGray)
		log.Println(prefixDebug+spaceElem, fmt.Sprint(v...))
	}
}

func Warnf(format string, v ...any) {
	if logLevel >= LevelWarning {
		log.SetPrefix(ColorYellow)
		log.Print(prefixWarning+spaceElem, fmt.Sprintf(format, v...))
	}
}

func Warnln(v ...any) {
	if logLevel >= LevelWarning {
		log.SetPrefix(ColorYellow)
		log.Println(prefixWarning+spaceElem, fmt.Sprint(v...))
	}
}

func Errorln(v ...any) {
	if logLevel >= LevelError {
		log.SetPrefix(ColorRed)
		log.Println(prefixError+spaceElem, fmt.Sprint(v...))
	}
}

func Errorf(format string, args ...any) {
	if logLevel >= LevelError {
		log.SetPrefix(ColorRed)
		log.Print(prefixWarning+spaceElem, fmt.Sprintf(format, args...))
	}
}

func Fatalln(v ...any) {
	log.SetPrefix(ColorRed)
	log.Fatalln(prefixError+spaceElem, fmt.Sprint(v...))
}
