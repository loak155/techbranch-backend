package db

import (
	"fmt"
	"log"

	"github.com/loak155/techbranch-backend/internal/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(url string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	return db
}

func NewUserDB(conf *config.Config) *gorm.DB {
	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		conf.User.DB.User,
		conf.User.DB.Password,
		conf.User.DB.Host,
		conf.User.DB.Port,
		conf.User.DB.Name,
	)
	return NewDB(url)
}

func NewArticleDB(conf *config.Config) *gorm.DB {
	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		conf.Article.DB.User,
		conf.Article.DB.Password,
		conf.Article.DB.Host,
		conf.Article.DB.Port,
		conf.Article.DB.Name,
	)
	return NewDB(url)
}

func NewBookmarkDB(conf *config.Config) *gorm.DB {
	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		conf.Bookmark.DB.User,
		conf.Bookmark.DB.Password,
		conf.Bookmark.DB.Host,
		conf.Bookmark.DB.Port,
		conf.Bookmark.DB.Name,
	)
	return NewDB(url)
}

func NewCommentDB(conf *config.Config) *gorm.DB {
	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		conf.Comment.DB.User,
		conf.Comment.DB.Password,
		conf.Comment.DB.Host,
		conf.Comment.DB.Port,
		conf.Comment.DB.Name,
	)
	return NewDB(url)
}

func CloseDB(db *gorm.DB) {
	sqlDB, _ := db.DB()
	if err := sqlDB.Close(); err != nil {
		log.Fatalln(err)
	}
}
