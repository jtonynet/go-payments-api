package service

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"gopkg.in/go-playground/assert.v1"

	"github.com/jtonynet/go-payments-api/internal/adapter/repository"
	"github.com/jtonynet/go-payments-api/internal/core/port"
)

var (
	timeoutSLAcfg = 100

	accountUIDtoTransact, _ = uuid.Parse("123e4567-e89b-12d3-a456-426614174000")

	balanceFoodAmount = decimal.NewFromFloat(205.11)
	foodCategoryID    = uint(1)

	balanceMealAmount = decimal.NewFromFloat(310.22)
	mealCategoryID    = uint(2)

	balanceCashAmount = decimal.NewFromFloat(415.33)
	cashCategoryID    = uint(3)

	correctFoodMCC = "5411"

	amountFoodFundsApproved         = decimal.NewFromFloat(100.10)
	amountFoodFundsRejected         = decimal.NewFromFloat(720.45)
	amountFoodFundsFallbackApproved = decimal.NewFromFloat(320.00)
)

type DBfake struct {
	Accounts     map[uint]port.AccountEntity
	Transactions map[uint]port.TransactionEntity
	Merchants    map[uint]port.MerchantEntity
}

func newDBfake() DBfake {
	db := DBfake{}

	db.Transactions = make(map[uint]port.TransactionEntity)

	categories := make(map[int]port.TransactionByCategoryEntity)
	foodCategoryUID, _ := uuid.Parse("32e04519-a979-4de2-a20e-77e8342d718f")
	categories[1] = port.TransactionByCategoryEntity{
		ID:       1,
		UID:      foodCategoryUID,
		Amount:   balanceFoodAmount,
		Category: port.CategoryEntity{ID: foodCategoryID, Name: "FOOD", MCCs: []string{"5411", "5412"}, Priority: 1},
	}

	mealCategoryUID, _ := uuid.Parse("e5ce3deb-7dea-4382-a1fd-1428c9888bdc")
	categories[2] = port.TransactionByCategoryEntity{
		ID:       2,
		UID:      mealCategoryUID,
		Amount:   balanceFoodAmount,
		Category: port.CategoryEntity{ID: mealCategoryID, Name: "MEAL", MCCs: []string{"5811", "5812"}, Priority: 2},
	}

	cashCategoryUID, _ := uuid.Parse("4cfdc9f0-a9d8-409d-ba8e-58a36e126ec1")
	categories[3] = port.TransactionByCategoryEntity{
		ID:       3,
		UID:      cashCategoryUID,
		Amount:   balanceFoodAmount,
		Category: port.CategoryEntity{ID: cashCategoryID, Name: "CASH", Priority: 3},
	}

	db.Accounts = map[uint]port.AccountEntity{
		1: {
			ID:  1,
			UID: accountUIDtoTransact,
			Balance: port.BalanceEntity{
				AmountTotal: decimal.NewFromFloat(930.66),
				Categories:  categories,
			},
		},
	}

	db.Merchants = map[uint]port.MerchantEntity{
		1: {
			Name: "UBER EATS                   SAO PAULO BR",
			MCC:  "5412",
		},
	}

	return db
}

func (dbf *DBfake) GetDB() *DBfake {
	return dbf
}

func (dbf *DBfake) AccountRepoFindByUID(_ context.Context, uid uuid.UUID) (port.AccountEntity, error) {
	for _, ae := range dbf.Accounts {
		if ae.UID == uid {
			return ae, nil
		}
	}

	return port.AccountEntity{}, fmt.Errorf("account with AccountUID %s not found", uid.String())
}

type AccountRepoFake struct {
	db DBfake
}

func newAccountRepoFake(db DBfake) port.AccountRepository {
	return &AccountRepoFake{
		db,
	}
}

func (arf *AccountRepoFake) FindByUID(_ context.Context, uid uuid.UUID) (port.AccountEntity, error) {
	accountEntity, err := arf.db.AccountRepoFindByUID(context.Background(), uid)
	return accountEntity, err
}

func (arf *AccountRepoFake) SaveTransactions(_ context.Context, transactions map[int]port.TransactionEntity) error {
	maxID := uint(1)

	for _, t := range transactions {
		arf.db.Transactions[maxID] = port.TransactionEntity{
			ID:           maxID,
			UID:          uuid.New(),
			AccountID:    t.AccountID,
			Amount:       t.Amount,
			MCC:          t.MCC,
			MerchantName: t.MerchantName,
			CategoryID:   t.CategoryID,
		}

		maxID = maxID + 1
	}

	return nil
}

