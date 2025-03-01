// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/repository/session_repository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	models "github.com/myacey/redditclone/internal/models"
)

// MockSessionRepository is a mock of SessionRepository interface.
type MockSessionRepository struct {
	ctrl     *gomock.Controller
	recorder *MockSessionRepositoryMockRecorder
}

// MockSessionRepositoryMockRecorder is the mock recorder for MockSessionRepository.
type MockSessionRepositoryMockRecorder struct {
	mock *MockSessionRepository
}

// NewMockSessionRepository creates a new mock instance.
func NewMockSessionRepository(ctrl *gomock.Controller) *MockSessionRepository {
	mock := &MockSessionRepository{ctrl: ctrl}
	mock.recorder = &MockSessionRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSessionRepository) EXPECT() *MockSessionRepositoryMockRecorder {
	return m.recorder
}

// CreateSession mocks base method.
func (m *MockSessionRepository) CreateSession(ctx context.Context, session *models.Session, username string, expirationTime time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSession", ctx, session, username, expirationTime)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateSession indicates an expected call of CreateSession.
func (mr *MockSessionRepositoryMockRecorder) CreateSession(ctx, session, username, expirationTime interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSession", reflect.TypeOf((*MockSessionRepository)(nil).CreateSession), ctx, session, username, expirationTime)
}

// GetSessionTokenByUsername mocks base method.
func (m *MockSessionRepository) GetSessionTokenByUsername(ctx context.Context, username string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSessionTokenByUsername", ctx, username)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSessionTokenByUsername indicates an expected call of GetSessionTokenByUsername.
func (mr *MockSessionRepositoryMockRecorder) GetSessionTokenByUsername(ctx, username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSessionTokenByUsername", reflect.TypeOf((*MockSessionRepository)(nil).GetSessionTokenByUsername), ctx, username)
}

// UpdateSessionToken mocks base method.
func (m *MockSessionRepository) UpdateSessionToken(ctx context.Context, newSession *models.Session, username string, expirationTime time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSessionToken", ctx, newSession, username, expirationTime)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateSessionToken indicates an expected call of UpdateSessionToken.
func (mr *MockSessionRepositoryMockRecorder) UpdateSessionToken(ctx, newSession, username, expirationTime interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSessionToken", reflect.TypeOf((*MockSessionRepository)(nil).UpdateSessionToken), ctx, newSession, username, expirationTime)
}
