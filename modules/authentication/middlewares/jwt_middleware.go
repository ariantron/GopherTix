package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gopher_tix/configs"
	errs "gopher_tix/packages/common/errors"
	"strings"
	"time"
)

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

const (
	authHeader   = "Authorization"
	bearerPrefix = "Bearer "
)

func GenerateToken(userID, email string) (string, error) {
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // 24 hours
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString(configs.SecretKey)
	if err != nil {
		return "", errs.NewInternalServerError("Failed to generate token")
	}
	return token, nil
}

func JwtProtected() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authHeader := ctx.Get(authHeader)
		if authHeader == "" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header is required",
			})
		}

		if !strings.HasPrefix(authHeader, bearerPrefix) {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header format",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, bearerPrefix)
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return configs.SecretKey, nil
		})

		if err != nil || !token.Valid {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		ctx.Locals("user_id", claims.UserID)
		ctx.Locals("email", claims.Email)

		return ctx.Next()
	}
}
