package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/myacey/redditclone/internal/handlers"
	"github.com/myacey/redditclone/internal/mocks"
	"github.com/myacey/redditclone/internal/models"
)

var ErrBasic = errors.New("some error")

var (
	mockToken   = "valid-token"
	mockUser    = models.NewUser("mockuser", "password")
	mockSession = models.NewSession("token")
	mockPost    = models.NewPost(mockUser, "music", "title", "text", "postText", "")
	mockPosts   = []*models.Post{mockPost, mockPost, mockPost}
)

func TestRegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockServiceInterface(ctrl)
	mockLogger := zap.NewNop().Sugar()
	mockTokenMaker := mocks.NewMockTokenMaker(ctrl)

	handler := handlers.NewHandler(mockService, mockLogger, mockTokenMaker)

	testCases := []struct {
		name           string
		reqBody        interface{}
		mockSetup      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Successful register",
			reqBody: handlers.RegisterRequest{
				Username: "testuser",
				Password: "qwerty123",
			},
			mockSetup: func() {
				mockService.EXPECT().CreateNewUser(gomock.Any(), gomock.Any()).Return(mockSession, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"token":"token"}`,
		},
		{
			name:    "Invalid JSON",
			reqBody: "invalid",
			mockSetup: func() {
				// mockService.EXPECT().CreateNewUser(gomock.Any(), gomock.Any()).Return(mockSession, nil)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"bad json"}`,
		},
		{
			name: "Service error",
			reqBody: handlers.RegisterRequest{
				Username: "testuser",
				Password: "qwerty123",
			},
			mockSetup: func() {
				mockService.EXPECT().CreateNewUser(gomock.Any(), gomock.Any()).Return(nil, ErrBasic)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"some error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			var reqBody bytes.Buffer
			err := json.NewEncoder(&reqBody).Encode(tc.reqBody)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/api/register", &reqBody)
			w := httptest.NewRecorder()
			handler.RegisterUser(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			var respBody map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&respBody)
			assert.NoError(t, err)

			expBody := map[string]interface{}{}
			err = json.Unmarshal([]byte(tc.expectedBody), &expBody)
			assert.NoError(t, err)
			assert.Equal(t, expBody, respBody)
		})
	}
}

func TestLoginRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockServiceInterface(ctrl)
	mockLogger := zap.NewNop().Sugar()
	mockTokenMaker := mocks.NewMockTokenMaker(ctrl)

	handler := handlers.NewHandler(mockService, mockLogger, mockTokenMaker)

	testCases := []struct {
		name           string
		reqBody        interface{}
		mockSetup      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Successful register",
			reqBody: handlers.LoginRequest{
				Username: "testuser",
				Password: "qwerty123",
			},
			mockSetup: func() {
				mockService.EXPECT().LoginUser(gomock.Any(), gomock.Any()).Return(mockSession, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"token":"token"}`,
		},
		{
			name:    "Invalid JSON",
			reqBody: "invalid",
			mockSetup: func() {
				// mockService.EXPECT().LoginUser(gomock.Any(), gomock.Any()).Return(mockSession, nil)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"bad json"}`,
		},
		{
			name: "Service error",
			reqBody: handlers.LoginRequest{
				Username: "testuser",
				Password: "qwerty123",
			},
			mockSetup: func() {
				mockService.EXPECT().LoginUser(gomock.Any(), gomock.Any()).Return(nil, ErrBasic)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"some error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			var reqBody bytes.Buffer
			err := json.NewEncoder(&reqBody).Encode(tc.reqBody)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/api/login", &reqBody)
			w := httptest.NewRecorder()
			handler.LoginUser(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			var respBody map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&respBody)
			assert.NoError(t, err)

			expBody := map[string]interface{}{}
			err = json.Unmarshal([]byte(tc.expectedBody), &expBody)
			assert.NoError(t, err)
			assert.Equal(t, expBody, respBody)
		})
	}
}

