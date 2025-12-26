package mem

import (
	_ "cmp"
	"context"
	_ "encoding/json"
	"errors"
	"slices"
	"time"

	"ecom-internship/internal/logger"
	"ecom-internship/internal/model"
)

func New(log logger.Logger) *MemDB {
	db := &MemDB{
		data: make([]model.ToDo, 0),
		log:  log,
	}

	return db
}

func (db *MemDB) GetAllToDos(ctx context.Context) ([]model.ToDo, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	res := make([]model.ToDo, len(db.data))
	copy(res, db.data)
	return res, nil
}

func (db *MemDB) GetToDoByID(ctx context.Context, id int) (model.ToDo, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	index, found := slices.BinarySearchFunc(
		db.data,
		model.ToDo{
			ID: id,
		},
		func(a, b model.ToDo) int {
			return a.ID - b.ID
		},
	)

	if found {
		return db.data[index], nil
	} else {
		return model.ToDo{}, errors.New("not found")
	}
}

func (db *MemDB) CreateToDo(ctx context.Context, todo model.ToDo) (int, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, err := db.GetToDoByID(ctx, todo.ID); err != nil {
		if len(db.data) == 0 {
			todo.ID = 1
		} else {
			todo.ID = db.data[len(db.data)-1].ID + 1
		}
	}

	todo.CreatedAt = time.Now()
	todo.UpdatedAt = time.Now()

	db.data = append(db.data, todo)
	return todo.ID, nil
}
