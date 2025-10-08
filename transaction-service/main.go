package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"transaction-service/grpc_client"
	"transaction-service/middleware"
	"transaction-service/model"
	"transaction-service/routes"
)

var (
	DB          *gorm.DB
	Redis       *redis.Client
	UserClient  *grpc_client.UserClient
)

func initDB() {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	pass := getEnv("DB_PASS", "postgres")
	name := getEnv("DB_NAME", "transactiondb")

	dsn := "host=" + host + " user=" + user + " password=" + pass + " dbname=" + name + " port=" + port + " sslmode=disable TimeZone=UTC"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect transaction db:", err)
	}

	if err := DB.AutoMigrate(&model.Transaction{}); err != nil {
		log.Fatal(err)
	}
	log.Println("✅ Connected to Postgres")
}

func initRedis() {
	addr := getEnv("REDIS_HOST", "localhost:6379")

	Redis = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	// optional: test ping
	if err := Redis.Ping(ctx()).Err(); err != nil {
		log.Fatal("failed to connect to redis:", err)
	}
	log.Println("✅ Connected to Redis at", addr)
}



func main() {
	initDB()
	initRedis()

	app := fiber.New()
	app.Use(logger.New())


	// Pass DB & Redis & middleware ke routes
	routes.RegisterTransactionRoutes(app, DB, Redis, middleware.AuthRequired)

	log.Println("🚀 Transaction service running on :3003")
	if err := app.Listen(":3003"); err != nil {
		log.Fatal(err)
	}
}

func getEnv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}

// helper untuk context redis
func ctx() context.Context {
	return context.Background()
}
