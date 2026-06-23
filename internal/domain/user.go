package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Role constants
const (
	RoleCustomer = "customer"
	RoleAdmin    = "admin"
)

// User represents a registered user
type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FirstName    string             `bson:"first_name" json:"first_name" validate:"required,min=2,max=50"`
	LastName     string             `bson:"last_name" json:"last_name" validate:"required,min=2,max=50"`
	Email        string             `bson:"email" json:"email" validate:"required,email"`
	Password     string             `bson:"password" json:"-"`
	Phone        string             `bson:"phone,omitempty" json:"phone,omitempty"`
	Role         string             `bson:"role" json:"role"`
	IsActive     bool               `bson:"is_active" json:"is_active"`
	Avatar       string             `bson:"avatar,omitempty" json:"avatar,omitempty"`
	Addresses    []Address          `bson:"addresses,omitempty" json:"addresses,omitempty"`
	Wishlist     []primitive.ObjectID `bson:"wishlist,omitempty" json:"wishlist,omitempty"`
	ResetToken   string             `bson:"reset_token,omitempty" json:"-"`
	ResetExpiry  time.Time          `bson:"reset_expiry,omitempty" json:"-"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

// Address represents a shipping/billing address
type Address struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Label      string             `bson:"label" json:"label"` // Home, Office, etc.
	FullName   string             `bson:"full_name" json:"full_name"`
	Phone      string             `bson:"phone" json:"phone"`
	Street     string             `bson:"street" json:"street"`
	City       string             `bson:"city" json:"city"`
	State      string             `bson:"state" json:"state"`
	ZipCode    string             `bson:"zip_code" json:"zip_code"`
	Country    string             `bson:"country" json:"country"`
	IsDefault  bool               `bson:"is_default" json:"is_default"`
}

// UserRepository defines data access methods for users
type UserRepository interface {
	Create(ctx context.Context, user *User) (*User, error)
	FindByID(ctx context.Context, id string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, id string, update map[string]interface{}) (*User, error)
	UpdatePassword(ctx context.Context, id string, hashedPassword string) error
	SetResetToken(ctx context.Context, email string, token string, expiry time.Time) error
	FindByResetToken(ctx context.Context, token string) (*User, error)
	ClearResetToken(ctx context.Context, id string) error
	AddToWishlist(ctx context.Context, userID string, productID string) error
	RemoveFromWishlist(ctx context.Context, userID string, productID string) error
	AddAddress(ctx context.Context, userID string, address Address) error
	UpdateAddress(ctx context.Context, userID string, addressID string, address Address) error
	DeleteAddress(ctx context.Context, userID string, addressID string) error
	List(ctx context.Context, page, limit int, search string) ([]*User, int64, error)
	Delete(ctx context.Context, id string) error
}

// RegisterInput represents registration request payload
type RegisterInput struct {
	FirstName string `json:"first_name" validate:"required,min=2,max=50"`
	LastName  string `json:"last_name" validate:"required,min=2,max=50"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	Phone     string `json:"phone,omitempty"`
}

// LoginInput represents login request payload
type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UpdateProfileInput for profile updates
type UpdateProfileInput struct {
	FirstName string `json:"first_name" validate:"omitempty,min=2,max=50"`
	LastName  string `json:"last_name" validate:"omitempty,min=2,max=50"`
	Phone     string `json:"phone,omitempty"`
	Avatar    string `json:"avatar,omitempty"`
}

// ChangePasswordInput for password changes
type ChangePasswordInput struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}

// AuthResponse returned after login/register
type AuthResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	User         *User  `json:"user"`
}
