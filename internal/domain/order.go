package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Order status constants
const (
	OrderStatusPending    = "pending"
	OrderStatusProcessing = "processing"
	OrderStatusShipped    = "shipped"
	OrderStatusDelivered  = "delivered"
	OrderStatusCancelled  = "cancelled"
	OrderStatusRefunded   = "refunded"
)

// Payment method constants
const (
	PaymentCOD    = "cod"
	PaymentOnline = "online"
)

// Payment status constants
const (
	PaymentPending  = "pending"
	PaymentPaid     = "paid"
	PaymentFailed   = "failed"
	PaymentRefunded = "refunded"
)

// Order represents a customer order
type Order struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OrderNumber     string             `bson:"order_number" json:"order_number"`
	UserID          primitive.ObjectID `bson:"user_id" json:"user_id"`
	User            *User              `bson:"user,omitempty" json:"user,omitempty"`
	Items           []OrderItem        `bson:"items" json:"items"`
	ShippingAddress Address            `bson:"shipping_address" json:"shipping_address"`
	Status          string             `bson:"status" json:"status"`
	PaymentMethod   string             `bson:"payment_method" json:"payment_method"`
	PaymentStatus   string             `bson:"payment_status" json:"payment_status"`
	Subtotal        float64            `bson:"subtotal" json:"subtotal"`
	ShippingCost    float64            `bson:"shipping_cost" json:"shipping_cost"`
	Discount        float64            `bson:"discount" json:"discount"`
	CouponCode      string             `bson:"coupon_code,omitempty" json:"coupon_code,omitempty"`
	Tax             float64            `bson:"tax" json:"tax"`
	Total           float64            `bson:"total" json:"total"`
	Note            string             `bson:"note,omitempty" json:"note,omitempty"`
	TrackingNumber  string             `bson:"tracking_number,omitempty" json:"tracking_number,omitempty"`
	StatusHistory   []StatusHistory    `bson:"status_history" json:"status_history"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
}

// OrderItem represents a line item in an order
type OrderItem struct {
	ProductID  primitive.ObjectID `bson:"product_id" json:"product_id"`
	VariantID  *primitive.ObjectID `bson:"variant_id,omitempty" json:"variant_id,omitempty"`
	Name       string             `bson:"name" json:"name"`
	Image      string             `bson:"image" json:"image"`
	SKU        string             `bson:"sku" json:"sku"`
	Attributes map[string]string  `bson:"attributes,omitempty" json:"attributes,omitempty"`
	Price      float64            `bson:"price" json:"price"`
	Quantity   int                `bson:"quantity" json:"quantity"`
	Total      float64            `bson:"total" json:"total"`
}

// StatusHistory tracks order status changes
type StatusHistory struct {
	Status    string    `bson:"status" json:"status"`
	Note      string    `bson:"note,omitempty" json:"note,omitempty"`
	ChangedAt time.Time `bson:"changed_at" json:"changed_at"`
}

// OrderRepository defines data access methods for orders
type OrderRepository interface {
	Create(ctx context.Context, order *Order) (*Order, error)
	FindByID(ctx context.Context, id string) (*Order, error)
	FindByOrderNumber(ctx context.Context, orderNumber string) (*Order, error)
	UpdateStatus(ctx context.Context, id string, status string, note string) error
	UpdatePaymentStatus(ctx context.Context, id string, status string) error
	SetTrackingNumber(ctx context.Context, id string, tracking string) error
	ListByUser(ctx context.Context, userID string, page, limit int) ([]*Order, int64, error)
	List(ctx context.Context, status string, page, limit int, search string) ([]*Order, int64, error)
	GetStats(ctx context.Context) (*OrderStats, error)
	GetRevenueByPeriod(ctx context.Context, days int) ([]RevenuePoint, error)
}

// OrderStats for dashboard analytics
type OrderStats struct {
	TotalOrders    int64   `json:"total_orders"`
	PendingOrders  int64   `json:"pending_orders"`
	TotalRevenue   float64 `json:"total_revenue"`
	TodayRevenue   float64 `json:"today_revenue"`
	TodayOrders    int64   `json:"today_orders"`
	AverageOrder   float64 `json:"average_order"`
}

// RevenuePoint for chart data
type RevenuePoint struct {
	Date    string  `json:"date"`
	Revenue float64 `json:"revenue"`
	Orders  int64   `json:"orders"`
}

// CreateOrderInput for order creation
type CreateOrderInput struct {
	Items           []CreateOrderItemInput `json:"items" validate:"required,min=1"`
	ShippingAddress Address               `json:"shipping_address" validate:"required"`
	PaymentMethod   string                `json:"payment_method" validate:"required,oneof=cod online"`
	CouponCode      string                `json:"coupon_code,omitempty"`
	Note            string                `json:"note,omitempty"`
}

// CreateOrderItemInput for order item creation
type CreateOrderItemInput struct {
	ProductID  string            `json:"product_id" validate:"required"`
	VariantID  string            `json:"variant_id,omitempty"`
	Quantity   int               `json:"quantity" validate:"required,min=1"`
	Attributes map[string]string `json:"attributes,omitempty"`
}
