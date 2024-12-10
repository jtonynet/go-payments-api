package gormRepos

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jtonynet/go-payments-api/internal/adapter/database"
	"github.com/jtonynet/go-payments-api/internal/adapter/model/gormModel"
	"github.com/jtonynet/go-payments-api/internal/core/port"
	"github.com/shopspring/decimal"

	"gorm.io/gorm"
)

type Account struct {
	gormConn database.Conn
	db       *gorm.DB
}

func NewGormAccount(conn database.Conn) (port.AccountRepository, error) {
	db, err := conn.GetDB(context.Background())
	if err != nil {
		return nil, fmt.Errorf("account repository failure on conn.GetDB()")
	}

	dbGorm, ok := db.(*gorm.DB)
	if !ok {
		return nil, fmt.Errorf("account repository failure to cast conn.GetDB() as gorm.DB")
	}

	return &Account{
		gormConn: conn,
		db:       dbGorm,
	}, nil
}

type accountResult struct {
	AccountID      uint
	AccountUID     uuid.UUID
	TransactionID  uint
	TransactionUID uuid.UUID
	Amount         decimal.Decimal
	CategoryID     uint
	CategoryName   string
	Priority       int
	Codes          sql.NullString
}

func (a *Account) FindByUID(ctx context.Context, uid uuid.UUID) (port.AccountEntity, error) {
	var account port.AccountEntity
	var balance port.BalanceEntity
	var results []accountResult

	firstFound := false
	transactionsByCategories := make(map[int]port.TransactionByCategoryEntity)

	/*
		https://gorm.io/docs/context.html#Context-Timeout
	*/
	err := a.db.WithContext(ctx).
		Table("accounts as a").
		Select(`
			a.id as account_id, 
			t.id as transaction_id, 
			t.amount as amount, 
			c.id as category_id, 
			c.name as category_name, 
			c.priority as priority,
			STRING_AGG(mc.mcc, ',') AS codes
		`).
		Joins("JOIN account_categories as ac ON ac.account_id = a.id").
		Joins("JOIN categories as c ON c.id = ac.category_id").
		Joins(`
        	JOIN transactions as t ON t.account_id = a.id AND 
				t.category_id = c.id AND 
				t.id = (
        	    	SELECT MAX(t2.id) 
        	    	FROM transactions t2 
        	    	WHERE t2.account_id = a.id AND t2.category_id = c.id
        		)
    	`).
		Joins("LEFT JOIN mccs as mc ON mc.category_id = c.id").
		Where("a.uid = ?", uid).
		Where(`
			a.deleted_at IS NULL
			AND ac.deleted_at IS NULL
			AND c.deleted_at IS NULL
		`).
		Group("a.id, a.uid, t.id, t.uid, t.amount, c.id, c.name, c.priority").
		Scan(&results).Error

	if err != nil {
		return account, fmt.Errorf("error retrying account:%s  err: %w", uid, err)
	}

	if len(results) > 0 {
		amountTotal := decimal.NewFromInt(0)

		for _, result := range results {
			mccs := []string{}
			if result.Codes.Valid {
				mccs = strings.Split(result.Codes.String, ",")
			}

			amountTotal = amountTotal.Add(result.Amount)

			if !firstFound {
				firstFound = true
				account.ID = result.AccountID
				account.UID = uid
			}

			transactionsByCategories[int(result.TransactionID)] = port.TransactionByCategoryEntity{
				ID:     result.TransactionID,
				Amount: result.Amount,
				Category: port.CategoryEntity{
					ID:       result.CategoryID,
					Name:     result.CategoryName,
					Priority: result.Priority,
					MCCs:     mccs,
				},
			}
		}

		balance.AmountTotal = amountTotal
		balance.Categories = transactionsByCategories
	}

	account.Balance = balance

	return account, nil
}

func (a *Account) SaveTransactions(ctx context.Context, transactions map[int]port.TransactionEntity) error {
	if len(transactions) == 0 {
		return fmt.Errorf("no transactions to save")
	}

	var tSlice []gormModel.Transaction

	for _, transaction := range transactions {
		tSlice = append(tSlice, gormModel.Transaction{
			UID:          uuid.New(),
			AccountID:    transaction.AccountID,
			CategoryID:   transaction.CategoryID,
			Amount:       transaction.Amount,
			MCC:          transaction.MCC,
			MerchantName: transaction.MerchantName,
		})
	}

	err := a.db.WithContext(ctx).Create(&tSlice).Error
	if err != nil {
		return fmt.Errorf("failed to save transactions: %w", err)
	}

	return nil
}
