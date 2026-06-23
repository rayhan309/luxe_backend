package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Review represents a product review
type Review struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ProductID  primitive.ObjectID `bson:"product_id" json:"product_id"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	User       *User              `bson:"user,omitempty" json:"user,omitempty"`
	Rating     int                `bson:"rating" json:"rating" validate:"required,min=1,max=5"`
	Title      string             `bson:"title,omitempty" json:"title,omitempty"`
	Comment    string             `bson:"comment" json:"comment" validate:"required,min=10"`
	Images     []string           `bson:"images,omitempty" json:"images,omitempty"`
	IsVerified bool               `bson:"is_verified" json:"is_verified"` // purchased product
	IsApproved bool               `bson:"is_approved" json:"is_approved"`
	Helpful    int                `bson:"helpful" json:"helpful"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}

// ReviewRepository defines data access methods for reviews
type ReviewRepository interface {
	Create(ctx context.Context, review *Review) (*Review, error)
	FindByID(ctx context.Context, id string) (*Review, error)
	FindByProductID(ctx context.Context, productID string, page, limit int, approved bool) ([]*Review, int64, error)
	FindByUserID(ctx context.Context, userID string, page, limit int) ([]*Review, int64, error)
	Approve(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, approved *bool, page, limit int) ([]*Review, int64, error)
	GetProductRatingStats(ctx context.Context, productID string) (float64, int, error)
	IncrementHelpful(ctx context.Context, id string) error
}

// CreateReviewInput for review creation
type CreateReviewInput struct {
	ProductID string   `json:"product_id" validate:"required"`
	Rating    int      `json:"rating" validate:"required,min=1,max=5"`
	Title     string   `json:"title,omitempty"`
	Comment   string   `json:"comment" validate:"required,min=10"`
	Images    []string `json:"images,omitempty"`
}
