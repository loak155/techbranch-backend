package repository

import (
	"regexp"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/loak155/techbranch-backend/internal/domain"
	"github.com/loak155/techbranch-backend/mock"
)

func testComment() *domain.Comment {
	return &domain.Comment{
		UserID:    1,
		ArticleID: 1,
		Content:   "test_content",
	}
}

func testComment2() *domain.Comment {
	return &domain.Comment{
		UserID:    2,
		ArticleID: 2,
		Content:   "test_content2",
	}
}

func TestCreateComment(t *testing.T) {
	testComment := testComment()

	db, mock, err := mock.NewDBMock()
	if err != nil {
		t.Errorf("Failed to initialize mock DB: %v", err)
	}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "comments" ("user_id","article_id","content","created_at","updated_at") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
		WillReturnRows(rows)
	mock.ExpectCommit()

	repo := NewCommentRepository(db)
	err = repo.CreateComment(testComment)
	if err != nil {
		t.Fatal(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Test Create Comment: %v", err)
	}
}

func TestListCommentsByUserID(t *testing.T) {
	testComment1 := testComment()
	testComment2 := testComment2()

	db, mock, err := mock.NewDBMock()
	if err != nil {
		t.Errorf("Failed to initialize mock DB: %v", err)
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "article_id", "content", "created_at", "updated_at"}).
		AddRow(1, testComment1.UserID, testComment1.ArticleID, testComment1.Content, time.Now(), time.Now()).
		AddRow(2, testComment2.UserID, testComment2.ArticleID, testComment2.Content, time.Now(), time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "comments" WHERE user_id=$1`)).
		WithArgs(1).
		WillReturnRows(rows)

	repo := NewCommentRepository(db)
	_, err = repo.ListCommentsByUserID(1)
	if err != nil {
		t.Fatalf("failed to list Comment: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Test Find Comment: %v", err)
	}
}

func TestListCommentsByArticleID(t *testing.T) {
	testComment1 := testComment()
	testComment2 := testComment2()

	db, mock, err := mock.NewDBMock()
	if err != nil {
		t.Errorf("Failed to initialize mock DB: %v", err)
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "article_id", "content", "created_at", "updated_at"}).
		AddRow(1, testComment1.UserID, testComment1.ArticleID, testComment1.Content, time.Now(), time.Now()).
		AddRow(2, testComment2.UserID, testComment2.ArticleID, testComment2.Content, time.Now(), time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "comments" WHERE article_id=$1`)).
		WithArgs(1).
		WillReturnRows(rows)

	repo := NewCommentRepository(db)
	_, err = repo.ListCommentsByArticleID(1)
	if err != nil {
		t.Fatalf("failed to list Comment: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Test Find Comment: %v", err)
	}
}

func TestDeleteComment(t *testing.T) {
	db, mock, err := mock.NewDBMock()
	if err != nil {
		t.Errorf("Failed to initialize mock DB: %v", err)
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM "comments" WHERE "comments"."id" = $1`)).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := NewCommentRepository(db)
	err = repo.DeleteComment(1)
	if err != nil {
		t.Fatalf("failed to delete article: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Test Find Article: %v", err)
	}
}

func TestDeleteCommentByUserIDAndArticleID(t *testing.T) {
	db, mock, err := mock.NewDBMock()
	if err != nil {
		t.Errorf("Failed to initialize mock DB: %v", err)
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM "comments" WHERE user_id=$1 AND article_id=$2`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := NewCommentRepository(db)
	err = repo.DeleteCommentByUserIDAndArticleID(1, 1)
	if err != nil {
		t.Fatalf("failed to list comment: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Test Find Comment: %v", err)
	}
}

func TestDeleteCommentByUserID(t *testing.T) {
	db, mock, err := mock.NewDBMock()
	if err != nil {
		t.Errorf("Failed to initialize mock DB: %v", err)
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM "comments" WHERE user_id=$1`)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := NewCommentRepository(db)
	err = repo.DeleteCommentByUserID(1)
	if err != nil {
		t.Fatalf("failed to list comment: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Test Find Comment: %v", err)
	}
}

func TestDeleteCommentByArticleID(t *testing.T) {
	db, mock, err := mock.NewDBMock()
	if err != nil {
		t.Errorf("Failed to initialize mock DB: %v", err)
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM "comments" WHERE article_id=$1`)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := NewCommentRepository(db)
	err = repo.DeleteCommentByArticleID(1)
	if err != nil {
		t.Fatalf("failed to list comment: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Test Find Comment: %v", err)
	}
}
