package routes

import (
	"category-service/controller"
	"category-service/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterCategoryRoutes(app *fiber.App, db *gorm.DB, authMiddleware fiber.Handler) {
	cc := &controller.CategoryController{DB: db}

	api := app.Group("/api")
	c := api.Group("/category") // pakai plural biar konsisten REST

	c.Get("/", authMiddleware, cc.List)                                       // list kategori milik user
	c.Get("/all", authMiddleware, middleware.RoleRequired("admin"), cc.List)  // list semua kategori (khusus admin)
	c.Get("/:id", authMiddleware, cc.Get)                                     // detail kategori
	c.Post("/", authMiddleware, cc.Create)                                    // buat kategori
	c.Put("/:id", authMiddleware, cc.Update)                                  // update kategori
	c.Delete("/:id", authMiddleware, cc.Delete)                               // hapus kategori
}
