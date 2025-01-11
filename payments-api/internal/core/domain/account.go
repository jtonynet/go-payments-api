package domain

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jtonynet/go-payments-api/internal/support/logger"
	"github.com/shopspring/decimal"
)

type Account struct {
	ID  uint
	UID uuid.UUID

	Balance

	Logger logger.Logger
}

func (a *Account) ApproveTransaction(ctx context.Context, tDomain Transaction) (map[int]Transaction, *CustomError) {
	transactions := make(map[int]Transaction)

	amountDebtRemaining := tDomain.Amount

	categoryMCC, err := a.Balance.TransactionByCategories.GetByMCC(tDomain.MCC)
	if err != nil {
		a.Logger.Debug(
			ctx,
			"Error retrieving category by MCC, attempting fallback.",
		)

	} else if categoryMCC.Amount.GreaterThanOrEqual(tDomain.Amount) {
		a.Logger.Debug(
			ctx,
			fmt.Sprintf(
				"Sufficient funds available in category '%s'",
				categoryMCC.Name,
			),
		)

		amountDebtRemaining = decimal.NewFromFloat(0)

		categoryMCC.Amount = categoryMCC.Amount.Sub(tDomain.Amount)
		transactions[categoryMCC.Priority] = a.mapCategoryToTransaction(categoryMCC, tDomain)

	} else if categoryMCC.Amount.IsPositive() {
		a.Logger.Debug(
			ctx,
			fmt.Sprintf(
				"Category '%s' has a positive balance but insufficient funds, attempting from fallback.",
				categoryMCC.Name,
			),
		)

		amountDebtRemaining = amountDebtRemaining.Sub(categoryMCC.Amount)

		categoryMCC.Amount = decimal.NewFromFloat(0)
		transactions[categoryMCC.Priority] = a.mapCategoryToTransaction(categoryMCC, tDomain)
	}

	if amountDebtRemaining.GreaterThan(decimal.Zero) {
		CategoryFallback, err := a.Balance.TransactionByCategories.GetFallback()
		if err != nil {
			a.Logger.Debug(
				ctx,
				"Error retrieving fallback category.",
			)

			fallbackNotFoundErr := fmt.Sprintf("Category Fallback not found for :%s", tDomain.AccountUID.String())
			return transactions, NewCustomError(CODE_REJECTED_INSUFICIENT_FUNDS, fallbackNotFoundErr)

		} else if CategoryFallback.Amount.GreaterThanOrEqual(amountDebtRemaining) {
			a.Logger.Debug(
				ctx,
				"Fallback category has sufficient funds available.",
			)

			CategoryFallback.Amount = CategoryFallback.Amount.Sub(amountDebtRemaining)
			transactions[CategoryFallback.Priority] = a.mapCategoryToTransaction(CategoryFallback, tDomain)

			amountDebtRemaining = decimal.NewFromFloat(0)
		} else {
			a.Logger.Debug(
				ctx,
				"Insufficient funds in the fallback category.",
			)

			amountDebtRemaining = tDomain.Amount
			transactions = make(map[int]Transaction)
		}
	}

	if amountDebtRemaining.GreaterThan(decimal.Zero) || len(transactions) == 0 {
		return transactions, NewCustomError(CODE_REJECTED_INSUFICIENT_FUNDS, "Insuficient funds for transaction")
	}

	return transactions, nil
}

func (a *Account) mapCategoryToTransaction(tc TransactionCategory, t Transaction) Transaction {
	return Transaction{
		UID:          t.UID,
		AccountID:    a.ID,
		AccountUID:   a.UID,
		CategoryID:   tc.CategoryID,
		Amount:       tc.Amount,
		MCC:          t.MCC,
		MerchantName: t.MerchantName,
	}
}
