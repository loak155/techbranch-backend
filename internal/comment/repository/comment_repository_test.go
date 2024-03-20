package repository

import (
	"regexp"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/loak155/techbranch-backend/internal/comment/domain"
	"github.com/loak155/techbranch-backend/mock"
)

func testComment() *domain.Comment {
	return &domain.Comment{
		UserID:    1,
		ArticleID: 1,
		Content:   "test_Content1",
	}
}

func testComment2() *domain.Comment {
	return &domain.Comment{
		UserID:    2,
		ArticleID: 2,
		Content:   "test_Content2",
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
		`INSERT INTO "comments" ("user_id","article_id","content","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`)).
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

func TestGetComment(t *testing.T) {
	testComment := testComment()

	db, mock, err := mock.NewDBMock()
	if err != nil {
		t.Errorf("Failed to initialize mock DB: %v", err)
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "article_id", "content", "created_at", "updated_at", "deleted_at"}).
		AddRow(1, testComment.UserID, testComment.ArticleID, testComment.Content, time.Now(), time.Now(), nil)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "comments" WHERE "comments"."id" = $1 AND "comments"."deleted_at" IS NULL ORDER BY "comments"."id" LIMIT $2`)).
		WithArgs(1, 1).
		WillReturnRows(rows)

	repo := NewCommentRepository(db)
	_, err = repo.GetComment(1)
	if err != nil {
		t.Fatalf("failed to get comment: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Test Find Comment: %v", err)
	}
}

func TestListComments(t *testing.T) {
	testComment1 := testComment()
	testComment2 := testComment2()

	db, mock, err := mock.NewDBMock()
	if err != nil {
		t.Errorf("Failed to initialize mock DB: %v", err)
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "article_id", "content", "created_at", "updated_at", "deleted_at"}).
		AddRow(1, testComment1.UserID, testComment1.ArticleID, testComment1.Content, time.Now(), time.Now(), nil).
		AddRow(2, testComment2.UserID, testComment2.ArticleID, testComment2.Content, time.Now(), time.Now(), nil)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "comments" WHERE article_id=$1 AND "comments"."deleted_at" IS NULL`)).
		WithArgs(1).
		WillReturnRows(rows)

	repo := NewCommentRepository(db)
	_, err = repo.ListCommentsByArticleID(1)
	if err != nil {
		t.Fatalf("failed to list comment: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Test Find Comment: %v", err)
	}
}

func TestUpdateComment(t *testing.T) {
	testComment := testComment()

	db, mock, err := mock.NewDBMock()
	if err != nil {
		t.Errorf("Failed to initialize mock DB: %v", err)
	}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "comments" ("user_id","article_id","content","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`)).
		WillReturnRows(rows)
	mock.ExpectCommit()

	repo := NewCommentRepository(db)
	err = repo.UpdateComment(testComment)
	if err != nil {
		t.Fatalf("failed to list comment: %s", err)
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
		`UPDATE "comments" SET "deleted_at"=$1 WHERE "comments"."id" = $2 AND "comments"."deleted_at" IS NULL`)).
		WithArgs(sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := NewCommentRepository(db)
	err = repo.DeleteComment(1)
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
		`UPDATE "comments" SET "deleted_at"=$1 WHERE user_id=$2 AND "comments"."deleted_at" IS NULL`)).
		WithArgs(sqlmock.AnyArg(), 1).
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

func TestUpdateByUserIDWithUnscoped(t *testing.T) {
	db, mock, err := mock.NewDBMock()
	if err != nil {
		t.Errorf("Failed to initialize mock DB: %v", err)
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "comments" SET "deleted_at"=$1,"updated_at"=$2 WHERE "user_id" = $3`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := NewCommentRepository(db)
	err = repo.UpdateByUserIDWithUnscoped(1)
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
		`UPDATE "comments" SET "deleted_at"=$1 WHERE article_id=$2 AND "comments"."deleted_at" IS NULL`)).
		WithArgs(sqlmock.AnyArg(), 1).
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

func TestUpdateByArticleIDWithUnscoped(t *testing.T) {
	db, mock, err := mock.NewDBMock()
	if err != nil {
		t.Errorf("Failed to initialize mock DB: %v", err)
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "comments" SET "deleted_at"=$1,"updated_at"=$2 WHERE "article_id" = $3`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := NewCommentRepository(db)
	err = repo.UpdateByArticleIDWithUnscoped(1)
	if err != nil {
		t.Fatalf("failed to list comment: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Test Find Comment: %v", err)
	}
}
