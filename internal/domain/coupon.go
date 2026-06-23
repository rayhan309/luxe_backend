package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Coupon discount type constants
const (
	DiscountTypeFixed      = "fixed"
	DiscountTypePercentage = "percentage"
)

// Coupon represents a discount coupon
type Coupon struct {
	ID              primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Code            string               `bson:"code" json:"code" validate:"required,uppercase"`
	Description     string               `bson:"description,omitempty" json:"description,omitempty"`
	DiscountType    string               `bson:"discount_type" json:"discount_type"` // fixed | percentage
	DiscountValue   float64              `bson:"discount_value" json:"discount_value" validate:"required,gt=0"`
	MinOrderAmount  float64              `bson:"min_order_amount" json:"min_order_amount"`
	MaxDiscount     float64              `bson:"max_discount,omitempty" json:"max_discount,omitempty"` // cap for percentage
	UsageLimit      int                  `bson:"usage_limit" json:"usage_limit"`           // 0 = unlimited
	UsedCount       int                  `bson:"used_count" json:"used_count"`
	UserUsageLimit  int                  `bson:"user_usage_limit" json:"user_usage_limit"` // per user
	AllowedUsers    []primitive.ObjectID `bson:"allowed_users,omitempty" json:"allowed_users,omitempty"`
	AllowedProducts []primitive.ObjectID `bson:"allowed_products,omitempty" json:"allowed_products,omitempty"`
	IsActive        bool                 `bson:"is_active" json:"is_active"`
	StartsAt        time.Time            `bson:"starts_at" json:"starts_at"`
	ExpiresAt       time.Time            `bson:"expires_at" json:"expires_at"`
	CreatedAt       time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time            `bson:"updated_at" json:"updated_at"`
}

// CouponRepository defines data access methods for coupons
type CouponRepository interface {
	Create(ctx context.Context, coupon *Coupon) (*Coupon, error)
	FindByCode(ctx context.Context, code string) (*Coupon, error)
	FindByID(ctx context.Context, id string) (*Coupon, error)
	Update(ctx context.Context, id string, update map[string]interface{}) (*Coupon, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, limit int) ([]*Coupon, int64, error)
	IncrementUsage(ctx context.Context, code string) error
}

// ValidateCouponInput for coupon validation request
type ValidateCouponInput struct {
	Code        string  `json:"code" validate:"required"`
	OrderAmount float64 `json:"order_amount" validate:"required,gt=0"`
}

// CouponValidationResult returned after coupon validation
type CouponValidationResult struct {
	Valid         bool    `json:"valid"`
	Coupon        *Coupon `json:"coupon,omitempty"`
	DiscountAmount float64 `json:"discount_amount"`
	Message       string  `json:"message"`
}
