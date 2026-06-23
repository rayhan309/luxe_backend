package http

import (
	"github.com/gin-gonic/gin"
	"github.com/luxe/backend/internal/domain"
	"github.com/luxe/backend/internal/pkg/pagination"
	"github.com/luxe/backend/internal/pkg/response"
	"github.com/luxe/backend/internal/pkg/validator"
	"github.com/luxe/backend/internal/usecase"
)

type OrderHandler struct {
	orderUC *usecase.OrderUseCase
}

func NewOrderHandler(uc *usecase.OrderUseCase) *OrderHandler {
	return &OrderHandler{orderUC: uc}
}

func (h *OrderHandler) Create(c *gin.Context) {
	userID := c.GetString("user_id")
	var input domain.CreateOrderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}
	if errs := validator.Validate(input); errs != nil {
		response.ValidationError(c, errs)
		return
	}

	order, err := h.orderUC.Create(c.Request.Context(), userID, input)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, "order placed successfully", order)
}

func (h *OrderHandler) GetMyOrders(c *gin.Context) {
	userID := c.GetString("user_id")
	pag := pagination.FromContext(c)

	orders, total, err := h.orderUC.GetUserOrders(c.Request.Context(), userID, pag.Page, pag.Limit)
	if err != nil {
		response.InternalError(c, "failed to fetch orders")
		return
	}
	response.Paginated(c, "orders fetched", orders, pag.Page, pag.Limit, total)
}

func (h *OrderHandler) GetByID(c *gin.Context) {
	order, err := h.orderUC.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.NotFound(c, "order not found")
		return
	}

	// Verify ownership (unless admin)
	userID := c.GetString("user_id")
	if order.UserID.Hex() != userID && c.GetString("user_role") != domain.RoleAdmin {
		response.Forbidden(c, "access denied")
		return
	}

	response.Success(c, "order fetched", order)
}

// Admin handlers
func (h *OrderHandler) AdminList(c *gin.Context) {
	pag := pagination.FromContext(c)
	orders, total, err := h.orderUC.List(
		c.Request.Context(),
		c.Query("status"),
		pag.Page,
		pag.Limit,
		c.Query("search"),
	)
	if err != nil {
		response.InternalError(c, "failed to fetch orders")
		return
	}
	response.Paginated(c, "orders fetched", orders, pag.Page, pag.Limit, total)
}

func (h *OrderHandler) AdminUpdateStatus(c *gin.Context) {
	var input struct {
		Status string `json:"status" validate:"required"`
		Note   string `json:"note"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	if err := h.orderUC.UpdateStatus(c.Request.Context(), c.Param("id"), input.Status, input.Note); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, "order status updated", nil)
}
