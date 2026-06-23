package repository

import (
	"context"
	"time"

	"github.com/luxe/backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type cartRepository struct {
	col *mongo.Collection
}

func NewCartRepository(db *Database) domain.CartRepository {
	return &cartRepository{col: db.Collection("carts")}
}

func (r *cartRepository) Create(ctx context.Context, cart *domain.Cart) (*domain.Cart, error) {
	cart.ID = primitive.NewObjectID()
	cart.UpdatedAt = time.Now()
	_, err := r.col.InsertOne(ctx, cart)
	return cart, err
}

func (r *cartRepository) GetByUserID(ctx context.Context, userID string) (*domain.Cart, error) {
	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	var cart domain.Cart
	err = r.col.FindOne(ctx, bson.M{"user_id": uid}).Decode(&cart)
	if err != nil {
		// Return empty cart on not found
		return &domain.Cart{UserID: uid, Items: []domain.CartItem{}}, nil
	}
	return &cart, nil
}

func (r *cartRepository) AddItem(ctx context.Context, userID string, item domain.CartItem) error {
	uid, _ := primitive.ObjectIDFromHex(userID)
	item.ID = primitive.NewObjectID()

	res, err := r.col.UpdateOne(ctx,
		bson.M{"user_id": uid},
		bson.M{
			"$push": bson.M{"items": item},
			"$set":  bson.M{"updated_at": time.Now()},
		},
	)
	if err != nil || res.MatchedCount == 0 {
		// Create cart
		cart := &domain.Cart{
			ID:        primitive.NewObjectID(),
			UserID:    uid,
			Items:     []domain.CartItem{item},
			UpdatedAt: time.Now(),
		}
		_, err = r.col.InsertOne(ctx, cart)
	}
	return err
}

func (r *cartRepository) UpdateItemQuantity(ctx context.Context, userID, itemID string, quantity int) error {
	uid, _ := primitive.ObjectIDFromHex(userID)
	iid, _ := primitive.ObjectIDFromHex(itemID)
	_, err := r.col.UpdateOne(ctx,
		bson.M{"user_id": uid, "items._id": iid},
		bson.M{
			"$set": bson.M{
				"items.$.quantity": quantity,
				"updated_at":       time.Now(),
			},
		},
	)
	return err
}

func (r *cartRepository) RemoveItem(ctx context.Context, userID, itemID string) error {
	uid, _ := primitive.ObjectIDFromHex(userID)
	iid, _ := primitive.ObjectIDFromHex(itemID)
	_, err := r.col.UpdateOne(ctx, bson.M{"user_id": uid}, bson.M{
		"$pull": bson.M{"items": bson.M{"_id": iid}},
		"$set":  bson.M{"updated_at": time.Now()},
	})
	return err
}

func (r *cartRepository) Clear(ctx context.Context, userID string) error {
	uid, _ := primitive.ObjectIDFromHex(userID)
	_, err := r.col.UpdateOne(ctx, bson.M{"user_id": uid}, bson.M{
		"$set": bson.M{"items": []domain.CartItem{}, "updated_at": time.Now()},
	})
	return err
}
