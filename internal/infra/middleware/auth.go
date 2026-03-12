package middleware

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware() echo.MiddlewareFunc {
	secret := os.Getenv("JWT_SECRET")

	return echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(secret),
		
		SuccessHandler: func(c echo.Context) {
			token, ok := c.Get("user").(*jwt.Token)
			if !ok {
				return
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				if userID, ok := claims["user_id"].(string); ok {
					// ハンドラー側で c.Get("userID") で取れるように
					c.Set("userID", userID)
				}
			}
		},
	})
}