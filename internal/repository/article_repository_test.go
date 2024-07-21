package repository

import (
	"regexp"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/loak155/techbranch-backend/internal/domain"
	"github.com/loak155/techbranch-backend/mock"
)

func testArticle() *domain.Article {
	return &domain.Article{
		Title: "test_title",
		Url:   "http://example.com",
		Image: "http://example.com/image",
	}
}

func testArticle2() *domain.Article {
	return &domain.Article{
		Title: "test_title2",
		Url:   "http://example2.com",
		Image: "http://example2.com/image",
	}
}

func TestCreateArticle(t *testing.T) {
	testArticle := testArticle()

	db, mock, err := mock.NewDBMock()
	if err != nil {
		t.Errorf("Failed to initialize mock DB: %v", err)
	}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "articles" ("title","url","image","created_at","updated_at") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
		WillReturnRows(rows)
	mock.ExpectCommit()

	repo := NewArticleRepository(db)
	err = repo.CreateArticle(testArticle)
	if err != nil {
		t.Fatal(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Test Create Article: %v", err)
	}
}

func TestGetArticle(t *testing.T) {
	testArticle := testArticle()

	db, mock, err := mock.NewDBMock()
	if err != nil {
		t.Errorf("Failed to initialize mock DB: %v", err)
	}

	rows := sqlmock.NewRows([]string{"id", "title", "url", "image", "created_at", "updated_at"}).
		AddRow(1, testArticle.Title, testArticle.Url, time.Now(), time.Now(), nil)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "articles" WHERE "articles"."id" = $1 ORDER BY "articles"."id" LIMIT $2`)).
		WithArgs(1, 1).
		WillReturnRows(rows)

	repo := NewArticleRepository(db)
	_, err = repo.GetArticle(1)
	if err != nil {
		t.Fatalf("failed to get article: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Test Find Article: %v", err)
	}
}

func TestListArticles(t *testing.T) {
	testArticle1 := testArticle()
	testArticle2 := testArticle2()

	db, mock, err := mock.NewDBMock()
	if err != nil {
		t.Errorf("Failed to initialize mock DB: %v", err)
	}

	rows := sqlmock.NewRows([]string{"id", "title", "url", "image", "created_at", "updated_at"}).
		AddRow(1, testArticle1.Title, testArticle1.Url, testArticle1.Image, time.Now(), time.Now()).
		AddRow(2, testArticle2.Title, testArticle2.Url, testArticle2.Image, time.Now(), time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "articles"`)).
		WillReturnRows(rows)

	repo := NewArticleRepository(db)
	_, err = repo.ListArticles(1, 2)
	if err != nil {
		t.Fatalf("failed to list article: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Test Find Article: %v", err)
	}
}

func TestUpdateArticle(t *testing.T) {
	testArticle := testArticle()

	db, mock, err := mock.NewDBMock()
	if err != nil {
		t.Errorf("Failed to initialize mock DB: %v", err)
	}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "articles" ("title","url","image","created_at","updated_at") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
		WillReturnRows(rows)
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		`UPDATE "articles" SET "title"=$1,"url"=$2,"image"=$3,"created_at"=$4,"updated_at"=$5 WHERE "id" = $6 RETURNING *`)).
		WillReturnRows(rows)
	mock.ExpectCommit()

	repo := NewArticleRepository(db)
	err = repo.CreateArticle(testArticle)
	if err != nil {
		t.Fatalf("failed to create article: %s", err)
	}
	err = repo.UpdateArticle(testArticle)
	if err != nil {
		t.Fatalf("failed to list article: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Test Find Article: %v", err)
	}
}

func TestDeleteArticle(t *testing.T) {
	db, mock, err := mock.NewDBMock()
	if err != nil {
		t.Errorf("Failed to initialize mock DB: %v", err)
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM "articles" WHERE "articles"."id" = $1`)).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := NewArticleRepository(db)
	err = repo.DeleteArticle(1)
	if err != nil {
		t.Fatalf("failed to delete article: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Test Find Article: %v", err)
	}
}
