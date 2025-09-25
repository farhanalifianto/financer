package model

import "time"

type Transaction struct {
	ID        	uint      `gorm:"primaryKey" json:"id"`
	Name      	string    `json:"name"`
	Desc      	string    `json:"desc"`
	OwnerID  	uint     `json:"owner_id"`
	CategoryID	string    `json:"category"`
	CreatedAt 	time.Time `json:"created_at"`
	Amount    	float64   `json:"amount"`
	Type	  	string    `json:"type"`
}
