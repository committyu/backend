package middleware

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

// AuthMiddleware はJWTの検証を行い、userIDをContextにセットします
func AuthMiddleware() echo.MiddlewareFunc {
	// 環境変数が空の場合のデフォルト値を設定するか、パニックを防ぐチェック
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// 開発時に気づけるようログを出すか、固定値を入れる（本番では厳禁）
		secret = "default_secret" 
	}

	return echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(secret),
		// TokenLookup: "header:Authorization", // デフォルトでこれになっています
		// AuthScheme:  "Bearer",               // デフォルトでこれになっています
		
		SuccessHandler: func(c echo.Context) {
			// echo-jwtがセットした "user" (token) を取得
			token, ok := c.Get("user").(*jwt.Token)
			if !ok {
				return
			}

			// Claimsから user_id を取り出す
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				if userID, ok := claims["user_id"].(string); ok {
					// ハンドラー側で c.Get("userID") で取れるように再セット
					c.Set("userID", userID)
				}
			}
		},
	})
}