package http

import (
	// "net/http"
	// "strings"

	"github.com/labstack/echo/v4"
	// "github.com/lulzshadowwalker/green-backend/internal"
)

func JWTAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// authHeader := c.Request().Header.Get("Authorization")
		// if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		// 	return c.JSON(http.StatusUnauthorized, echo.Map{"error": "missing or invalid Authorization header"})
		// }
		// token := strings.TrimPrefix(authHeader, "Bearer ")
		// claims, err := internal.ParseJWT(token)
		// if err != nil {
		// 	return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid or expired token"})
		// }
		// // Optionally set user info in context
		// c.Set("user_id", claims.UserID)
		// c.Set("username", claims.Username)
		return next(c)
	}
}
