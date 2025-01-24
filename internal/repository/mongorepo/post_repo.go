package mongorepo

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/myacey/redditclone/internal/models"
	"github.com/myacey/redditclone/internal/repository"
)

type MongoPostRepository struct {
	postsCollection *mongo.Collection
	commentRepo     repository.CommentRepository
}

func NewMongoPostRepository(client *mongo.Client, dbName string, commentRepo repository.CommentRepository) repository.PostRepository {
	return &MongoPostRepository{
		postsCollection: client.Database(dbName).Collection("posts"),
		commentRepo:     commentRepo,
	}
}

func (r *MongoPostRepository) CreatePost(ctx context.Context, newPost *models.Post) error {
	_, err := r.postsCollection.InsertOne(ctx, newPost)
	return err
}

func (r *MongoPostRepository) GetPostByID(ctx context.Context, postID string) (*models.Post, error) {
	filter := bson.M{"_id": postID}
	post := models.Post{}
	err := r.postsCollection.FindOne(ctx, filter).Decode(&post)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &post, err
}

func (r *MongoPostRepository) GetAllPosts(ctx context.Context) ([]*models.Post, error) {
	res, err := r.postsCollection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	posts := []*models.Post{}
	err = res.All(ctx, &posts)

	return posts, err
}

func (r *MongoPostRepository) UpdatePostInfo(ctx context.Context, newPost *models.Post) error {
	filter := bson.M{"_id": newPost.ID}

	opts := options.Replace().SetUpsert(false)

	_, err := r.postsCollection.ReplaceOne(ctx, filter, newPost, opts)
	return err
}

func (r *MongoPostRepository) DeletePost(ctx context.Context, postID string) error {
	filter := bson.M{"_id": postID}

	_, err := r.postsCollection.DeleteOne(ctx, filter)
	return err
}

// vote logic (maybe to different file?....)
func (r *MongoPostRepository) VotePost(ctx context.Context, postID string, newVote *models.Vote) (*models.Post, error) {
	filter := bson.M{"_id": postID, "votes.user": newVote.UserID}
	update := bson.M{
		"$set": bson.M{
			"votes.$.vote": newVote.Vote,
		},
	}

	after := options.After
	opts := options.FindOneAndUpdate().SetReturnDocument(after)

	res := r.postsCollection.FindOneAndUpdate(ctx, filter, update, opts)
	if errors.Is(res.Err(), mongo.ErrNoDocuments) {
		// create new vote

		update = bson.M{
			"$push": bson.M{
				"votes": newVote,
			},
		}
		_, err := r.postsCollection.UpdateOne(ctx, bson.M{"_id": postID}, update)
		if err != nil {
			return nil, err
		}

		return r.GetPostByID(ctx, postID)
	} else if res.Err() != nil {
		return nil, res.Err()
	}

	// was no error
	var updatedPost models.Post
	if err := res.Decode(&updatedPost); err != nil {
		return nil, err
	}

	return &updatedPost, nil
}

func (r *MongoPostRepository) UnvotePost(ctx context.Context, postID, userID string) (*models.Post, error) {
	filter := bson.M{"_id": postID}
	update := bson.M{
		"$pull": bson.M{
			"votes": bson.M{"user": userID},
		},
	}

	after := options.After
	opts := options.FindOneAndUpdate().SetReturnDocument(after)

	res := r.postsCollection.FindOneAndUpdate(ctx, filter, update, opts)
	if errors.Is(res.Err(), mongo.ErrNoDocuments) {
		// error request, wtf he did?
		return r.GetPostByID(ctx, postID)
	} else if res.Err() != nil {
		return nil, res.Err()
	}

	// was no error
	var updatedPost models.Post
	if err := res.Decode(&updatedPost); err != nil {
		return nil, err
	}

	return &updatedPost, nil
}
