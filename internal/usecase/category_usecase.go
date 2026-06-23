package usecase

import (
	"context"

	"github.com/luxe/backend/internal/domain"
)

type CategoryUseCase struct {
	categoryRepo domain.CategoryRepository
}

func NewCategoryUseCase(repo domain.CategoryRepository) *CategoryUseCase {
	return &CategoryUseCase{categoryRepo: repo}
}

func (uc *CategoryUseCase) Create(ctx context.Context, input domain.CreateCategoryInput) (*domain.Category, error) {
	cat := &domain.Category{
		Name:        input.Name,
		Description: input.Description,
		Image:       input.Image,
		IsActive:    input.IsActive,
		SortOrder:   input.SortOrder,
		MetaTitle:   input.MetaTitle,
		MetaDesc:    input.MetaDesc,
	}
	return uc.categoryRepo.Create(ctx, cat)
}

func (uc *CategoryUseCase) GetBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	return uc.categoryRepo.FindBySlug(ctx, slug)
}

func (uc *CategoryUseCase) Update(ctx context.Context, id string, update map[string]interface{}) (*domain.Category, error) {
	return uc.categoryRepo.Update(ctx, id, update)
}

func (uc *CategoryUseCase) Delete(ctx context.Context, id string) error {
	return uc.categoryRepo.Delete(ctx, id)
}

func (uc *CategoryUseCase) List(ctx context.Context, activeOnly bool) ([]*domain.Category, error) {
	return uc.categoryRepo.List(ctx, activeOnly)
}

func (uc *CategoryUseCase) GetTree(ctx context.Context) ([]*domain.Category, error) {
	return uc.categoryRepo.GetTree(ctx)
}
