package domain

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Merchant struct {
	Name string
	MCC  string
}

func (m *Merchant) NewTransaction(
	transactionUID uuid.UUID,
	mcc string,
	totalAmount decimal.Decimal,
	name string,

	account Account,
) Transaction {
	correctMCC := mcc
	if m.MCC != "" {
		correctMCC = m.MCC
	}

	return Transaction{
		UID:          transactionUID,
		AccountID:    account.ID,
		AccountUID:   account.UID,
		MCC:          correctMCC,
		Amount:       totalAmount,
		MerchantName: name,
	}
}
