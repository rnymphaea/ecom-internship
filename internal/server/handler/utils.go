package handler

import (
	"fmt"
	"net/http"
)

func getRequestID(r *http.Request) string {
	if id, ok := r.Context().Value("request_id").(string); ok {
		return id
	}

	return ""
}

func buildLocation(r *http.Request, id int) string {
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
