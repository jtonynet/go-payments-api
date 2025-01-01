package port

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type TransactionPaymentRequest struct {
	AccountUID     uuid.UUID       `json:"account" validate:"required,uuid" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000"`
	TransactionUID uuid.UUID       `json:"-" swaggerignore:"true"`
	TotalAmount    decimal.Decimal `json:"totalAmount" validate:"required,min=0.01" binding:"required" example:"100.09"`
	MCC            string          `json:"mcc" validate:"required,min=4,max=4" binding:"required" example:"5411"`
	Merchant       string          `json:"merchant" validate:"required,min=3,max=255" binding:"required" example:"PADARIA DO ZE              SAO PAULO BR"`
}

type TransactionPaymentResponse struct {
	Code string `json:"code" example:"00"`
}

type TransactionEntity struct {
	ID           uint
	UID          uuid.UUID
	AccountID    uint
	CategoryID   uint
	Amount       decimal.Decimal
	MCC          string
	MerchantName string
}

type TransactionByCategoryEntity struct {
	ID       uint
	UID      uuid.UUID
	Amount   decimal.Decimal
	Category CategoryEntity
}
