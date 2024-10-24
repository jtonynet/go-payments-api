package service

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jtonynet/go-payments-api/internal/core/domain"
	"github.com/jtonynet/go-payments-api/internal/core/port"
	"github.com/shopspring/decimal"
)

func mapParamsToTransactionDomain(
	accountID uint,
	accountUID uuid.UUID,
	MccCode string,
	totalAmount decimal.Decimal,
	merchant string) (*domain.Transaction, error) {

	if !totalAmount.GreaterThan(decimal.NewFromInt(0)) {
		return &domain.Transaction{}, fmt.Errorf("transaction totalAmount must be a positive value %v", totalAmount)
	}

	return &domain.Transaction{
		AccountID:   accountID,
		AccountUID:  accountUID,
		MccCode:     MccCode,
		TotalAmount: totalAmount,
		Merchant:    merchant,
	}, nil
}

func mapAccountEntityToDomain(aEntity port.AccountEntity) (domain.Account, error) {
	return domain.Account{
		ID:  aEntity.ID,
		UID: aEntity.UID,
	}, nil
}

func mapBalanceEntityToDomain(be port.BalanceEntity) (*domain.Balance, error) {
	categoryItens := make(map[int]domain.Category)
	for _, ce := range be.Categories {
		category := domain.Category{
			ID:       ce.ID,
			Name:     ce.Category.Name,
			Amount:   ce.Amount,
			MccCodes: ce.Category.MccCodes,
			Order:    ce.Category.Order,
		}
		categoryItens[ce.Category.Order] = category
	}

	if len(categoryItens) == 0 {
		return &domain.Balance{}, fmt.Errorf("failed to map Balance with AccountID %v, Cotegories not found", be.AccountID)
	}

	categories := domain.Categories{Itens: categoryItens}

	b := &domain.Balance{
		AccountID:   be.AccountID,
		AmountTotal: be.AmountTotal,
		Categories:  categories,
	}

	return b, nil
}

func mapBalanceDomainToEntity(dBalance *domain.Balance) port.BalanceEntity {
	bCategories := make(map[int]port.BalanceByCategoryEntity)
	for _, categoryItem := range dBalance.Categories.Itens {
		bCategory := port.BalanceByCategoryEntity{
			ID:        categoryItem.ID,
			AccountID: dBalance.AccountID,
			Amount:    categoryItem.Amount,
			Category:  port.Categories[categoryItem.Name],
		}

		bCategories[categoryItem.Order] = bCategory
	}

	return port.BalanceEntity{
		AccountID:   dBalance.AccountID,
		AmountTotal: dBalance.AmountTotal,
		Categories:  bCategories,
	}
}

func mapTransactionDomainToEntity(tDomain *domain.Transaction) port.TransactionEntity {
	return port.TransactionEntity{
		AccountID:   tDomain.AccountID,
		TotalAmount: tDomain.TotalAmount,
		MccCode:     tDomain.MccCode,
		Merchant:    tDomain.Merchant,
	}
}