func TestGetPosts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockServiceInterface(ctrl)
	mockLogger := zap.NewNop().Sugar()
	mockTokenMaker := mocks.NewMockTokenMaker(ctrl)

	handler := handlers.NewHandler(mockService, mockLogger, mockTokenMaker)

	allPosts, err := json.Marshal(mockPosts)
	assert.NoError(t, err)

	testCases := []struct {
		name           string
		mockSetup      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Successful get posts",
			mockSetup: func() {
				mockService.EXPECT().GetAllPosts(gomock.Any()).Return(mockPosts, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   string(allPosts),
		},
		{
			name: "Invalid Service",
			mockSetup: func() {
				mockService.EXPECT().GetAllPosts(gomock.Any()).Return(nil, ErrBasic)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"some error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			req := httptest.NewRequest(http.MethodGet, "/api/posts", nil)
			w := httptest.NewRecorder()
			handler.GetPosts(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.JSONEq(t, tc.expectedBody, string(body))
		})
	}
}

func TestAddPosts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockServiceInterface(ctrl)
	mockLogger := zap.NewNop().Sugar()
	mockTokenMaker := mocks.NewMockTokenMaker(ctrl)

	handler := handlers.NewHandler(mockService, mockLogger, mockTokenMaker)

	expPost, err := json.Marshal(mockPost)
	assert.NoError(t, err)

	testCases := []struct {
		name           string
		token          string
		reqBody        interface{}
		mockSetup      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:  "Successful add posts",
			token: mockToken,
			reqBody: handlers.AddPostRequest{
				Category: mockPost.Category,
				Title:    mockPost.Title,
				Type:     mockPost.Type,
				Text:     mockPost.Text,
				URL:      mockPost.URL,
			},
			mockSetup: func() {
				mockService.EXPECT().AddPost(gomock.Any(), gomock.Any()).Return(nil)
				mockService.EXPECT().GetUserFromDBByID(gomock.Any(), gomock.Any()).Return(mockUser, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   string(expPost),
		},
		{
			name:    "Bad JSON",
			token:   mockToken,
			reqBody: "invalid",
			mockSetup: func() {
				// mockService.EXPECT().AddPost(gomock.Any(), gomock.Any()).Return(nil)
				mockService.EXPECT().GetUserFromDBByID(gomock.Any(), gomock.Any()).Return(mockUser, nil)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"bad json"}`,
		},
		{
			name:  "Service error",
			token: mockToken,
			reqBody: handlers.AddPostRequest{
				Category: mockPost.Category,
				Title:    mockPost.Title,
				Type:     mockPost.Type,
				Text:     mockPost.Text,
				URL:      mockPost.URL,
			},
			mockSetup: func() {
				mockService.EXPECT().AddPost(gomock.Any(), gomock.Any()).Return(ErrBasic)
				mockService.EXPECT().GetUserFromDBByID(gomock.Any(), gomock.Any()).Return(mockUser, nil)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"failed to add post"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			var reqBody bytes.Buffer
			err := json.NewEncoder(&reqBody).Encode(tc.reqBody)
			assert.NoError(t, err)

			ctx := context.Background()
			ctx = context.WithValue(ctx, handlers.UserIDCtxKeyValue, mockToken)
			req := httptest.NewRequest(http.MethodPost, "/api/posts", &reqBody).WithContext(ctx)

			w := httptest.NewRecorder()
			handler.AddPost(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			if strings.Contains(tc.expectedBody, `"message"`) {
				// Ожидается ответ с полем "message" (например, ошибка)
				var gotMessage map[string]string
				err := json.Unmarshal(body, &gotMessage)
				assert.NoError(t, err)
				assert.JSONEq(t, tc.expectedBody, string(body), "expected error message response")
			} else {
				// Ожидается ответ типа post
				var gotPost models.Post
				err := json.Unmarshal(body, &gotPost)
				assert.NoError(t, err)

				// Проверка значимых полей post с учетом CreatedAt и ID
				if time.Since(gotPost.CreatedAt) < time.Minute {
					gotPost.CreatedAt = mockPost.CreatedAt
					gotPost.ID = mockPost.ID
				}

				b, err := json.Marshal(gotPost)
				assert.NoError(t, err)
				assert.JSONEq(t, tc.expectedBody, string(b), "expected post response")
			}
		})
	}
}

func TestGetSinglePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockServiceInterface(ctrl)
	mockLogger := zap.NewNop().Sugar()
	mockTokenMaker := mocks.NewMockTokenMaker(ctrl)

	handler := handlers.NewHandler(mockService, mockLogger, mockTokenMaker)

	marshalledPost, err := json.Marshal(mockPost)
	assert.NoError(t, err)

	testCases := []struct {
		name           string
		mockSetup      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Successful get post",
			mockSetup: func() {
				mockService.EXPECT().GetPostByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockPost, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   string(marshalledPost),
		},
		{
			name: "Service err",
			mockSetup: func() {
				mockService.EXPECT().GetPostByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, ErrBasic)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"some error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			req := httptest.NewRequest(http.MethodGet, "/api/post/", nil)
			w := httptest.NewRecorder()
			handler.GetPost(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			assert.JSONEq(t, tc.expectedBody, string(body))
		})
	}
}

func TestGetPostsByCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockServiceInterface(ctrl)
	mockLogger := zap.NewNop().Sugar()
	mockTokenMaker := mocks.NewMockTokenMaker(ctrl)

	handler := handlers.NewHandler(mockService, mockLogger, mockTokenMaker)

	allPosts, err := json.Marshal(mockPosts)
	assert.NoError(t, err)

	testCases := []struct {
		name           string
		mockSetup      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Successful get posts",
			mockSetup: func() {
				mockService.EXPECT().GetPostsByCategory(gomock.Any(), gomock.Any()).Return(mockPosts, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   string(allPosts),
		},
		{
			name: "Service error",
			mockSetup: func() {
				mockService.EXPECT().GetPostsByCategory(gomock.Any(), gomock.Any()).Return(nil, ErrBasic)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"some error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			req := httptest.NewRequest(http.MethodGet, "/api/posts/music", nil)
			w := httptest.NewRecorder()
			handler.GetPostsByCategory(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.JSONEq(t, tc.expectedBody, string(body))
		})
	}
}

func TestDeletePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockServiceInterface(ctrl)
	mockLogger := zap.NewNop().Sugar()
	mockTokenMaker := mocks.NewMockTokenMaker(ctrl)

	handler := handlers.NewHandler(mockService, mockLogger, mockTokenMaker)

	testCases := []struct {
		name           string
		mockSetup      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Successful delete",
			mockSetup: func() {
				mockService.EXPECT().DeletePostWithID(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"success"}`,
		},
		{
			name: "Service error",
			mockSetup: func() {
				mockService.EXPECT().DeletePostWithID(gomock.Any(), gomock.Any()).Return(ErrBasic)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"some error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			req := httptest.NewRequest(http.MethodDelete, "/api/post/0", nil)
			w := httptest.NewRecorder()
			handler.DeletePost(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.JSONEq(t, tc.expectedBody, string(body))
		})
	}
}

func TestGetUserPosts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockServiceInterface(ctrl)
	mockLogger := zap.NewNop().Sugar()
	mockTokenMaker := mocks.NewMockTokenMaker(ctrl)

	handler := handlers.NewHandler(mockService, mockLogger, mockTokenMaker)

	allPosts, err := json.Marshal(mockPosts)
	assert.NoError(t, err)

	testCases := []struct {
		name             string
		mockSetup        func()
		usernameToSearch string
		expectedStatus   int
		expectedBody     string
	}{
		{
			name:             "Successful get posts",
			usernameToSearch: mockUser.Username,
			mockSetup: func() {
				mockService.EXPECT().GetPostsByAuthor(gomock.Any(), gomock.Any()).Return(mockPosts, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   string(allPosts),
		},
		{
			name:             "Invalid username",
			usernameToSearch: "",
			mockSetup: func() {
				// mockService.EXPECT().GetPostsByAuthor(gomock.Any(), gomock.Any()).Return(mockPosts, nil)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"invalid username"}`,
		},
		{
			name:             "Service error",
			usernameToSearch: mockUser.Username,
			mockSetup: func() {
				mockService.EXPECT().GetPostsByAuthor(gomock.Any(), gomock.Any()).Return(nil, ErrBasic)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"some error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			req := httptest.NewRequest(http.MethodGet, "/api/user/{username}", nil)
			req = mux.SetURLVars(req, map[string]string{ // we have checkout that username==""?
				"username": tc.usernameToSearch,
			})
			w := httptest.NewRecorder()
			handler.GetUserPosts(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.JSONEq(t, tc.expectedBody, string(body))
		})
	}
}
