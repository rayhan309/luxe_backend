package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/luxe/backend/internal/domain"
)

type CouponUseCase struct {
	couponRepo domain.CouponRepository
}

func NewCouponUseCase(repo domain.CouponRepository) *CouponUseCase {
	return &CouponUseCase{couponRepo: repo}
}

func (uc *CouponUseCase) Create(ctx context.Context, coupon *domain.Coupon) (*domain.Coupon, error) {
	return uc.couponRepo.Create(ctx, coupon)
}

func (uc *CouponUseCase) Update(ctx context.Context, id string, update map[string]interface{}) (*domain.Coupon, error) {
	return uc.couponRepo.Update(ctx, id, update)
}

func (uc *CouponUseCase) Delete(ctx context.Context, id string) error {
	return uc.couponRepo.Delete(ctx, id)
}

func (uc *CouponUseCase) List(ctx context.Context, page, limit int) ([]*domain.Coupon, int64, error) {
	return uc.couponRepo.List(ctx, page, limit)
}

func (uc *CouponUseCase) Validate(ctx context.Context, input domain.ValidateCouponInput) (*domain.CouponValidationResult, error) {
	coupon, err := uc.couponRepo.FindByCode(ctx, input.Code)
	if err != nil {
		return &domain.CouponValidationResult{Valid: false, Message: "Coupon not found"}, nil
	}

	now := time.Now()
	if !coupon.IsActive {
		return &domain.CouponValidationResult{Valid: false, Message: "Coupon is inactive"}, nil
	}
	if now.Before(coupon.StartsAt) {
		return &domain.CouponValidationResult{Valid: false, Message: "Coupon has not started yet"}, nil
	}
	if now.After(coupon.ExpiresAt) {
		return &domain.CouponValidationResult{Valid: false, Message: "Coupon has expired"}, nil
	}
	if coupon.UsageLimit > 0 && coupon.UsedCount >= coupon.UsageLimit {
		return &domain.CouponValidationResult{Valid: false, Message: "Coupon usage limit reached"}, nil
	}
	if input.OrderAmount < coupon.MinOrderAmount {
		return &domain.CouponValidationResult{
			Valid:   false,
			Message: errors.New("minimum order amount not met").Error(),
		}, nil
	}

	var discountAmount float64
	if coupon.DiscountType == domain.DiscountTypeFixed {
		discountAmount = coupon.DiscountValue
	} else {
		discountAmount = input.OrderAmount * (coupon.DiscountValue / 100)
		if coupon.MaxDiscount > 0 && discountAmount > coupon.MaxDiscount {
			discountAmount = coupon.MaxDiscount
		}
	}

	return &domain.CouponValidationResult{
		Valid:          true,
		Coupon:         coupon,
		DiscountAmount: discountAmount,
		Message:        "Coupon applied successfully",
	}, nil
}
