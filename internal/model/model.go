// Package model defines the application's models.
package model

import (
	"time"
)

// ToDo represents a task (item).
//
//nolint:godox
type ToDo struct {
	ID          int       `json:"id"`
	Caption     string    `json:"caption"`
	Description string    `json:"description"`
	IsCompleted bool      `json:"is_completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
