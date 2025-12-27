package server

import (
	"net/http"

	"ecom-internship/internal/database"
	"ecom-internship/internal/logger"
	"ecom-internship/internal/server/handler"
)

func NewRouter(log logger.Logger, db database.Database) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /todos", handler.GetAllToDos(log, db))
	mux.HandleFunc("GET /todos/{id}", handler.GetToDoByID(log, db))

	mux.HandleFunc("POST /todos", handler.CreateToDo(log, db))

	mux.HandleFunc("PUT /todos/{id}", handler.UpdateToDo(log, db))

	mux.HandleFunc("DELETE /todos/{id}", handler.DeleteToDo(log, db))

	return mux
}
