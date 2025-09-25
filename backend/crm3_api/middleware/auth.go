package middleware

import (
	"net/http"
	"strings"

	"crm3_api/dto"
	"crm3_api/utils"

	"github.com/labstack/echo/v4"
)

// JWTMiddleware JWT認証ミドルウェア
func JWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Authorizationヘッダーを取得
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
					Error:   "missing_token",
					Message: "認証トークンがありません",
				})
			}

			// "Bearer " プレフィックスをチェック
			if !strings.HasPrefix(authHeader, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
					Error:   "invalid_token_format",
					Message: "トークンの形式が正しくありません",
				})
			}

			// トークンを抽出
			token := strings.TrimPrefix(authHeader, "Bearer ")

			// トークンを検証
			claims, err := utils.ValidateToken(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
					Error:   "invalid_token",
					Message: "無効な認証トークンです",
				})
			}

			// ユーザー情報をコンテキストに設定
			c.Set("user_id", claims.UserID)
			c.Set("email", claims.Email)

			return next(c)
		}
	}
}