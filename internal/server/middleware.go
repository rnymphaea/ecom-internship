package server

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"ecom-internship/internal/logger"
	"ecom-internship/internal/pkg/httputils"
)

func loggingMiddleware(log logger.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := generateRequestID()
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

func generateRequestID() string {
	extra := make([]byte, 2)

	//nolint:errcheck,gosec
	rand.Read(extra)

	requestID := fmt.Sprintf("%d-%s", time.Now().UnixNano(), hex.EncodeToString(extra))

	return requestID
}
