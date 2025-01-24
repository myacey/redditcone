package mongorepo_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"

	"github.com/myacey/redditclone/internal/models"
	"github.com/myacey/redditclone/internal/repository/mongorepo"
)

var (
	ErrBasic  = errors.New("some error")
	mockUser  = models.NewUser("testuser", "qwerty123")
	mockPost  = models.NewPost(mockUser, "music", "mock title", "text", "mock text", "")
	mockPosts = []*models.Post{mockPost, mockPost, mockPost}
)

func setupMockDB(t *testing.T) *mtest.T {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	return mt
}

func TestCreatePost(t *testing.T) {
	mt := setupMockDB(t)

	mt.Run("Test Cases", func(mt *mtest.T) {
		repo := mongorepo.NewMongoPostRepository(mt.Client, "testDB", nil)

		testCases := []struct {
			name         string
			post         *models.Post
			mockBehavior func()
			expErr       error
		}{
			{
				name: "Success",
				post: mockPost,
				mockBehavior: func() {
					mt.AddMockResponses(mtest.CreateSuccessResponse())
				},
				expErr: nil,
			},
			{
				name: "Error",
				post: mockPost,
				mockBehavior: func() {
					mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
						Message: ErrBasic.Error(),
					}))
				},
				expErr: ErrBasic,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				tc.mockBehavior()

				err := repo.CreatePost(context.Background(), tc.post)
				if tc.expErr == nil {
					assert.NoError(t, err)
				} else {
					var cmdErr mongo.CommandError
					if errors.As(err, &cmdErr) {
						assert.Equal(t, tc.expErr.Error(), cmdErr.Message)
					} else {
						assert.EqualError(t, err, tc.expErr.Error())
					}
				}
			})
		}
	})
}

func TestGetPostByID(t *testing.T) {
	mt := setupMockDB(t)

	mt.Run("Test Cases", func(mt *mtest.T) {
		repo := mongorepo.NewMongoPostRepository(mt.Client, "testDB", nil)

		testCases := []struct {
			name         string
			postID       string
			mockBehavior func()
			expRes       *models.Post
			expErr       error
		}{
			{
				name:   "Success",
				postID: mockPost.ID,
				mockBehavior: func() {
					mt.AddMockResponses(mtest.CreateCursorResponse(1, "testDB.posts", mtest.FirstBatch, bson.D{
						{Key: "_id", Value: mockPost.ID},
						{Key: "title", Value: mockPost.Title},
					}))
				},
				expRes: &models.Post{ID: mockPost.ID, Title: mockPost.Title},
				expErr: nil,
			},
			{
				name:   "Err no docs",
				postID: mockPost.ID,
				mockBehavior: func() {
					mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{Message: ErrBasic.Error()}))
				},
				expRes: nil,
				expErr: ErrBasic,
			},
			{
				name:   "Error No Documents",
				postID: mockPost.ID,
				mockBehavior: func() {
					mt.AddMockResponses(mtest.CreateCursorResponse(0, "testDB.posts", mtest.FirstBatch))
				},
				expErr: nil, // bcs GetPostByID returns nil nil
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				tc.mockBehavior()

				res, err := repo.GetPostByID(context.Background(), tc.postID)
				assert.Equal(t, tc.expRes, res)
				if tc.expErr == nil {
					assert.NoError(t, err)
				} else {
					var cmdErr mongo.CommandError
					if errors.As(err, &cmdErr) {
						assert.Equal(t, tc.expErr.Error(), cmdErr.Message)
					} else {
						assert.EqualError(t, err, tc.expErr.Error())
					}
				}
			})
		}
	})
}

func toDoc(v interface{}, t *testing.T) (doc bson.D) {
	data, err := bson.Marshal(v)
	assert.NoError(t, err)

	err = bson.Unmarshal(data, &doc)
	assert.NoError(t, err)
	return doc
}

func toDocSlice(vs []*models.Post, t *testing.T) (docs []bson.D) {
	for _, v := range vs {
		doc := toDoc(v, t)
		docs = append(docs, doc)
	}
	return
}

