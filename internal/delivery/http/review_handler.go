package http

import (
	"github.com/gin-gonic/gin"
	"github.com/luxe/backend/internal/domain"
	"github.com/luxe/backend/internal/pkg/pagination"
	"github.com/luxe/backend/internal/pkg/response"
	"github.com/luxe/backend/internal/pkg/validator"
	"github.com/luxe/backend/internal/usecase"
)

type ReviewHandler struct {
	reviewUC *usecase.ReviewUseCase
}

func NewReviewHandler(uc *usecase.ReviewUseCase) *ReviewHandler {
	return &ReviewHandler{reviewUC: uc}
}

func (h *ReviewHandler) GetProductReviews(c *gin.Context) {
	pag := pagination.FromContext(c)
	reviews, total, err := h.reviewUC.GetProductReviews(c.Request.Context(), c.Param("productId"), pag.Page, pag.Limit)
	if err != nil {
		response.InternalError(c, "failed to fetch reviews")
		return
	}
	response.Paginated(c, "reviews fetched", reviews, pag.Page, pag.Limit, total)
}

func (h *ReviewHandler) Create(c *gin.Context) {
	var input domain.CreateReviewInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}
	if errs := validator.Validate(input); errs != nil {
		response.ValidationError(c, errs)
		return
	}

	review, err := h.reviewUC.Create(c.Request.Context(), c.GetString("user_id"), input)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, "review submitted (pending approval)", review)
}

// Admin handlers
func (h *ReviewHandler) AdminList(c *gin.Context) {
	pag := pagination.FromContext(c)
	var approved *bool
	if v := c.Query("approved"); v == "true" {
		t := true
		approved = &t
	} else if v == "false" {
		f := false
		approved = &f
	}

	reviews, total, err := h.reviewUC.List(c.Request.Context(), approved, pag.Page, pag.Limit)
	if err != nil {
		response.InternalError(c, "failed to fetch reviews")
		return
	}
	response.Paginated(c, "reviews fetched", reviews, pag.Page, pag.Limit, total)
}

func (h *ReviewHandler) AdminApprove(c *gin.Context) {
	if err := h.reviewUC.Approve(c.Request.Context(), c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, "review approved", nil)
}

func (h *ReviewHandler) AdminDelete(c *gin.Context) {
	if err := h.reviewUC.Delete(c.Request.Context(), c.Param("id")); err != nil {
		response.InternalError(c, "failed to delete review")
		return
	}
	response.Success(c, "review deleted", nil)
}
