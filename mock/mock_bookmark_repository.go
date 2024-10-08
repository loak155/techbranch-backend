// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/repository/bookmark_repository.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/loak155/techbranch-backend/internal/domain"
)

// MockIBookmarkRepository is a mock of IBookmarkRepository interface.
type MockIBookmarkRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIBookmarkRepositoryMockRecorder
}

// MockIBookmarkRepositoryMockRecorder is the mock recorder for MockIBookmarkRepository.
type MockIBookmarkRepositoryMockRecorder struct {
	mock *MockIBookmarkRepository
}

// NewMockIBookmarkRepository creates a new mock instance.
func NewMockIBookmarkRepository(ctrl *gomock.Controller) *MockIBookmarkRepository {
	mock := &MockIBookmarkRepository{ctrl: ctrl}
	mock.recorder = &MockIBookmarkRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIBookmarkRepository) EXPECT() *MockIBookmarkRepositoryMockRecorder {
	return m.recorder
}

// CreateBookmark mocks base method.
func (m *MockIBookmarkRepository) CreateBookmark(bookmark *domain.Bookmark) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateBookmark", bookmark)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateBookmark indicates an expected call of CreateBookmark.
func (mr *MockIBookmarkRepositoryMockRecorder) CreateBookmark(bookmark interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateBookmark", reflect.TypeOf((*MockIBookmarkRepository)(nil).CreateBookmark), bookmark)
}

// DeleteBookmarkByArticleID mocks base method.
func (m *MockIBookmarkRepository) DeleteBookmarkByArticleID(ArticleID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBookmarkByArticleID", ArticleID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteBookmarkByArticleID indicates an expected call of DeleteBookmarkByArticleID.
func (mr *MockIBookmarkRepositoryMockRecorder) DeleteBookmarkByArticleID(ArticleID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBookmarkByArticleID", reflect.TypeOf((*MockIBookmarkRepository)(nil).DeleteBookmarkByArticleID), ArticleID)
}

// DeleteBookmarkByUserID mocks base method.
func (m *MockIBookmarkRepository) DeleteBookmarkByUserID(UserID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBookmarkByUserID", UserID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteBookmarkByUserID indicates an expected call of DeleteBookmarkByUserID.
func (mr *MockIBookmarkRepositoryMockRecorder) DeleteBookmarkByUserID(UserID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBookmarkByUserID", reflect.TypeOf((*MockIBookmarkRepository)(nil).DeleteBookmarkByUserID), UserID)
}

// DeleteBookmarkByUserIDAndArticleID mocks base method.
func (m *MockIBookmarkRepository) DeleteBookmarkByUserIDAndArticleID(userID, articleID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBookmarkByUserIDAndArticleID", userID, articleID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteBookmarkByUserIDAndArticleID indicates an expected call of DeleteBookmarkByUserIDAndArticleID.
func (mr *MockIBookmarkRepositoryMockRecorder) DeleteBookmarkByUserIDAndArticleID(userID, articleID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBookmarkByUserIDAndArticleID", reflect.TypeOf((*MockIBookmarkRepository)(nil).DeleteBookmarkByUserIDAndArticleID), userID, articleID)
}

// GetBookmarkCountByArticleID mocks base method.
func (m *MockIBookmarkRepository) GetBookmarkCountByArticleID(articleID int) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBookmarkCountByArticleID", articleID)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBookmarkCountByArticleID indicates an expected call of GetBookmarkCountByArticleID.
func (mr *MockIBookmarkRepositoryMockRecorder) GetBookmarkCountByArticleID(articleID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBookmarkCountByArticleID", reflect.TypeOf((*MockIBookmarkRepository)(nil).GetBookmarkCountByArticleID), articleID)
}

// ListBookmarksByArticleID mocks base method.
func (m *MockIBookmarkRepository) ListBookmarksByArticleID(articleID int) (*[]domain.Bookmark, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListBookmarksByArticleID", articleID)
	ret0, _ := ret[0].(*[]domain.Bookmark)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListBookmarksByArticleID indicates an expected call of ListBookmarksByArticleID.
func (mr *MockIBookmarkRepositoryMockRecorder) ListBookmarksByArticleID(articleID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListBookmarksByArticleID", reflect.TypeOf((*MockIBookmarkRepository)(nil).ListBookmarksByArticleID), articleID)
}

// ListBookmarksByUserID mocks base method.
func (m *MockIBookmarkRepository) ListBookmarksByUserID(userID int) (*[]domain.Bookmark, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListBookmarksByUserID", userID)
	ret0, _ := ret[0].(*[]domain.Bookmark)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListBookmarksByUserID indicates an expected call of ListBookmarksByUserID.
func (mr *MockIBookmarkRepositoryMockRecorder) ListBookmarksByUserID(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListBookmarksByUserID", reflect.TypeOf((*MockIBookmarkRepository)(nil).ListBookmarksByUserID), userID)
}
