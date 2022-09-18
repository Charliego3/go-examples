package main

import (
	"fmt"
	"github.com/whimthen/kits/logger"
)
import "time"

// import "david/util"

var StdLogger = &Logger{
	level: -1,
}

type Logger struct {
	level int // 1个level对应4个空格
}

// LevelUp
func (l *Logger) LevelUp() {
	l.level += 1
}

// LevelDown
func (l *Logger) LevelDown() {
	l.level -= 1
	if l.level < 0 {
		fmt.Errorf("l.level should not less than 0 %d\n", l.level)
	}
}

// Debug print debug level
func (l *Logger) Debug(format string, a ...interface{}) {
	logger.Debug(format, a...)
}

// Info print debug level
func (l *Logger) Info(format string, a ...interface{}) {
	logger.Info(format, a...)
}

// Warn print debug level
func (l *Logger) Warn(format string, a ...interface{}) {
	logger.Warn(format, a...)
}

// Error print debug level
func (l *Logger) Error(format string, a ...interface{}) {
	logger.Error(format, a...)
}

// nowTime hh:mm:ss
func NowTime() string {
	now := time.Now()
	return fmt.Sprintf("%02d:%02d:%02d", now.Hour(), now.Minute(), now.Second())
}
