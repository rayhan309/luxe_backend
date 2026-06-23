package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Category represents a product category
type Category struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name" validate:"required,min=2,max=100"`
	Slug        string             `bson:"slug" json:"slug"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	Image       string             `bson:"image,omitempty" json:"image,omitempty"`
	ParentID    *primitive.ObjectID `bson:"parent_id,omitempty" json:"parent_id,omitempty"`
	Children    []*Category        `bson:"children,omitempty" json:"children,omitempty"`
	IsActive    bool               `bson:"is_active" json:"is_active"`
	SortOrder   int                `bson:"sort_order" json:"sort_order"`
	MetaTitle   string             `bson:"meta_title,omitempty" json:"meta_title,omitempty"`
	MetaDesc    string             `bson:"meta_desc,omitempty" json:"meta_desc,omitempty"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// CategoryRepository defines data access methods for categories
type CategoryRepository interface {
	Create(ctx context.Context, category *Category) (*Category, error)
	FindByID(ctx context.Context, id string) (*Category, error)
	FindBySlug(ctx context.Context, slug string) (*Category, error)
	Update(ctx context.Context, id string, update map[string]interface{}) (*Category, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, activeOnly bool) ([]*Category, error)
	GetTree(ctx context.Context) ([]*Category, error)
}

// CreateCategoryInput for category creation
type CreateCategoryInput struct {
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Description string `json:"description,omitempty"`
	Image       string `json:"image,omitempty"`
	ParentID    string `json:"parent_id,omitempty"`
	IsActive    bool   `json:"is_active"`
	SortOrder   int    `json:"sort_order"`
	MetaTitle   string `json:"meta_title,omitempty"`
	MetaDesc    string `json:"meta_desc,omitempty"`
}
