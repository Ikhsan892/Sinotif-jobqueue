package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"
)

type Log struct {
	log           *log.Logger
	file          *os.File
	defaultLogger *log.Logger
}

type ILog interface {
	Debug(appName string, data ...any)
	Error(appName string, data ...any)
	Info(appName string, data ...any)
	ErrorF(appName string, format string, v ...any)
	Warn(appName string, data ...any)
}

func (l *Log) NewLogger(prefix string, appName string) {
	// logging
	currentTime := time.Now()
	formatLog := fmt.Sprintf("/logs/actions/app_%s.log", currentTime.Format("2006-01-02"))
	logDirectory := GetWorkingDirectoryContent(formatLog)

	_, errDirectory := os.Stat(logDirectory)
	if os.IsNotExist(errDirectory) {
		errMKDIR := os.MkdirAll(GetWorkingDirectoryContent(fmt.Sprintf("/logs/actions")), os.ModePerm)
		if errMKDIR != nil {
			log.Println(errMKDIR)
		}
	}
	f, err := os.OpenFile(logDirectory, os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			f, _ = os.Create(logDirectory)
			f, err = os.OpenFile(logDirectory, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		}
	}

	l.file = f
	l.log = log.New(f, fmt.Sprintf("[ %s - %s ] : ", prefix, appName), log.LstdFlags|log.Lshortfile|log.Lmsgprefix)
	l.defaultLogger = log.New(os.Stdout, fmt.Sprintf("[ %s - %s ] : ", prefix, appName), log.LstdFlags|log.Lshortfile|log.Lmsgprefix)
}

func (l *Log) print(prefix string, appName string, msg ...any) {
	l.NewLogger(prefix, appName)
	l.log.Println(msg)
	l.defaultLogger.Println(msg)
	l.file.Close()
}

func (l *Log) printf(prefix string, appName string, format string, msg ...any) {
	l.NewLogger(prefix, appName)
	l.log.Printf(format, msg)
	l.defaultLogger.Printf(format, msg)
	l.file.Close()
}

func Debug(appName string, msg ...any) {
	l := Log{}
	l.print("DEBUG", appName, msg)
}

func Info(appName string, msg ...any) {
	l := Log{}
	l.print("INFO", appName, msg)
}

func Error(appName string, msg ...any) {
	l := Log{}
	l.print("ERROR", appName, msg)
}

func Warn(appName string, msg ...any) {
	l := Log{}
	l.print("WARN", appName, msg)
}

func ErrorF(appName string, format string, v ...any) {
	l := Log{}
	l.printf("ERROR", appName, format, v)
}
