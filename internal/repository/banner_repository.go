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

type bannerRepository struct {
	col *mongo.Collection
}

func NewBannerRepository(db *Database) domain.BannerRepository {
	return &bannerRepository{col: db.Collection("banners")}
}

func (r *bannerRepository) Create(ctx context.Context, banner *domain.Banner) (*domain.Banner, error) {
	banner.ID = primitive.NewObjectID()
	banner.CreatedAt = time.Now()
	banner.UpdatedAt = time.Now()
	_, err := r.col.InsertOne(ctx, banner)
	return banner, err
}

func (r *bannerRepository) FindByID(ctx context.Context, id string) (*domain.Banner, error) {
	oid, _ := primitive.ObjectIDFromHex(id)
	var banner domain.Banner
	err := r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&banner)
	return &banner, err
}

func (r *bannerRepository) Update(ctx context.Context, id string, update map[string]interface{}) (*domain.Banner, error) {
	oid, _ := primitive.ObjectIDFromHex(id)
	update["updated_at"] = time.Now()
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated domain.Banner
	err := r.col.FindOneAndUpdate(ctx, bson.M{"_id": oid}, bson.M{"$set": update}, opts).Decode(&updated)
	return &updated, err
}

func (r *bannerRepository) Delete(ctx context.Context, id string) error {
	oid, _ := primitive.ObjectIDFromHex(id)
	_, err := r.col.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}

func (r *bannerRepository) List(ctx context.Context, activeOnly bool, position string) ([]*domain.Banner, error) {
	filter := bson.M{}
	if activeOnly {
		filter["is_active"] = true
		filter["$or"] = bson.A{
			bson.M{"expires_at": bson.M{"$gt": time.Now()}},
			bson.M{"expires_at": nil},
		}
	}
	if position != "" {
		filter["position"] = position
	}
	opts := options.Find().SetSort(bson.M{"sort_order": 1})
	cursor, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var banners []*domain.Banner
	return banners, cursor.All(ctx, &banners)
}
