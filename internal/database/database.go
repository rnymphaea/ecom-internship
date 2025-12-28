package database

import (
	"context"
	"errors"

	"ecom-internship/internal/model"
)

type Database interface {
	GetAllToDos(ctx context.Context) ([]model.ToDo, error)
	GetToDoByID(ctx context.Context, id int) (model.ToDo, error)
	CreateToDo(ctx context.Context, todo model.ToDo) (int, error)
	UpdateToDo(ctx context.Context, todo model.ToDo) error
	DeleteToDo(ctx context.Context, id int) error
}

var (
	ErrNotFound        = errors.New("todo not found")
	ErrIDAlreadyExists = errors.New("todo with provided id already exists")
)
