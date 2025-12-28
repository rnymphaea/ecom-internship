package mem

import (
	"context"
	"errors"
	"testing"

	"ecom-internship/internal/database"
	"ecom-internship/internal/logger/std"
	"ecom-internship/internal/model"
)

func TestMemDB_CreateToDo(t *testing.T) {
	logger := std.New("debug")
	db := New(logger)
	ctx := context.Background()

	todo1 := model.ToDo{
		Caption:     "Test Todo 1",
		Description: "Test Description 1",
	}

	id, err := db.CreateToDo(ctx, todo1)
	if err != nil {
		t.Fatalf("CreateToDo failed: %v", err)
	}
	if id != 1 {
		t.Errorf("Expected ID 1, got %d", id)
	}

	todo2 := model.ToDo{
		Caption:     "Test Todo 2",
		Description: "Test Description 2",
	}

	id2, err := db.CreateToDo(ctx, todo2)
	if err != nil {
		t.Fatalf("CreateToDo failed: %v", err)
	}
	if id2 != 2 {
		t.Errorf("Expected ID 2, got %d", id2)
	}

	todo3 := model.ToDo{
		ID:          1,
		Caption:     "Test Todo 3",
		Description: "Test Description 3",
	}

	_, err = db.CreateToDo(ctx, todo3)
	if !errors.Is(err, database.ErrIDAlreadyExists) {
		t.Errorf("Expected ErrIDAlreadyExists, got %v", err)
	}

	todo4 := model.ToDo{
		ID:          100,
		Caption:     "Test Todo 100",
		Description: "Test Description 100",
	}

	id4, err := db.CreateToDo(ctx, todo4)
	if err != nil {
		t.Fatalf("CreateToDo failed: %v", err)
	}
	if id4 != 100 {
		t.Errorf("Expected ID 100, got %d", id4)
	}
}

func TestMemDB_GetAllToDos(t *testing.T) {
	logger := std.New("debug")
	db := New(logger)
	ctx := context.Background()

	todos, err := db.GetAllToDos(ctx)
	if err != nil {
		t.Fatalf("GetAllToDos failed: %v", err)
	}
	if len(todos) != 0 {
		t.Errorf("Expected 0 todos, got %d", len(todos))
	}

	_, err = db.CreateToDo(ctx, model.ToDo{Caption: "Todo 1", Description: "Desc 1"})
	if err != nil {
		t.Fatalf("CreateToDo failed: %v", err)
	}

	_, err = db.CreateToDo(ctx, model.ToDo{Caption: "Todo 2", Description: "Desc 2"})
	if err != nil {
		t.Fatalf("CreateToDo failed: %v", err)
	}

	todos, err = db.GetAllToDos(ctx)
	if err != nil {
		t.Fatalf("GetAllToDos failed: %v", err)
	}
	if len(todos) != 2 {
		t.Errorf("Expected 2 todos, got %d", len(todos))
	}
}

func TestMemDB_GetToDoByID(t *testing.T) {
	logger := std.New("debug")
	db := New(logger)
	ctx := context.Background()

	todo := model.ToDo{
		Caption:     "Test Todo",
		Description: "Test Description",
	}
	id, err := db.CreateToDo(ctx, todo)
	if err != nil {
		t.Fatalf("CreateToDo failed: %v", err)
	}

	retrieved, err := db.GetToDoByID(ctx, id)
	if err != nil {
		t.Fatalf("GetToDoByID failed: %v", err)
	}
	if retrieved.Caption != todo.Caption {
		t.Errorf("Expected caption %s, got %s", todo.Caption, retrieved.Caption)
	}

	_, err = db.GetToDoByID(ctx, 999)
	if !errors.Is(err, database.ErrNotFound) {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

func TestMemDB_UpdateToDo(t *testing.T) {
	logger := std.New("debug")
	db := New(logger)
	ctx := context.Background()

	original := model.ToDo{
		Caption:     "Original",
		Description: "Original Description",
	}
	id, err := db.CreateToDo(ctx, original)
	if err != nil {
		t.Fatalf("CreateToDo failed: %v", err)
	}

	updated := model.ToDo{
		ID:          id,
		Caption:     "Updated",
		Description: "Updated Description",
		IsCompleted: true,
	}

	err = db.UpdateToDo(ctx, updated)
	if err != nil {
		t.Fatalf("UpdateToDo failed: %v", err)
	}

	retrieved, err := db.GetToDoByID(ctx, id)
	if err != nil {
		t.Fatalf("GetToDoByID failed: %v", err)
	}
	if retrieved.Caption != updated.Caption {
		t.Errorf("Expected caption %s, got %s", updated.Caption, retrieved.Caption)
	}
	if !retrieved.IsCompleted {
		t.Error("Expected IsCompleted to be true")
	}

	newTodo := model.ToDo{
		ID:          999,
		Caption:     "New via Update",
		Description: "Created via update",
		IsCompleted: false,
	}

	err = db.UpdateToDo(ctx, newTodo)
	if !errors.Is(err, database.ErrNotFound) {
		t.Errorf("Expected ErrNotFound for non-existent todo, got %v", err)
	}

	_, err = db.GetToDoByID(ctx, 999)
	if !errors.Is(err, database.ErrNotFound) {
		t.Errorf("Expected ErrNotFound for non-existent todo, got %v", err)
	}
}

func TestMemDB_DeleteToDo(t *testing.T) {
	logger := std.New("debug")
	db := New(logger)
	ctx := context.Background()

	todo := model.ToDo{
		Caption:     "Test Todo",
		Description: "Test Description",
	}
	id, err := db.CreateToDo(ctx, todo)
	if err != nil {
		t.Fatalf("CreateToDo failed: %v", err)
	}

	err = db.DeleteToDo(ctx, id)
	if err != nil {
		t.Fatalf("DeleteToDo failed: %v", err)
	}

	_, err = db.GetToDoByID(ctx, id)
	if !errors.Is(err, database.ErrNotFound) {
		t.Errorf("Expected ErrNotFound after deletion, got %v", err)
	}

	err = db.DeleteToDo(ctx, 999)
	if !errors.Is(err, database.ErrNotFound) {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

func TestMemDB_ConcurrentAccess(t *testing.T) {
	logger := std.New("debug")
	db := New(logger)
	ctx := context.Background()

	const numGoroutines = 10
	const opsPerGoroutine = 100
	done := make(chan bool)

	for i := range numGoroutines {
		go func(goroutineID int) {
			for j := range opsPerGoroutine {
				id := goroutineID*1000 + j

				todo := model.ToDo{
					ID:          id,
					Caption:     "Todo",
					Description: "Description",
				}

				//nolint:errcheck,gosec
				db.CreateToDo(ctx, todo)
			}
			done <- true
		}(i)
	}

	for range numGoroutines {
		<-done
	}

	todos, err := db.GetAllToDos(ctx)
	if err != nil {
		t.Fatalf("GetAllToDos failed: %v", err)
	}
	expectedCount := numGoroutines * opsPerGoroutine
	if len(todos) != expectedCount {
		t.Errorf("Expected %d todos, got %d", expectedCount, len(todos))
	}
}
