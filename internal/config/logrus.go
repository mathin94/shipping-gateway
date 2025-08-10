package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type CustomConsoleFormatter struct {
	Identifier string
}

func (f *CustomConsoleFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := entry.Time.Format(time.RFC3339)
	level := strings.ToUpper(entry.Level.String()[:1])
	file, line := getCaller(8)
	msg := entry.Message

	var fields []string
	for k, v := range entry.Data {
		if k == "traceId" {
			continue // Skip traceId from fields
		}
		fields = append(fields, fmt.Sprintf("%s=%v", k, v))
	}
	fieldStr := ""
	if len(fields) > 0 {
		fieldStr = "[" + strings.Join(fields, ", ") + "] "
	}

	traceID := ""
	if id, ok := entry.Data["traceId"]; ok {
		traceID = fmt.Sprintf("[%v] ", id)
	}

	return []byte(fmt.Sprintf("%s [%s] [%s:%d] %s%s%s\n", timestamp, level, file, line, traceID, fieldStr, msg)), nil
}

func getCaller(depth int) (string, int) {
	pc := make([]uintptr, depth+10)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	for {
		frame, more := frames.Next()
		if !strings.Contains(frame.File, "sirupsen/logrus") &&
			!strings.Contains(frame.File, "config/logrus.go") &&
			!strings.Contains(frame.File, "config/logger.go") {
			return trimFilePath(frame.File, 2), frame.Line
		}
		if !more {
			break
		}
	}
	return "???", 0
}

func trimFilePath(path string, levels int) string {
	parts := strings.Split(filepath.ToSlash(path), "/")
	if len(parts) < levels {
		return path
	}
	return strings.Join(parts[len(parts)-levels:], "/")
}

func NewLogger(v *viper.Viper) *logrus.Logger {
	level := v.GetString("log.level")
	filePath := v.GetString("log.file_path")
	maxSize := v.GetInt("log.max_size")
	maxBackups := v.GetInt("log.max_backups")
	maxAge := v.GetInt("log.max_age")
	consoleOutput := v.GetBool("log.console_enabled")

	logger := logrus.New()

	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	logger.SetLevel(logLevel)

	// Setup file writer
	fileWriter := &lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   true,
	}

	if consoleOutput {
		// Write to both console and file
		logger.SetOutput(io.MultiWriter(os.Stdout, fileWriter))
		logger.SetFormatter(&CustomConsoleFormatter{})
	} else {
		// Only write to file
		logger.SetOutput(fileWriter)
		logger.SetFormatter(&logrus.JSONFormatter{DisableHTMLEscape: true, TimestampFormat: time.RFC3339})
	}

	return logger
}