type MerchantRepoFake struct {
	db DBfake
}

func newMerchantRepoFake(db DBfake) port.MerchantRepository {
	return &MerchantRepoFake{
		db,
	}
}

func (m *MerchantRepoFake) FindByName(_ context.Context, name string) (*port.MerchantEntity, error) {
	MerchantEntity, err := m.db.MerchantRepoFindByName(name)
	return MerchantEntity, err
}

func (dbf *DBfake) MerchantRepoFindByName(Name string) (*port.MerchantEntity, error) {
	for _, m := range dbf.Merchants {
		if m.Name == Name {
			return &m, nil
		}
	}

	return nil, nil
}

type InMemoryDBfake struct {
	Lock map[string]string
}

func newInMemoryDBfake() InMemoryDBfake {
	imdbf := InMemoryDBfake{}
	imdbf.Lock = make(map[string]string)

	return imdbf
}

func (suite *PaymentSuite) getInMemoryDBfake() InMemoryDBfake {
	inMemorydbFake := newInMemoryDBfake()
	return inMemorydbFake
}

func (imdbf *InMemoryDBfake) MemoryLockRepoLock(_ context.Context, mle port.MemoryLockEntity) (port.MemoryLockEntity, error) {
	imdbf.Lock[mle.Key] = strconv.FormatInt(mle.Timestamp, 10)
	return mle, nil
}

func (imdbf *InMemoryDBfake) MemoryLockRepoUnlock(_ context.Context, key string) error {
	if _, exists := imdbf.Lock[key]; exists {
		delete(imdbf.Lock, key)
		return nil
	}
	return nil
}

type MemoryLockRepoFake struct {
	memoryDB InMemoryDBfake
}

func newMemoryLockRepoFake(memoryDB InMemoryDBfake) port.MemoryLockRepository {
	return &MemoryLockRepoFake{
		memoryDB,
	}
}

func (m *MemoryLockRepoFake) Lock(_ context.Context, timeoutSLA port.TimeoutSLA, mle port.MemoryLockEntity) (port.MemoryLockEntity, error) {
	return m.memoryDB.MemoryLockRepoLock(context.Background(), mle)
}

func (m *MemoryLockRepoFake) Unlock(_ context.Context, key string) error {
	return m.memoryDB.MemoryLockRepoUnlock(context.Background(), key)
}

type PaymentSuite struct {
	suite.Suite
}

func (suite *PaymentSuite) getDBfake() *DBfake {
	dbFake := newDBfake()
	return &dbFake
}

func (suite *PaymentSuite) getMemoryLockRepoFake(memoryDB InMemoryDBfake) port.MemoryLockRepository {
	return newMemoryLockRepoFake(memoryDB)
}

func (suite *PaymentSuite) getAllRepositories(dbFake *DBfake) *repository.AllRepos {
	allRepos := repository.AllRepos{}
	allRepos.Account = newAccountRepoFake(*dbFake)
	allRepos.Merchant = newMerchantRepoFake(*dbFake)

	return &allRepos
}

func (suite *PaymentSuite) TestL1PaymentExecuteGenericRejected() {
	//Arrange
	timeoutSLA := port.TimeoutSLA(
		time.Duration(timeoutSLAcfg) * time.Millisecond,
	)

	dbFake := DBfake{}
	allRepos := suite.getAllRepositories(&dbFake)

	inMemoryDBfake := suite.getInMemoryDBfake()
	memoryLockRepo := suite.getMemoryLockRepoFake(inMemoryDBfake)

	tRequest := port.TransactionPaymentRequest{
		AccountUID:  accountUIDtoTransact,
		TotalAmount: amountFoodFundsApproved,
		MCC:         correctFoodMCC,
		Merchant:    "PADARIA DO ZE               SAO PAULO BR",
	}

	//Act
	paymentService := NewPayment(
		timeoutSLA,
		allRepos.Account,
		allRepos.Merchant,
		memoryLockRepo,
		nil,
	)

	returnCode, _ := paymentService.Execute(tRequest)

	//Assert
	codeRejected := "07" // domain.CODE_REJECTED_GENERIC
	assert.Equal(suite.T(), returnCode, codeRejected)
}

