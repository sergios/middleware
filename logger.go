package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/urfave/negroni"
)

const (
	requestWarning    = 500 * time.Millisecond
	requestMsg        = "Completed request"
	requestWarningMsg = "Completed request alert"
)

type Logger struct {
	*logrus.Logger
	AppName string
}

func NewLogger(logLevel string, appName string) *Logger {
	logrusLevel, err := logrus.ParseLevel(logLevel)
	if err != nil {
		log.Fatal(err)
	}

	logr := &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.999",
	}
	return NewLoggerCustom(logrusLevel, logr, appName)
}

func NewLoggerCustom(level logrus.Level, formatter logrus.Formatter, appName string) *Logger {
	log := logrus.New()
	log.Level = level
	log.Formatter = formatter

	return &Logger{Logger: log, AppName: appName}
}

func (l *Logger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	start := time.Now()
	next(rw, r)
	latency := time.Since(start)

	res := rw.(negroni.ResponseWriter)
	entry := l.WithFields(logrus.Fields{
		"app":          l.AppName,
		"method":       r.Method,
		"request":      r.RequestURI,
		"status":       res.Status(),
		"size":         res.Size(),
		"took":         latency,
		"cachecontrol": rw.Header().Get("Cache-Control"),
	})

	n := res.Status()
	switch {
	case n >= 500:
		entry.Error(requestMsg)
	case n >= 400:
		entry.Warn(requestMsg)
	default:
		if latency > requestWarning {
			entry.Warn(requestWarningMsg)
		} else {
			entry.Info(requestMsg)
		}
	}
}
