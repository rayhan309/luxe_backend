package http

import (
	"github.com/gin-gonic/gin"
	"github.com/luxe/backend/internal/domain"
	"github.com/luxe/backend/internal/pkg/response"
	"github.com/luxe/backend/internal/pkg/validator"
	"github.com/luxe/backend/internal/usecase"
)

type AuthHandler struct {
	authUC *usecase.AuthUseCase
}

func NewAuthHandler(authUC *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUC: authUC}
}

// Register godoc
// @Summary Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param body body domain.RegisterInput true "Registration data"
// @Success 201 {object} domain.AuthResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var input domain.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}
	if errs := validator.Validate(input); errs != nil {
		response.ValidationError(c, errs)
		return
	}

	result, err := h.authUC.Register(c.Request.Context(), input)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, "account created successfully", result)
}

// Login godoc
// @Summary Login
// @Tags auth
// @Accept json
// @Produce json
// @Param body body domain.LoginInput true "Login credentials"
// @Success 200 {object} domain.AuthResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var input domain.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}
	if errs := validator.Validate(input); errs != nil {
		response.ValidationError(c, errs)
		return
	}

	result, err := h.authUC.Login(c.Request.Context(), input)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}
	response.Success(c, "login successful", result)
}

// ForgotPassword godoc
// @Summary Request password reset
// @Tags auth
// @Accept json
// @Produce json
// @Router /auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var input struct {
		Email string `json:"email" validate:"required,email"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}
	if errs := validator.Validate(input); errs != nil {
		response.ValidationError(c, errs)
		return
	}

	token, err := h.authUC.ForgotPassword(c.Request.Context(), input.Email)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, "reset link generated", gin.H{"reset_token": token})
}

// ResetPassword godoc
// @Summary Reset password with token
// @Tags auth
// @Accept json
// @Produce json
// @Router /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var input struct {
		Token       string `json:"token" validate:"required"`
		NewPassword string `json:"new_password" validate:"required,min=8"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	if err := h.authUC.ResetPassword(c.Request.Context(), input.Token, input.NewPassword); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, "password reset successfully", nil)
}
