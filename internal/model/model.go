package model

import (
	"time"
)

type ToDo struct {
	ID          int       `json:"id"`
	Caption     string    `json:"caption"`
	Description string    `json:"description"`
	IsCompleted bool      `json:"is_completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
