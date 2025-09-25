package controller

import (
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
	var response []map[string]interface{}
	for _, t := range transactions {
		UserInfo, _ := userClient.GetUserEmail(t.OwnerID)
		response = append(response, map[string]interface{}{
			"id":         t.ID,
			"name":       t.Name,
			"desc":       t.Desc,
			"category":   t.CategoryID,
			"created_at": t.CreatedAt,
			"amount":     t.Amount,
			"type":       t.Type,
			"owner":      UserInfo.Email,
			"owner_name": UserInfo.Name,
			 
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
	var transaction model.Transaction
	if err := tc.DB.Where("owner_id = ?", userID).First(&transaction, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	if err := tc.DB.First(&transaction, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(transaction)
}

// Update transaction

type EditRequest struct {
	Name 		string 		`json:"name"`
	Amount    	float64  	`json:"amount"`
	Desc      	string   	`json:"desc"`
	Type	  	string   	`json:"type"`
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
    transaction.Type = body.Type

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