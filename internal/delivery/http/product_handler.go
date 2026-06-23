package http

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/luxe/backend/internal/domain"
	"github.com/luxe/backend/internal/pkg/pagination"
	"github.com/luxe/backend/internal/pkg/response"
	"github.com/luxe/backend/internal/pkg/validator"
	"github.com/luxe/backend/internal/usecase"
)

type ProductHandler struct {
	productUC *usecase.ProductUseCase
}

func NewProductHandler(productUC *usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{productUC: productUC}
}

func (h *ProductHandler) List(c *gin.Context) {
	pag := pagination.FromContext(c)

	minPrice, _ := strconv.ParseFloat(c.Query("min_price"), 64)
	maxPrice, _ := strconv.ParseFloat(c.Query("max_price"), 64)

	var ids []string
	if idsStr := c.Query("ids"); idsStr != "" {
		ids = strings.Split(idsStr, ",")
	}

	filter := domain.ProductFilter{
		CategoryID: c.Query("category"),
		MinPrice:   minPrice,
		MaxPrice:   maxPrice,
		Search:     c.Query("search"),
		Status:     domain.ProductStatusActive,
		SortBy:     c.DefaultQuery("sort", "newest"),
		Page:       pag.Page,
		Limit:      pag.Limit,
		IDs:        ids,
	}

	products, total, err := h.productUC.List(c.Request.Context(), filter)
	if err != nil {
		response.InternalError(c, "failed to fetch products")
		return
	}
	response.Paginated(c, "products fetched", products, pag.Page, pag.Limit, total)
}

func (h *ProductHandler) GetBySlug(c *gin.Context) {
	product, err := h.productUC.GetBySlug(c.Request.Context(), c.Param("slug"))
	if err != nil {
		response.NotFound(c, "product not found")
		return
	}
	response.Success(c, "product fetched", product)
}

func (h *ProductHandler) GetFeatured(c *gin.Context) {
	products, err := h.productUC.GetFeatured(c.Request.Context(), 8)
	if err != nil {
		response.InternalError(c, "failed to fetch featured products")
		return
	}
	response.Success(c, "featured products fetched", products)
}

func (h *ProductHandler) GetBestSellers(c *gin.Context) {
	products, err := h.productUC.GetBestSellers(c.Request.Context(), 8)
	if err != nil {
		response.InternalError(c, "failed to fetch best sellers")
		return
	}
	response.Success(c, "best sellers fetched", products)
}

func (h *ProductHandler) GetNewArrivals(c *gin.Context) {
	products, err := h.productUC.GetNewArrivals(c.Request.Context(), 8)
	if err != nil {
		response.InternalError(c, "failed to fetch new arrivals")
		return
	}
	response.Success(c, "new arrivals fetched", products)
}

func (h *ProductHandler) GetRelated(c *gin.Context) {
	products, err := h.productUC.GetRelated(c.Request.Context(), c.Query("category"), c.Param("slug"), 6)
	if err != nil {
		response.InternalError(c, "failed to fetch related products")
		return
	}
	response.Success(c, "related products fetched", products)
}

// Admin handlers

func (h *ProductHandler) AdminCreate(c *gin.Context) {
	var input domain.CreateProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}
	if errs := validator.Validate(input); errs != nil {
		response.ValidationError(c, errs)
		return
	}

	product, err := h.productUC.Create(c.Request.Context(), input)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, "product created", product)
}

func (h *ProductHandler) AdminUpdate(c *gin.Context) {
	var update map[string]interface{}
	if err := c.ShouldBindJSON(&update); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	product, err := h.productUC.Update(c.Request.Context(), c.Param("id"), update)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, "product updated", product)
}

func (h *ProductHandler) AdminDelete(c *gin.Context) {
	if err := h.productUC.Delete(c.Request.Context(), c.Param("id")); err != nil {
		response.InternalError(c, "failed to delete product")
		return
	}
	response.Success(c, "product deleted", nil)
}

func (h *ProductHandler) AdminList(c *gin.Context) {
	pag := pagination.FromContext(c)
	filter := domain.ProductFilter{
		Search: c.Query("search"),
		Status: c.Query("status"),
		Page:   pag.Page,
		Limit:  pag.Limit,
	}
	products, total, err := h.productUC.List(c.Request.Context(), filter)
	if err != nil {
		response.InternalError(c, "failed to fetch products")
		return
	}
	response.Paginated(c, "products fetched", products, pag.Page, pag.Limit, total)
}
