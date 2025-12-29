// Package mem provides an in-memory storage implementation for ToDo items.
package mem

import (
	"sync"

	"ecom-internship/internal/logger"
	"ecom-internship/internal/model"
)

// MemDB represents an in-memory ToDo storage.
//
//nolint:revive
type MemDB struct {
	data []model.ToDo
	log  logger.Logger
	mu   sync.RWMutex
	last int
}
