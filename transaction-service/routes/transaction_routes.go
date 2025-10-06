package routes

import (
	"transaction-service/controller"
	"transaction-service/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func RegisterTransactionRoutes(app *fiber.App, db *gorm.DB,rdb *redis.Client, authMiddleware fiber.Handler) {
	tc := &controller.TransactionController{DB: db}

	api := app.Group("/api")
	t := api.Group("/transaction")

	t.Get("/",authMiddleware, tc.List)
	t.Get("/all",authMiddleware,middleware.RoleRequired("admin"), tc.ListAll)
	t.Get("/balance",authMiddleware, tc.GetBalance)
	t.Get("/balance/category",authMiddleware, tc.GetBalanceCategory)
	t.Get("/budget",authMiddleware, tc.GetBudgetStatus)
	t.Get("/:id",authMiddleware, tc.Get)
	t.Post("/", authMiddleware, tc.CreateFiltered)
	t.Put("/:id", authMiddleware, tc.Update)
	t.Delete("/:id", authMiddleware, tc.Delete)
}
