package routes

import (
	"transaction-service/controller"
	"transaction-service/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterTransactionRoutes(app *fiber.App, db *gorm.DB, authMiddleware fiber.Handler) {
	tc := &controller.TransactionController{DB: db}

	api := app.Group("/api")
	t := api.Group("/transaction")

	t.Get("/",authMiddleware, tc.List)
	t.Get("/all",authMiddleware,middleware.RoleRequired("admin"), tc.ListAll)
	t.Get("/balance",authMiddleware, tc.GetBalance)
	t.Get("/:id",authMiddleware, tc.Get)
	t.Post("/", authMiddleware, tc.Create)
	t.Put("/:id", authMiddleware, tc.Update)
	t.Delete("/:id", authMiddleware, tc.Delete)
}
