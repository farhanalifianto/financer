package model

import "time"

type Category struct {
	ID        	uint      `gorm:"primaryKey" json:"id"`
	Name      	string    `json:"name"`
	OwnerID  	uint     `json:"owner_id"`
	CreatedAt 	time.Time `json:"created_at"`
	Type		string    `json:"type"` // "expense" or "income"
	Budget		float64   `json:"budget"`
}
