package service

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/myacey/redditclone/internal/mocks"
	"github.com/myacey/redditclone/internal/models"
)

var ErrBasic = errors.New("some error")

var (
	mockUser       = models.NewUser("testuser", "qwerty123")
	mockSession    = models.NewSession("token")
	mockSinglePost = models.NewPost(mockUser, "music", "mock title", "text", "mock text", "")
	mockPosts      = []*models.Post{mockSinglePost, mockSinglePost, mockSinglePost}
)

func TestGetUserFromDBByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockPostRepo := mocks.NewMockPostRepository(ctrl)
	mockCommentRepo := mocks.NewMockCommentRepository(ctrl)
	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
	mockTokenMaker := mocks.NewMockTokenMaker(ctrl)
	mockLogger := zap.NewNop().Sugar()

	service := &Service{
		userRepo:    mockUserRepo,
		postRepo:    mockPostRepo,
		commentRepo: mockCommentRepo,
		sessionRepo: mockSessionRepo,
		tokenMaker:  mockTokenMaker,
		logger:      mockLogger,
	}

	testCases := []struct {
		name       string
		userID     string
		mockSetup  func()
		expRes     interface{}
		wantErrMsg string
	}{
		{
			name:   "Success",
			userID: mockUser.ID,
			mockSetup: func() {
				mockUserRepo.EXPECT().GetUserByID(gomock.Any(), mockUser.ID).Return(mockUser, nil)
			},
			expRes:     mockUser,
			wantErrMsg: "",
		},
		{
			name:   "Error",
			userID: mockUser.ID,
			mockSetup: func() {
				mockUserRepo.EXPECT().GetUserByID(gomock.Any(), mockUser.ID).Return(nil, ErrBasic)
			},
			expRes:     nil,
			wantErrMsg: ErrBasic.Error(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			res, err := service.GetUserFromDBByID(context.Background(), tc.userID)
			if tc.wantErrMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.wantErrMsg)
			}
			if tc.expRes == nil {
				assert.Nil(t, res)
			} else {
				assert.Equal(t, tc.expRes, res)
			}
		})
	}
}

func TestGetUserFromDBByUsername(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockPostRepo := mocks.NewMockPostRepository(ctrl)
	mockCommentRepo := mocks.NewMockCommentRepository(ctrl)
	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
	mockTokenMaker := mocks.NewMockTokenMaker(ctrl)
	mockLogger := zap.NewNop().Sugar()

	service := &Service{
		userRepo:    mockUserRepo,
		postRepo:    mockPostRepo,
		commentRepo: mockCommentRepo,
		sessionRepo: mockSessionRepo,
		tokenMaker:  mockTokenMaker,
		logger:      mockLogger,
	}

	testCases := []struct {
		name       string
		username   string
		mockSetup  func()
		expRes     interface{}
		wantErrMsg string
	}{
		{
			name:     "Success",
			username: mockUser.Username,
			mockSetup: func() {
				mockUserRepo.EXPECT().GetUserByUsername(gomock.Any(), mockUser.Username).Return(mockUser, nil)
			},
			expRes:     mockUser,
			wantErrMsg: "",
		},
		{
			name:     "Error",
			username: mockUser.Username,
			mockSetup: func() {
				mockUserRepo.EXPECT().GetUserByUsername(gomock.Any(), mockUser.Username).Return(nil, ErrBasic)
			},
			expRes:     nil,
			wantErrMsg: ErrBasic.Error(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			res, err := service.GetUserFromDBByUsername(context.Background(), tc.username)
			if tc.expRes == nil {
				assert.Nil(t, res)
			} else {
				assert.Equal(t, tc.expRes, res)
			}
			if tc.wantErrMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.wantErrMsg)
			}
		})
	}
}

func TestCreateNewUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockPostRepo := mocks.NewMockPostRepository(ctrl)
	mockCommentRepo := mocks.NewMockCommentRepository(ctrl)
	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
	mockTokenMaker := mocks.NewMockTokenMaker(ctrl)
	mockLogger := zap.NewNop().Sugar()

	service := &Service{
		userRepo:    mockUserRepo,
		postRepo:    mockPostRepo,
		commentRepo: mockCommentRepo,
		sessionRepo: mockSessionRepo,
		tokenMaker:  mockTokenMaker,
		logger:      mockLogger,
	}

	testCases := []struct {
		name       string
		userToAdd  *models.User
		mockSetup  func()
		expRes     interface{}
		wantErrMsg string
	}{
		{
			name:      "Success",
			userToAdd: mockUser,
			mockSetup: func() {
				mockUserRepo.EXPECT().CreateUser(gomock.Any(), mockUser).Return(nil)
				mockTokenMaker.EXPECT().CreateToken(gomock.Any()).Return(mockSession.Token, nil)
				mockSessionRepo.EXPECT().CreateSession(gomock.Any(), mockSession, mockUser.ID, gomock.Any()).Return(nil)
			},
			expRes:     mockSession,
			wantErrMsg: "",
		},
		{
			name:      "Repo Error",
			userToAdd: mockUser,
			mockSetup: func() {
				mockUserRepo.EXPECT().CreateUser(gomock.Any(), mockUser).Return(ErrBasic)
				// mockTokenMaker.EXPECT().CreateToken(gomock.Any()).Return(mockSession.Token, nil)
				// mockSessionRepo.EXPECT().CreateSession(gomock.Any(), mockSession, mockUser.ID, gomock.Any()).Return(nil)
			},
			expRes:     nil,
			wantErrMsg: "internal error",
		},
		{
			name:      "Token Maker Error",
			userToAdd: mockUser,
			mockSetup: func() {
				mockUserRepo.EXPECT().CreateUser(gomock.Any(), mockUser).Return(nil)
				mockTokenMaker.EXPECT().CreateToken(gomock.Any()).Return("", ErrBasic)
				// mockSessionRepo.EXPECT().CreateSession(gomock.Any(), mockSession, mockUser.ID, gomock.Any()).Return(nil)
			},
			expRes:     nil,
			wantErrMsg: "internal error",
		},
		{
			name:      "Session Error",
			userToAdd: mockUser,
			mockSetup: func() {
				mockUserRepo.EXPECT().CreateUser(gomock.Any(), mockUser).Return(nil)
				mockTokenMaker.EXPECT().CreateToken(gomock.Any()).Return(mockSession.Token, nil)
				mockSessionRepo.EXPECT().CreateSession(gomock.Any(), mockSession, mockUser.ID, gomock.Any()).Return(ErrBasic)
			},
			expRes:     nil,
			wantErrMsg: "internal error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			res, err := service.CreateNewUser(context.Background(), tc.userToAdd)
			if tc.expRes == nil {
				assert.Nil(t, res)
			} else {
				assert.Equal(t, tc.expRes, res)
			}
			if tc.wantErrMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.wantErrMsg)
			}
		})
	}
}

func TestLoginUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockPostRepo := mocks.NewMockPostRepository(ctrl)
	mockCommentRepo := mocks.NewMockCommentRepository(ctrl)
	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
	mockTokenMaker := mocks.NewMockTokenMaker(ctrl)
	mockLogger := zap.NewNop().Sugar()

	service := &Service{
		userRepo:    mockUserRepo,
		postRepo:    mockPostRepo,
		commentRepo: mockCommentRepo,
		sessionRepo: mockSessionRepo,
		tokenMaker:  mockTokenMaker,
		logger:      mockLogger,
	}

	testCases := []struct {
		name       string
		username   string
		mockSetup  func()
		expRes     interface{}
		wantErrMsg string
	}{
		{
			name:     "Success",
			username: mockUser.Username,
			mockSetup: func() {
				mockUserRepo.EXPECT().GetUserByUsername(gomock.Any(), mockUser.Username).Return(mockUser, nil)
				mockTokenMaker.EXPECT().CreateToken(mockUser).Return("token", nil)
				mockSessionRepo.EXPECT().UpdateSessionToken(gomock.Any(), mockSession, mockUser.ID, gomock.Any()).Return(nil)
			},
			expRes:     mockSession,
			wantErrMsg: "",
		},
		{
			name:     "Err user repo",
			username: mockUser.Username,
			mockSetup: func() {
				mockUserRepo.EXPECT().GetUserByUsername(gomock.Any(), mockUser.Username).Return(nil, ErrBasic)
				// mockTokenMaker.EXPECT().CreateToken(mockUser).Return("token", nil)
				// mockSessionRepo.EXPECT().UpdateSessionToken(gomock.Any(), mockSession, mockUser.ID, gomock.Any()).Return(nil)
			},
			expRes:     nil,
			wantErrMsg: ErrBasic.Error(),
		},
		{
			name:     "Err token maker",
			username: mockUser.Username,
			mockSetup: func() {
				mockUserRepo.EXPECT().GetUserByUsername(gomock.Any(), mockUser.Username).Return(mockUser, nil)
				mockTokenMaker.EXPECT().CreateToken(mockUser).Return("", ErrBasic)
				// mockSessionRepo.EXPECT().UpdateSessionToken(gomock.Any(), mockSession, mockUser.ID, gomock.Any()).Return(nil)
			},
			expRes:     nil,
			wantErrMsg: ErrBasic.Error(),
		},
		{
			name:     "Err session repo",
			username: mockUser.Username,
			mockSetup: func() {
				mockUserRepo.EXPECT().GetUserByUsername(gomock.Any(), mockUser.Username).Return(mockUser, nil)
				mockTokenMaker.EXPECT().CreateToken(mockUser).Return("token", nil)
				mockSessionRepo.EXPECT().UpdateSessionToken(gomock.Any(), mockSession, mockUser.ID, gomock.Any()).Return(ErrBasic)
			},
			expRes:     nil,
			wantErrMsg: ErrBasic.Error(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			res, err := service.LoginUser(context.Background(), tc.username)
			if tc.expRes == nil {
				assert.Nil(t, res)
			} else {
				assert.Equal(t, tc.expRes, res)
			}
			if tc.wantErrMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.wantErrMsg)
			}
		})
	}
}

