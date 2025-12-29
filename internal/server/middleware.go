package server

import (
	"net/http"
	"time"

	"ecom-internship/internal/httputils"
	"ecom-internship/internal/logger"
)

func loggingMiddleware(log logger.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := httputils.GenerateRequestID()
		ctx := httputils.WithRequestID(r.Context(), requestID)
		r = r.WithContext(ctx)

		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)

		log.Info("request completed",
			"request_id", requestID,
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
			"duration", duration.String(),
		)
	})
}

func panicRecoveryMiddleware(log logger.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				requestID := httputils.RequestID(r)

				log.Error("recovered from panic",
					"request_id", requestID,
					"error", err,
					"path", r.URL.Path,
					"method", r.Method,
					"remote_addr", r.RemoteAddr,
				)

				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func chain(log logger.Logger,
	h http.Handler,
	middlewares ...func(logger.Logger, http.Handler) http.Handler,
) http.Handler {
	for _, m := range middlewares {
		h = m(log, h)
	}

	return h
}
