package usecase

import (
	"context"

	"github.com/luxe/backend/internal/domain"
)

type BannerUseCase struct {
	bannerRepo domain.BannerRepository
}

func NewBannerUseCase(repo domain.BannerRepository) *BannerUseCase {
	return &BannerUseCase{bannerRepo: repo}
}

func (uc *BannerUseCase) Create(ctx context.Context, input domain.CreateBannerInput) (*domain.Banner, error) {
	banner := &domain.Banner{
		Title:       input.Title,
		Subtitle:    input.Subtitle,
		Image:       input.Image,
		MobileImage: input.MobileImage,
		Link:        input.Link,
		ButtonText:  input.ButtonText,
		Position:    input.Position,
		SortOrder:   input.SortOrder,
		IsActive:    input.IsActive,
		StartsAt:    input.StartsAt,
		ExpiresAt:   input.ExpiresAt,
	}
	return uc.bannerRepo.Create(ctx, banner)
}

func (uc *BannerUseCase) Update(ctx context.Context, id string, update map[string]interface{}) (*domain.Banner, error) {
	return uc.bannerRepo.Update(ctx, id, update)
}

func (uc *BannerUseCase) Delete(ctx context.Context, id string) error {
	return uc.bannerRepo.Delete(ctx, id)
}

func (uc *BannerUseCase) List(ctx context.Context, activeOnly bool, position string) ([]*domain.Banner, error) {
	return uc.bannerRepo.List(ctx, activeOnly, position)
}
