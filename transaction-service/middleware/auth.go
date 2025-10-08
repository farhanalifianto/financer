package middleware

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"transaction-service/grpc_client"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/status"
)

var (
	ctx = context.Background()
	rdb *redis.Client
)

// Struct untuk data yang disimpan di cache
type CachedUser struct {
	ID   uint `json:"id"`
	Role string `json:"role"`
}

// InitRedis dipanggil di main.go sekali saja
func InitRedis(client *redis.Client) {
	rdb = client
}

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

	cacheKey := "auth:" + token

	// 1️⃣ Cek cache Redis dulu
	if rdb != nil {
		cached, err := rdb.Get(ctx, cacheKey).Result()
		if err == nil && cached != "" {
			var cu CachedUser
			if err := json.Unmarshal([]byte(cached), &cu); err == nil {
				c.Locals("user_id", cu.ID)
				c.Locals("user_role", cu.Role)
				return c.Next()
			}
		}
	}

	// 2️⃣ Kalau tidak ada cache, panggil user-service
	user, err := UserClient.GetMe(token)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			log.Printf("gRPC error - code: %v, message: %s", st.Code(), st.Message())
		} else {
			log.Printf("Unknown gRPC error: %v", err)
		}
		return c.Status(401).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// 3️⃣ Simpan hasil ke Redis
	if rdb != nil {
		cu := CachedUser{
			ID:   user.Id,
			Role: user.Role,
		}
		jsonVal, _ := json.Marshal(cu)
		// TTL disesuaikan dengan sisa masa berlaku JWT, di sini contoh 10 menit
		rdb.Set(ctx, cacheKey, jsonVal, 10*time.Minute)
	}

	// 4️⃣ Simpan ke Fiber locals
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
