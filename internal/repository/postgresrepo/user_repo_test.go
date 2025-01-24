package postgresrepo_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/myacey/redditclone/internal/models"
	"github.com/myacey/redditclone/internal/repository/postgresrepo"
)

var ErrBasic = errors.New("some error")

var mockUser = models.NewUser("testuser", "qwerty123")

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:                 mockDB,
		PreferSimpleProtocol: true,
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	return db, mock
}

func TestCreateUser(t *testing.T) {
	db, mock := setupMockDB(t)

	repo := postgresrepo.NewPostgresUserRepository(db)
	ctx := context.TODO()

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "users"`).
		WithArgs(mockUser.ID, mockUser.Username, mockUser.Password).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.CreateUser(ctx, mockUser)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByID(t *testing.T) {
	db, mock := setupMockDB(t)

	repo := postgresrepo.NewPostgresUserRepository(db)
	ctx := context.TODO()

	testCases := []struct {
		name         string
		userID       string
		mockBehavior func(userID string)
		expUser      *models.User
		expErr       error
	}{
		{
			name:   "Success",
			userID: mockUser.ID,
			mockBehavior: func(userID string) {
				rows := sqlmock.NewRows([]string{"id", "username", "password"}).
					AddRow(mockUser.ID, mockUser.Username, mockUser.Password)
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE id = \$1 ORDER BY "users"\."id" LIMIT \$2`).
					WithArgs(userID, 1).
					WillReturnRows(rows)
			},
			expUser: mockUser,
			expErr:  nil,
		},
		{
			name:   "Error sql",
			userID: mockUser.ID,
			mockBehavior: func(userID string) {
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE id = \$1 ORDER BY "users"\."id" LIMIT \$2`).
					WithArgs(userID, 1).
					WillReturnError(ErrBasic)
			},
			expUser: nil,
			expErr:  ErrBasic,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.userID)

			user, err := repo.GetUserByID(ctx, tc.userID)
			assert.Equal(t, tc.expUser, user)
			assert.Equal(t, tc.expErr, err)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetUserByUsername(t *testing.T) {
	db, mock := setupMockDB(t)

	repo := postgresrepo.NewPostgresUserRepository(db)
	ctx := context.TODO()

	testCases := []struct {
		name         string
		username     string
		mockBehavior func(username string)
		expUser      *models.User
		expErr       error
	}{
		{
			name:     "Success",
			username: mockUser.Username,
			mockBehavior: func(username string) {
				rows := sqlmock.NewRows([]string{"id", "username", "password"}).
					AddRow(mockUser.ID, mockUser.Username, mockUser.Password)
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE username = \$1 ORDER BY "users"\."id" LIMIT \$2`).
					WithArgs(username, 1).
					WillReturnRows(rows)
			},
			expUser: mockUser,
			expErr:  nil,
		},
		{
			name:     "Err sql",
			username: mockUser.Username,
			mockBehavior: func(username string) {
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE username = \$1 ORDER BY "users"\."id" LIMIT \$2`).
					WithArgs(username, 1).
					WillReturnError(ErrBasic)
			},
			expUser: nil,
			expErr:  ErrBasic,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.username)

			user, err := repo.GetUserByUsername(ctx, tc.username)
			assert.Equal(t, tc.expUser, user)
			assert.Equal(t, tc.expErr, err)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
