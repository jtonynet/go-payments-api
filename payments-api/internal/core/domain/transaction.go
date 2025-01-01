package domain

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Transaction struct {
	UID          uuid.UUID
	AccountID    uint
	AccountUID   uuid.UUID
	CategoryID   uint
	MCC          string
	Amount       decimal.Decimal
	MerchantName string
}
