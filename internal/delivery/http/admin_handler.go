package http

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/luxe/backend/internal/pkg/pagination"
	"github.com/luxe/backend/internal/pkg/response"
	"github.com/luxe/backend/internal/usecase"
)

type AdminHandler struct {
	orderUC *usecase.OrderUseCase
	userUC  *usecase.UserUseCase
}

func NewAdminHandler(orderUC *usecase.OrderUseCase, userUC *usecase.UserUseCase) *AdminHandler {
	return &AdminHandler{orderUC: orderUC, userUC: userUC}
}

// Dashboard returns aggregated analytics
func (h *AdminHandler) Dashboard(c *gin.Context) {
	ctx := c.Request.Context()

	stats, err := h.orderUC.GetStats(ctx)
	if err != nil {
		response.InternalError(c, "failed to fetch stats")
		return
	}

	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	if days <= 0 || days > 365 {
		days = 30
	}

	chart, _ := h.orderUC.GetRevenueChart(ctx, days)

	response.Success(c, "dashboard data fetched", gin.H{
		"stats": stats,
		"chart": chart,
	})
}

// ListCustomers returns paginated customer list
func (h *AdminHandler) ListCustomers(c *gin.Context) {
	pag := pagination.FromContext(c)
	users, total, err := h.userUC.ListUsers(c.Request.Context(), pag.Page, pag.Limit, c.Query("search"))
	if err != nil {
		response.InternalError(c, "failed to fetch customers")
		return
	}
	response.Paginated(c, "customers fetched", users, pag.Page, pag.Limit, total)
}
