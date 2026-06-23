package http

import (
	"github.com/gin-gonic/gin"
	"github.com/luxe/backend/internal/domain"
	"github.com/luxe/backend/internal/pkg/response"
	"github.com/luxe/backend/internal/pkg/validator"
	"github.com/luxe/backend/internal/usecase"
)

type CartHandler struct {
	cartUC *usecase.CartUseCase
}

func NewCartHandler(uc *usecase.CartUseCase) *CartHandler {
	return &CartHandler{cartUC: uc}
}

func (h *CartHandler) GetCart(c *gin.Context) {
	cart, err := h.cartUC.GetCart(c.Request.Context(), c.GetString("user_id"))
	if err != nil {
		response.InternalError(c, "failed to fetch cart")
		return
	}
	response.Success(c, "cart fetched", cart)
}

func (h *CartHandler) AddItem(c *gin.Context) {
	var input domain.AddToCartInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}
	if errs := validator.Validate(input); errs != nil {
		response.ValidationError(c, errs)
		return
	}

	cart, err := h.cartUC.AddItem(c.Request.Context(), c.GetString("user_id"), input)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, "item added to cart", cart)
}

func (h *CartHandler) UpdateItem(c *gin.Context) {
	var input domain.UpdateCartItemInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	cart, err := h.cartUC.UpdateItem(c.Request.Context(), c.GetString("user_id"), c.Param("itemId"), input.Quantity)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, "cart updated", cart)
}

func (h *CartHandler) RemoveItem(c *gin.Context) {
	cart, err := h.cartUC.RemoveItem(c.Request.Context(), c.GetString("user_id"), c.Param("itemId"))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, "item removed", cart)
}

func (h *CartHandler) Clear(c *gin.Context) {
	if err := h.cartUC.Clear(c.Request.Context(), c.GetString("user_id")); err != nil {
		response.InternalError(c, "failed to clear cart")
		return
	}
	response.Success(c, "cart cleared", nil)
}
