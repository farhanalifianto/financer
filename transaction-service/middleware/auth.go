package middleware

import (
	"log"
	"strings"

	"transaction-service/grpc_client"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc/status"
)



func AuthRequired(c *fiber.Ctx) error {
	UserClient := grpc_client.NewUserClient()
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(401).JSON(fiber.Map{"error": "missing auth"})
	}

	token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	if token == "" {
		return c.Status(401).JSON(fiber.Map{"error": "invalid auth header"})
	}
	

	user, err := UserClient.GetMe(token)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			log.Printf("gRPC error - code: %v, message: %s", st.Code(), st.Message())
		} else {
			log.Printf("Unknown gRPC error: %v", err)
		}
		return c.Status(401).JSON(fiber.Map{
			"error": st.Message(), // atau pakai err.Error() jika bukan gRPC error
		})
	}

	c.Locals("user_id", user.Id)
	c.Locals("user_role", user.Role)

	return c.Next()
}

func RoleRequired(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("user_role")
		if userRole == nil {
			return c.Status(403).JSON(fiber.Map{"error": "no role"})
		}

		role := userRole.(string)
		for _, r := range roles {
			if role == r {
				return c.Next()
			}
		}
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}
}
