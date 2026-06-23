package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Cart represents a user's shopping cart
type Cart struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Items     []CartItem         `bson:"items" json:"items"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// CartItem represents a single item in the cart
type CartItem struct {
	ID         primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	ProductID  primitive.ObjectID  `bson:"product_id" json:"product_id"`
	VariantID  *primitive.ObjectID `bson:"variant_id,omitempty" json:"variant_id,omitempty"`
	Product    *Product            `bson:"product,omitempty" json:"product,omitempty"`
	Name       string              `bson:"name" json:"name"`
	Image      string              `bson:"image" json:"image"`
	SKU        string              `bson:"sku" json:"sku"`
	Attributes map[string]string   `bson:"attributes,omitempty" json:"attributes,omitempty"`
	Price      float64             `bson:"price" json:"price"`
	Quantity   int                 `bson:"quantity" json:"quantity"`
}

// CartRepository defines data access methods for carts
type CartRepository interface {
	GetByUserID(ctx context.Context, userID string) (*Cart, error)
	AddItem(ctx context.Context, userID string, item CartItem) error
	UpdateItemQuantity(ctx context.Context, userID string, itemID string, quantity int) error
	RemoveItem(ctx context.Context, userID string, itemID string) error
	Clear(ctx context.Context, userID string) error
	Create(ctx context.Context, cart *Cart) (*Cart, error)
}

// AddToCartInput for adding item to cart
type AddToCartInput struct {
	ProductID  string            `json:"product_id" validate:"required"`
	VariantID  string            `json:"variant_id,omitempty"`
	Quantity   int               `json:"quantity" validate:"required,min=1"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

// UpdateCartItemInput for updating cart item quantity
type UpdateCartItemInput struct {
	Quantity int `json:"quantity" validate:"required,min=1"`
}
