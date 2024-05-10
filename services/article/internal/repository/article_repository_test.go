package repository

import (
	"regexp"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/loak155/techbranch-backend/pkg/db"
	"github.com/loak155/techbranch-backend/services/article/internal/domain"
)

func testArticle() *domain.Article {
	return &domain.Article{
		Title: "test_title",
		Url:   "http://example.com",
	}
}

func testArticle2() *domain.Article {
	return &domain.Article{
		Title: "test_title2",
		Url:   "http://example2.com",
	}
}

func TestCreateArticle(t *testing.T) {
	testArticle := testArticle()

	db, mock, err := db.NewDBMock()
	if err != nil {
		t.Errorf("Failed to initialize mock DB: %v", err)
	}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "articles" ("title","url","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
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

	db, mock, err := db.NewDBMock()
	if err != nil {
		t.Errorf("Failed to initialize mock DB: %v", err)
	}

	rows := sqlmock.NewRows([]string{"id", "title", "url", "created_at", "updated_at", "deleted_at"}).
		AddRow(1, testArticle.Title, testArticle.Url, time.Now(), time.Now(), nil)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "articles" WHERE "articles"."id" = $1 AND "articles"."deleted_at" IS NULL ORDER BY "articles"."id" LIMIT $2`)).
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

	db, mock, err := db.NewDBMock()
	if err != nil {
		t.Errorf("Failed to initialize mock DB: %v", err)
	}

	rows := sqlmock.NewRows([]string{"id", "title", "url", "created_at", "updated_at", "deleted_at"}).
		AddRow(1, testArticle1.Title, testArticle1.Url, time.Now(), time.Now(), nil).
		AddRow(2, testArticle2.Title, testArticle2.Url, time.Now(), time.Now(), nil)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "articles" WHERE "articles"."deleted_at" IS NULL`)).
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

	db, mock, err := db.NewDBMock()
	if err != nil {
		t.Errorf("Failed to initialize mock DB: %v", err)
	}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "articles" ("title","url","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
		WillReturnRows(rows)
	mock.ExpectCommit()

	repo := NewArticleRepository(db)
	err = repo.UpdateArticle(testArticle)
	if err != nil {
		t.Fatalf("failed to list article: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Test Find Article: %v", err)
	}
}

func TestDeleteArticle(t *testing.T) {
	db, mock, err := db.NewDBMock()
	if err != nil {
		t.Errorf("Failed to initialize mock DB: %v", err)
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "articles" SET "deleted_at"=$1 WHERE "articles"."id" = $2 AND "articles"."deleted_at" IS NULL`)).
		WithArgs(sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := NewArticleRepository(db)
	err = repo.DeleteArticle(1)
	if err != nil {
		t.Fatalf("failed to list article: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Test Find Article: %v", err)
	}
}