func TestGetAllPosts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockPostRepo := mocks.NewMockPostRepository(ctrl)
	mockCommentRepo := mocks.NewMockCommentRepository(ctrl)
	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
	mockTokenMaker := mocks.NewMockTokenMaker(ctrl)
	mockLogger := zap.NewNop().Sugar()

	service := &Service{
		userRepo:    mockUserRepo,
		postRepo:    mockPostRepo,
		commentRepo: mockCommentRepo,
		sessionRepo: mockSessionRepo,
		tokenMaker:  mockTokenMaker,
		logger:      mockLogger,
	}

	testCases := []struct {
		name       string
		mockSetup  func()
		expRes     interface{}
		wantErrMsg string
	}{
		{
			name: "Success",
			mockSetup: func() {
				mockPostRepo.EXPECT().GetAllPosts(gomock.Any()).Return(mockPosts, nil)
			},
			expRes:     mockPosts,
			wantErrMsg: "",
		},
		{
			name: "Err post repo",
			mockSetup: func() {
				mockPostRepo.EXPECT().GetAllPosts(gomock.Any()).Return(nil, ErrBasic)
			},
			expRes:     nil,
			wantErrMsg: ErrBasic.Error(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			res, err := service.GetAllPosts(context.Background())
			if tc.expRes == nil {
				assert.Nil(t, res)
			} else {
				assert.Equal(t, tc.expRes, res)
			}
			if tc.wantErrMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.wantErrMsg)
			}
		})
	}
}

func TestAddPost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockPostRepo := mocks.NewMockPostRepository(ctrl)
	mockCommentRepo := mocks.NewMockCommentRepository(ctrl)
	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
	mockTokenMaker := mocks.NewMockTokenMaker(ctrl)
	mockLogger := zap.NewNop().Sugar()

	service := &Service{
		userRepo:    mockUserRepo,
		postRepo:    mockPostRepo,
		commentRepo: mockCommentRepo,
		sessionRepo: mockSessionRepo,
		tokenMaker:  mockTokenMaker,
		logger:      mockLogger,
	}

	testCases := []struct {
		name       string
		newPost    *models.Post
		mockSetup  func()
		wantErrMsg string
	}{
		{
			name:    "Success",
			newPost: mockSinglePost,
			mockSetup: func() {
				mockPostRepo.EXPECT().CreatePost(gomock.Any(), mockSinglePost).Return(nil)
			},
			wantErrMsg: "",
		},
		{
			name:    "Err invalid post data",
			newPost: &models.Post{Category: "invalid"},
			mockSetup: func() {
				// mockPostRepo.EXPECT().CreatePost(gomock.Any(), mockSinglePost).Return(nil)
			},
			wantErrMsg: ErrInvalidPostData.Error(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			err := service.AddPost(context.Background(), tc.newPost)
			if tc.wantErrMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.wantErrMsg)
			}
		})
	}
}

func TestGetPostByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockPostRepo := mocks.NewMockPostRepository(ctrl)
	mockCommentRepo := mocks.NewMockCommentRepository(ctrl)
	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
	mockTokenMaker := mocks.NewMockTokenMaker(ctrl)
	mockLogger := zap.NewNop().Sugar()

	service := &Service{
		userRepo:    mockUserRepo,
		postRepo:    mockPostRepo,
		commentRepo: mockCommentRepo,
		sessionRepo: mockSessionRepo,
		tokenMaker:  mockTokenMaker,
		logger:      mockLogger,
	}

	testCases := []struct {
		name         string
		postID       string
		increaseVote bool
		mockSetup    func()
		expRes       *models.Post
		wantErrMsg   string
	}{
		{
			name:         "Success",
			postID:       mockSinglePost.ID,
			increaseVote: true,
			mockSetup: func() {
				mockPostRepo.EXPECT().GetPostByID(gomock.Any(), mockSinglePost.ID).Return(mockSinglePost, nil)
				mockPostRepo.EXPECT().UpdatePostInfo(gomock.Any(), mockSinglePost).Return(nil) // increate vote
				mockCommentRepo.EXPECT().GetCommentsByPostID(gomock.Any(), mockSinglePost.ID).Return(nil, nil)
			},
			expRes:     mockSinglePost,
			wantErrMsg: "",
		},
		{
			name:         "Err post repo get user",
			postID:       mockSinglePost.ID,
			increaseVote: true,
			mockSetup: func() {
				mockPostRepo.EXPECT().GetPostByID(gomock.Any(), mockSinglePost.ID).Return(nil, ErrBasic)
				// mockPostRepo.EXPECT().UpdatePostInfo(gomock.Any(), mockSinglePost).Return(nil) // increate vote
				// mockCommentRepo.EXPECT().GetCommentsByPostID(gomock.Any(), mockSinglePost.ID).Return(nil, nil)
			},
			expRes:     nil,
			wantErrMsg: "cant find post",
		},
		{
			name:         "Err increase votes",
			postID:       mockSinglePost.ID,
			increaseVote: true,
			mockSetup: func() {
				mockPostRepo.EXPECT().GetPostByID(gomock.Any(), mockSinglePost.ID).Return(mockSinglePost, nil)
				mockPostRepo.EXPECT().UpdatePostInfo(gomock.Any(), mockSinglePost).Return(ErrBasic) // increate vote
				// mockCommentRepo.EXPECT().GetCommentsByPostID(gomock.Any(), mockSinglePost.ID).Return(nil, nil)
			},
			expRes:     nil,
			wantErrMsg: ErrBasic.Error(),
		},
		{
			name:         "Err comment repo",
			postID:       mockSinglePost.ID,
			increaseVote: true,
			mockSetup: func() {
				mockPostRepo.EXPECT().GetPostByID(gomock.Any(), mockSinglePost.ID).Return(mockSinglePost, nil)
				mockPostRepo.EXPECT().UpdatePostInfo(gomock.Any(), mockSinglePost).Return(nil) // increate vote
				mockCommentRepo.EXPECT().GetCommentsByPostID(gomock.Any(), mockSinglePost.ID).Return(nil, ErrBasic)
			},
			expRes:     nil,
			wantErrMsg: "cant find comments",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			res, err := service.GetPostByID(context.Background(), tc.postID, tc.increaseVote)
			if tc.expRes == nil {
				assert.Nil(t, res)
			} else {
				assert.Equal(t, tc.expRes, res)
			}
			if tc.wantErrMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.wantErrMsg)
			}
		})
	}
}

func TestGetPostsByAuthor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockPostRepo := mocks.NewMockPostRepository(ctrl)
	mockCommentRepo := mocks.NewMockCommentRepository(ctrl)
	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
	mockTokenMaker := mocks.NewMockTokenMaker(ctrl)
	mockLogger := zap.NewNop().Sugar()

	service := &Service{
		userRepo:    mockUserRepo,
		postRepo:    mockPostRepo,
		commentRepo: mockCommentRepo,
		sessionRepo: mockSessionRepo,
		tokenMaker:  mockTokenMaker,
		logger:      mockLogger,
	}

	testCases := []struct {
		name       string
		username   string
		mockSetup  func()
		expRes     interface{}
		wantErrMsg string
	}{
		{
			name:     "Success",
			username: mockUser.Username,
			mockSetup: func() {
				mockUserRepo.EXPECT().GetUserByUsername(gomock.Any(), mockUser.Username).Return(mockUser, nil)
				mockPostRepo.EXPECT().GetAllPosts(gomock.Any()).Return(mockPosts, nil)
			},
			expRes:     mockPosts,
			wantErrMsg: "",
		},
		{
			name:     "Err user repo",
			username: mockUser.Username,
			mockSetup: func() {
				mockUserRepo.EXPECT().GetUserByUsername(gomock.Any(), mockUser.Username).Return(nil, ErrBasic)
				// mockPostRepo.EXPECT().GetAllPosts(gomock.Any()).Return(mockPosts, nil)
			},
			expRes:     nil,
			wantErrMsg: ErrBasic.Error(),
		},
		{
			name:     "Err post repo",
			username: mockUser.Username,
			mockSetup: func() {
				mockUserRepo.EXPECT().GetUserByUsername(gomock.Any(), mockUser.Username).Return(mockUser, nil)
				mockPostRepo.EXPECT().GetAllPosts(gomock.Any()).Return(nil, ErrBasic)
			},
			expRes:     nil,
			wantErrMsg: ErrBasic.Error(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			res, err := service.GetPostsByAuthor(context.Background(), tc.username)
			if tc.expRes == nil {
				assert.Nil(t, res)
			} else {
				assert.Equal(t, tc.expRes, res)
			}
			if tc.wantErrMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.wantErrMsg)
			}
		})
	}
}

func TestGetPostsByCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockPostRepo := mocks.NewMockPostRepository(ctrl)
	mockCommentRepo := mocks.NewMockCommentRepository(ctrl)
	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
	mockTokenMaker := mocks.NewMockTokenMaker(ctrl)
	mockLogger := zap.NewNop().Sugar()

	service := &Service{
		userRepo:    mockUserRepo,
		postRepo:    mockPostRepo,
		commentRepo: mockCommentRepo,
		sessionRepo: mockSessionRepo,
		tokenMaker:  mockTokenMaker,
		logger:      mockLogger,
	}

	mockPostMusic := models.NewPost(mockUser, "music", "mock title", "text", "mock text", "")
	mockPostsMusic := []*models.Post{mockPostMusic, mockPostMusic, mockPostMusic}
	mockPostVideo := models.NewPost(mockUser, "video", "mock title", "text", "mock text", "")
	mockPostsVideo := []*models.Post{mockPostVideo, mockPostVideo, mockPostVideo}

	allMockedPosts := append([]*models.Post{}, mockPostsMusic...)
	allMockedPosts = append(allMockedPosts, mockPostsVideo...)

	testCases := []struct {
		name       string
		category   string
		mockSetup  func()
		expRes     interface{}
		wantErrMsg string
	}{
		{
			name:     "Success",
			category: "music",
			mockSetup: func() {
				mockPostRepo.EXPECT().GetAllPosts(gomock.Any()).Return(allMockedPosts, nil)
			},
			expRes:     mockPostsMusic,
			wantErrMsg: "",
		},
		{
			name:     "Err invalid category",
			category: "invalid",
			mockSetup: func() {
				// mockPostRepo.EXPECT().GetAllPosts(gomock.Any()).Return(allMockedPosts, nil)
			},
			expRes:     nil,
			wantErrMsg: "invalid category",
		},
		{
			name:     "Err post repo",
			category: "music",
			mockSetup: func() {
				mockPostRepo.EXPECT().GetAllPosts(gomock.Any()).Return(nil, ErrBasic)
			},
			expRes:     nil,
			wantErrMsg: ErrBasic.Error(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			res, err := service.GetPostsByCategory(context.Background(), tc.category)
			if tc.expRes == nil {
				assert.Nil(t, res)
			} else {
				assert.Equal(t, tc.expRes, res)
			}
			if tc.wantErrMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.wantErrMsg)
			}
		})
	}
}

func TestDeletePostWithID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockPostRepo := mocks.NewMockPostRepository(ctrl)
	mockCommentRepo := mocks.NewMockCommentRepository(ctrl)
	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
	mockTokenMaker := mocks.NewMockTokenMaker(ctrl)
	mockLogger := zap.NewNop().Sugar()

	service := &Service{
		userRepo:    mockUserRepo,
		postRepo:    mockPostRepo,
		commentRepo: mockCommentRepo,
		sessionRepo: mockSessionRepo,
		tokenMaker:  mockTokenMaker,
		logger:      mockLogger,
	}

	testCases := []struct {
		name         string
		postID       string
		increaseVote bool
		mockSetup    func()
		wantErrMsg   string
	}{
		{
			name:         "Success",
			postID:       mockSinglePost.ID,
			increaseVote: true,
			mockSetup: func() {
				mockPostRepo.EXPECT().DeletePost(gomock.Any(), mockSinglePost.ID).Return(nil)
			},
			wantErrMsg: "",
		},
		{
			name:         "Error",
			postID:       mockSinglePost.ID,
			increaseVote: true,
			mockSetup: func() {
				mockPostRepo.EXPECT().DeletePost(gomock.Any(), mockSinglePost.ID).Return(ErrBasic)
			},
			wantErrMsg: ErrBasic.Error(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			err := service.DeletePostWithID(context.Background(), tc.postID)
			if tc.wantErrMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.wantErrMsg)
			}
		})
	}
}
