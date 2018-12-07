package gobatis

import (
	"testing"
)

func TestInfo(t *testing.T) {
	SetLevel(LogLevelDebug)
	Info("test info -> level debug")

	SetLevel(LogLevelOff)
	Info("test info -> level off")
}

func TestFileLogger_Info(t *testing.T)  {
	logger := NewFileLog("d:/logs/nohup.out", LogLevelDebug)
	logger.Info("test file logger info")
}