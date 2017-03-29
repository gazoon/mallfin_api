package tracing

import (
	"context"
	"net/http"

	"github.com/satori/go.uuid"
)

type ContextKey int

const (
	RequestIDHeader = "X-Request-ID"
	requestIDCtxKey = ContextKey(228)
)

func NewRequestID() string {
	return uuid.NewV4().String()
}

func NewContext(ctx context.Context, reqID string) context.Context {
	return context.WithValue(ctx, requestIDCtxKey, reqID)
}

func FromContext(ctx context.Context) string {
	reqID, _ := ctx.Value(requestIDCtxKey).(string)
	return reqID
}

func Middleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	requestID := r.Header.Get(RequestIDHeader)
	if requestID == "" {
		requestID = NewRequestID()
		r.Header.Set(RequestIDHeader, requestID)
	}
	ctx := NewContext(r.Context(), requestID)
	w.Header().Set(RequestIDHeader, requestID)
	next(w, r.WithContext(ctx))
}
