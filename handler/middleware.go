package handler

import (
	"time"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gofiber/fiber/v2"
	"github.com/iqbaleff214/easynote-backend-go/helper"
	"github.com/iqbaleff214/easynote-backend-go/user"
)

func AuthMiddleware(jwtSecret string, userService user.Service) func(c *fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(jwtSecret)},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(
				helper.APIResponse("Invalid or expired JWT", "error", fiber.StatusUnauthorized, nil),
			)
		},
		SuccessHandler: func(c *fiber.Ctx) error {
			t := c.Locals("user").(*jwt.Token)
			claims := t.Claims.(jwt.MapClaims)

			// Check if token has expired date
			expiredAt, err := time.Parse(time.RFC822, claims["expired_at"].(string))
			if err != nil {
				return c.Status(fiber.StatusUnauthorized).JSON(
					helper.APIResponse("Invalid or expired JWT", "error", fiber.StatusUnauthorized, nil),
				)
			}

			// Check if token has expired
			if time.Now().After(expiredAt) {
				return c.Status(fiber.StatusUnauthorized).JSON(
					helper.APIResponse("Expired JWT", "error", fiber.StatusUnauthorized, nil),
				)
			}

			// TODO: is it necessary?
			userID := int(claims["user_id"].(float64))
			user, err := userService.GetUserByID(userID)
			if err != nil {
				return c.Status(fiber.StatusUnauthorized).JSON(
					helper.APIResponse("User not found", "error", fiber.StatusUnauthorized, nil),
				)
			}
			c.Locals("currentUser", user)

			return c.Next()
		},
	})
}
