package controller

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"transaction-service/grpc_client"

	"transaction-service/model"

	"github.com/gofiber/fiber/v2"

	"gorm.io/gorm"
)

// List all transactions

type TransactionController struct {
	DB       *gorm.DB
}

// func (tc *TransactionController) List(c *fiber.Ctx) error {
// 	var transactions []model.Transaction

// 	userID := c.Locals("user_id").(uint)
// 	role := c.Locals("user_role").(string)
// 	if role == "admin" {
// 		if err := tc.DB.Find(&transactions).Error; err != nil {
// 			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
// 		}
// 		return c.JSON(transactions)
// 	}

// 	if err := tc.DB.Where("owner_id = ?", userID).Find(&transactions).Error; err != nil {
// 		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
// 	}

// 	return c.JSON(transactions)
// }
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
		// Ambil info user
		userInfo, err := userClient.GetUserInfo(t.OwnerID)
		if err != nil {
			log.Printf("failed to get user info for ID %d: %v", t.OwnerID, err)
			continue
		}

		// Ambil info kategori
		categoryInfo, err := categoryClient.GetCategoryInfo(t.CategoryID)
		if err != nil {
			log.Printf("failed to get category info for ID %d: %v", t.CategoryID, err)
			continue
		}

		// Bentuk response
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

// Get transaction by ID
func (tc *TransactionController) Get(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	userID := c.Locals("user_id").(uint)
	var transaction []model.Transaction
	if err := tc.DB.Where("owner_id = ?", userID).First(&transaction, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	if err := tc.DB.First(&transaction, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	userClient := grpc_client.NewUserClient()
	categoryClient := grpc_client.NewCategoryClient()

	var response []map[string]interface{}

	for _, t := range transaction {
		// Ambil info user
		userInfo, err := userClient.GetUserInfo(t.OwnerID)
		if err != nil {
			log.Printf("failed to get user info for ID %d: %v", t.OwnerID, err)
			continue
		}

		// Ambil info kategori
		categoryInfo, err := categoryClient.GetCategoryInfo(t.CategoryID)
		if err != nil {
			log.Printf("failed to get category info for ID %d: %v", t.CategoryID, err)
			continue
		}

		// Bentuk response
		response = append(response, map[string]interface{}{
			"id":           t.ID,
			"name":         t.Name,
			"desc":         t.Desc,
			"category_id":  t.CategoryID,
			"category_name": categoryInfo.Name,
			"created_at":   t.CreatedAt,
			"amount":       t.Amount,
			"owner_email":  userInfo.Email,
			"owner_name":   userInfo.Name,
		})
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

func (tc *TransactionController) GetBalancerr(c *fiber.Ctx) error {
	// Ambil semua transaksi dari DB
	var transactions []model.Transaction

	userID := c.Locals("user_id").(uint)
	
	if err := tc.DB.Where("owner_id = ?", userID).Find(&transactions).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	categoryClient := grpc_client.NewCategoryClient() // client gRPC ke category service

	var totalIncome float64
	var totalExpense float64

	for _, t := range transactions {
		// Ambil detail kategori lewat gRPC
		category, err := categoryClient.GetCategoryInfo(t.CategoryID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "failed to fetch category info",
			})
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
func (tc *TransactionController) GetBalance(c *fiber.Ctx) error {
	var transactions []model.Transaction
	userID := c.Locals("user_id").(uint)

	if err := tc.DB.Where("owner_id = ?", userID).Find(&transactions).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	categoryClient := grpc_client.NewCategoryClient()

	var totalIncome float64
	var totalExpense float64

	// 🧠 Cache lokal kategori supaya tidak gRPC berkali-kali
	categoryCache := make(map[uint]*grpc_client.CategoryInfo)

	for _, t := range transactions {
		var category *grpc_client.CategoryInfo
		var ok bool

		if category, ok = categoryCache[t.CategoryID]; !ok {
			// Kalau belum ada di cache → ambil lewat gRPC
			cat, err := categoryClient.GetCategoryInfo(t.CategoryID)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{
					"error": fmt.Sprintf("failed to fetch category info for ID %d", t.CategoryID),
				})
			}
			category = cat
			categoryCache[t.CategoryID] = category
		}

		// Hitung total berdasarkan tipe
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

