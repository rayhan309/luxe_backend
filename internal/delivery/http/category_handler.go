package http

import (
	"github.com/gin-gonic/gin"
	"github.com/luxe/backend/internal/domain"
	"github.com/luxe/backend/internal/pkg/response"
	"github.com/luxe/backend/internal/pkg/validator"
	"github.com/luxe/backend/internal/usecase"
)

type CategoryHandler struct {
	categoryUC *usecase.CategoryUseCase
}

func NewCategoryHandler(uc *usecase.CategoryUseCase) *CategoryHandler {
	return &CategoryHandler{categoryUC: uc}
}

func (h *CategoryHandler) List(c *gin.Context) {
	cats, err := h.categoryUC.List(c.Request.Context(), true)
	if err != nil {
		response.InternalError(c, "failed to fetch categories")
		return
	}
	response.Success(c, "categories fetched", cats)
}

func (h *CategoryHandler) GetTree(c *gin.Context) {
	tree, err := h.categoryUC.GetTree(c.Request.Context())
	if err != nil {
		response.InternalError(c, "failed to fetch category tree")
		return
	}
	response.Success(c, "category tree fetched", tree)
}

func (h *CategoryHandler) GetBySlug(c *gin.Context) {
	cat, err := h.categoryUC.GetBySlug(c.Request.Context(), c.Param("slug"))
	if err != nil {
		response.NotFound(c, "category not found")
		return
	}
	response.Success(c, "category fetched", cat)
}

func (h *CategoryHandler) AdminCreate(c *gin.Context) {
	var input domain.CreateCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}
	if errs := validator.Validate(input); errs != nil {
		response.ValidationError(c, errs)
		return
	}
	cat, err := h.categoryUC.Create(c.Request.Context(), input)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, "category created", cat)
}

func (h *CategoryHandler) AdminUpdate(c *gin.Context) {
	var update map[string]interface{}
	if err := c.ShouldBindJSON(&update); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}
	cat, err := h.categoryUC.Update(c.Request.Context(), c.Param("id"), update)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, "category updated", cat)
}

func (h *CategoryHandler) AdminDelete(c *gin.Context) {
	if err := h.categoryUC.Delete(c.Request.Context(), c.Param("id")); err != nil {
		response.InternalError(c, "failed to delete category")
		return
	}
	response.Success(c, "category deleted", nil)
}

func (h *CategoryHandler) AdminList(c *gin.Context) {
	cats, err := h.categoryUC.List(c.Request.Context(), false)
	if err != nil {
		response.InternalError(c, "failed to fetch categories")
		return
	}
	response.Success(c, "categories fetched", cats)
}
