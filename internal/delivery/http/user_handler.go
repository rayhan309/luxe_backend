package http

import (
	"github.com/gin-gonic/gin"
	"github.com/luxe/backend/internal/domain"
	"github.com/luxe/backend/internal/pkg/response"
	"github.com/luxe/backend/internal/pkg/validator"
	"github.com/luxe/backend/internal/usecase"
)

type UserHandler struct {
	userUC *usecase.UserUseCase
}

func NewUserHandler(uc *usecase.UserUseCase) *UserHandler {
	return &UserHandler{userUC: uc}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	user, err := h.userUC.GetProfile(c.Request.Context(), c.GetString("user_id"))
	if err != nil {
		response.NotFound(c, "user not found")
		return
	}
	response.Success(c, "profile fetched", user)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	var input domain.UpdateProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}
	if errs := validator.Validate(input); errs != nil {
		response.ValidationError(c, errs)
		return
	}
	user, err := h.userUC.UpdateProfile(c.Request.Context(), c.GetString("user_id"), input)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, "profile updated", user)
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	var input domain.ChangePasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}
	if err := h.userUC.ChangePassword(c.Request.Context(), c.GetString("user_id"), input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, "password changed", nil)
}

func (h *UserHandler) GetWishlist(c *gin.Context) {
	user, err := h.userUC.GetWishlist(c.Request.Context(), c.GetString("user_id"))
	if err != nil {
		response.InternalError(c, "failed to fetch wishlist")
		return
	}
	response.Success(c, "wishlist fetched", user.Wishlist)
}

func (h *UserHandler) ToggleWishlist(c *gin.Context) {
	if err := h.userUC.ToggleWishlist(c.Request.Context(), c.GetString("user_id"), c.Param("productId")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, "wishlist updated", nil)
}

func (h *UserHandler) AddAddress(c *gin.Context) {
	var address domain.Address
	if err := c.ShouldBindJSON(&address); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}
	if err := h.userUC.AddAddress(c.Request.Context(), c.GetString("user_id"), address); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, "address added", nil)
}

func (h *UserHandler) UpdateAddress(c *gin.Context) {
	var address domain.Address
	if err := c.ShouldBindJSON(&address); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}
	if err := h.userUC.UpdateAddress(c.Request.Context(), c.GetString("user_id"), c.Param("addressId"), address); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, "address updated", nil)
}

func (h *UserHandler) DeleteAddress(c *gin.Context) {
	if err := h.userUC.DeleteAddress(c.Request.Context(), c.GetString("user_id"), c.Param("addressId")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, "address deleted", nil)
}
