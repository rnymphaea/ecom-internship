package mem

import (
	_ "cmp"
	"encoding/json"
	"errors"
	"slices"

	"ecom-internship/internal/model"
)

func New() *MemDB {
	db := &MemDB{
		data: make([]model.ToDo, 0),
		last: 0,
	}

	return db
}

func (db *MemDB) GetAllToDos() ([]model.ToDo, error) {
	return data, nil
}

func (db *MemDB) GetToDoByID(id int) (model.ToDo, error) {
	index, found := slices.BinarySearchFunc(
		db.data,
		id,
		func(a, b model.ToDo) int {
			return a.id - b.id
		},
	)

	if found {
		return db.data[index]
	} else {
		return errors.New("not found")
	}
}
