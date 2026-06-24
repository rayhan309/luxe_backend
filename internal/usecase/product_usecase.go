package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/luxe/backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ProductUseCase handles product business logic
type ProductUseCase struct {
	productRepo  domain.ProductRepository
	categoryRepo domain.CategoryRepository
}

func NewProductUseCase(productRepo domain.ProductRepository, categoryRepo domain.CategoryRepository) *ProductUseCase {
	return &ProductUseCase{productRepo: productRepo, categoryRepo: categoryRepo}
}

// Create creates a new product
func (uc *ProductUseCase) Create(ctx context.Context, input domain.CreateProductInput) (*domain.Product, error) {
	// Validate category
	if input.CategoryID != "" {
		if _, err := uc.categoryRepo.FindByID(ctx, input.CategoryID); err != nil {
			return nil, errors.New("category not found")
		}
	}

	product := &domain.Product{
		Name:          input.Name,
		Description:   input.Description,
		ShortDesc:     input.ShortDesc,
		Price:         input.Price,
		ComparePrice:  input.ComparePrice,
		CostPrice:     input.CostPrice,
		Images:        input.Images,
		Thumbnail:     input.Thumbnail,
		SKU:           input.SKU,
		Stock:         input.Stock,
		LowStockAlert: input.LowStockAlert,
		Weight:        input.Weight,
		Tags:          input.Tags,
		Variants:      input.Variants,
		Attributes:    input.Attributes,
		Status:        input.Status,
		IsFeatured:    input.IsFeatured,
		IsBestSeller:  input.IsBestSeller,
		IsNewArrival:  input.IsNewArrival,
		MetaTitle:     input.MetaTitle,
		MetaDesc:      input.MetaDesc,
	}

	if input.CategoryID != "" {
		categoryOID, err := primitive.ObjectIDFromHex(input.CategoryID)
		if err != nil {
			return nil, errors.New("invalid category id")
		}
		product.CategoryID = categoryOID
	}

	return uc.productRepo.Create(ctx, product)
}

// GetBySlug retrieves a product by slug
func (uc *ProductUseCase) GetBySlug(ctx context.Context, slug string) (*domain.Product, error) {
	product, err := uc.productRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, errors.New("product not found")
	}
	return product, nil
}

// GetByID retrieves a product by ID
func (uc *ProductUseCase) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	product, err := uc.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("product not found")
	}
	return product, nil
}

// List returns paginated filtered products
func (uc *ProductUseCase) List(ctx context.Context, filter domain.ProductFilter) ([]*domain.Product, int64, error) {
	return uc.productRepo.List(ctx, filter)
}

// Update updates a product
func (uc *ProductUseCase) Update(ctx context.Context, id string, input map[string]interface{}) (*domain.Product, error) {
	input["updated_at"] = time.Now()
	return uc.productRepo.Update(ctx, id, input)
}

// Delete deletes a product
func (uc *ProductUseCase) Delete(ctx context.Context, id string) error {
	return uc.productRepo.Delete(ctx, id)
}

// GetFeatured returns featured products
func (uc *ProductUseCase) GetFeatured(ctx context.Context, limit int) ([]*domain.Product, error) {
	return uc.productRepo.GetFeatured(ctx, limit)
}

// GetBestSellers returns best-selling products
func (uc *ProductUseCase) GetBestSellers(ctx context.Context, limit int) ([]*domain.Product, error) {
	return uc.productRepo.GetBestSellers(ctx, limit)
}

// GetNewArrivals returns new arrival products
func (uc *ProductUseCase) GetNewArrivals(ctx context.Context, limit int) ([]*domain.Product, error) {
	return uc.productRepo.GetNewArrivals(ctx, limit)
}

// GetRelated returns related products
func (uc *ProductUseCase) GetRelated(ctx context.Context, categoryID, excludeSlug string, limit int) ([]*domain.Product, error) {
	excludeID := excludeSlug
	if product, err := uc.productRepo.FindBySlug(ctx, excludeSlug); err == nil {
		excludeID = product.ID.Hex()
	}
	return uc.productRepo.GetRelated(ctx, categoryID, excludeID, limit)
}
