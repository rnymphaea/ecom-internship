package mem

import (
	"context"
	"testing"

	"ecom-internship/internal/database"
	"ecom-internship/internal/logger/std"
	"ecom-internship/internal/model"
)

func TestMemDB_CreateToDo(t *testing.T) {
	logger := std.New("debug")
	db := New(logger)
	ctx := context.Background()

	// Test 1: Create todo without ID
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

	// Test 2: Create todo with ID 0 (should auto-generate)
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

	// Test 3: Create todo with existing ID (should fail)
	todo3 := model.ToDo{
		ID:          1,
		Caption:     "Test Todo 3",
		Description: "Test Description 3",
	}

	_, err = db.CreateToDo(ctx, todo3)
	if err != database.ErrIDAlreadyExists {
		t.Errorf("Expected ErrIDAlreadyExists, got %v", err)
	}

	// Test 4: Create todo with custom valid ID
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

	// Empty database
	todos, err := db.GetAllToDos(ctx)
	if err != nil {
		t.Fatalf("GetAllToDos failed: %v", err)
	}
	if len(todos) != 0 {
		t.Errorf("Expected 0 todos, got %d", len(todos))
	}

	// Add some todos
	db.CreateToDo(ctx, model.ToDo{Caption: "Todo 1", Description: "Desc 1"})
	db.CreateToDo(ctx, model.ToDo{Caption: "Todo 2", Description: "Desc 2"})

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

	// Create a todo
	todo := model.ToDo{
		Caption:     "Test Todo",
		Description: "Test Description",
	}
	id, err := db.CreateToDo(ctx, todo)
	if err != nil {
		t.Fatalf("CreateToDo failed: %v", err)
	}

	// Test 1: Get existing todo
	retrieved, err := db.GetToDoByID(ctx, id)
	if err != nil {
		t.Fatalf("GetToDoByID failed: %v", err)
	}
	if retrieved.Caption != todo.Caption {
		t.Errorf("Expected caption %s, got %s", todo.Caption, retrieved.Caption)
	}

	// Test 2: Get non-existing todo
	_, err = db.GetToDoByID(ctx, 999)
	if err != database.ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

func TestMemDB_UpdateToDo(t *testing.T) {
	logger := std.New("debug")
	db := New(logger)
	ctx := context.Background()

	// Create a todo
	original := model.ToDo{
		Caption:     "Original",
		Description: "Original Description",
	}
	id, err := db.CreateToDo(ctx, original)
	if err != nil {
		t.Fatalf("CreateToDo failed: %v", err)
	}

	// Test 1: Update existing todo
	updated := model.ToDo{
		ID:          id,
		Caption:     "Updated",
		Description: "Updated Description",
		IsCompleted: true,
	}

	found, err := db.UpdateToDo(ctx, updated)
	if err != nil {
		t.Fatalf("UpdateToDo failed: %v", err)
	}
	if !found {
		t.Error("Expected todo to be found for update")
	}

	// Verify update
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

	// Test 2: Update non-existing todo (should create)
	newTodo := model.ToDo{
		ID:          999,
		Caption:     "New via Update",
		Description: "Created via update",
		IsCompleted: false,
	}

	found, err = db.UpdateToDo(ctx, newTodo)
	if err != nil {
		t.Fatalf("UpdateToDo failed: %v", err)
	}
	if found {
		t.Error("Expected todo not to be found (should create new)")
	}

	// Verify new todo was created
	retrieved, err = db.GetToDoByID(ctx, 999)
	if err != nil {
		t.Fatalf("GetToDoByID failed: %v", err)
	}
	if retrieved.Caption != newTodo.Caption {
		t.Errorf("Expected caption %s, got %s", newTodo.Caption, retrieved.Caption)
	}
}

func TestMemDB_DeleteToDo(t *testing.T) {
	logger := std.New("debug")
	db := New(logger)
	ctx := context.Background()

	// Create a todo
	todo := model.ToDo{
		Caption:     "Test Todo",
		Description: "Test Description",
	}
	id, err := db.CreateToDo(ctx, todo)
	if err != nil {
		t.Fatalf("CreateToDo failed: %v", err)
	}

	// Test 1: Delete existing todo
	err = db.DeleteToDo(ctx, id)
	if err != nil {
		t.Fatalf("DeleteToDo failed: %v", err)
	}

	// Verify deletion
	_, err = db.GetToDoByID(ctx, id)
	if err != database.ErrNotFound {
		t.Errorf("Expected ErrNotFound after deletion, got %v", err)
	}

	// Test 2: Delete non-existing todo
	err = db.DeleteToDo(ctx, 999)
	if err != database.ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

func TestMemDB_ConcurrentAccess(t *testing.T) {
	logger := std.New("debug")
	db := New(logger)
	ctx := context.Background()

	// Start multiple goroutines
	const numGoroutines = 10
	const opsPerGoroutine = 100
	done := make(chan bool)

	for i := 0; i < numGoroutines; i++ {
		go func(goroutineID int) {
			for j := 0; j < opsPerGoroutine; j++ {
				id := goroutineID*1000 + j

				todo := model.ToDo{
					ID:          id,
					Caption:     "Todo",
					Description: "Description",
				}
				db.CreateToDo(ctx, todo)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// Verify data integrity
	todos, err := db.GetAllToDos(ctx)
	if err != nil {
		t.Fatalf("GetAllToDos failed: %v", err)
	}
	expectedCount := numGoroutines * opsPerGoroutine
	if len(todos) != expectedCount {
		t.Errorf("Expected %d todos, got %d", expectedCount, len(todos))
	}
}
