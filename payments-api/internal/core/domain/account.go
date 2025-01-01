package domain

import (
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

func (a *Account) ApproveTransaction(tDomain Transaction) (map[int]Transaction, *CustomError) {
	transactions := make(map[int]Transaction)

	amountDebtRemaining := tDomain.Amount

	categoryMCC, err := a.Balance.TransactionByCategories.GetByMCC(tDomain.MCC)
	if err != nil {
		a.debugLog(
			fmt.Sprintf(
				"Error retrieving category by MCC, attempting to debit remaining amount (%s) directly from fallback.",
				amountDebtRemaining,
			),
		)

	} else if categoryMCC.Amount.GreaterThanOrEqual(tDomain.Amount) {
		a.debugLog(
			fmt.Sprintf(
				"Sufficient funds available in category '%s' to cover the transaction amount.",
				categoryMCC.Name,
			),
		)

		amountDebtRemaining = decimal.NewFromFloat(0)

		categoryMCC.Amount = categoryMCC.Amount.Sub(tDomain.Amount)
		transactions[categoryMCC.Priority] = a.mapCategoryToTransaction(categoryMCC, tDomain)

	} else if categoryMCC.Amount.IsPositive() {
		a.debugLog(
			fmt.Sprintf(
				"Category '%s' has a positive balance but insufficient funds for the transaction. Check fallback category.",
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
			a.debugLog(
				"Error retrieving fallback category. Rolling back changes and returning an error.",
			)

			return transactions, NewCustomError(CODE_REJECTED_GENERIC, "Category Fallback not found")

		} else if CategoryFallback.Amount.GreaterThanOrEqual(amountDebtRemaining) {
			a.debugLog(
				"Fallback category has sufficient funds available for the transaction.",
			)

			CategoryFallback.Amount = CategoryFallback.Amount.Sub(amountDebtRemaining)
			transactions[CategoryFallback.Priority] = a.mapCategoryToTransaction(CategoryFallback, tDomain)

			amountDebtRemaining = decimal.NewFromFloat(0)
		} else {
			a.debugLog(
				"Insufficient funds in the fallback category. Rolling back changes and returning an error.",
			)

			amountDebtRemaining = tDomain.Amount
			transactions = make(map[int]Transaction)
		}
	}

	if amountDebtRemaining.GreaterThan(decimal.Zero) || len(transactions) == 0 {
		return transactions, NewCustomError(CODE_REJECTED_INSUFICIENT_FUNDS, "transaction category has insuficient funds")
	}

	return transactions, nil
}

func (a *Account) debugLog(msg string) {
	if a.Logger != nil {
		a.Logger.Debug(msg)
	}
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
