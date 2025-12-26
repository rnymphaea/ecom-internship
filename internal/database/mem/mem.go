package mem

import (
	"sync"

	"ecom-internship/internal/logger"
	"ecom-internship/internal/model"
)

type MemDB struct {
	data []model.ToDo
	log  logger.Logger
	mu   sync.RWMutex
}
