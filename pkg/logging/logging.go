package logging

import (
	"log"
	"os"
)

var (
	warnLogger  *log.Logger
	infoLogger  *log.Logger
	errorLogger *log.Logger
	fatalLogger *log.Logger
	traceLogger *log.Logger
)

type logs struct {
	warnLogger  *log.Logger
	infoLogger  *log.Logger
	errorLogger *log.Logger
	fatalLogger *log.Logger
	traceLogger *log.Logger
}

func New() *logs {
	out := os.Stdout

	warnLogger = log.New(out, " [WARN]:  ", log.Ldate|log.Ltime)
	infoLogger = log.New(out, " [INFO]:  ", log.Ldate|log.Ltime)
	errorLogger = log.New(out, "[ERROR]:  ", log.Ldate|log.Ltime)
	fatalLogger = log.New(out, "[FATAL]:  ", log.Ldate|log.Ltime)
	traceLogger = log.New(out, "[TRACE]:  ", log.Ldate|log.Ltime)

	return &logs{warnLogger, infoLogger, errorLogger, fatalLogger, traceLogger}
}

func (l logs) Info(message any) {
	l.infoLogger.Println(message)
}

func (l logs) Infof(format string, v ...any) {
	l.infoLogger.Printf(format, v...)
}

func (l logs) Warn(message any) {
	l.warnLogger.Println(message)
}

func (l logs) Warnf(format string, v ...any) {
	l.warnLogger.Printf(format, v...)
}

func (l logs) Error(message any) {
	l.errorLogger.Println(message)
}

func (l logs) Errorf(format string, v ...any) {
	l.errorLogger.Printf(format, v...)
}

func (l logs) Fatal(message any) {
	l.fatalLogger.Fatal(message)
}

func (l logs) Fatalf(format string, v ...any) {
	l.fatalLogger.Printf(format, v...)
}

func (l logs) Trace(message any) {
	l.traceLogger.Println(message)
}

func (l logs) Tracef(format string, v ...any) {
	l.traceLogger.Printf(format, v...)
}

type Logger struct {
	*logs
}

func GetLogger() *Logger {
	return &logger
}

var logger Logger

func init() {
	l := New()
	logger = Logger{l}
}
