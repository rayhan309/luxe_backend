package http

import (
	"github.com/gin-gonic/gin"
	"github.com/luxe/backend/internal/domain"
	"github.com/luxe/backend/internal/pkg/response"
	"github.com/luxe/backend/internal/usecase"
)

type BannerHandler struct {
	bannerUC *usecase.BannerUseCase
}

func NewBannerHandler(uc *usecase.BannerUseCase) *BannerHandler {
	return &BannerHandler{bannerUC: uc}
}

func (h *BannerHandler) List(c *gin.Context) {
	banners, err := h.bannerUC.List(c.Request.Context(), true, c.Query("position"))
	if err != nil {
		response.InternalError(c, "failed to fetch banners")
		return
	}
	response.Success(c, "banners fetched", banners)
}

func (h *BannerHandler) AdminList(c *gin.Context) {
	banners, err := h.bannerUC.List(c.Request.Context(), false, "")
	if err != nil {
		response.InternalError(c, "failed to fetch banners")
		return
	}
	response.Success(c, "banners fetched", banners)
}

func (h *BannerHandler) AdminCreate(c *gin.Context) {
	var input domain.CreateBannerInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}
	banner, err := h.bannerUC.Create(c.Request.Context(), input)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, "banner created", banner)
}

func (h *BannerHandler) AdminUpdate(c *gin.Context) {
	var update map[string]interface{}
	if err := c.ShouldBindJSON(&update); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}
	banner, err := h.bannerUC.Update(c.Request.Context(), c.Param("id"), update)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, "banner updated", banner)
}

func (h *BannerHandler) AdminDelete(c *gin.Context) {
	if err := h.bannerUC.Delete(c.Request.Context(), c.Param("id")); err != nil {
		response.InternalError(c, "failed to delete banner")
		return
	}
	response.Success(c, "banner deleted", nil)
}
