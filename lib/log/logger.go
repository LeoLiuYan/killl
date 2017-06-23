package log

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"os"
	"runtime"
	"strings"
)

const (
	SKIP       = 2
	InfoLevel  = logrus.InfoLevel
	DebugLevel = logrus.DebugLevel
	FatalLevel = logrus.FatalLevel
	ErrorLevel = logrus.ErrorLevel
)

var logger *logrus.Logger = logrus.New()

func InitLogger(verbose bool) {
	logger.Out = os.Stdout
	if verbose {
		logger.Level = DebugLevel
		return
	}
	logger.Level = InfoLevel
	return
}

func Info(args ...interface{}) {
	if logger.Level >= InfoLevel {
		item := logger.WithFields(logrus.Fields{})
		item.Data["__file__"] = lineInfo(SKIP)
		item.Info(args...)
	}
}

func Infof(format string, args ...interface{}) {
	if logger.Level >= InfoLevel {
		item := logger.WithFields(logrus.Fields{})
		item.Data["_file__"] = lineInfo(SKIP)
		item.Infof(format, args...)
	}
}

func Debug(args ...interface{}) {
	if logger.Level >= DebugLevel {
		item := logger.WithFields(logrus.Fields{})
		item.Data["__file__"] = lineInfo(SKIP)
		item.Debug(args...)
	}
}

func Debugf(format string, args ...interface{}) {
	if logger.Level >= DebugLevel {
		item := logger.WithFields(logrus.Fields{})
		item.Data["__file__"] = lineInfo(SKIP)
		item.Debugf(format, args...)
	}
}

func Fatal(args ...interface{}) {
	if logger.Level >= FatalLevel {
		item := logger.WithFields(logrus.Fields{})
		item.Data["__file__"] = lineInfo(SKIP)
		item.Fatal(args...)
	}
}

func Fatalf(format string, args ...interface{}) {
	if logger.Level >= FatalLevel {
		item := logger.WithFields(logrus.Fields{})
		item.Data["__file__"] = lineInfo(SKIP)
		item.Fatalf(format, args...)
	}
}

func Error(args ...interface{}) {
	if logger.Level >= ErrorLevel {
		item := logger.WithFields(logrus.Fields{})
		item.Data["__file__"] = lineInfo(SKIP)
		item.Error(args...)
	}
}

func Errorf(format string, args ...interface{}) {
	if logger.Level >= ErrorLevel {
		item := logger.WithFields(logrus.Fields{})
		item.Data["__file__"] = lineInfo(SKIP)
		item.Errorf(format, args...)
	}
}

func lineInfo(skip int) (res string) {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return fmt.Sprintf("%s:%d", "?", 0)
	}
	idx := strings.LastIndex(file, "/")
	if idx >= 0 {
		file = file[idx+1:]
	}
	return fmt.Sprintf("%s:%d", file, line)
}
