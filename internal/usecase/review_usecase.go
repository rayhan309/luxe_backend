package usecase

import (
	"context"

	"github.com/luxe/backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReviewUseCase struct {
	reviewRepo  domain.ReviewRepository
	productRepo domain.ProductRepository
	orderRepo   domain.OrderRepository
}

func NewReviewUseCase(reviewRepo domain.ReviewRepository, productRepo domain.ProductRepository, orderRepo domain.OrderRepository) *ReviewUseCase {
	return &ReviewUseCase{reviewRepo: reviewRepo, productRepo: productRepo, orderRepo: orderRepo}
}

func (uc *ReviewUseCase) Create(ctx context.Context, userID string, input domain.CreateReviewInput) (*domain.Review, error) {
	pid, _ := primitive.ObjectIDFromHex(input.ProductID)
	uid, _ := primitive.ObjectIDFromHex(userID)

	review := &domain.Review{
		ProductID:  pid,
		UserID:     uid,
		Rating:     input.Rating,
		Title:      input.Title,
		Comment:    input.Comment,
		Images:     input.Images,
		IsApproved: false, // pending admin approval
	}

	created, err := uc.reviewRepo.Create(ctx, review)
	if err != nil {
		return nil, err
	}

	// Update product rating stats
	avg, count, err := uc.reviewRepo.GetProductRatingStats(ctx, input.ProductID)
	if err == nil {
		_ = uc.productRepo.UpdateRating(ctx, input.ProductID, avg, count)
	}

	return created, nil
}

func (uc *ReviewUseCase) GetProductReviews(ctx context.Context, productID string, page, limit int) ([]*domain.Review, int64, error) {
	return uc.reviewRepo.FindByProductID(ctx, productID, page, limit, true)
}

func (uc *ReviewUseCase) List(ctx context.Context, approved *bool, page, limit int) ([]*domain.Review, int64, error) {
	return uc.reviewRepo.List(ctx, approved, page, limit)
}

func (uc *ReviewUseCase) Approve(ctx context.Context, id string) error {
	if err := uc.reviewRepo.Approve(ctx, id); err != nil {
		return err
	}

	// Get the review to update product rating
	review, err := uc.reviewRepo.FindByID(ctx, id)
	if err != nil {
		return nil
	}

	avg, count, err := uc.reviewRepo.GetProductRatingStats(ctx, review.ProductID.Hex())
	if err == nil {
		_ = uc.productRepo.UpdateRating(ctx, review.ProductID.Hex(), avg, count)
	}

	return nil
}

func (uc *ReviewUseCase) Delete(ctx context.Context, id string) error {
	return uc.reviewRepo.Delete(ctx, id)
}
