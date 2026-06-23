package usecase

import (
	"context"
	"errors"

	"github.com/luxe/backend/internal/domain"
	"github.com/luxe/backend/internal/pkg/password"
)

type UserUseCase struct {
	userRepo domain.UserRepository
}

func NewUserUseCase(repo domain.UserRepository) *UserUseCase {
	return &UserUseCase{userRepo: repo}
}

func (uc *UserUseCase) GetProfile(ctx context.Context, id string) (*domain.User, error) {
	return uc.userRepo.FindByID(ctx, id)
}

func (uc *UserUseCase) UpdateProfile(ctx context.Context, id string, input domain.UpdateProfileInput) (*domain.User, error) {
	update := map[string]interface{}{}
	if input.FirstName != "" {
		update["first_name"] = input.FirstName
	}
	if input.LastName != "" {
		update["last_name"] = input.LastName
	}
	if input.Phone != "" {
		update["phone"] = input.Phone
	}
	if input.Avatar != "" {
		update["avatar"] = input.Avatar
	}
	return uc.userRepo.Update(ctx, id, update)
}

func (uc *UserUseCase) ChangePassword(ctx context.Context, id string, input domain.ChangePasswordInput) error {
	user, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		return errors.New("user not found")
	}

	if !password.Compare(user.Password, input.CurrentPassword) {
		return errors.New("current password is incorrect")
	}

	hashed, err := password.Hash(input.NewPassword)
	if err != nil {
		return errors.New("failed to process password")
	}

	return uc.userRepo.UpdatePassword(ctx, id, hashed)
}

func (uc *UserUseCase) GetWishlist(ctx context.Context, userID string) (*domain.User, error) {
	return uc.userRepo.FindByID(ctx, userID)
}

func (uc *UserUseCase) ToggleWishlist(ctx context.Context, userID, productID string) error {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Check if product is already in wishlist
	for _, id := range user.Wishlist {
		if id.Hex() == productID {
			return uc.userRepo.RemoveFromWishlist(ctx, userID, productID)
		}
	}

	return uc.userRepo.AddToWishlist(ctx, userID, productID)
}

func (uc *UserUseCase) AddAddress(ctx context.Context, userID string, address domain.Address) error {
	return uc.userRepo.AddAddress(ctx, userID, address)
}

func (uc *UserUseCase) UpdateAddress(ctx context.Context, userID, addressID string, address domain.Address) error {
	return uc.userRepo.UpdateAddress(ctx, userID, addressID, address)
}

func (uc *UserUseCase) DeleteAddress(ctx context.Context, userID, addressID string) error {
	return uc.userRepo.DeleteAddress(ctx, userID, addressID)
}

func (uc *UserUseCase) ListUsers(ctx context.Context, page, limit int, search string) ([]*domain.User, int64, error) {
	return uc.userRepo.List(ctx, page, limit, search)
}
