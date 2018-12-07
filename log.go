package gobatis

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

// 测试日志类, 性能低下, 不建议使用
func now() string {
	date := time.Now().Format("2006-01-02 15:04:06")
	return date
}

func getCallers() []string {
	callers := make([]string, 0)
	for i := 0; true; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		id := strings.LastIndex(file, "/") + 1
		caller := fmt.Sprintf("%s:%d", (string)(([]byte(file))[id:]), line)
		callers = append(callers, caller)
	}

	return callers
}

type ILogger interface {
	SetLevel(level LogLevel)
	SetFileName(fileName string)
	Info(format string, v ...interface{})
	Debug(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Error(format string, v ...interface{})
	Fatal(format string, v ...interface{})
}

type LogLevel int

// ALL < DEBUG < INFO < WARN < ERROR < FATAL < OFF
const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
	LogLevelOff
)

type OutType int

const (
	OutTypeFile OutType = iota
	OutTypeStd
)

type iOut interface {
	getOutType() OutType
	println(msg string)
	Close()
}

type logger struct {
	out           iOut
	logLevel      LogLevel
	mu            sync.Mutex
	callStepDepth int
}

var defaultLogLevel = LogLevelDebug

type stdLogger struct{ mu sync.Mutex }

func (this *stdLogger) println(v string) {
	this.mu.Lock()
	defer this.mu.Unlock()

	fmt.Println(v)
}

func (this *stdLogger) getOutType() OutType {
	return OutTypeStd
}

func (this *stdLogger) Close() {}

type fileLogger struct {
	file    *os.File
	mu      sync.Mutex
	isClose bool
}

func (this *fileLogger) println(v string) {
	this.mu.Lock()
	defer this.mu.Unlock()

	if nil != this.file {
		this.file.WriteString(v)
		this.file.WriteString("\n")
	}
}

func (this *fileLogger) getOutType() OutType {
	return OutTypeFile
}

func (this *fileLogger) Close() {
	this.mu.Lock()
	defer this.mu.Unlock()

	if this.isClose {
		return
	}

	this.file.Close()
	this.file = nil
	this.isClose = true
}

var defLog = &logger{logLevel: defaultLogLevel, out: &stdLogger{}, callStepDepth: 1}

func NewFileLog(fileName string, level LogLevel) ILogger {
	logger := &logger{
		logLevel: level,
	}

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if nil == err {
		logger.out = &fileLogger{
			file: file,
		}

		return logger
	}

	return nil
}

func (this *logger) getPrefix(flag string) string {
	prefix := fmt.Sprintf("%s [%5s] - ", now(), flag)
	callers := getCallers()
	if len(callers) >= (3 + this.callStepDepth + 1) {
		prefix = fmt.Sprintf("%s [%5s] [%s] - ", now(), flag, callers[3+this.callStepDepth])
	}

	return prefix
}

func (this *logger) SetCallStepDepth(stepDepth int) {
	this.mu.Lock()
	defer this.mu.Unlock()

	this.callStepDepth = stepDepth
}

func (this *logger) SetFileName(fileName string) {
	if this.out.getOutType() == OutTypeFile {
		this.mu.Lock()
		defer this.mu.Unlock()

		file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if nil == err {
			this.out = &fileLogger{
				file: file,
			}
		}
	}
}

func (this *logger) SetLevel(level LogLevel) {
	this.mu.Lock()
	defer this.mu.Unlock()

	this.logLevel = level
}

func (this *logger) Info(format string, v ...interface{}) {
	if this.logLevel <= LogLevelInfo {
		logStr := fmt.Sprintf(this.getPrefix("INFO")+format, v...)

		this.out.println(logStr)
	}
}

func (this *logger) Debug(format string, v ...interface{}) {
	if this.logLevel <= LogLevelDebug {
		logStr := fmt.Sprintf(this.getPrefix("DEBUG")+format, v...)
		this.out.println(logStr)
	}
}

func (this *logger) Warn(format string, v ...interface{}) {
	if this.logLevel <= LogLevelWarn {
		logStr := fmt.Sprintf(this.getPrefix("WARN")+format, v...)
		this.out.println(logStr)
	}
}

func (this *logger) Error(format string, v ...interface{}) {
	if this.logLevel <= LogLevelError {
		logStr := fmt.Sprintf(this.getPrefix("ERROR")+format, v...)
		this.out.println(logStr)
	}
}

func (this *logger) Fatal(format string, v ...interface{}) {
	if this.logLevel <= LogLevelFatal {
		logStr := fmt.Sprintf(this.getPrefix("FATAL")+format, v...)
		this.out.println(logStr)
	}
}

func SetLevel(lv LogLevel) {
	defLog.SetLevel(lv)
}

func Info(format string, v ...interface{}) {
	defLog.Info(format, v...)
}

func Debug(format string, v ...interface{}) {
	defLog.Debug(format, v...)
}

func Warn(format string, v ...interface{}) {
	defLog.Warn(format, v...)
}

func Error(format string, v ...interface{}) {
	defLog.Error(format, v...)
}

func Fatal(format string, v ...interface{}) {
	defLog.Fatal(format, v...)
}
