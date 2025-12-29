// Package handler provides HTTP handlers for the ToDo API.
package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"ecom-internship/internal/database"
	"ecom-internship/internal/httputils"
	"ecom-internship/internal/logger"
	"ecom-internship/internal/model"
)

type apiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(apiError{
		Code:    status,
		Message: message,
	})

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

type allToDosResponse struct {
	ToDos []model.ToDo `json:"todos"`
}

// GetAllToDos returns a handler for retrieving all ToDo items.
func GetAllToDos(log logger.Logger, db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		requestID := httputils.RequestID(r)

		toDos, err := db.GetAllToDos(r.Context())
		if err != nil {
			log.Error("failed get all todos",
				"request_id", requestID,
				"error", err)
			writeError(w, http.StatusInternalServerError, "Internal server error")

			return
		}

		response := allToDosResponse{
			ToDos: toDos,
		}

		if err = json.NewEncoder(w).Encode(response); err != nil {
			log.Error("failed to encode response",
				"request_id", requestID,
				"error", err)
			writeError(w, http.StatusInternalServerError, "Internal server error")
		}
	}
}

// GetToDoByID returns a handler for retrieving a ToDo item by ID.
func GetToDoByID(log logger.Logger, db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		requestID := httputils.RequestID(r)

		idFromPath := r.PathValue("id")
		id, err := strconv.Atoi(idFromPath)
		if err != nil {
			log.Error("invalid id",
				"request_id", requestID,
				"error", err,
				"id", idFromPath)
			writeError(w, http.StatusBadRequest, "Invalid id")

			return
		}

		toDo, err := db.GetToDoByID(r.Context(), id)
		if err != nil {
			if errors.Is(err, database.ErrNotFound) {
				log.Debug("invalid id",
					"request_id", requestID,
					"error", err)
				writeError(w, http.StatusNotFound, "ToDo id not found")
			} else {
				log.Error("error get todo by id",
					"request_id", requestID,
					"error", err)

				writeError(w, http.StatusInternalServerError, "Internal server error")
			}

			return
		}

		if err = json.NewEncoder(w).Encode(toDo); err != nil {
			log.Error("failed to encode response",
				"request_id", requestID,
				"error", err)

			writeError(w, http.StatusInternalServerError, "Internal server error")
		}
	}
}

// CreateToDo returns a handler for creating a new ToDo item.
func CreateToDo(log logger.Logger, db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		requestID := httputils.RequestID(r)

		var toDo model.ToDo
		if err := json.NewDecoder(r.Body).Decode(&toDo); err != nil {
			log.Error("failed to decode request",
				"request_id", requestID,
				"error", err)
			writeError(w, http.StatusBadRequest, "Internal request body")

			return
		}

		if len(toDo.Caption) == 0 {
			log.Debug("empty caption",
				"request_id", requestID)
			writeError(w, http.StatusBadRequest, "Empty caption provided")

			return
		}

		id, err := db.CreateToDo(r.Context(), toDo)
		if err != nil {
			if errors.Is(err, database.ErrIDAlreadyExists) {
				writeError(w, http.StatusConflict, "ToDo with this ID already exists")
			} else {
				log.Error("error create todo",
					"request_id", requestID,
					"error", err)
				writeError(w, http.StatusInternalServerError, "Internal server error")
			}

			return
		}

		location := httputils.BuildLocation(r, id)
		w.Header().Add("Location", location)
		w.WriteHeader(http.StatusCreated)
	}
}

type updateToDoRequest struct {
	Caption     string `json:"caption"`
	Description string `json:"description"`
	IsCompleted bool   `json:"is_completed"`
}

// UpdateToDo returns a handler for updating an existing ToDo item.
func UpdateToDo(log logger.Logger, db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		requestID := httputils.RequestID(r)

		idFromPath := r.PathValue("id")
		id, err := strconv.Atoi(idFromPath)
		if err != nil {
			log.Error("invalid id",
				"request_id", requestID,
				"error", err,
				"id", idFromPath)
			writeError(w, http.StatusBadRequest, "Invalid id")

			return
		}

		var update updateToDoRequest
		if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
			log.Error("failed to decode request",
				"request_id", requestID,
				"error", err)
			writeError(w, http.StatusBadRequest, "Invalid request body")

			return
		}

		if len(update.Caption) == 0 {
			log.Debug("empty caption",
				"request_id", requestID)
			writeError(w, http.StatusBadRequest, "Empty caption provided")

			return
		}

		todo := model.ToDo{
			ID:          id,
			Caption:     update.Caption,
			Description: update.Description,
			IsCompleted: update.IsCompleted,
		}

		err = db.UpdateToDo(r.Context(), todo)
		if err != nil {
			if errors.Is(err, database.ErrNotFound) {
				log.Debug("invalid id",
					"request_id", requestID,
					"error", err)
				writeError(w, http.StatusNotFound, "ToDo id not found")
			} else {
				log.Error("failed to update todo",
					"request_id", requestID,
					"error", err)
				writeError(w, http.StatusInternalServerError, "Internal server error")
			}

			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// DeleteToDo returns a handler for deleting a ToDo item by ID.
func DeleteToDo(log logger.Logger, db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		requestID := httputils.RequestID(r)

		idFromPath := r.PathValue("id")
		id, err := strconv.Atoi(idFromPath)
		if err != nil {
			log.Error("invalid id",
				"request_id", requestID,
				"error", err,
				"id", idFromPath)

			writeError(w, http.StatusBadRequest, "Invalid id")

			return
		}

		if err = db.DeleteToDo(r.Context(), id); err != nil {
			if errors.Is(err, database.ErrNotFound) {
				log.Debug("invalid id",
					"request_id", requestID,
					"error", err)
				writeError(w, http.StatusNotFound, "ToDo id not found")
			} else {
				log.Error("failed to update todo",
					"request_id", requestID,
					"error", err)
				writeError(w, http.StatusInternalServerError, "Internal server error")
			}

			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
