package logutil

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

var (
	defaultFormat = &TextFormatter{
		TextFormatter: &logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05.000"}}
	defaultLogger = &logrus.Logger{
		Out:       os.Stderr,
		Level:     logrus.InfoLevel,
		Formatter: defaultFormat,
	}
)

type TextFormatter struct {
	*logrus.TextFormatter
}

func (f *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	keys := make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		keys = append(keys, k)
	}

	if !f.DisableSorting {
		sort.Strings(keys)
	}
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = time.RFC3339
	}

	levelText := strings.ToUpper(entry.Level.String())[0:4]
	fmt.Fprintf(b, "%s[%s] %-44s ", levelText, entry.Time.Format(timestampFormat), entry.Message)
	for _, k := range keys {
		v := entry.Data[k]
		fmt.Fprintf(b, " %s", k)
		stringVal, ok := v.(string)
		if !ok {
			stringVal = fmt.Sprint(v)
		}

		if !f.needsQuoting(stringVal) {
			b.WriteString(stringVal)
		} else {
			b.WriteString(fmt.Sprintf("%q", stringVal))
		}
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *TextFormatter) needsQuoting(text string) bool {
	if f.QuoteEmptyFields && len(text) == 0 {
		return true
	}
	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.' || ch == '_' || ch == '/' || ch == '@' || ch == '^' || ch == '+') {
			return true
		}
	}
	return false
}

type Logger interface {
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Fatal(args ...interface{})
	Printf(format string, v ...interface{})
}

func DefaultLogger() Logger {
	return defaultLogger
}

func SetLogFile(filename string) error {
	var opts []rotatelogs.Option
	opts = append(opts, rotatelogs.WithMaxAge(time.Hour*24*3))
	opts = append(opts, rotatelogs.WithRotationTime(time.Hour*6))
	if runtime.GOOS == "linux" {
		opts = append(opts, rotatelogs.WithLinkName(filename))
	}

	writer, err := rotatelogs.New(
		filename+"_%Y%m%d%H%M.logutil",
		opts...,
	)
	if err == nil {
		defaultLogger.SetOutput(writer)
	}
	return nil
}

func SetLogLevel(l logrus.Level) {
	defaultLogger.SetLevel(l)
}

func Printf(format string, v ...interface{}) {
	defaultLogger.Printf(format, v...)
}

func LogWarn(args ...interface{}) {
	defaultLogger.Warn(args...)
}

func Warn(args ...interface{}) {
	defaultLogger.Warn(args...)
}

func LogWarnf(format string, args ...interface{}) {
	defaultLogger.Warnf(format, args...)
}

func Warnf(format string, args ...interface{}) {
	defaultLogger.Warnf(format, args...)
}

func LogInfo(args ...interface{}) {
	defaultLogger.Info(args...)
}

func Info(args ...interface{}) {
	defaultLogger.Info(args...)
}

func LogInfof(format string, args ...interface{}) {
	defaultLogger.Infof(format, args...)
}

func Infof(format string, args ...interface{}) {
	defaultLogger.Infof(format, args...)
}

func LogError(args ...interface{}) {
	defaultLogger.Error(args...)
}

func Error(args ...interface{}) {
	defaultLogger.Error(args...)
}

func LogErrorf(format string, args ...interface{}) {
	defaultLogger.Errorf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	defaultLogger.Errorf(format, args...)
}

func LogDebug(args ...interface{}) {
	defaultLogger.Debug(args...)
}

func Debug(args ...interface{}) {
	defaultLogger.Debug(args...)
}

func LogDebugf(format string, args ...interface{}) {
	defaultLogger.Debugf(format, args...)
}

func Debugf(format string, args ...interface{}) {
	defaultLogger.Debugf(format, args...)
}

func LogFatal(args ...interface{}) {
	defaultLogger.Fatal(args...)
}

func Fatal(args ...interface{}) {
	defaultLogger.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	defaultLogger.Fatalf(format, args...)
}