func (suite *PaymentSuite) TestL1PaymentExecuteCorrectMCCWithFundsRejected() {
	//Arrange
	timeoutSLA := port.TimeoutSLA(
		time.Duration(timeoutSLAcfg) * time.Millisecond,
	)

	dbFake := suite.getDBfake()
	allRepos := suite.getAllRepositories(dbFake)

	inMemoryDBfake := suite.getInMemoryDBfake()
	memoryLockRepo := suite.getMemoryLockRepoFake(inMemoryDBfake)

	tRequest := port.TransactionPaymentRequest{
		AccountUID:  accountUIDtoTransact,
		TotalAmount: amountFoodFundsRejected,
		MCC:         correctFoodMCC,
		Merchant:    "PADARIA DO ZE               SAO PAULO BR",
	}

	//Act
	paymentService := NewPayment(
		timeoutSLA,
		allRepos.Account,
		allRepos.Merchant,
		memoryLockRepo,
		nil,
	)

	returnCode, err := paymentService.Execute(tRequest)

	//Assert
	codeRejected := "51" // domain.CODE_REJECTED_INSUFICIENT_FUNDS
	insuficientFundsError := "failed to approve transaction: transaction category has insuficient funds"

	assert.Equal(suite.T(), returnCode, codeRejected)
	assert.Equal(suite.T(), err.Error(), insuficientFundsError)
}

func (suite *PaymentSuite) TestL1PaymentExecuteCorrectMCCWithFundsApproved() {
	//Arrange
	timeoutSLA := port.TimeoutSLA(
		time.Duration(timeoutSLAcfg) * time.Millisecond,
	)

	dbFake := suite.getDBfake()
	allRepos := suite.getAllRepositories(dbFake)

	inMemoryDBfake := suite.getInMemoryDBfake()
	memoryLockRepo := suite.getMemoryLockRepoFake(inMemoryDBfake)

	tRequest := port.TransactionPaymentRequest{
		AccountUID:  accountUIDtoTransact,
		TotalAmount: amountFoodFundsApproved,
		MCC:         correctFoodMCC,
		Merchant:    "PADARIA DO ZE               SAO PAULO BR",
	}

	//Act
	paymentService := NewPayment(
		timeoutSLA,
		allRepos.Account,
		allRepos.Merchant,
		memoryLockRepo,
		nil,
	)
	returnCode, err := paymentService.Execute(tRequest)

	//Assert
	codeApproved := "00" // domain.CODE_APPROVED

	assert.Equal(suite.T(), returnCode, codeApproved)
	assert.Equal(suite.T(), err, nil)

	foodTransaction, err := getLastTransaction(dbFake.Transactions, port.TransactionEntity{AccountID: 1, CategoryID: foodCategoryID})
	assert.Equal(suite.T(), foodTransaction.Amount, decimal.NewFromFloat(105.01))
	assert.Equal(suite.T(), err, nil)
}

func (suite *PaymentSuite) TestL2PaymentExecuteCorrectMCCFallbackApproved() {
	//Arrange
	timeoutSLA := port.TimeoutSLA(
		time.Duration(timeoutSLAcfg) * time.Millisecond,
	)

	dbFake := suite.getDBfake()
	allRepos := suite.getAllRepositories(dbFake)

	inMemoryDBfake := suite.getInMemoryDBfake()
	memoryLockRepo := suite.getMemoryLockRepoFake(inMemoryDBfake)

	tRequest := port.TransactionPaymentRequest{
		AccountUID:  accountUIDtoTransact,
		TotalAmount: amountFoodFundsFallbackApproved,
		MCC:         correctFoodMCC,
		Merchant:    "PADARIA DO ZE               SAO PAULO BR",
	}

	//Act
	paymentService := NewPayment(
		timeoutSLA,
		allRepos.Account,
		allRepos.Merchant,
		memoryLockRepo,
		nil,
	)
	returnCode, err := paymentService.Execute(tRequest)

	//Assert
	codeApproved := "00" // domain.CODE_APPROVED

	assert.Equal(suite.T(), returnCode, codeApproved)
	assert.Equal(suite.T(), err, nil)

	foodTransaction, _ := getLastTransaction(dbFake.Transactions, port.TransactionEntity{AccountID: 1, CategoryID: foodCategoryID})
	assert.Equal(suite.T(), err, nil)
	assert.Equal(suite.T(), foodTransaction.Amount, decimal.NewFromFloat(0))

	cashTransaction, _ := getLastTransaction(dbFake.Transactions, port.TransactionEntity{AccountID: 1, CategoryID: cashCategoryID})
	assert.Equal(suite.T(), err, nil)
	assert.Equal(suite.T(), cashTransaction.Amount, decimal.NewFromFloat(90.22))
}

