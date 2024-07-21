package domain

import (
	"time"
)

type Article struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Url       string    `json:"url"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
