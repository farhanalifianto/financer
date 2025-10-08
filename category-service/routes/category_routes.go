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
	c := api.Group("/category") 

	c.Get("/list", authMiddleware, cc.List)                                       
	c.Get("/all", authMiddleware, middleware.RoleRequired("admin"), cc.List)  
	c.Get("/:id", authMiddleware, cc.Get)                                     
	c.Post("/", authMiddleware, cc.Create)                                    
	c.Put("/:id", authMiddleware, cc.Update)                                  
	c.Delete("/:id", authMiddleware, cc.Delete)                               
}