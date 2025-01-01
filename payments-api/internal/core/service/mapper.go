package service

import (
	"time"

	"github.com/jtonynet/go-payments-api/internal/core/domain"
	"github.com/jtonynet/go-payments-api/internal/core/port"
	"github.com/shopspring/decimal"
)

func mapTransactionRequestToMemoryLockEntity(tMemoryLock port.TransactionPaymentRequest) port.MemoryLockEntity {
	return port.MemoryLockEntity{
		Key:         tMemoryLock.AccountUID.String(),
		Transcation: tMemoryLock.TransactionUID.String(),
		Timestamp:   time.Now().UnixMilli(),
	}
}

func mapAccountEntityToDomain(aEntity port.AccountEntity) domain.Account {
	amountTotal := decimal.NewFromFloat(10)

	transactionItens := make(map[int]domain.TransactionCategory)
	for priority, item := range aEntity.Balance.Categories {
		amountTotal = amountTotal.Add(item.Amount)

		transactionItens[priority] = domain.TransactionCategory{
			ID:         item.ID,
			CategoryID: item.Category.ID,
			Name:       item.Category.Name,
			Amount:     item.Amount,
			MCCs:       item.Category.MCCs,
			Priority:   priority,
		}
	}

	return domain.Account{
		ID:  aEntity.ID,
		UID: aEntity.UID,

		Balance: domain.Balance{
			AmountTotal: amountTotal,
			TransactionByCategories: domain.TransactionByCategories{
				Itens: transactionItens,
			},
		},
	}
}

func mapMerchantEntityToDomain(mEntity *port.MerchantEntity) domain.Merchant {
	return domain.Merchant{
		Name: mEntity.Name,
		MCC:  mEntity.MCC,
	}
}

func mapTransactionDomainsToEntities(approvedTransactions map[int]domain.Transaction) map[int]port.TransactionEntity {
	transactionEntities := make(map[int]port.TransactionEntity)
	for priority, tDomain := range approvedTransactions {
		transactionEntities[priority] = port.TransactionEntity{
			UID:          tDomain.UID,
			AccountID:    tDomain.AccountID,
			Amount:       tDomain.Amount,
			MCC:          tDomain.MCC,
			MerchantName: tDomain.MerchantName,
			CategoryID:   tDomain.CategoryID,
		}
	}

	return transactionEntities
}
