package usecase

import (
	"context"
	"errors"

	"github.com/luxe/backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartUseCase struct {
	cartRepo    domain.CartRepository
	productRepo domain.ProductRepository
}

func NewCartUseCase(cartRepo domain.CartRepository, productRepo domain.ProductRepository) *CartUseCase {
	return &CartUseCase{cartRepo: cartRepo, productRepo: productRepo}
}

func (uc *CartUseCase) GetCart(ctx context.Context, userID string) (*domain.Cart, error) {
	return uc.cartRepo.GetByUserID(ctx, userID)
}

func (uc *CartUseCase) AddItem(ctx context.Context, userID string, input domain.AddToCartInput) (*domain.Cart, error) {
	product, err := uc.productRepo.FindByID(ctx, input.ProductID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	if product.Stock < input.Quantity {
		return nil, errors.New("insufficient stock")
	}

	pid, _ := primitive.ObjectIDFromHex(input.ProductID)
	price := product.Price
	image := product.Thumbnail
	sku := product.SKU

	item := domain.CartItem{
		ProductID:  pid,
		Name:       product.Name,
		Image:      image,
		SKU:        sku,
		Price:      price,
		Quantity:   input.Quantity,
		Attributes: input.Attributes,
	}

	if input.VariantID != "" {
		vid, _ := primitive.ObjectIDFromHex(input.VariantID)
		item.VariantID = &vid
		for _, v := range product.Variants {
			if v.ID.Hex() == input.VariantID {
				item.Price = v.Price
				item.SKU = v.SKU
				if len(v.Images) > 0 {
					item.Image = v.Images[0]
				}
				break
			}
		}
	}

	if err := uc.cartRepo.AddItem(ctx, userID, item); err != nil {
		return nil, err
	}

	return uc.cartRepo.GetByUserID(ctx, userID)
}

func (uc *CartUseCase) UpdateItem(ctx context.Context, userID, itemID string, quantity int) (*domain.Cart, error) {
	if quantity < 1 {
		return nil, errors.New("quantity must be at least 1")
	}
	if err := uc.cartRepo.UpdateItemQuantity(ctx, userID, itemID, quantity); err != nil {
		return nil, err
	}
	return uc.cartRepo.GetByUserID(ctx, userID)
}

func (uc *CartUseCase) RemoveItem(ctx context.Context, userID, itemID string) (*domain.Cart, error) {
	if err := uc.cartRepo.RemoveItem(ctx, userID, itemID); err != nil {
		return nil, err
	}
	return uc.cartRepo.GetByUserID(ctx, userID)
}

func (uc *CartUseCase) Clear(ctx context.Context, userID string) error {
	return uc.cartRepo.Clear(ctx, userID)
}
