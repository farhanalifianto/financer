package model

type TransactionResponse struct {
	ID         uint    `json:"id"`
	Name       string  `json:"name"`
	Desc       string  `json:"desc"`
	Amount     float64 `json:"amount"`
	Category   string  `json:"category"`
	OwnerID    uint    `json:"owner_id"`
	OwnerEmail string  `json:"owner_email"`
	Type       string  `json:"type"` // tambahan dari gRPC User Service
}
