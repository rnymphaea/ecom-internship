package mem

import (
	"context"
	_ "slices"
	"time"

	"ecom-internship/internal/database"
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
	const funcName = "GetAllToDos"

	db.mu.RLock()
	defer db.mu.RUnlock()

	res := make([]model.ToDo, len(db.data))
	copy(res, db.data)

	if err := ctx.Err(); err != nil {
		db.log.Info("Context cancelled", "func", funcName)

		return nil, err
	}

	return res, nil
}

func (db *MemDB) GetToDoByID(ctx context.Context, id int) (model.ToDo, error) {
	const funcName = "GetToDoByID"

	db.mu.RLock()
	defer db.mu.RUnlock()

	index, found := db.find(id)

	if err := ctx.Err(); err != nil {
		db.log.Info("Context cancelled", "func", funcName)

		return model.ToDo{}, err
	}

	if found {
		return db.data[index], nil
	} else {
		return model.ToDo{}, database.ErrNotFound
	}
}

func (db *MemDB) find(id int) (int, bool) {
	for ind, value := range db.data {
		if value.ID == id {
			return ind, true
		}
	}

	return -1, false
}

func (db *MemDB) CreateToDo(ctx context.Context, todo model.ToDo) (int, error) {
	const funcName = "CreateToDo"

	db.mu.Lock()
	defer db.mu.Unlock()

	_, found := db.find(todo.ID)
	if found {
		return -1, database.ErrIDAlreadyExists
	}

	if todo.ID == 0 {
		if len(db.data) == 0 {
			todo.ID = 1
		} else {
			todo.ID = db.data[len(db.data)-1].ID + 1
		}
	}

	createdAt := time.Now()
	todo.CreatedAt = createdAt
	todo.UpdatedAt = createdAt

	if err := ctx.Err(); err != nil {
		db.log.Info("Context cancelled", "func", funcName)

		return -1, err
	}

	db.data = append(db.data, todo)

	return todo.ID, nil
}

func (db *MemDB) UpdateToDo(ctx context.Context, todo model.ToDo) error {
	const funcName = "UpdateToDo"

	db.mu.Lock()
	defer db.mu.Unlock()

	index, found := db.find(todo.ID)
	if !found {
		return database.ErrNotFound
	}

	todo.CreatedAt = db.data[index].CreatedAt
	todo.UpdatedAt = time.Now()

	if err := ctx.Err(); err != nil {
		db.log.Info("Context cancelled", "func", funcName)

		return err
	}

	db.data[index] = todo

	return nil
}

func (db *MemDB) DeleteToDo(ctx context.Context, id int) error {
	const funcName = "DeleteToDo"

	db.mu.Lock()
	defer db.mu.Unlock()

	index, found := db.find(id)
	if !found {
		return database.ErrNotFound
	} else {
		if err := ctx.Err(); err != nil {
			db.log.Info("Context cancelled", "func", funcName)

			return err
		}

		db.data = append(db.data[:index], db.data[index+1:]...)

		return nil
	}
}
