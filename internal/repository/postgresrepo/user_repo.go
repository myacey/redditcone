package postgresrepo

import (
	"context"

	"gorm.io/gorm"

	"github.com/myacey/redditclone/internal/models"
	"github.com/myacey/redditclone/internal/repository"
)

// PostgresUserRepository should implement UserRepository interface.
// Use PostgreSQL and GORM
type PostgresUserRepository struct {
	db *gorm.DB
}

func NewPostgresUserRepository(db *gorm.DB) repository.UserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) CreateUser(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *PostgresUserRepository) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	var usr models.User
	err := r.db.WithContext(ctx).Where("id = ?", userID).First(&usr).Error
	if err != nil {
		return nil, err
	}

	return &usr, nil
}

func (r *PostgresUserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var usr models.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&usr).Error
	if err != nil {
		return nil, err
	}

	return &usr, nil
}
