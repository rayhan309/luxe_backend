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

type userRepository struct {
	col *mongo.Collection
}

// NewUserRepository creates a new MongoDB-backed user repository
func NewUserRepository(db *Database) domain.UserRepository {
	return &userRepository{col: db.Collection("users")}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.IsActive = true
	if user.Role == "" {
		user.Role = domain.RoleCustomer
	}

	_, err := r.col.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var user domain.User
	err = r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.col.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, id string, update map[string]interface{}) (*domain.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	update["updated_at"] = time.Now()

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated domain.User
	err = r.col.FindOneAndUpdate(ctx, bson.M{"_id": oid}, bson.M{"$set": update}, opts).Decode(&updated)
	return &updated, err
}

func (r *userRepository) UpdatePassword(ctx context.Context, id string, hashedPassword string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.col.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{
		"$set": bson.M{"password": hashedPassword, "updated_at": time.Now()},
	})
	return err
}

func (r *userRepository) SetResetToken(ctx context.Context, email, token string, expiry time.Time) error {
	_, err := r.col.UpdateOne(ctx, bson.M{"email": email}, bson.M{
		"$set": bson.M{
			"reset_token":  token,
			"reset_expiry": expiry,
			"updated_at":   time.Now(),
		},
	})
	return err
}

func (r *userRepository) FindByResetToken(ctx context.Context, token string) (*domain.User, error) {
	var user domain.User
	err := r.col.FindOne(ctx, bson.M{
		"reset_token":  token,
		"reset_expiry": bson.M{"$gt": time.Now()},
	}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) ClearResetToken(ctx context.Context, id string) error {
	oid, _ := primitive.ObjectIDFromHex(id)
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{
		"$unset": bson.M{"reset_token": "", "reset_expiry": ""},
		"$set":   bson.M{"updated_at": time.Now()},
	})
	return err
}

func (r *userRepository) AddToWishlist(ctx context.Context, userID, productID string) error {
	uid, _ := primitive.ObjectIDFromHex(userID)
	pid, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return err
	}
	_, err = r.col.UpdateOne(ctx, bson.M{"_id": uid}, bson.M{
		"$addToSet": bson.M{"wishlist": pid},
		"$set":      bson.M{"updated_at": time.Now()},
	})
	return err
}

func (r *userRepository) RemoveFromWishlist(ctx context.Context, userID, productID string) error {
	uid, _ := primitive.ObjectIDFromHex(userID)
	pid, _ := primitive.ObjectIDFromHex(productID)
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": uid}, bson.M{
		"$pull": bson.M{"wishlist": pid},
		"$set":  bson.M{"updated_at": time.Now()},
	})
	return err
}

func (r *userRepository) AddAddress(ctx context.Context, userID string, address domain.Address) error {
	uid, _ := primitive.ObjectIDFromHex(userID)
	address.ID = primitive.NewObjectID()
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": uid}, bson.M{
		"$push": bson.M{"addresses": address},
		"$set":  bson.M{"updated_at": time.Now()},
	})
	return err
}

func (r *userRepository) UpdateAddress(ctx context.Context, userID, addressID string, address domain.Address) error {
	uid, _ := primitive.ObjectIDFromHex(userID)
	aid, _ := primitive.ObjectIDFromHex(addressID)
	_, err := r.col.UpdateOne(ctx,
		bson.M{"_id": uid, "addresses._id": aid},
		bson.M{"$set": bson.M{
			"addresses.$":  address,
			"updated_at":   time.Now(),
		}},
	)
	return err
}

func (r *userRepository) DeleteAddress(ctx context.Context, userID, addressID string) error {
	uid, _ := primitive.ObjectIDFromHex(userID)
	aid, _ := primitive.ObjectIDFromHex(addressID)
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": uid}, bson.M{
		"$pull": bson.M{"addresses": bson.M{"_id": aid}},
		"$set":  bson.M{"updated_at": time.Now()},
	})
	return err
}

func (r *userRepository) List(ctx context.Context, page, limit int, search string) ([]*domain.User, int64, error) {
	filter := bson.M{}
	if search != "" {
		filter["$or"] = bson.A{
			bson.M{"first_name": bson.M{"$regex": search, "$options": "i"}},
			bson.M{"last_name": bson.M{"$regex": search, "$options": "i"}},
			bson.M{"email": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	total, err := r.col.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSkip(int64((page - 1) * limit)).
		SetLimit(int64(limit)).
		SetSort(bson.M{"created_at": -1})

	cursor, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var users []*domain.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	oid, _ := primitive.ObjectIDFromHex(id)
	_, err := r.col.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}
