package database

import (
	"context"

	"ecom-internship/internal/model"
)

type Database interface {
	GetAllToDos(ctx context.Context) ([]model.ToDo, error)
	GetToDoByID(ctx context.Context, id int) (model.ToDo, error)
	CreateToDo(ctx context.Context, todo model.ToDo) (int, error)
	UpdateToDo(ctx context.Context, todo model.ToDo) (bool, error)
	DeleteToDo(ctx context.Context, id int) error
}
