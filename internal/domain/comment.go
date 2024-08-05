package domain

import (
	"time"
)

type Comment struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	ArticleID uint      `json:"article_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
