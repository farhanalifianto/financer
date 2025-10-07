package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"
	"transaction-service/grpc_client"

	"transaction-service/model"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"

	"gorm.io/gorm"
)

// List all transactions

type TransactionController struct {
	DB       *gorm.DB
	Redis 	*redis.Client
}

type UserInfo struct {
    Email string
    Name  string
}
func (tc *TransactionController) List(c *fiber.Ctx) error {
	var transactions []model.Transaction

	userID := c.Locals("user_id").(uint)
	role := c.Locals("user_role").(string)

	if role == "admin" {
		if err := tc.DB.Find(&transactions).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
	} else {
		if err := tc.DB.Where("owner_id = ?", userID).Find(&transactions).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
	}

	// get email from grpc
	userClient := grpc_client.NewUserClient()
	categoryClient := grpc_client.NewCategoryClient()

	var response []map[string]interface{}

	for _, t := range transactions {
		userInfo, err := userClient.GetUserInfo(t.OwnerID)
		if err != nil {
			log.Printf("failed to get user info for ID %d: %v", t.OwnerID, err)
			continue
		}

		categoryInfo, err := categoryClient.GetCategoryInfo(t.CategoryID)
		if err != nil {
			log.Printf("failed to get category info for ID %d: %v", t.CategoryID, err)
			continue
		}

		response = append(response, map[string]interface{}{
			"id":           t.ID,
			"name":         t.Name,
			"desc":         t.Desc,
			"category_id":  t.CategoryID,
			"category_name": categoryInfo.Name,
			"category_type": categoryInfo.Type,
			"created_at":   t.CreatedAt,
			"amount":       t.Amount,
			"owner_email":  userInfo.Email,
			"owner_name":   userInfo.Name,
		})
	}

	return c.JSON(response)
	}

func (tc *TransactionController) ListAll(c *fiber.Ctx) error {
	var transactions []model.Transaction
	if err := tc.DB.Find(&transactions).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if len(transactions) == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "transactions not found"})
	} else {
		return c.JSON(transactions)
	}
}

// Create new transaction
func (tc *TransactionController) Create(c *fiber.Ctx) error {
	var transaction model.Transaction
	if err := c.BodyParser(&transaction); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	userID := c.Locals("user_id").(uint)
	
	in := model.Transaction{}
	if err := c.BodyParser(&in); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid payload"})
	}
	in.OwnerID = userID

	transaction.CreatedAt = time.Now()

	if err := tc.DB.Create(&in).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(in)
}

func (tc *TransactionController) CreateFiltered(c *fiber.Ctx) error {
	var in model.Transaction
	if err := c.BodyParser(&in); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid payload"})
	}

	userID := c.Locals("user_id").(uint)
	in.OwnerID = userID
	in.CreatedAt = time.Now()

	categoryClient := grpc_client.NewCategoryClient()
	categoryInfo, err := categoryClient.GetCategoryInfo(in.CategoryID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": fmt.Sprintf("invalid category ID: %d", in.CategoryID)})
	}
	if categoryInfo == nil {
		return c.Status(400).JSON(fiber.Map{"error": "category not found"})
	}

	if uint64(categoryInfo.OwnerID) != uint64(userID) {
		return c.Status(403).JSON(fiber.Map{
			"error": "you do not own this category",
		})
	}

	if err := tc.DB.Create(&in).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(in)
}

// Get transaction by id
 
func (tc *TransactionController) Get(c *fiber.Ctx) error {
	// 1. Ambil dan validasi ID dari params
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}

	// 2. Ambil user ID dari JWT middleware
	userID := c.Locals("user_id").(uint)

	// 3. Ambil transaksi dari DB milik user tersebut
	var transaction model.Transaction
	if err := tc.DB.Where("id = ? AND owner_id = ?", id, userID).First(&transaction).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "transaction not found"})
	}

	// 4. Inisialisasi gRPC client (sebaiknya ini nanti dipindah ke controller struct)
	userClient := grpc_client.NewUserClient()
	categoryClient := grpc_client.NewCategoryClient()

	// 5. Ambil info user dari user-service
	userInfo, err := userClient.GetUserInfo(transaction.OwnerID)
	if err != nil {
		log.Printf("failed to get user info for ID %d: %v", transaction.OwnerID, err)
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch user info"})
	}

	// 6. Ambil info category dari category-service
	categoryInfo, err := categoryClient.GetCategoryInfo(transaction.CategoryID)
	if err != nil {
		log.Printf("failed to get category info for ID %d: %v", transaction.CategoryID, err)
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch category info"})
	}

	// 7. Buat response
	response := map[string]interface{}{
		"id":            transaction.ID,
		"name":          transaction.Name,
		"desc":          transaction.Desc,
		"category_id":   transaction.CategoryID,
		"category_name": categoryInfo.Name,
		"created_at":    transaction.CreatedAt,
		"amount":        transaction.Amount,
		"owner_email":   userInfo.Email,
		"owner_name":    userInfo.Name,
	}

	return c.JSON(response)
}


// Update transaction

type EditRequest struct {
	Name 		string 		`json:"name"`
	Amount    	float64  	`json:"amount"`
	Desc      	string   	`json:"desc"`
}

func (tc *TransactionController) Update(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
    }
    userID := c.Locals("user_id").(uint)

    var transaction model.Transaction
    if err := tc.DB.Where("id = ? AND owner_id = ?", id, userID).First(&transaction).Error; err != nil {
        return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
    }

    var body EditRequest
    if err := c.BodyParser(&body); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": err.Error()})
    }

   
    transaction.Name = body.Name
    transaction.Amount = body.Amount
    transaction.Desc = body.Desc

    if err := tc.DB.Save(&transaction).Error; err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(transaction)
}

