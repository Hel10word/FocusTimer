package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// LogLevel 定义日志级别
type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

// Logger 提供日志功能
type Logger struct {
	level  LogLevel
	logger *log.Logger
	file   *os.File
}

// New 创建新的日志器
func New(levelStr string) *Logger {
	var level LogLevel
	switch levelStr {
	case "debug":
		level = LevelDebug
	case "info":
		level = LevelInfo
	case "warn":
		level = LevelWarn
	case "error":
		level = LevelError
	case "fatal":
		level = LevelFatal
	default:
		level = LevelInfo
	}

	// 创建日志文件
	logDir := "logs"
	os.MkdirAll(logDir, 0755)

	logFile := filepath.Join(logDir, fmt.Sprintf("FocusTimer%s.log",
		time.Now().Format("2025-05-21")))

	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Failed to open log file: %v", err)
		file = nil
	}

	logger := log.New(os.Stdout, "", log.LstdFlags)

	return &Logger{
		level:  level,
		logger: logger,
		file:   file,
	}
}

// log 记录一条日志
func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	var levelStr string
	switch level {
	case LevelDebug:
		levelStr = "DEBUG"
	case LevelInfo:
		levelStr = "INFO"
	case LevelWarn:
		levelStr = "WARN"
	case LevelError:
		levelStr = "ERROR"
	case LevelFatal:
		levelStr = "FATAL"
	}

	msg := fmt.Sprintf(format, args...)
	logMsg := fmt.Sprintf("[%s] %s", levelStr, msg)

	// 输出到控制台
	l.logger.Println(logMsg)

	// 输出到文件
	if l.file != nil {
		fmt.Fprintln(l.file, logMsg)
	}

	// 如果是致命错误，程序终止
	if level == LevelFatal {
		os.Exit(1)
	}
}

// Debug 记录调试日志
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(LevelDebug, format, args...)
}

// Info 记录信息日志
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(LevelInfo, format, args...)
}

// Warn 记录警告日志
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(LevelWarn, format, args...)
}

// Error 记录错误日志
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(LevelError, format, args...)
}

// Fatal 记录致命错误日志并退出程序
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(LevelFatal, format, args...)
}

// Close 关闭日志文件
func (l *Logger) Close() {
	if l.file != nil {
		l.file.Close()
	}
}
