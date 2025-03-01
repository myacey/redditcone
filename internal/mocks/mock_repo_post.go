// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/repository/post_repository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	models "github.com/myacey/redditclone/internal/models"
)

// MockPostRepository is a mock of PostRepository interface.
type MockPostRepository struct {
	ctrl     *gomock.Controller
	recorder *MockPostRepositoryMockRecorder
}

// MockPostRepositoryMockRecorder is the mock recorder for MockPostRepository.
type MockPostRepositoryMockRecorder struct {
	mock *MockPostRepository
}

// NewMockPostRepository creates a new mock instance.
func NewMockPostRepository(ctrl *gomock.Controller) *MockPostRepository {
	mock := &MockPostRepository{ctrl: ctrl}
	mock.recorder = &MockPostRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPostRepository) EXPECT() *MockPostRepositoryMockRecorder {
	return m.recorder
}

// CreatePost mocks base method.
func (m *MockPostRepository) CreatePost(ctx context.Context, newPost *models.Post) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePost", ctx, newPost)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreatePost indicates an expected call of CreatePost.
func (mr *MockPostRepositoryMockRecorder) CreatePost(ctx, newPost interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePost", reflect.TypeOf((*MockPostRepository)(nil).CreatePost), ctx, newPost)
}

// DeletePost mocks base method.
func (m *MockPostRepository) DeletePost(ctx context.Context, postID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePost", ctx, postID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePost indicates an expected call of DeletePost.
func (mr *MockPostRepositoryMockRecorder) DeletePost(ctx, postID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePost", reflect.TypeOf((*MockPostRepository)(nil).DeletePost), ctx, postID)
}

// GetAllPosts mocks base method.
func (m *MockPostRepository) GetAllPosts(ctx context.Context) ([]*models.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllPosts", ctx)
	ret0, _ := ret[0].([]*models.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllPosts indicates an expected call of GetAllPosts.
func (mr *MockPostRepositoryMockRecorder) GetAllPosts(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllPosts", reflect.TypeOf((*MockPostRepository)(nil).GetAllPosts), ctx)
}

// GetPostByID mocks base method.
func (m *MockPostRepository) GetPostByID(ctx context.Context, postID string) (*models.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPostByID", ctx, postID)
	ret0, _ := ret[0].(*models.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPostByID indicates an expected call of GetPostByID.
func (mr *MockPostRepositoryMockRecorder) GetPostByID(ctx, postID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPostByID", reflect.TypeOf((*MockPostRepository)(nil).GetPostByID), ctx, postID)
}

// UpdatePostInfo mocks base method.
func (m *MockPostRepository) UpdatePostInfo(ctx context.Context, updatedPost *models.Post) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePostInfo", ctx, updatedPost)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdatePostInfo indicates an expected call of UpdatePostInfo.
func (mr *MockPostRepositoryMockRecorder) UpdatePostInfo(ctx, updatedPost interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePostInfo", reflect.TypeOf((*MockPostRepository)(nil).UpdatePostInfo), ctx, updatedPost)
}
