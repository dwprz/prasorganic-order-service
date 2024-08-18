package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func (m *Middleware) VerifyAdmin(c *fiber.Ctx) error {
	userData := c.Locals("user_data").(jwt.MapClaims)

	role := userData["role"].(string)

	if role != "SUPER_ADMIN" && role != "ADMIN" {
		return c.Status(403).JSON(fiber.Map{"errors": "access denied"})
	}

	return c.Next()
}
