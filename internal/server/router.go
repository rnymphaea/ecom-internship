package server

import (
	"net/http"

	"ecom-internship/internal/database"
	"ecom-internship/internal/logger"
	"ecom-internship/internal/server/handler"
)

// NewRouter creates and configures the HTTP router with middleware.
func NewRouter(log logger.Logger, db database.Database) *http.ServeMux {
	mux := http.NewServeMux()

	middlewares := []func(logger.Logger, http.Handler) http.Handler{
		panicRecoveryMiddleware,
		loggingMiddleware,
	}

	mux.Handle("GET /todos", chain(log, handler.GetAllToDos(log, db), middlewares...))
	mux.Handle("GET /todos/{id}", chain(log, handler.GetToDoByID(log, db), middlewares...))

	mux.Handle("POST /todos", chain(log, handler.CreateToDo(log, db), middlewares...))

	mux.Handle("PUT /todos/{id}", chain(log, handler.UpdateToDo(log, db), middlewares...))

	mux.Handle("DELETE /todos/{id}", chain(log, handler.DeleteToDo(log, db), middlewares...))

	return mux
}
