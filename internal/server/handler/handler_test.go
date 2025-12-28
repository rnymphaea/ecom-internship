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

func (m *mockDB) GetAllToDos(ctx context.Context) ([]model.ToDo, error) {
	if m.shouldErr {
		return nil, errors.New("database error")
	}
	todos := make([]model.ToDo, 0, len(m.todos))
	for _, todo := range m.todos {
		todos = append(todos, todo)
	}
	return todos, nil
}

func (m *mockDB) GetToDoByID(ctx context.Context, id int) (model.ToDo, error) {
	if m.shouldErr {
		return model.ToDo{}, errors.New("database error")
	}
	todo, ok := m.todos[id]
	if !ok {
		return model.ToDo{}, database.ErrNotFound
	}
	return todo, nil
}

func (m *mockDB) CreateToDo(ctx context.Context, todo model.ToDo) (int, error) {
	if m.shouldErr {
		return 0, errors.New("database error")
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

func (m *mockDB) UpsertToDo(ctx context.Context, todo model.ToDo) (bool, error) {
	if m.shouldErr {
		return false, errors.New("database error")
	}
	_, exists := m.todos[todo.ID]
	m.todos[todo.ID] = todo
	return exists, nil
}

func (m *mockDB) DeleteToDo(ctx context.Context, id int) error {
	if m.shouldErr {
		return errors.New("database error")
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
	req := httptest.NewRequest("GET", "/todos", nil)
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

	req := httptest.NewRequest("GET", "/todos/1", nil)
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

	req = httptest.NewRequest("GET", "/todos/abc", nil)
	req.SetPathValue("id", "abc")
	w = httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid ID, got %d", w.Code)
	}

	req = httptest.NewRequest("GET", "/todos/999", nil)
	req.SetPathValue("id", "999")
	w = httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 for non-existent ID, got %d", w.Code)
	}

	db.shouldErr = true
	req = httptest.NewRequest("GET", "/todos/1", nil)
	req.SetPathValue("id", "1")
	w = httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500 on database error, got %d", w.Code)
	}
}

func TestCreateToDo(t *testing.T) {
	logger := std.New("debug")
	db := &mockDB{todos: make(map[int]model.ToDo)}

	handler := CreateToDo(logger, db)

	todo := model.ToDo{
		Caption:     "New Todo",
		Description: "New Description",
	}
	body, _ := json.Marshal(todo)

	req := httptest.NewRequest("POST", "/todos", bytes.NewReader(body))
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
	body, _ = json.Marshal(todo)

	req = httptest.NewRequest("POST", "/todos", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for empty caption, got %d", w.Code)
	}

	req = httptest.NewRequest("POST", "/todos", bytes.NewReader([]byte("{invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid JSON, got %d", w.Code)
	}

	db.todos[5] = model.ToDo{ID: 5, Caption: "Existing"}
	todo = model.ToDo{ID: 5, Caption: "Duplicate"}
	body, _ = json.Marshal(todo)

	req = httptest.NewRequest("POST", "/todos", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("Expected status 409 for duplicate ID, got %d", w.Code)
	}

	db.shouldErr = true
	todo = model.ToDo{Caption: "Should Fail"}
	body, _ = json.Marshal(todo)

	req = httptest.NewRequest("POST", "/todos", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500 on database error, got %d", w.Code)
	}
}

func TestUpsertToDo(t *testing.T) {
	logger := std.New("debug")
	db := &mockDB{
		todos: map[int]model.ToDo{
			1: {ID: 1, Caption: "Original"},
		},
	}

	handler := UpsertToDo(logger, db)

	update := upsertToDoRequest{
		Caption:     "Upsertd",
		Description: "Upsertd Description",
		IsCompleted: true,
	}
	body, _ := json.Marshal(update)

	req := httptest.NewRequest("PUT", "/todos/1", bytes.NewReader(body))
	req.SetPathValue("id", "1")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204 for update, got %d", w.Code)
	}

	update.Caption = "New via PUT"
	body, _ = json.Marshal(update)

	req = httptest.NewRequest("PUT", "/todos/2", bytes.NewReader(body))
	req.SetPathValue("id", "2")
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201 for create via update, got %d", w.Code)
	}

	update.Caption = ""
	body, _ = json.Marshal(update)

	req = httptest.NewRequest("PUT", "/todos/1", bytes.NewReader(body))
	req.SetPathValue("id", "1")
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for empty caption, got %d", w.Code)
	}

	req = httptest.NewRequest("PUT", "/todos/abc", bytes.NewReader(body))
	req.SetPathValue("id", "abc")
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid ID, got %d", w.Code)
	}

	db.shouldErr = true
	update.Caption = "Should Fail"
	body, _ = json.Marshal(update)

	req = httptest.NewRequest("PUT", "/todos/1", bytes.NewReader(body))
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

	req := httptest.NewRequest("DELETE", "/todos/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}

	req = httptest.NewRequest("DELETE", "/todos/999", nil)
	req.SetPathValue("id", "999")
	w = httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 for non-existent todo, got %d", w.Code)
	}

	req = httptest.NewRequest("DELETE", "/todos/abc", nil)
	req.SetPathValue("id", "abc")
	w = httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid ID, got %d", w.Code)
	}

	db.shouldErr = true
	db.todos[2] = model.ToDo{ID: 2, Caption: "Should Fail"}
	req = httptest.NewRequest("DELETE", "/todos/2", nil)
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

	var decoded map[string]interface{}
	if err := json.Unmarshal(body, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if _, ok := decoded["todos"]; !ok {
		t.Error("Response should have 'todos' field")
	}
}
