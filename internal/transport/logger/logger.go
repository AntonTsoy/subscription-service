package logger

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		reqID := uuid.New().String()
		ctx := context.WithValue(r.Context(), "ReqID", reqID)
		r = r.WithContext(ctx)

		next.ServeHTTP(lrw, r)

		duration := time.Since(start)

		log.Printf("RequestID=%s Method=%s Path=%s Status=%d Duration=%s",
			reqID, r.Method, r.URL.String(), lrw.statusCode, duration)
	})
}
