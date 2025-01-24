package mongorepo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/myacey/redditclone/internal/models"
	"github.com/myacey/redditclone/internal/repository"
)

type MongoCommentRepo struct {
	commentCollection *mongo.Collection
}

func NewMongoCommentRepo(client *mongo.Client, dbName string) repository.CommentRepository {
	return &MongoCommentRepo{
		commentCollection: client.Database(dbName).Collection("comments"),
	}
}

func (r *MongoCommentRepo) GetCommentsByPostID(ctx context.Context, postID string) ([]*models.Comment, error) {
	filter := bson.M{"post_id": postID}

	comments := []*models.Comment{}
	cursor, err := r.commentCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	err = cursor.All(ctx, &comments)

	return comments, err
}

func (r *MongoCommentRepo) GetCommentByID(ctx context.Context, commentID string) (*models.Comment, error) {
	filter := bson.M{"_id": commentID}

	var post *models.Comment
	err := r.commentCollection.FindOne(ctx, filter).Decode(&post)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return post, nil
}

func (r *MongoCommentRepo) CreateComment(ctx context.Context, newComment *models.Comment) error {
	_, err := r.commentCollection.InsertOne(ctx, newComment)
	return err
}

func (r *MongoCommentRepo) DeleteComment(ctx context.Context, commentID string) error {
	filter := bson.M{"_id": commentID}

	_, err := r.commentCollection.DeleteOne(ctx, filter)
	return err
}
