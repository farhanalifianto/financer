package controller

import (
	"log"
	"strconv"
	"time"

	"category-service/grpc_client"
	"category-service/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CategoryController struct {
	DB *gorm.DB
}

// List categories
func (cc *CategoryController) List(c *fiber.Ctx) error {
	
	var categories []model.Category

	userID := c.Locals("user_id").(uint)
	role := c.Locals("user_role").(string)

	var err error
	if role == "admin" {
		err = cc.DB.Find(&categories).Error
	} else {
		err = cc.DB.Where("owner_id = ?", userID).Find(&categories).Error
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Siapkan user client gRPC untuk ambil email pemilik kategori
	userClient := grpc_client.NewUserClient()

	// Buat response final
	var response []map[string]interface{}
	for _, cat := range categories {
		userInfo, userErr := userClient.GetUserEmail(cat.OwnerID)
		if userErr != nil {
			log.Printf("failed to get user email for owner_id %d: %v", cat.OwnerID, userErr)
		}

		response = append(response, map[string]interface{}{
			"id":     cat.ID,
			"name":   cat.Name,
			"owner":  userInfo.Email,
			"type":   cat.Type,
			"budget": cat.Budget,
		})
	}

	return c.JSON(response)
}

// Create category
func (cc *CategoryController) Create(c *fiber.Ctx) error {
	var category model.Category
	if err := c.BodyParser(&category); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if category.Type != "income" && category.Type != "expense" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Type should be 'income' or 'expense'",
		})
	}

	if category.Type == "income" {
		category.Budget = 0
	}

	userID := c.Locals("user_id").(uint)
	category.OwnerID = userID
	category.CreatedAt = time.Now()

	if err := cc.DB.Create(&category).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(category)
}


// gget category by id
func (cc *CategoryController) Get(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}

	userID := c.Locals("user_id").(uint)
	role := c.Locals("user_role").(string)

	var category model.Category
	if role == "admin" {
		if err := cc.DB.First(&category, id).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "not found"})
		}
	} else {
		if err := cc.DB.Where("id = ? AND owner_id = ?", id, userID).First(&category).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "not found"})
		}
	}

	userClient := grpc_client.NewUserClient()
	UserInfo, _ := userClient.GetUserEmail(category.OwnerID)

	response := map[string]interface{}{
		"id":     category.ID,
		"name":   category.Name,
		"owner":  UserInfo.Email,
		"type":   category.Type,
		"budget": category.Budget,
	}

	return c.JSON(response)
}

// update category

type EditRequest struct {
	Name 		string 		`json:"name"`
	Budget    	float64  	`json:"budget"`
	Type		string		`json:"type"`
}

func (cc *CategoryController) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}

	userID := c.Locals("user_id").(uint)

	var category model.Category
	if err := cc.DB.Where("id = ? AND owner_id = ?", id, userID).First(&category).Error; err != nil {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}

	var body EditRequest	
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	category.Name = body.Name

	if err := cc.DB.Save(&category).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(category)
}

// Delete category
func (cc *CategoryController) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}

	userID := c.Locals("user_id").(uint)

	var category model.Category
	if err := cc.DB.Where("id = ? AND owner_id = ?", id, userID).First(&category).Error; err != nil {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}

	if err := cc.DB.Delete(&category).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "category deleted"})
}
