package repository

import (
	"context"
	"time"

	"github.com/luxe/backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type reviewRepository struct {
	col *mongo.Collection
}

func NewReviewRepository(db *Database) domain.ReviewRepository {
	return &reviewRepository{col: db.Collection("reviews")}
}

func (r *reviewRepository) Create(ctx context.Context, review *domain.Review) (*domain.Review, error) {
	review.ID = primitive.NewObjectID()
	review.CreatedAt = time.Now()
	review.UpdatedAt = time.Now()
	_, err := r.col.InsertOne(ctx, review)
	return review, err
}

func (r *reviewRepository) FindByID(ctx context.Context, id string) (*domain.Review, error) {
	oid, _ := primitive.ObjectIDFromHex(id)
	var review domain.Review
	err := r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&review)
	return &review, err
}

func (r *reviewRepository) FindByProductID(ctx context.Context, productID string, page, limit int, approved bool) ([]*domain.Review, int64, error) {
	pid, _ := primitive.ObjectIDFromHex(productID)
	filter := bson.M{"product_id": pid}
	if approved {
		filter["is_approved"] = true
	}
	total, _ := r.col.CountDocuments(ctx, filter)
	opts := options.Find().
		SetSkip(int64((page - 1) * limit)).
		SetLimit(int64(limit)).
		SetSort(bson.M{"created_at": -1})
	cursor, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	var reviews []*domain.Review
	return reviews, total, cursor.All(ctx, &reviews)
}

func (r *reviewRepository) FindByUserID(ctx context.Context, userID string, page, limit int) ([]*domain.Review, int64, error) {
	uid, _ := primitive.ObjectIDFromHex(userID)
	filter := bson.M{"user_id": uid}
	total, _ := r.col.CountDocuments(ctx, filter)
	opts := options.Find().
		SetSkip(int64((page - 1) * limit)).
		SetLimit(int64(limit)).
		SetSort(bson.M{"created_at": -1})
	cursor, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	var reviews []*domain.Review
	return reviews, total, cursor.All(ctx, &reviews)
}

func (r *reviewRepository) Approve(ctx context.Context, id string) error {
	oid, _ := primitive.ObjectIDFromHex(id)
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{
		"$set": bson.M{"is_approved": true, "updated_at": time.Now()},
	})
	return err
}

func (r *reviewRepository) Delete(ctx context.Context, id string) error {
	oid, _ := primitive.ObjectIDFromHex(id)
	_, err := r.col.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}

func (r *reviewRepository) List(ctx context.Context, approved *bool, page, limit int) ([]*domain.Review, int64, error) {
	filter := bson.M{}
	if approved != nil {
		filter["is_approved"] = *approved
	}
	total, _ := r.col.CountDocuments(ctx, filter)
	opts := options.Find().
		SetSkip(int64((page - 1) * limit)).
		SetLimit(int64(limit)).
		SetSort(bson.M{"created_at": -1})
	cursor, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	var reviews []*domain.Review
	return reviews, total, cursor.All(ctx, &reviews)
}

func (r *reviewRepository) GetProductRatingStats(ctx context.Context, productID string) (float64, int, error) {
	pid, _ := primitive.ObjectIDFromHex(productID)
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"product_id": pid, "is_approved": true}}},
		{{Key: "$group", Value: bson.M{
			"_id":   nil,
			"avg":   bson.M{"$avg": "$rating"},
			"count": bson.M{"$sum": 1},
		}}},
	}
	cursor, err := r.col.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, 0, err
	}
	defer cursor.Close(ctx)
	var result []struct {
		Avg   float64 `bson:"avg"`
		Count int     `bson:"count"`
	}
	if err := cursor.All(ctx, &result); err != nil || len(result) == 0 {
		return 0, 0, err
	}
	return result[0].Avg, result[0].Count, nil
}

func (r *reviewRepository) IncrementHelpful(ctx context.Context, id string) error {
	oid, _ := primitive.ObjectIDFromHex(id)
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{
		"$inc": bson.M{"helpful": 1},
	})
	return err
}