func TestGetAllPosts(t *testing.T) {
	mt := setupMockDB(t)

	mt.Run("Test Cases", func(mt *mtest.T) {
		repo := mongorepo.NewMongoPostRepository(mt.Client, "testDB", nil)

		testCases := []struct {
			name         string
			mockBehavior func()
			expRes       []*models.Post
			expErr       error
		}{
			{
				name: "Success",
				mockBehavior: func() {
					first := mtest.CreateCursorResponse(1, "testDB.posts", mtest.FirstBatch, toDoc(mockPosts[0], t))
					getMore := mtest.CreateCursorResponse(1, "testDB.posts", mtest.NextBatch, toDocSlice(mockPosts[1:], t)...)
					killCursors := mtest.CreateCursorResponse(0, "testDB.posts", mtest.NextBatch)

					mt.AddMockResponses(first, getMore, killCursors)
				},
				expRes: mockPosts,
				expErr: nil,
			},
			{
				name: "Err post collection",
				mockBehavior: func() {
					mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{Message: ErrBasic.Error()}))
				},
				expRes: nil,
				expErr: ErrBasic,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				tc.mockBehavior()

				res, err := repo.GetAllPosts(context.Background())
				if tc.expRes != nil && res == nil && assert.ObjectsAreEqualValues(tc.expRes, res) {
					t.Fatalf("objects dont equal:\n %v\n%v\n----", tc.expErr, res)
				}
				if tc.expErr == nil {
					assert.NoError(t, err)
				} else {
					var cmdErr mongo.CommandError
					if errors.As(err, &cmdErr) {
						assert.Equal(t, tc.expErr.Error(), cmdErr.Message)
					} else {
						assert.EqualError(t, err, tc.expErr.Error())
					}
				}
			})
		}
	})
}

func TestUpdatePostInfo(t *testing.T) {
	mt := setupMockDB(t)

	mt.Run("Test Cases", func(mt *mtest.T) {
		repo := mongorepo.NewMongoPostRepository(mt.Client, "testDB", nil)

		testCases := []struct {
			name         string
			newPost      *models.Post
			mockBehavior func()
			expErr       error
		}{
			{
				name:    "Success",
				newPost: mockPost,
				mockBehavior: func() {
					mt.AddMockResponses(mtest.CreateSuccessResponse())
				},
				expErr: nil,
			},
			{
				name:    "Error",
				newPost: mockPost,
				mockBehavior: func() {
					mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{Message: ErrBasic.Error()}))
				},
				expErr: ErrBasic,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				tc.mockBehavior()

				err := repo.UpdatePostInfo(context.Background(), tc.newPost)
				if tc.expErr == nil {
					assert.NoError(t, err)
				} else {
					var cmdErr mongo.CommandError
					if errors.As(err, &cmdErr) {
						assert.Equal(t, tc.expErr.Error(), cmdErr.Message)
					} else {
						assert.EqualError(t, err, tc.expErr.Error())
					}
				}
			})
		}
	})
}

func TestDeletePost(t *testing.T) {
	mt := setupMockDB(t)

	mt.Run("Test Cases", func(mt *mtest.T) {
		repo := mongorepo.NewMongoPostRepository(mt.Client, "testDB", nil)

		testCases := []struct {
			name         string
			postID       string
			mockBehavior func()
			expErr       error
		}{
			{
				name:   "Success",
				postID: mockPost.ID,
				mockBehavior: func() {
					mt.AddMockResponses(mtest.CreateSuccessResponse())
				},
				expErr: nil,
			},
			{
				name:   "Error",
				postID: mockPost.ID,
				mockBehavior: func() {
					mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{Message: ErrBasic.Error()}))
				},
				expErr: ErrBasic,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				tc.mockBehavior()

				err := repo.DeletePost(context.Background(), tc.postID)
				if tc.expErr == nil {
					assert.NoError(t, err)
				} else {
					var cmdErr mongo.CommandError
					if errors.As(err, &cmdErr) {
						assert.Equal(t, tc.expErr.Error(), cmdErr.Message)
					} else {
						assert.EqualError(t, err, tc.expErr.Error())
					}
				}
			})
		}
	})
}
