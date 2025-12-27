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

	index, found := db.find(ctx, id)
	if found {
		return db.data[index], nil
	} else {
		return model.ToDo{}, errors.New("not found")
	}
}

func (db *MemDB) find(ctx context.Context, id int) (int, bool) {
	index, found := slices.BinarySearchFunc(
		db.data,
		model.ToDo{
			ID: id,
		},
		func(a, b model.ToDo) int {
			return a.ID - b.ID
		},
	)
	return index, found
}

func (db *MemDB) CreateToDo(ctx context.Context, todo model.ToDo) (int, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	_, found := db.find(ctx, todo.ID)
	if todo.ID == 0 || found {
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

func (db *MemDB) UpdateToDo(ctx context.Context, todo model.ToDo) (bool, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	index, found := db.find(ctx, todo.ID)
	if found {
		todo.CreatedAt = db.data[index].CreatedAt
		todo.UpdatedAt = time.Now()
		db.data[index] = todo
	} else {
		todo.CreatedAt = time.Now()
		todo.UpdatedAt = time.Now()
		db.data = append(db.data, todo)
	}

	return found, nil
}

func (db *MemDB) DeleteToDo(ctx context.Context, id int) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	index, found := db.find(ctx, id)
	if !found {
		return errors.New("not found")
	} else {
		db.data = append(db.data[:index], db.data[index+1:]...)
		return nil
	}
}
