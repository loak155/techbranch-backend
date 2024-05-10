package domain

import (
	"time"

	"gorm.io/gorm"
)

type Article struct {
	ID        uint           `json:"id"`
	Title     string         `json:"title"`
	Url       string         `json:"url"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
