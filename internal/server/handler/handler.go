package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"ecom-internship/internal/database"
	"ecom-internship/internal/logger"
	"ecom-internship/internal/model"
)

type allToDosResponse struct {
	ToDos []model.ToDo `json:"todos"`
}

func GetAllToDos(log logger.Logger, db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		requestID := getRequestID(r)

		toDos, err := db.GetAllToDos(r.Context())
		if err != nil {
			log.Error("failed get all todos",
				"request_id", requestID,
				"error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		response := allToDosResponse{
			ToDos: toDos,
		}

		if err = json.NewEncoder(w).Encode(response); err != nil {
			log.Error("failed to encode response",
				"request_id", requestID,
				"error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

func GetToDoByID(log logger.Logger, db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		requestID := getRequestID(r)

		idFromPath := r.PathValue("id")
		id, err := strconv.Atoi(idFromPath)
		if err != nil {
			log.Error("invalid id",
				"request_id", requestID,
				"error", err,
				"id", idFromPath)
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		toDo, err := db.GetToDoByID(r.Context(), id)
		if err != nil {
			if errors.Is(err, database.ErrNotFound) {
				log.Error("invalid id",
					"request_id", requestID,
					"error", err)
				w.WriteHeader(http.StatusNotFound)
			} else {
				log.Error("error get todo by id",
					"request_id", requestID,
					"error", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		if err = json.NewEncoder(w).Encode(toDo); err != nil {
			log.Error("failed to encode response",
				"request_id", requestID,
				"error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

func CreateToDo(log logger.Logger, db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		requestID := getRequestID(r)

		var toDo model.ToDo
		if err := json.NewDecoder(r.Body).Decode(&toDo); err != nil {
			log.Error("failed to decode request",
				"request_id", requestID,
				"error", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if len(toDo.Caption) == 0 {
			log.Error("empty caption",
				"request_id", requestID)
			http.Error(w, "Empty caption provided", http.StatusBadRequest)
			return
		}

		id, err := db.CreateToDo(r.Context(), toDo)
		if err != nil {
			if errors.Is(err, database.ErrIDAlreadyExists) {
				http.Error(w, "ToDo with this ID already exists", http.StatusConflict)
			} else {
				log.Error("error create todo",
					"request_id", requestID,
					"error", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		location := buildLocation(r, id)
		w.Header().Add("Location", location)
		w.WriteHeader(http.StatusCreated)
	}
}

type upsertToDoRequest struct {
	Caption     string `json:"caption"`
	Description string `json:"description"`
	IsCompleted bool   `json:"is_completed"`
}

func UpsertToDo(log logger.Logger, db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		requestID := getRequestID(r)

		idFromPath := r.PathValue("id")
		id, err := strconv.Atoi(idFromPath)
		if err != nil {
			log.Error("invalid id",
				"request_id", requestID,
				"error", err,
				"id", idFromPath)
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		var update upsertToDoRequest
		if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
			log.Error("failed to decode request",
				"request_id", requestID,
				"error", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if len(update.Caption) == 0 {
			log.Error("empty caption",
				"request_id", requestID)
			http.Error(w, "Empty caption provided", http.StatusBadRequest)
			return
		}

		todo := model.ToDo{
			ID:          id,
			Caption:     update.Caption,
			Description: update.Description,
			IsCompleted: update.IsCompleted,
		}

		found, err := db.UpsertToDo(r.Context(), todo)
		if err != nil {
			log.Error("failed to update todo",
				"request_id", requestID,
				"error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if found {
			w.WriteHeader(http.StatusNoContent)
		} else {
			location := buildLocation(r, id)
			w.Header().Add("Location", location)
			w.WriteHeader(http.StatusCreated)
		}
	}
}

func DeleteToDo(log logger.Logger, db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		requestID := getRequestID(r)

		idFromPath := r.PathValue("id")
		id, err := strconv.Atoi(idFromPath)
		if err != nil {
			log.Error("invalid id",
				"request_id", requestID,
				"error", err,
				"id", idFromPath)
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		if err = db.DeleteToDo(r.Context(), id); err != nil {
			if errors.Is(err, database.ErrNotFound) {
				log.Error("invalid id",
					"request_id", requestID,
					"error", err)
				w.WriteHeader(http.StatusNotFound)
			} else {
				log.Error("failed to update todo",
					"request_id", requestID,
					"error", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
