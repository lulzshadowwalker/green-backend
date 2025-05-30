package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lulzshadowwalker/green-backend/internal"
	"github.com/lulzshadowwalker/green-backend/internal/service"
)

type LoginHandler struct {
	userService *service.UserService
}

func NewLoginHandler(userService *service.UserService) *LoginHandler {
	return &LoginHandler{userService: userService}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

func (h *LoginHandler) RegisterRoutes(e *echo.Echo) {
	e.POST("/api/login", h.Login)
}

func (h *LoginHandler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request"})
	}
	user, err := h.userService.Authenticate(c.Request().Context(), req.Username, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid username or password"})
	}

	// Generate JWT token
	token, err := internal.GenerateJWT(user.ID, user.Username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to generate token"})
	}

	return c.JSON(http.StatusOK, LoginResponse{AccessToken: token})
}
