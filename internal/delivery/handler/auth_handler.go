package handler

import (
	"net/http"

	"github.com/LuuDinhTheTai/tzone/internal/dto"
	"github.com/LuuDinhTheTai/tzone/internal/service"
	"github.com/LuuDinhTheTai/tzone/util/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService}
}

// register endpoint
func (h *AuthHandler) Register(c *gin.Context) {

	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err := h.authService.Register(req.Email, req.Password)

	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, "register success", nil)
}

// login endpoint
func (h *AuthHandler) Login(c *gin.Context) {

	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	token, user, err := h.authService.Login(req.Email, req.Password)

	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "login success", gin.H{
		"token": token,
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
		},
	})
}
