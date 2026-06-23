package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ProductStatus constants
const (
	ProductStatusActive   = "active"
	ProductStatusInactive = "inactive"
	ProductStatusDraft    = "draft"
)

// Product represents a store product
type Product struct {
	ID             primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Name           string               `bson:"name" json:"name" validate:"required,min=2,max=200"`
	Slug           string               `bson:"slug" json:"slug"`
	Description    string               `bson:"description" json:"description"`
	ShortDesc      string               `bson:"short_desc" json:"short_desc"`
	CategoryID     primitive.ObjectID   `bson:"category_id" json:"category_id"`
	Category       *Category            `bson:"category,omitempty" json:"category,omitempty"`
	Price          float64              `bson:"price" json:"price" validate:"required,gt=0"`
	ComparePrice   float64              `bson:"compare_price,omitempty" json:"compare_price,omitempty"`
	CostPrice      float64              `bson:"cost_price,omitempty" json:"cost_price,omitempty"`
	Images         []string             `bson:"images" json:"images"`
	Thumbnail      string               `bson:"thumbnail" json:"thumbnail"`
	SKU            string               `bson:"sku" json:"sku"`
	Stock          int                  `bson:"stock" json:"stock"`
	LowStockAlert  int                  `bson:"low_stock_alert" json:"low_stock_alert"`
	Weight         float64              `bson:"weight,omitempty" json:"weight,omitempty"`
	Tags           []string             `bson:"tags,omitempty" json:"tags,omitempty"`
	Variants       []ProductVariant     `bson:"variants,omitempty" json:"variants,omitempty"`
	Attributes     map[string][]string  `bson:"attributes,omitempty" json:"attributes,omitempty"` // e.g. {"color": ["Black","White"], "size": ["S","M","L"]}
	Status         string               `bson:"status" json:"status"`
	IsFeatured     bool                 `bson:"is_featured" json:"is_featured"`
	IsBestSeller   bool                 `bson:"is_best_seller" json:"is_best_seller"`
	IsNewArrival   bool                 `bson:"is_new_arrival" json:"is_new_arrival"`
	MetaTitle      string               `bson:"meta_title,omitempty" json:"meta_title,omitempty"`
	MetaDesc       string               `bson:"meta_desc,omitempty" json:"meta_desc,omitempty"`
	AverageRating  float64              `bson:"average_rating" json:"average_rating"`
	ReviewCount    int                  `bson:"review_count" json:"review_count"`
	SoldCount      int                  `bson:"sold_count" json:"sold_count"`
	CreatedAt      time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time            `bson:"updated_at" json:"updated_at"`
}

// ProductVariant represents a product SKU variant
type ProductVariant struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name       string             `bson:"name" json:"name"` // e.g. "Black / M"
	SKU        string             `bson:"sku" json:"sku"`
	Price      float64            `bson:"price" json:"price"`
	Stock      int                `bson:"stock" json:"stock"`
	Attributes map[string]string  `bson:"attributes" json:"attributes"` // e.g. {"color": "Black", "size": "M"}
	Images     []string           `bson:"images,omitempty" json:"images,omitempty"`
}

// ProductFilter for list queries
type ProductFilter struct {
	CategoryID  string
	MinPrice    float64
	MaxPrice    float64
	Search      string
	Tags        []string
	Status      string
	IsFeatured  *bool
	IsBestSeller *bool
	IsNewArrival *bool
	SortBy      string // price_asc, price_desc, newest, rating, popular
	Page        int
	Limit       int
	IDs         []string
}

// ProductRepository defines data access methods for products
type ProductRepository interface {
	Create(ctx context.Context, product *Product) (*Product, error)
	FindByID(ctx context.Context, id string) (*Product, error)
	FindBySlug(ctx context.Context, slug string) (*Product, error)
	Update(ctx context.Context, id string, update map[string]interface{}) (*Product, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter ProductFilter) ([]*Product, int64, error)
	UpdateStock(ctx context.Context, id string, quantity int) error
	UpdateRating(ctx context.Context, id string, avgRating float64, count int) error
	IncrementSoldCount(ctx context.Context, id string, qty int) error
	GetFeatured(ctx context.Context, limit int) ([]*Product, error)
	GetBestSellers(ctx context.Context, limit int) ([]*Product, error)
	GetNewArrivals(ctx context.Context, limit int) ([]*Product, error)
	GetRelated(ctx context.Context, categoryID string, excludeID string, limit int) ([]*Product, error)
}

// CreateProductInput for product creation
type CreateProductInput struct {
	Name          string              `json:"name" validate:"required,min=2,max=200"`
	Description   string              `json:"description" validate:"required"`
	ShortDesc     string              `json:"short_desc"`
	CategoryID    string              `json:"category_id" validate:"required"`
	Price         float64             `json:"price" validate:"required,gt=0"`
	ComparePrice  float64             `json:"compare_price,omitempty"`
	CostPrice     float64             `json:"cost_price,omitempty"`
	Images        []string            `json:"images" validate:"required,min=1"`
	Thumbnail     string              `json:"thumbnail"`
	SKU           string              `json:"sku" validate:"required"`
	Stock         int                 `json:"stock" validate:"gte=0"`
	LowStockAlert int                 `json:"low_stock_alert"`
	Weight        float64             `json:"weight,omitempty"`
	Tags          []string            `json:"tags,omitempty"`
	Variants      []ProductVariant    `json:"variants,omitempty"`
	Attributes    map[string][]string `json:"attributes,omitempty"`
	Status        string              `json:"status"`
	IsFeatured    bool                `json:"is_featured"`
	IsBestSeller  bool                `json:"is_best_seller"`
	IsNewArrival  bool                `json:"is_new_arrival"`
	MetaTitle     string              `json:"meta_title,omitempty"`
	MetaDesc      string              `json:"meta_desc,omitempty"`
}
