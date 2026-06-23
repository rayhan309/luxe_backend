package usecase

import (
	"context"
	"errors"

	"github.com/luxe/backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderUseCase struct {
	orderRepo   domain.OrderRepository
	productRepo domain.ProductRepository
	cartRepo    domain.CartRepository
	couponRepo  domain.CouponRepository
}

func NewOrderUseCase(
	orderRepo domain.OrderRepository,
	productRepo domain.ProductRepository,
	cartRepo domain.CartRepository,
	couponRepo domain.CouponRepository,
) *OrderUseCase {
	return &OrderUseCase{
		orderRepo:   orderRepo,
		productRepo: productRepo,
		cartRepo:    cartRepo,
		couponRepo:  couponRepo,
	}
}

// Create builds and saves an order from the input
func (uc *OrderUseCase) Create(ctx context.Context, userID string, input domain.CreateOrderInput) (*domain.Order, error) {
	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	var items []domain.OrderItem
	var subtotal float64

	for _, i := range input.Items {
		product, err := uc.productRepo.FindByID(ctx, i.ProductID)
		if err != nil {
			return nil, errors.New("product not found: " + i.ProductID)
		}
		if product.Stock < i.Quantity {
			return nil, errors.New("insufficient stock for: " + product.Name)
		}

		price := product.Price
		image := product.Thumbnail
		sku := product.SKU

		// Check variant
		if i.VariantID != "" {
			for _, v := range product.Variants {
				if v.ID.Hex() == i.VariantID {
					price = v.Price
					sku = v.SKU
					if len(v.Images) > 0 {
						image = v.Images[0]
					}
					break
				}
			}
		}

		total := price * float64(i.Quantity)
		subtotal += total

		item := domain.OrderItem{
			ProductID:  product.ID,
			Name:       product.Name,
			Image:      image,
			SKU:        sku,
			Price:      price,
			Quantity:   i.Quantity,
			Total:      total,
			Attributes: i.Attributes,
		}
		if i.VariantID != "" {
			vid, _ := primitive.ObjectIDFromHex(i.VariantID)
			item.VariantID = &vid
		}
		items = append(items, item)
	}

	// Apply coupon
	var discount float64
	var couponCode string
	if input.CouponCode != "" {
		coupon, err := uc.couponRepo.FindByCode(ctx, input.CouponCode)
		if err == nil && coupon.IsActive {
			if coupon.DiscountType == domain.DiscountTypeFixed {
				discount = coupon.DiscountValue
			} else {
				discount = subtotal * (coupon.DiscountValue / 100)
				if coupon.MaxDiscount > 0 && discount > coupon.MaxDiscount {
					discount = coupon.MaxDiscount
				}
			}
			couponCode = coupon.Code
			_ = uc.couponRepo.IncrementUsage(ctx, coupon.Code)
		}
	}

	// Shipping cost (flat rate logic; can be extended)
	shippingCost := 0.0
	if subtotal < 100 {
		shippingCost = 9.99
	}

	total := subtotal - discount + shippingCost

	order := &domain.Order{
		UserID:          uid,
		Items:           items,
		ShippingAddress: input.ShippingAddress,
		PaymentMethod:   input.PaymentMethod,
		Subtotal:        subtotal,
		ShippingCost:    shippingCost,
		Discount:        discount,
		CouponCode:      couponCode,
		Total:           total,
		Note:            input.Note,
	}

	created, err := uc.orderRepo.Create(ctx, order)
	if err != nil {
		return nil, errors.New("failed to create order")
	}

	// Decrement stock for each item
	for _, i := range input.Items {
		_ = uc.productRepo.IncrementSoldCount(ctx, i.ProductID, i.Quantity)
	}

	// Clear cart
	_ = uc.cartRepo.Clear(ctx, userID)

	return created, nil
}

func (uc *OrderUseCase) GetByID(ctx context.Context, id string) (*domain.Order, error) {
	return uc.orderRepo.FindByID(ctx, id)
}

func (uc *OrderUseCase) GetUserOrders(ctx context.Context, userID string, page, limit int) ([]*domain.Order, int64, error) {
	return uc.orderRepo.ListByUser(ctx, userID, page, limit)
}

func (uc *OrderUseCase) UpdateStatus(ctx context.Context, id, status, note string) error {
	return uc.orderRepo.UpdateStatus(ctx, id, status, note)
}

func (uc *OrderUseCase) List(ctx context.Context, status string, page, limit int, search string) ([]*domain.Order, int64, error) {
	return uc.orderRepo.List(ctx, status, page, limit, search)
}

func (uc *OrderUseCase) GetStats(ctx context.Context) (*domain.OrderStats, error) {
	return uc.orderRepo.GetStats(ctx)
}

func (uc *OrderUseCase) GetRevenueChart(ctx context.Context, days int) ([]domain.RevenuePoint, error) {
	return uc.orderRepo.GetRevenueByPeriod(ctx, days)
}
