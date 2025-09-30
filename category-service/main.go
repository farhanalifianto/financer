package main

import (
	"log"
	"net"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	pb "category-service/proto/category"

	"category-service/grpc_server"
	"category-service/middleware"
	"category-service/model"
	"category-service/routes"
)

var DB *gorm.DB

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
		log.Fatal("failed to connect product db:", err)
	}

	if err := DB.AutoMigrate(&model.Category{}); err != nil {
		log.Fatal(err)
	}
}

func main() {
	initDB()

	go func(){
		app := fiber.New()
		app.Use(logger.New())

		// inject DB & middleware ke routes
		routes.RegisterCategoryRoutes(app, DB, middleware.AuthRequired)
		if err := app.Listen(":3004"); err != nil {
            log.Fatalf("failed to start REST server: %v", err)
        }
	}()
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterCategoryServiceServer(grpcServer, &grpc_server.CategoryServer{DB: DB})

	log.Println("gRPC CategoryService running on :50052")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	
}

func getEnv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
