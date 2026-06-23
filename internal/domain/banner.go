package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Banner represents a promotional banner
type Banner struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title" validate:"required"`
	Subtitle    string             `bson:"subtitle,omitempty" json:"subtitle,omitempty"`
	Image       string             `bson:"image" json:"image" validate:"required"`
	MobileImage string             `bson:"mobile_image,omitempty" json:"mobile_image,omitempty"`
	Link        string             `bson:"link,omitempty" json:"link,omitempty"`
	ButtonText  string             `bson:"button_text,omitempty" json:"button_text,omitempty"`
	Position    string             `bson:"position" json:"position"` // hero, mid, bottom
	SortOrder   int                `bson:"sort_order" json:"sort_order"`
	IsActive    bool               `bson:"is_active" json:"is_active"`
	StartsAt    *time.Time         `bson:"starts_at,omitempty" json:"starts_at,omitempty"`
	ExpiresAt   *time.Time         `bson:"expires_at,omitempty" json:"expires_at,omitempty"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// BannerRepository defines data access methods for banners
type BannerRepository interface {
	Create(ctx context.Context, banner *Banner) (*Banner, error)
	FindByID(ctx context.Context, id string) (*Banner, error)
	Update(ctx context.Context, id string, update map[string]interface{}) (*Banner, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, activeOnly bool, position string) ([]*Banner, error)
}

// CreateBannerInput for banner creation
type CreateBannerInput struct {
	Title       string     `json:"title" validate:"required"`
	Subtitle    string     `json:"subtitle,omitempty"`
	Image       string     `json:"image" validate:"required"`
	MobileImage string     `json:"mobile_image,omitempty"`
	Link        string     `json:"link,omitempty"`
	ButtonText  string     `json:"button_text,omitempty"`
	Position    string     `json:"position" validate:"required,oneof=hero mid bottom"`
	SortOrder   int        `json:"sort_order"`
	IsActive    bool       `json:"is_active"`
	StartsAt    *time.Time `json:"starts_at,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}
