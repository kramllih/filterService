package logger

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func init() {
	log := logrus.New()

	log.Level = logrus.InfoLevel
	log.Formatter = &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339Nano}

	filePath, err := getWorkingPath()
	if err != nil {
		log.Fatal(err)
	}

	filename := path.Join(filePath, "log")

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		mw := io.MultiWriter(os.Stdout, file)
		log.SetOutput(mw)
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	Log = log
}

type Logger struct {
	logger *logrus.Logger
	entry  *logrus.Entry
	fields map[string]interface{}
}

func NewLogger(service string) *Logger {

	logger := Log

	log := &Logger{logger, nil, make(map[string]interface{})}

	log.fields["service"] = service

	log.entry = logger.WithFields(log.fields)

	return log
}

func (l *Logger) WithField(key string, value interface{}) *Logger {

	l.fields[key] = value

	entry := l.logger.WithFields(l.fields)

	l.entry = entry

	return l
}

func (l *Logger) WithFields(fields logrus.Fields) *Logger {

	newFields := mergeFields(l.fields, fields)
	entry := l.logger.WithFields(newFields)

	l.entry = entry
	return l
}

func (l *Logger) Info(args ...interface{}) {
	l.entry.Info(args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.entry.Infof(format, args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.entry.Warn(args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.entry.Warnf(format, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.entry.Error(args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.entry.Errorf(format, args...)
}

func (l *Logger) Writer() *io.PipeWriter {
	return l.logger.WriterLevel(logrus.ErrorLevel)
}

func getWorkingPath() (string, error) {

	fullexecpath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}

	dir, _ := filepath.Split(fullexecpath)

	return dir, nil

}

func mergeFields(ms ...logrus.Fields) logrus.Fields {
	res := logrus.Fields{}
	for _, m := range ms {
		for k, v := range m {
			res[k] = v
		}
	}
	return res
}
