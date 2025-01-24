package service

import (
	"context"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/myacey/redditclone/internal/models"
	"github.com/myacey/redditclone/internal/repository"
	"github.com/myacey/redditclone/internal/repository/mongorepo"
	"github.com/myacey/redditclone/internal/repository/postgresrepo"
	"github.com/myacey/redditclone/internal/repository/redisrepo"
	"github.com/myacey/redditclone/internal/token"
)

type ServiceInterface interface {
	// user
	GetUserFromDBByID(ctx context.Context, userID string) (*models.User, error)
	GetUserFromDBByUsername(ctx context.Context, username string) (*models.User, error)
	CreateNewUser(ctx context.Context, user *models.User) (*models.Session, error)
	LoginUser(ctx context.Context, username string) (*models.Session, error)

	// post
	GetAllPosts(ctx context.Context) ([]*models.Post, error)
	AddPost(ctx context.Context, newPost *models.Post) error
	GetPostByID(ctx context.Context, postID string, increateVote bool) (*models.Post, error)
	GetPostsByAuthor(ctx context.Context, username string) ([]*models.Post, error)
	GetPostsByCategory(ctx context.Context, category string) ([]*models.Post, error)
	DeletePostWithID(ctx context.Context, postID string) error

	// comment
	RemoveComment(ctx context.Context, postID, commentID string) (*models.Post, error)
	AddCommentToPost(ctx context.Context, postID string, newComment models.Comment) (*models.Post, error)

	// session
	CheckUserSession(ctx context.Context, userID, token string) error

	// vote
	VotePostWithID(ctx context.Context, postID string, newVote *models.Vote) (*models.Post, error)
	UnvotePostWithID(ctx context.Context, postID, userID string) (*models.Post, error)
}

type Service struct {
	userRepo    repository.UserRepository
	postRepo    repository.PostRepository
	commentRepo repository.CommentRepository
	sessionRepo repository.SessionRepository

	tokenMaker token.TokenMaker

	logger *zap.SugaredLogger
}

func NewService(db *gorm.DB,
	mongoClient *mongo.Client,
	mongoDatabaseName string,
	redisPool *redis.Client,
	tokenMaker token.TokenMaker,
	lg *zap.SugaredLogger,
) ServiceInterface {
	userRepo := postgresrepo.NewPostgresUserRepository(db)
	commentRepo := mongorepo.NewMongoCommentRepo(mongoClient, mongoDatabaseName)
	postRepo := mongorepo.NewMongoPostRepository(mongoClient, mongoDatabaseName, commentRepo)
	sessionRepo := redisrepo.NewRedisSessionRepo(redisPool)

	return &Service{
		userRepo:    userRepo,
		postRepo:    postRepo,
		commentRepo: commentRepo,
		sessionRepo: sessionRepo,

		tokenMaker: tokenMaker,

		logger: lg,
	}
}
