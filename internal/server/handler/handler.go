package handler

import (
	"encoding/json"
	_ "errors"
	"net/http"
	"strconv"

	"ecom-internship/internal/database"
	"ecom-internship/internal/logger"
	"ecom-internship/internal/model"
)

const BASE_URL = "http://localhost:8080/todos/"

type allToDosResponse struct {
	ToDos []model.ToDo `json:"todos"`
}

func GetAllToDos(log logger.Logger, db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		toDos, err := db.GetAllToDos(r.Context())
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		response := allToDosResponse{
			ToDos: toDos,
		}

		if err = json.NewEncoder(w).Encode(response); err != nil {
			log.Error("failed to encode response", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

func GetToDoByID(log logger.Logger, db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		idFromPath := r.PathValue("id")
		id, err := strconv.Atoi(idFromPath)
		if err != nil {
			log.Error("invalid id", "id", idFromPath)
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		toDo, err := db.GetToDoByID(r.Context(), id)
		if err != nil {
			log.Error("error get todo by id", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if err = json.NewEncoder(w).Encode(toDo); err != nil {
			log.Error("failed to encode response", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

func CreateToDo(log logger.Logger, db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		var toDo model.ToDo
		if err := json.NewDecoder(r.Body).Decode(&toDo); err != nil {
			log.Error("failed to decode request", "error", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if len(toDo.Caption) == 0 {
			log.Error("empty caption")
			http.Error(w, "Empty caption provided", http.StatusBadRequest)
			return
		}

		id, err := db.CreateToDo(r.Context(), toDo)
		if err != nil {
			log.Error("error create todo", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		url := BASE_URL + strconv.Itoa(id)
		w.Header().Add("Location", url)
		w.WriteHeader(http.StatusCreated)
	}
}

type updateToDoRequest struct {
	Caption     string `json:"caption"`
	Description string `json:"description"`
	IsCompleted bool   `json:"is_completed"`
}

func UpdateToDo(log logger.Logger, db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		idFromPath := r.PathValue("id")
		id, err := strconv.Atoi(idFromPath)
		if err != nil {
			log.Error("invalid id", "id", idFromPath)
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		var update updateToDoRequest
		if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
			log.Error("failed to decode request", "error", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if len(update.Caption) == 0 {
			log.Error("empty caption")
			http.Error(w, "Empty caption provided", http.StatusBadRequest)
			return
		}

		todo := model.ToDo{
			ID:          id,
			Caption:     update.Caption,
			Description: update.Description,
			IsCompleted: update.IsCompleted,
		}

		found, err := db.UpdateToDo(r.Context(), todo)
		if err != nil {
			log.Error("failed to update todo", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if found {
			w.WriteHeader(http.StatusNoContent)
		} else {
			url := BASE_URL + strconv.Itoa(id)
			w.Header().Add("Location", url)
			w.WriteHeader(http.StatusCreated)
		}
	}
}

func DeleteToDo(log logger.Logger, db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		idFromPath := r.PathValue("id")
		id, err := strconv.Atoi(idFromPath)
		if err != nil {
			log.Error("invalid id", "id", idFromPath)
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		if err = db.DeleteToDo(r.Context(), id); err != nil {
			log.Error("failed to delete todo", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
