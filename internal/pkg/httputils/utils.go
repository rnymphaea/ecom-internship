package httputils

import (
	"context"
	"fmt"
	"net/http"
)

type KeyType int

const (
	RequestIDKey KeyType = iota
)

func RequestID(r *http.Request) string {
	if id, ok := r.Context().Value(RequestIDKey).(string); ok {
		return id
	}

	return ""
}

func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, RequestIDKey, id)
}

func BuildLocation(r *http.Request, id int) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	return fmt.Sprintf(
		"%s://%s%s/%d",
		scheme,
		r.Host,
		r.URL.Path,
		id,
	)
}