func (suite *PaymentSuite) TestL3PaymentExecuteNameMCCWithFundsApproved() {
	//Arrange
	timeoutSLA := port.TimeoutSLA(
		time.Duration(timeoutSLAcfg) * time.Millisecond,
	)

	dbFake := suite.getDBfake()
	allRepos := suite.getAllRepositories(dbFake)

	inMemoryDBfake := suite.getInMemoryDBfake()
	memoryLockRepo := suite.getMemoryLockRepoFake(inMemoryDBfake)

	tRequest := port.TransactionPaymentRequest{
		AccountUID:  accountUIDtoTransact,
		TotalAmount: amountFoodFundsApproved,
		MCC:         "5555",
		Merchant:    "UBER EATS                   SAO PAULO BR",
	}

	//Act
	paymentService := NewPayment(
		timeoutSLA,
		allRepos.Account,
		allRepos.Merchant,
		memoryLockRepo,
		nil,
	)
	returnCode, err := paymentService.Execute(tRequest)

	//Assert
	codeApproved := "00" // domain.CODE_APPROVED

	assert.Equal(suite.T(), returnCode, codeApproved)
	assert.Equal(suite.T(), err, nil)

	foodTransaction, err := getLastTransaction(dbFake.Transactions, port.TransactionEntity{AccountID: 1, CategoryID: foodCategoryID})
	assert.Equal(suite.T(), foodTransaction.Amount, decimal.NewFromFloat(105.01))
	assert.Equal(suite.T(), err, nil)
}

func (suite *PaymentSuite) TestL3PaymentExecuteNameMCCFallbackApproved() {
	//Arrange
	timeoutSLA := port.TimeoutSLA(
		time.Duration(timeoutSLAcfg) * time.Millisecond,
	)

	dbFake := suite.getDBfake()
	allRepos := suite.getAllRepositories(dbFake)

	inMemoryDBfake := suite.getInMemoryDBfake()
	memoryLockRepo := suite.getMemoryLockRepoFake(inMemoryDBfake)

	tRequest := port.TransactionPaymentRequest{
		AccountUID:  accountUIDtoTransact,
		TotalAmount: amountFoodFundsFallbackApproved,
		MCC:         "5555",
		Merchant:    "UBER EATS                   SAO PAULO BR",
	}

	//Act
	paymentService := NewPayment(
		timeoutSLA,
		allRepos.Account,
		allRepos.Merchant,
		memoryLockRepo,
		nil,
	)
	returnCode, err := paymentService.Execute(tRequest)

	//Assert
	codeApproved := "00" // domain.CODE_APPROVED

	assert.Equal(suite.T(), returnCode, codeApproved)
	assert.Equal(suite.T(), err, nil)

	foodTransaction, _ := getLastTransaction(dbFake.Transactions, port.TransactionEntity{AccountID: 1, CategoryID: foodCategoryID})
	assert.Equal(suite.T(), err, nil)
	assert.Equal(suite.T(), foodTransaction.Amount, decimal.NewFromFloat(0))

	cashTransaction, _ := getLastTransaction(dbFake.Transactions, port.TransactionEntity{AccountID: 1, CategoryID: cashCategoryID})
	assert.Equal(suite.T(), err, nil)
	assert.Equal(suite.T(), cashTransaction.Amount, decimal.NewFromFloat(90.22))
}

func getLastTransaction(transactions map[uint]port.TransactionEntity, tParams port.TransactionEntity) (*port.TransactionEntity, error) {
	var transaction port.TransactionEntity
	var maxKey uint
	var found bool

	found = false
	for key, t := range transactions {
		if tParams.AccountID == t.AccountID && tParams.CategoryID == t.CategoryID {
			if !found || key > maxKey {
				maxKey = key
				transaction = t
				found = true
			}
		}
	}

	if found {
		return &transaction, nil
	}

	return &transaction, fmt.Errorf(
		"transaction with AccountID: %v and CategoryID: %v not found",
		tParams.AccountID,
		tParams.CategoryID,
	)
}

func TestPaymentSuite(t *testing.T) {
	suite.Run(t, new(PaymentSuite))
}
