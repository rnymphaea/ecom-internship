package handler

import (
	"encoding/json"
	"net/http"

	"ecom-internship/internal/database"
	"ecom-internship/internal/logger"
	"ecom-internship/internal/model"
)

type allToDosResponse struct {
	toDos []model.ToDo `json: "todos"`
}

func GetAllToDos(log logger.Logger, db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		toDos, err := db.GetAllToDos(r.Context())
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		response := allToDosResponse{
			toDos: toDos,
		}

		if err = json.NewEncoder(w).Encode(response); err != nil {
			log.Error("failed to encode response: %w", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}
