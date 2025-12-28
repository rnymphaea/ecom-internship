package server

import (
	"net/http"

	"ecom-internship/internal/database"
	"ecom-internship/internal/logger"
	"ecom-internship/internal/server/handler"
)

func NewRouter(log logger.Logger, db database.Database) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("GET /todos", loggingMiddleware(log, handler.GetAllToDos(log, db)))
	mux.Handle("GET /todos/{id}", loggingMiddleware(log, handler.GetToDoByID(log, db)))

	mux.Handle("POST /todos", loggingMiddleware(log, handler.CreateToDo(log, db)))

	mux.Handle("PUT /todos/{id}", loggingMiddleware(log, handler.UpdateToDo(log, db)))

	mux.Handle("DELETE /todos/{id}", loggingMiddleware(log, handler.DeleteToDo(log, db)))

	return mux
}
