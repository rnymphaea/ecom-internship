package mem

import (
	"ecom-internship/internal/model"
)

type MemDB struct {
	data []model.ToDo
	last int
}
