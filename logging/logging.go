package logging

import (
	"mallfin_api/config"
	"strings"

	"context"
	"mallfin_api/tracing"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/urfave/negroni"
)

type ContextKey int

const (
	ServiceNameField = "service_name"
	ServerIDField    = "server_id"
	RequestIDField   = "request_id"
	PackageField     = "package"

	loggerCtxKey = ContextKey(1)
)

type customFormatter struct {
	logFormatter     log.Formatter
	additionalFields log.Fields
}

func WithPackage(packageName string) *log.Entry {
	return log.WithField(PackageField, packageName)
}

func (cf *customFormatter) Format(e *log.Entry) ([]byte, error) {
	for field, value := range cf.additionalFields {
		e.Data[field] = value
	}
	return cf.logFormatter.Format(e)
}

func FromContext(ctx context.Context) *log.Entry {
	logger, ok := ctx.Value(loggerCtxKey).(*log.Entry)
	if !ok {
		logger = log.NewEntry(log.StandardLogger())
	}
	return logger
}

func NewContext(ctx context.Context, logger *log.Entry) context.Context {
	return context.WithValue(ctx, loggerCtxKey, logger)
}

func Initialization() {
	var logLevel log.Level
	switch strings.ToLower(config.LogLevel()) {
	case "debug":
		logLevel = log.DebugLevel
	case "info":
		logLevel = log.InfoLevel
	case "warning":
		logLevel = log.WarnLevel
	case "error":
		logLevel = log.ErrorLevel
	default:
		logLevel = log.DebugLevel
	}
	formatter := &customFormatter{logFormatter: &log.TextFormatter{}, additionalFields: log.Fields{
		ServiceNameField: config.ServiceName(),
		ServerIDField:    config.ServerID(),
	}}
	log.SetLevel(logLevel)
	log.SetFormatter(formatter)
}

func Middleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ctx := r.Context()
	requestID := tracing.FromContext(ctx)
	logger := log.WithFields(log.Fields{RequestIDField: requestID})
	ctx = NewContext(ctx, logger)

	logger = logger.WithFields(log.Fields{"path": r.URL.Path, "method": r.Method})

	userIP := r.Header.Get("X-Real-IP")
	if userIP == "" {
		userIP = r.RemoteAddr
	}
	logger.WithFields(log.Fields{"user_ip": userIP, "user_agent": r.UserAgent()}).Debug("Request started")

	next(w, r.WithContext(ctx))

	res := w.(negroni.ResponseWriter)
	logger.WithField("status", res.Status()).Info("Request finished")
}
