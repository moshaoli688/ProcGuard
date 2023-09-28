package utils

import (
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log"
	"os"
	"path/filepath"
)

var (
	Logger *log.Logger
)
var (
	StdoutLogger *log.Logger
	StderrLogger *log.Logger
)

func InitLogger(logDir string, maxsize int, maxage int, filename string) error {
	if err := createLogDir(logDir); err != nil {
		return err
	}

	logFilePath := filepath.Join(logDir, filename+".log")

	lumberjackLogger := &lumberjack.Logger{
		Filename:  logFilePath,
		MaxSize:   maxsize,
		MaxAge:    maxage,
		LocalTime: true,
		Compress:  true,
	}

	multiWriter := io.MultiWriter(lumberjackLogger, os.Stdout)

	Logger = log.New(multiWriter, "", log.Ldate|log.Ltime)

	return nil
}

func InitLoggers(logDir string, maxsize int, maxage int, stdoutLogFile, stderrLogFile string) error {
	if err := createLogDir(logDir); err != nil {
		return err
	}

	stdoutLogFilePath := filepath.Join(logDir, stdoutLogFile)
	stdoutLogger := &lumberjack.Logger{
		Filename:  stdoutLogFilePath,
		MaxSize:   maxsize,
		MaxAge:    maxage,
		LocalTime: true,
		Compress:  true,
	}
	stdoutMultiWriter := io.MultiWriter(stdoutLogger, os.Stdout)
	StdoutLogger = log.New(stdoutMultiWriter, "", log.Ldate|log.Ltime)

	stderrLogFilePath := filepath.Join(logDir, stderrLogFile)
	stderrLogger := &lumberjack.Logger{
		Filename:  stderrLogFilePath,
		MaxSize:   maxsize, // 日志文件最大尺寸，单位：MB
		MaxAge:    maxage,  // 日志文件最大保存天数
		LocalTime: true,
		Compress:  true, // 是否压缩旧日志文件
	}
	stderrMultiWriter := io.MultiWriter(stderrLogger, os.Stderr)
	StderrLogger = log.New(stderrMultiWriter, "", log.Ldate|log.Ltime)

	return nil
}

// createLogDir 用于创建日志目录
func createLogDir(logDir string) error {
	_, err := os.Stat(logDir)
	if os.IsNotExist(err) {
		err := os.MkdirAll(logDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create log directory: %v", err)
		}
	}
	return nil
}

// Error 用于记录错误日志
func Error(format string, v ...interface{}) {
	Logger.Printf("ERROR: "+format, v...)
}

// Info 用于记录信息日志
func Info(format string, v ...interface{}) {
	Logger.Printf("INFO: "+format, v...)
}

// Debug 用于记录调试日志
func Debug(format string, v ...interface{}) {
	Logger.Printf("DEBUG: "+format, v...)
}

// Warn 用于记录警告日志
func Warn(format string, v ...interface{}) {
	Logger.Printf("WARN: "+format, v...)
}
