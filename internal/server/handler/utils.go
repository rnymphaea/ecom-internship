package handler

import (
	"net/http"
)

func getRequestID(r *http.Request) string {
	if id, ok := r.Context().Value("request_id").(string); ok {
		return id
	}
	return ""
}
