// Package httputils provides HTTP-related utils.
package httputils

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"
)

// KeyType represents the type of context key.
type KeyType int

const (
	// RequestIDKey is used to store request ID in context.
	// Using custom type to avoid collisions with other packages.
	RequestIDKey KeyType = iota
)

// RequestID extracts the request ID from the context.
func RequestID(r *http.Request) string {
	if id, ok := r.Context().Value(RequestIDKey).(string); ok {
		return id
	}

	return ""
}

// WithRequestID adds a request ID to the context.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, RequestIDKey, id)
}

// BuildLocation creates a URL for a newly created resource.
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

// GenerateRequestID generates a request ID based on time and random bytes to avoid dublication.
func GenerateRequestID() string {
	extra := make([]byte, 2)

	//nolint:errcheck,gosec
	rand.Read(extra)

	requestID := fmt.Sprintf("%d-%s", time.Now().UnixNano(), hex.EncodeToString(extra))

	return requestID
}
