package gobatis

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"
)

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

// 如果想定制logger可以实现此接口，否则日志将使用默认打印
type ILogger interface {
	SetLevel(level LogLevel)
	Info(format string, v ...interface{})
	Debug(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Error(format string, v ...interface{})
	Fatal(format string, v ...interface{})
}

type LogLevel int

// ALL < DEBUG < INFO < WARN < ERROR < FATAL < OFF
const (
	LOG_LEVEL_DEBUG LogLevel = iota
	LOG_LEVEL_INFO
	LOG_LEVEL_WARN
	LOG_LEVEL_ERROR
	LOG_LEVEL_FATAL
	LOG_LEVEL_OFF
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

var defaultLogLevel = LOG_LEVEL_DEBUG

type stdLogger struct{ mu sync.Mutex }

func (sl *stdLogger) println(v string) {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	fmt.Println(v)
}

func (sl *stdLogger) getOutType() OutType {
	return OutTypeStd
}

func (sl *stdLogger) Close() {}

var defLog = &logger{logLevel: defaultLogLevel, out: &stdLogger{}, callStepDepth: 0}

func (l *logger) getPrefix(flag string) string {
	prefix := fmt.Sprintf("%s [%5s] - ", now(), flag)
	callers := getCallers()
	if len(callers) >= 6 {
		prefix = fmt.Sprintf("%s [%5s] [%s] - ", now(), flag, callers[3+l.callStepDepth])
	}

	return prefix
}

func (l *logger) SetCallStepDepth(stepDepth int) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.callStepDepth = stepDepth
}

func (l *logger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.logLevel = level
}

func (l *logger) Info(format string, v ...interface{}) {
	if l.logLevel <= LOG_LEVEL_INFO {
		logStr := fmt.Sprintf(l.getPrefix("INFO")+format, v...)

		l.out.println(logStr)
	}
}

func (l *logger) Debug(format string, v ...interface{}) {
	if l.logLevel <= LOG_LEVEL_DEBUG {
		logStr := fmt.Sprintf(l.getPrefix("DEBUG")+format, v...)
		l.out.println(logStr)
	}
}

func (l *logger) Warn(format string, v ...interface{}) {
	if l.logLevel <= LOG_LEVEL_WARN {
		logStr := fmt.Sprintf(l.getPrefix("WARN")+format, v...)
		l.out.println(logStr)
	}
}

func (l *logger) Error(format string, v ...interface{}) {
	if l.logLevel <= LOG_LEVEL_ERROR {
		logStr := fmt.Sprintf(l.getPrefix("ERROR")+format, v...)
		l.out.println(logStr)
	}
}

func (l *logger) Fatal(format string, v ...interface{}) {
	if l.logLevel <= LOG_LEVEL_FATAL {
		logStr := fmt.Sprintf(l.getPrefix("FATAL")+format, v...)
		l.out.println(logStr)
	}
}