// Delete transaction
func (tc *TransactionController) Delete(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
    }

    userID := c.Locals("user_id").(uint)

    var transaction model.Transaction
    // cari transaksi yang sesuai id & owner_id
    if err := tc.DB.Where("id = ? AND owner_id = ?", id, userID).First(&transaction).Error; err != nil {
        return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
    }

    if err := tc.DB.Delete(&transaction).Error; err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(fiber.Map{"message": "transaction deleted"})
}

func (tc *TransactionController) GetBalance(c *fiber.Ctx) error {
	ctx := context.Background()
	userID := c.Locals("user_id").(uint)

	var transactions []model.Transaction
	if err := tc.DB.Where("owner_id = ?", userID).Find(&transactions).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	categoryClient := grpc_client.NewCategoryClient()

	var totalIncome float64
	var totalExpense float64

	for _, t := range transactions {
		cacheKey := fmt.Sprintf("category:%d", t.CategoryID)

		var category *grpc_client.CategoryInfo

		// check redis
		if tc.Redis != nil {
		if cachedData, err := tc.Redis.Get(ctx, cacheKey).Result(); err == nil {
		var cachedCat grpc_client.CategoryInfo
		if json.Unmarshal([]byte(cachedData), &cachedCat) == nil {
			category = &cachedCat
		} else {
			_ = tc.Redis.Del(ctx, cacheKey).Err()
		}
		} else if err != redis.Nil {
			log.Printf("redis get error for key %s: %v", cacheKey, err)
		}
}
		// get category info if nill
		if category == nil {
			cat, err := categoryClient.GetCategoryInfo(t.CategoryID)
			if err != nil {
				// log error dan skip transaksi ini agar tidak crash seluruh endpoint
				log.Printf("failed to fetch category info for ID %d: %v", t.CategoryID, err)
				continue
			}
			if cat == nil {
				// skip category if nill
				log.Printf("category info nil for ID %d", t.CategoryID)
				continue
			}
			category = cat

			// save to redis
			if tc.Redis != nil {
				if jsonCat, err := json.Marshal(category); err == nil {
					if err := tc.Redis.Set(ctx, cacheKey, jsonCat, 10*time.Minute).Err(); err != nil {
						log.Printf("failed to set redis key %s: %v", cacheKey, err)
					}
				}
			}
		}

		if category == nil {
			continue
		}

		if category.Type == "income" {
			totalIncome += float64(t.Amount)
		} else if category.Type == "expense" {
			totalExpense += float64(t.Amount)
		}
	}

	balance := totalIncome - totalExpense

	return c.JSON(fiber.Map{
		"total_income":  totalIncome,
		"total_expense": totalExpense,
		"balance":       balance,
	})
}


func (tc *TransactionController) GetBalanceCategory(c *fiber.Ctx) error {
	var transactions []model.Transaction
	userID := c.Locals("user_id").(uint)

	if err := tc.DB.Where("owner_id = ?", userID).Find(&transactions).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	categoryClient := grpc_client.NewCategoryClient()
	categoryCache := make(map[uint]*grpc_client.CategoryInfo)

	type CategorySummary struct {
		CategoryID   uint    `json:"category_id"`
		CategoryName string  `json:"category_name"`
		Type         string  `json:"type"`
		Total        float64 `json:"total"`
	}

	categoryTotals := make(map[uint]*CategorySummary)

	for _, t := range transactions {
		category, ok := categoryCache[t.CategoryID]
		if !ok {
			cat, err := categoryClient.GetCategoryInfo(t.CategoryID)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{
					"error": fmt.Sprintf("failed to fetch category info for ID %d", t.CategoryID),
				})
			}
			category = cat
			categoryCache[t.CategoryID] = category
		}

		if _, exists := categoryTotals[t.CategoryID]; !exists {
			categoryTotals[t.CategoryID] = &CategorySummary{
				CategoryID:   t.CategoryID,
				CategoryName: category.Name,
				Type:         category.Type,
				Total:        0,
			}
		}

		categoryTotals[t.CategoryID].Total += float64(t.Amount)
	}

	result := make([]CategorySummary, 0, len(categoryTotals))
	for _, summary := range categoryTotals {
		result = append(result, *summary)
	}

	return c.JSON(fiber.Map{
		"categories": result,
	})
}

func (tc *TransactionController) GetBudgetStatus(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var transactions []model.Transaction
	if err := tc.DB.Where("owner_id = ?", userID).Find(&transactions).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	categoryClient := grpc_client.NewCategoryClient()

	categoryCache := make(map[uint]*grpc_client.CategoryInfo)

	budgetUsage := make(map[uint]float64)

	for _, t := range transactions {
		var category *grpc_client.CategoryInfo
		var ok bool
		if category, ok = categoryCache[t.CategoryID]; !ok {
			cat, err := categoryClient.GetCategoryInfo(t.CategoryID)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{
					"error": fmt.Sprintf("failed to fetch category info for ID %d", t.CategoryID),
				})
			}
			category = cat
			categoryCache[t.CategoryID] = category
		}

		if category.Type == "expense" {
			budgetUsage[t.CategoryID] += float64(t.Amount)
		}
	}

	results := []fiber.Map{}
	for categoryID, cat := range categoryCache {
		if cat.Type == "expense" {
			used := budgetUsage[categoryID]
			remaining := cat.Budget - used
			status := "within budget"
			if remaining < 0 {
				status = "over budget"
			}

			results = append(results, fiber.Map{
				"category_id": categoryID,
				"name":        cat.Name,
				"budget":      cat.Budget,
				"used":        used,
				"remaining":   remaining,
				"status":      status,
			})
		}
	}

	return c.JSON(results)
}
