package database

import (
	"context"

	"ecom-internship/internal/model"
)

type Database interface {
	GetAllToDos() ([]model.ToDo, error)
	GetToDoByID(id int) (model.ToDo, error)
}
