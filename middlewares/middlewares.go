package middlewares

import (
	"mallfin_api/logging"
	"mallfin_api/tracing"
	"net/http"
	"runtime/debug"

	log "github.com/Sirupsen/logrus"
	"github.com/urfave/negroni"
)

func RecoveryMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer func() {
		if err := recover(); err != nil {
			logger := logging.FromContext(r.Context())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal server error"))
			logger.Errorf("Panic recovered: %s", err)
			debug.PrintStack()
		}
	}()
	next(w, r)
}

func TracingMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	requestID := tracing.InitializeHeaders(w, r)
	ctx := tracing.NewContext(r.Context(), requestID)
	next(w, r.WithContext(ctx))
}

func LoggerMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ctx := r.Context()
	requestID := tracing.FromContext(ctx)
	logger := log.WithFields(log.Fields{logging.RequestIDField: requestID})
	ctx = logging.NewContext(ctx, logger)

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
