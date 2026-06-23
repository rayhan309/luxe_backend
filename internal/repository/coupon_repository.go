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

type couponRepository struct {
	col *mongo.Collection
}

func NewCouponRepository(db *Database) domain.CouponRepository {
	return &couponRepository{col: db.Collection("coupons")}
}

func (r *couponRepository) Create(ctx context.Context, coupon *domain.Coupon) (*domain.Coupon, error) {
	coupon.ID = primitive.NewObjectID()
	coupon.CreatedAt = time.Now()
	coupon.UpdatedAt = time.Now()
	_, err := r.col.InsertOne(ctx, coupon)
	return coupon, err
}

func (r *couponRepository) FindByCode(ctx context.Context, code string) (*domain.Coupon, error) {
	var coupon domain.Coupon
	err := r.col.FindOne(ctx, bson.M{"code": code}).Decode(&coupon)
	return &coupon, err
}

func (r *couponRepository) FindByID(ctx context.Context, id string) (*domain.Coupon, error) {
	oid, _ := primitive.ObjectIDFromHex(id)
	var coupon domain.Coupon
	err := r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&coupon)
	return &coupon, err
}

func (r *couponRepository) Update(ctx context.Context, id string, update map[string]interface{}) (*domain.Coupon, error) {
	oid, _ := primitive.ObjectIDFromHex(id)
	update["updated_at"] = time.Now()
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated domain.Coupon
	err := r.col.FindOneAndUpdate(ctx, bson.M{"_id": oid}, bson.M{"$set": update}, opts).Decode(&updated)
	return &updated, err
}

func (r *couponRepository) Delete(ctx context.Context, id string) error {
	oid, _ := primitive.ObjectIDFromHex(id)
	_, err := r.col.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}

func (r *couponRepository) List(ctx context.Context, page, limit int) ([]*domain.Coupon, int64, error) {
	total, _ := r.col.CountDocuments(ctx, bson.M{})
	opts := options.Find().
		SetSkip(int64((page - 1) * limit)).
		SetLimit(int64(limit)).
		SetSort(bson.M{"created_at": -1})
	cursor, err := r.col.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	var coupons []*domain.Coupon
	return coupons, total, cursor.All(ctx, &coupons)
}

func (r *couponRepository) IncrementUsage(ctx context.Context, code string) error {
	_, err := r.col.UpdateOne(ctx, bson.M{"code": code}, bson.M{
		"$inc": bson.M{"used_count": 1},
		"$set": bson.M{"updated_at": time.Now()},
	})
	return err
}
