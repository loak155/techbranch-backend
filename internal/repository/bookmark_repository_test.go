package repository

import (
	"regexp"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/loak155/techbranch-backend/internal/domain"
	"github.com/loak155/techbranch-backend/mock"
)

func testBookmark() *domain.Bookmark {
	return &domain.Bookmark{
		UserID:    1,
		ArticleID: 1,
	}
}

func testBookmark2() *domain.Bookmark {
	return &domain.Bookmark{
		UserID:    2,
		ArticleID: 2,
	}
}

func TestCreateBookmark(t *testing.T) {
	testBookmark := testBookmark()

	db, mock, err := mock.NewDBMock()
	if err != nil {
		t.Errorf("Failed to initialize mock DB: %v", err)
	}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "bookmarks" ("user_id","article_id","created_at","updated_at") VALUES ($1,$2,$3,$4) RETURNING "id"`)).
		WillReturnRows(rows)
	mock.ExpectCommit()

	repo := NewBookmarkRepository(db)
	err = repo.CreateBookmark(testBookmark)
	if err != nil {
		t.Fatal(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Test Create Bookmark: %v", err)
	}
}

func TestGetBookmarkCountByArticleID(t *testing.T) {
	db, mock, err := mock.NewDBMock()
	if err != nil {
		t.Errorf("Failed to initialize mock DB: %v", err)
	}

	rows := sqlmock.NewRows([]string{"count"}).AddRow(1)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT count(*) FROM "bookmarks" WHERE article_id=$1`)).
		WithArgs(1).
		WillReturnRows(rows)

	repo := NewBookmarkRepository(db)
	_, err = repo.GetBookmarkCountByArticleID(1)
	if err != nil {
		t.Fatalf("failed to get bookmark count: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Test Get Bookmark Count: %v", err)
	}
}

func TestListBookmarksByUserID(t *testing.T) {
	testBookmark1 := testBookmark()
	testBookmark2 := testBookmark2()

	db, mock, err := mock.NewDBMock()
	if err != nil {
		t.Errorf("Failed to initialize mock DB: %v", err)
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "article_id", "created_at", "updated_at"}).
		AddRow(1, testBookmark1.UserID, testBookmark1.ArticleID, time.Now(), time.Now()).
		AddRow(2, testBookmark2.UserID, testBookmark2.ArticleID, time.Now(), time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "bookmarks" WHERE user_id=$1`)).
		WithArgs(1).
		WillReturnRows(rows)

	repo := NewBookmarkRepository(db)
	_, err = repo.ListBookmarksByUserID(1)
	if err != nil {
		t.Fatalf("failed to list Bookmark: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Test Find Bookmark: %v", err)
	}
}

func TestListBookmarksByArticleID(t *testing.T) {
	testBookmark1 := testBookmark()
	testBookmark2 := testBookmark2()

	db, mock, err := mock.NewDBMock()
	if err != nil {
		t.Errorf("Failed to initialize mock DB: %v", err)
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "article_id", "created_at", "updated_at"}).
		AddRow(1, testBookmark1.UserID, testBookmark1.ArticleID, time.Now(), time.Now()).
		AddRow(2, testBookmark2.UserID, testBookmark2.ArticleID, time.Now(), time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "bookmarks" WHERE article_id=$1`)).
		WithArgs(1).
		WillReturnRows(rows)

	repo := NewBookmarkRepository(db)
	_, err = repo.ListBookmarksByArticleID(1)
	if err != nil {
		t.Fatalf("failed to list Bookmark: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Test Find Bookmark: %v", err)
	}
}

func TestDeleteBookmarkByUserIDAndArticleID(t *testing.T) {
	db, mock, err := mock.NewDBMock()
	if err != nil {
		t.Errorf("Failed to initialize mock DB: %v", err)
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM "bookmarks" WHERE user_id=$1 AND article_id=$2`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := NewBookmarkRepository(db)
	err = repo.DeleteBookmarkByUserIDAndArticleID(1, 1)
	if err != nil {
		t.Fatalf("failed to list bookmark: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Test Find Bookmark: %v", err)
	}
}

func TestDeleteBookmarkByUserID(t *testing.T) {
	db, mock, err := mock.NewDBMock()
	if err != nil {
		t.Errorf("Failed to initialize mock DB: %v", err)
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM "bookmarks" WHERE user_id=$1`)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := NewBookmarkRepository(db)
	err = repo.DeleteBookmarkByUserID(1)
	if err != nil {
		t.Fatalf("failed to list bookmark: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Test Find Bookmark: %v", err)
	}
}

func TestDeleteBookmarkByArticleID(t *testing.T) {
	db, mock, err := mock.NewDBMock()
	if err != nil {
		t.Errorf("Failed to initialize mock DB: %v", err)
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM "bookmarks" WHERE article_id=$1`)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := NewBookmarkRepository(db)
	err = repo.DeleteBookmarkByArticleID(1)
	if err != nil {
		t.Fatalf("failed to list bookmark: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Test Find Bookmark: %v", err)
	}
}
