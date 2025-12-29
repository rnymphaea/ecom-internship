// Package database defines interfaces and errors for data storage operations.
package database

import (
	"context"
	"errors"

	"ecom-internship/internal/model"
)

// Database defines the interface for ToDo storage operations.
type Database interface {
	GetAllToDos(ctx context.Context) ([]model.ToDo, error)
	GetToDoByID(ctx context.Context, id int) (model.ToDo, error)
	CreateToDo(ctx context.Context, todo model.ToDo) (int, error)
	UpdateToDo(ctx context.Context, todo model.ToDo) error
	DeleteToDo(ctx context.Context, id int) error
}

var (
	// ErrNotFound is returned when a ToDo is not found.
	ErrNotFound = errors.New("todo not found")

	// ErrIDAlreadyExists is returned when creating a ToDo with an existing ID.
	ErrIDAlreadyExists = errors.New("todo with provided id already exists")
)
