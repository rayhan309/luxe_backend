package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/luxe/backend/internal/domain"
	jwtpkg "github.com/luxe/backend/internal/pkg/jwt"
	"github.com/luxe/backend/internal/pkg/password"
)

// AuthUseCase handles authentication business logic
type AuthUseCase struct {
	userRepo   domain.UserRepository
	jwtManager *jwtpkg.Manager
}

func NewAuthUseCase(userRepo domain.UserRepository, jwtManager *jwtpkg.Manager) *AuthUseCase {
	return &AuthUseCase{userRepo: userRepo, jwtManager: jwtManager}
}

// Register creates a new user account
func (uc *AuthUseCase) Register(ctx context.Context, input domain.RegisterInput) (*domain.AuthResponse, error) {
	// Check if email already exists
	existing, _ := uc.userRepo.FindByEmail(ctx, input.Email)
	if existing != nil {
		return nil, errors.New("email already in use")
	}

	hashed, err := password.Hash(input.Password)
	if err != nil {
		return nil, errors.New("failed to process password")
	}

	user := &domain.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Password:  hashed,
		Phone:     input.Phone,
		Role:      domain.RoleCustomer,
	}

	created, err := uc.userRepo.Create(ctx, user)
	if err != nil {
		return nil, errors.New("failed to create account")
	}

	return uc.generateAuthResponse(created)
}

// Login authenticates a user and returns tokens
func (uc *AuthUseCase) Login(ctx context.Context, input domain.LoginInput) (*domain.AuthResponse, error) {
	user, err := uc.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !password.Compare(user.Password, input.Password) {
		return nil, errors.New("invalid email or password")
	}

	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	return uc.generateAuthResponse(user)
}

// ForgotPassword generates a reset token
func (uc *AuthUseCase) ForgotPassword(ctx context.Context, email string) (string, error) {
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		// Don't reveal if email exists
		return "If that email is registered, you will receive a reset link.", nil
	}
	_ = user

	token := uuid.NewString()
	expiry := time.Now().Add(1 * time.Hour)

	if err := uc.userRepo.SetResetToken(ctx, email, token, expiry); err != nil {
		return "", errors.New("failed to generate reset token")
	}

	// In production: send email with token
	// For now, return token directly (client should display/email it)
	return token, nil
}

// ResetPassword sets a new password using a valid reset token
func (uc *AuthUseCase) ResetPassword(ctx context.Context, token, newPassword string) error {
	user, err := uc.userRepo.FindByResetToken(ctx, token)
	if err != nil {
		return errors.New("invalid or expired reset token")
	}

	hashed, err := password.Hash(newPassword)
	if err != nil {
		return errors.New("failed to process password")
	}

	if err := uc.userRepo.UpdatePassword(ctx, user.ID.Hex(), hashed); err != nil {
		return errors.New("failed to update password")
	}

	return uc.userRepo.ClearResetToken(ctx, user.ID.Hex())
}

func (uc *AuthUseCase) generateAuthResponse(user *domain.User) (*domain.AuthResponse, error) {
	token, err := uc.jwtManager.Generate(user.ID.Hex(), user.Email, user.Role)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	refreshToken, err := uc.jwtManager.GenerateRefresh(user.ID.Hex(), user.Email, user.Role)
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	return &domain.AuthResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}
