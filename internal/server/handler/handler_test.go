package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"ecom-internship/internal/database"
	"ecom-internship/internal/logger/std"
	"ecom-internship/internal/model"
)

type mockDB struct {
	todos     map[int]model.ToDo
	nextID    int
	shouldErr bool
}

var ErrDb = errors.New("database error")

//nolint:revive
func (m *mockDB) GetAllToDos(ctx context.Context) ([]model.ToDo, error) {
	if m.shouldErr {
		return nil, ErrDb
	}
	todos := make([]model.ToDo, 0, len(m.todos))
	for _, todo := range m.todos {
		todos = append(todos, todo)
	}

	return todos, nil
}

//nolint:revive
func (m *mockDB) GetToDoByID(ctx context.Context, id int) (model.ToDo, error) {
	if m.shouldErr {
		return model.ToDo{}, ErrDb
	}
	todo, ok := m.todos[id]
	if !ok {
		return model.ToDo{}, database.ErrNotFound
	}

	return todo, nil
}

//nolint:revive
func (m *mockDB) CreateToDo(ctx context.Context, todo model.ToDo) (int, error) {
	if m.shouldErr {
		return 0, ErrDb
	}
	if todo.ID != 0 {
		if _, exists := m.todos[todo.ID]; exists {
			return 0, database.ErrIDAlreadyExists
		}
	} else {
		m.nextID++
		todo.ID = m.nextID
	}
	m.todos[todo.ID] = todo

	return todo.ID, nil
}

//nolint:revive
func (m *mockDB) UpdateToDo(ctx context.Context, todo model.ToDo) error {
	if m.shouldErr {
		return ErrDb
	}
	if _, exists := m.todos[todo.ID]; !exists {
		return database.ErrNotFound
	}
	m.todos[todo.ID] = todo

	return nil
}

//nolint:revive
func (m *mockDB) DeleteToDo(ctx context.Context, id int) error {
	if m.shouldErr {
		return ErrDb
	}
	if _, exists := m.todos[id]; !exists {
		return database.ErrNotFound
	}
	delete(m.todos, id)

	return nil
}

func TestGetAllToDos(t *testing.T) {
	logger := std.New("debug")
	db := &mockDB{
		todos: map[int]model.ToDo{
			1: {ID: 1, Caption: "Todo 1"},
			2: {ID: 2, Caption: "Todo 2"},
		},
	}

	handler := GetAllToDos(logger, db)
	req := httptest.NewRequest(http.MethodGet, "/todos", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response allToDosResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(response.ToDos) != 2 {
		t.Errorf("Expected 2 todos, got %d", len(response.ToDos))
	}

	db.shouldErr = true
	w = httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500 on error, got %d", w.Code)
	}
}

func TestGetToDoByID(t *testing.T) {
	logger := std.New("debug")
	db := &mockDB{
		todos: map[int]model.ToDo{
			1: {ID: 1, Caption: "Todo 1"},
		},
	}

	handler := GetToDoByID(logger, db)

	req := httptest.NewRequest(http.MethodGet, "/todos/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var todo model.ToDo
	if err := json.Unmarshal(w.Body.Bytes(), &todo); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if todo.Caption != "Todo 1" {
		t.Errorf("Expected caption 'Todo 1', got %s", todo.Caption)
	}

	req = httptest.NewRequest(http.MethodGet, "/todos/abc", nil)
	req.SetPathValue("id", "abc")
	w = httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid ID, got %d", w.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/todos/999", nil)
	req.SetPathValue("id", "999")
	w = httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 for non-existent ID, got %d", w.Code)
	}

	db.shouldErr = true
	req = httptest.NewRequest(http.MethodGet, "/todos/1", nil)
	req.SetPathValue("id", "1")
	w = httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500 on database error, got %d", w.Code)
	}
}

