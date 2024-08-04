package repository

import (
	"regexp"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/loak155/techbranch-backend/internal/domain"
	"github.com/loak155/techbranch-backend/mock"
)

func testUser() *domain.User {
	return &domain.User{
		Username: "test_username",
		Email:    "test@example.com",
		Password: "test_password",
	}
}

func testUser2() *domain.User {
	return &domain.User{
		Username: "test_username2",
		Email:    "test2@example.com",
		Password: "test_password",
	}
}

func TestCreateUser(t *testing.T) {
	testUser := testUser()

	db, mock, err := mock.NewDBMock()
	if err != nil {
		t.Errorf("Failed to initialize mock DB: %v", err)
	}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "users" ("username","email","password","created_at","updated_at") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
		WillReturnRows(rows)
	mock.ExpectCommit()

	repo := NewUserRepository(db)
	err = repo.CreateUser(testUser)
	if err != nil {
		t.Fatal(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Test Create User: %v", err)
	}
}

func TestGetUser(t *testing.T) {
	testUser := testUser()

	db, mock, err := mock.NewDBMock()
	if err != nil {
		t.Errorf("Failed to initialize mock DB: %v", err)
	}

	rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "created_at", "updated_at"}).
		AddRow(1, testUser.Username, testUser.Email, testUser.Password, time.Now(), time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users" WHERE "users"."id" = $1 ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(1, 1).
		WillReturnRows(rows)

	repo := NewUserRepository(db)
	_, err = repo.GetUser(1)
	if err != nil {
		t.Fatalf("failed to get user: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Test Find User: %v", err)
	}
}

func TestGetUserByEmail(t *testing.T) {
	testUser := testUser()

	db, mock, err := mock.NewDBMock()
	if err != nil {
		t.Errorf("Failed to initialize mock DB: %v", err)
	}

	rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "created_at", "updated_at"}).
		AddRow(1, testUser.Username, testUser.Email, testUser.Password, time.Now(), time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users" WHERE email=$1 ORDER BY "users"."id" LIMIT $2`)).
		WithArgs("test@example.com", 1).
		WillReturnRows(rows)

	repo := NewUserRepository(db)
	_, err = repo.GetUserByEmail("test@example.com")
	if err != nil {
		t.Fatalf("failed to get user by email: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Test Find User by email: %v", err)
	}
}

func TestListUsers(t *testing.T) {
	testUser1 := testUser()
	testUser2 := testUser2()

	db, mock, err := mock.NewDBMock()
	if err != nil {
		t.Errorf("Failed to initialize mock DB: %v", err)
	}

	rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "created_at", "updated_at"}).
		AddRow(1, testUser1.Username, testUser1.Email, testUser1.Password, time.Now(), time.Now()).
		AddRow(2, testUser2.Username, testUser2.Email, testUser2.Password, time.Now(), time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users"`)).
		WillReturnRows(rows)

	repo := NewUserRepository(db)
	_, err = repo.ListUsers(1, 2)
	if err != nil {
		t.Fatalf("failed to list user: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Test Find User: %v", err)
	}
}

func TestUpdateUser(t *testing.T) {
	testUser := testUser()

	db, mock, err := mock.NewDBMock()
	if err != nil {
		t.Errorf("Failed to initialize mock DB: %v", err)
	}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "users" ("username","email","password","created_at","updated_at") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
		WillReturnRows(rows)
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		`UPDATE "users" SET "username"=$1,"email"=$2,"password"=$3,"created_at"=$4,"updated_at"=$5 WHERE "id" = $6 RETURNING *`)).
		WillReturnRows(rows)
	mock.ExpectCommit()

	repo := NewUserRepository(db)
	err = repo.CreateUser(testUser)
	if err != nil {
		t.Fatalf("failed to create user: %s", err)
	}
	err = repo.UpdateUser(testUser)
	if err != nil {
		t.Fatalf("failed to list user: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Test Find User: %v", err)
	}
}

func TestDeleteUser(t *testing.T) {
	db, mock, err := mock.NewDBMock()
	if err != nil {
		t.Errorf("Failed to initialize mock DB: %v", err)
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM "users" WHERE "users"."id" = $1`)).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := NewUserRepository(db)
	err = repo.DeleteUser(1)
	if err != nil {
		t.Fatalf("failed to delete user: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Test Find User: %v", err)
	}
}
