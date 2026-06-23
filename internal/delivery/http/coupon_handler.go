package http

import (
	"github.com/gin-gonic/gin"
	"github.com/luxe/backend/internal/domain"
	"github.com/luxe/backend/internal/pkg/response"
	"github.com/luxe/backend/internal/pkg/validator"
	"github.com/luxe/backend/internal/usecase"
)

type CouponHandler struct {
	couponUC *usecase.CouponUseCase
}

func NewCouponHandler(uc *usecase.CouponUseCase) *CouponHandler {
	return &CouponHandler{couponUC: uc}
}

func (h *CouponHandler) Validate(c *gin.Context) {
	var input domain.ValidateCouponInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	result, err := h.couponUC.Validate(c.Request.Context(), input)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result.Message, result)
}

func (h *CouponHandler) AdminCreate(c *gin.Context) {
	var coupon domain.Coupon
	if err := c.ShouldBindJSON(&coupon); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}
	if errs := validator.Validate(coupon); errs != nil {
		response.ValidationError(c, errs)
		return
	}
	created, err := h.couponUC.Create(c.Request.Context(), &coupon)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, "coupon created", created)
}

func (h *CouponHandler) AdminList(c *gin.Context) {
	coupons, total, err := h.couponUC.List(c.Request.Context(), 1, 50)
	if err != nil {
		response.InternalError(c, "failed to fetch coupons")
		return
	}
	response.Paginated(c, "coupons fetched", coupons, 1, 50, total)
}

func (h *CouponHandler) AdminUpdate(c *gin.Context) {
	var update map[string]interface{}
	if err := c.ShouldBindJSON(&update); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}
	coupon, err := h.couponUC.Update(c.Request.Context(), c.Param("id"), update)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, "coupon updated", coupon)
}

func (h *CouponHandler) AdminDelete(c *gin.Context) {
	if err := h.couponUC.Delete(c.Request.Context(), c.Param("id")); err != nil {
		response.InternalError(c, "failed to delete coupon")
		return
	}
	response.Success(c, "coupon deleted", nil)
}
