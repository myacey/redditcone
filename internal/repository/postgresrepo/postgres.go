package postgresrepo

import (
	"fmt"
	"os"

	"github.com/myacey/redditclone/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConfigurePostgres() (*gorm.DB, error) {
	dsn := os.Getenv("POSTGRES_DSN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("cant connect to postgres: %v", err)
	}
	if err = db.AutoMigrate(&models.User{}); err != nil {
		return nil, err
	}

	return db, nil
}