//nolint:funlen,cyclop
func TestCreateToDo(t *testing.T) {
	logger := std.New("debug")
	db := &mockDB{todos: make(map[int]model.ToDo)}

	handler := CreateToDo(logger, db)

	todo := model.ToDo{
		Caption:     "New Todo",
		Description: "New Description",
	}
	body, err := json.Marshal(todo)
	if err != nil {
		t.Fatalf("Failed to marshal todo: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}
	if location := w.Header().Get("Location"); location == "" {
		t.Error("Expected Location header")
	}

	todo.Caption = ""
	body, err = json.Marshal(todo)
	if err != nil {
		t.Fatalf("Failed to marshal todo: %v", err)
	}

	req = httptest.NewRequest(http.MethodPost, "/todos", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for empty caption, got %d", w.Code)
	}

	req = httptest.NewRequest(http.MethodPost, "/todos", bytes.NewReader([]byte("{invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid JSON, got %d", w.Code)
	}

	db.todos[5] = model.ToDo{ID: 5, Caption: "Existing"}
	todo = model.ToDo{ID: 5, Caption: "Duplicate"}
	body, err = json.Marshal(todo)
	if err != nil {
		t.Fatalf("Failed to marshal todo: %v", err)
	}

	req = httptest.NewRequest(http.MethodPost, "/todos", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("Expected status 409 for duplicate ID, got %d", w.Code)
	}

	db.shouldErr = true
	todo = model.ToDo{Caption: "Should Fail"}
	body, err = json.Marshal(todo)
	if err != nil {
		t.Fatalf("Failed to marshal todo: %v", err)
	}

	req = httptest.NewRequest(http.MethodPost, "/todos", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500 on database error, got %d", w.Code)
	}
}

//nolint:funlen,cyclop
func TestUpdateToDo(t *testing.T) {
	logger := std.New("debug")
	db := &mockDB{
		todos: map[int]model.ToDo{
			1: {ID: 1, Caption: "Original"},
		},
	}

	handler := UpdateToDo(logger, db)

	update := updateToDoRequest{
		Caption:     "Updated",
		Description: "Updated Description",
		IsCompleted: true,
	}
	body, err := json.Marshal(update)
	if err != nil {
		t.Fatalf("Failed to marshal update: %v", err)
	}

	req := httptest.NewRequest(http.MethodPut, "/todos/1", bytes.NewReader(body))
	req.SetPathValue("id", "1")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204 for update, got %d", w.Code)
	}

	updatedTodo, err := db.GetToDoByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("Failed to get updated todo: %v", err)
	}
	if updatedTodo.Caption != "Updated" {
		t.Errorf("Expected caption 'Updated', got %s", updatedTodo.Caption)
	}

	update.Caption = "New via PUT"
	body, err = json.Marshal(update)
	if err != nil {
		t.Fatalf("Failed to marshal update: %v", err)
	}

	req = httptest.NewRequest(http.MethodPut, "/todos/999", bytes.NewReader(body))
	req.SetPathValue("id", "999")
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 for non-existent todo, got %d", w.Code)
	}

	update.Caption = ""
	body, err = json.Marshal(update)
	if err != nil {
		t.Fatalf("Failed to marshal update: %v", err)
	}

	req = httptest.NewRequest(http.MethodPut, "/todos/1", bytes.NewReader(body))
	req.SetPathValue("id", "1")
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for empty caption, got %d", w.Code)
	}

	req = httptest.NewRequest(http.MethodPut, "/todos/abc", bytes.NewReader(body))
	req.SetPathValue("id", "abc")
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid ID, got %d", w.Code)
	}

	req = httptest.NewRequest(http.MethodPut, "/todos/1", bytes.NewReader([]byte("{invalid json")))
	req.SetPathValue("id", "1")
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid JSON, got %d", w.Code)
	}

	db.shouldErr = true
	update.Caption = "Should Fail"
	body, err = json.Marshal(update)
	if err != nil {
		t.Fatalf("Failed to marshal update: %v", err)
	}

	req = httptest.NewRequest(http.MethodPut, "/todos/1", bytes.NewReader(body))
	req.SetPathValue("id", "1")
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500 on database error, got %d", w.Code)
	}
}

func TestDeleteToDo(t *testing.T) {
	logger := std.New("debug")
	db := &mockDB{
		todos: map[int]model.ToDo{
			1: {ID: 1, Caption: "To Delete"},
		},
	}

	handler := DeleteToDo(logger, db)

	req := httptest.NewRequest(http.MethodDelete, "/todos/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}

	req = httptest.NewRequest(http.MethodDelete, "/todos/999", nil)
	req.SetPathValue("id", "999")
	w = httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 for non-existent todo, got %d", w.Code)
	}

	req = httptest.NewRequest(http.MethodDelete, "/todos/abc", nil)
	req.SetPathValue("id", "abc")
	w = httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid ID, got %d", w.Code)
	}

	db.shouldErr = true
	db.todos[2] = model.ToDo{ID: 2, Caption: "Should Fail"}
	req = httptest.NewRequest(http.MethodDelete, "/todos/2", nil)
	req.SetPathValue("id", "2")
	w = httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500 on database error, got %d", w.Code)
	}
}

func TestResponseStructures(t *testing.T) {
	response := allToDosResponse{
		ToDos: []model.ToDo{{ID: 1, Caption: "Test"}},
	}

	body, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	var decoded map[string]any
	if err := json.Unmarshal(body, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if _, ok := decoded["todos"]; !ok {
		t.Error("Response should have 'todos' field")
	}
}
