package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jtonynet/go-payments-api/internal/core/domain"
	"github.com/jtonynet/go-payments-api/internal/core/port"
	"github.com/jtonynet/go-payments-api/internal/support/logger"
)

type Payment struct {
	timeoutSLA           port.TimeoutSLA
	accountRepository    port.AccountRepository
	merchantRepository   port.MerchantRepository
	memoryLockRepository port.MemoryLockRepository

	logger            logger.Logger
	transactionLocked port.MemoryLockEntity
}

func NewPayment(
	timeoutSLA port.TimeoutSLA,

	aRepository port.AccountRepository,
	mRepository port.MerchantRepository,
	mlRepository port.MemoryLockRepository,

	logger logger.Logger,
) *Payment {
	return &Payment{
		timeoutSLA:           timeoutSLA,
		accountRepository:    aRepository,
		merchantRepository:   mRepository,
		memoryLockRepository: mlRepository,

		logger: logger,
	}
}

func (p *Payment) Execute(tpr port.TransactionPaymentRequest) (string, error) {
	timeout := time.Duration(p.timeoutSLA)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	transactionLocked, err := p.memoryLockRepository.Lock(
		ctx,
		mapTransactionRequestToMemoryLockEntity(tpr),
	)
	if err != nil {
		return p.rejectedGenericErr(fmt.Errorf("failed concurrent transaction locked: %w", err))
	}
	p.transactionLocked = transactionLocked

	accountEntity, err := p.accountRepository.FindByUID(ctx, tpr.AccountUID)
	if err != nil {
		return p.rejectedGenericErr(fmt.Errorf("failed to retrieve account entity: %w", err))
	}

	account := mapAccountEntityToDomain(accountEntity)
	account.Logger = p.logger

	var merchant domain.Merchant
	merchantEntity, err := p.merchantRepository.FindByName(ctx, tpr.Merchant)
	if err != nil {
		return p.rejectedGenericErr(fmt.Errorf("failed to retrieve merchant entity with name %s", tpr.Merchant))
	}

	if merchantEntity != nil {
		merchant = mapMerchantEntityToDomain(merchantEntity)
	}

	transaction := merchant.NewTransaction(
		tpr.MCC,
		tpr.TotalAmount,
		tpr.Merchant,
		account,
	)

	approvedTransactions, cErr := account.ApproveTransaction(transaction)
	if cErr != nil {
		return p.rejectedCustomErr(cErr)
	}

	err = p.accountRepository.SaveTransactions(
		ctx,
		mapTransactionDomainsToEntities(approvedTransactions),
	)
	if err != nil {
		return p.rejectedGenericErr(fmt.Errorf("failed to save transaction entity: %w", err))
	}

	_ = p.memoryLockRepository.Unlock(context.Background(), p.transactionLocked.Key)

	return domain.CODE_APPROVED, nil
}

func (p *Payment) rejectedGenericErr(err error) (string, error) {
	p.debugLog(err.Error())

	_ = p.memoryLockRepository.Unlock(context.Background(), p.transactionLocked.Key)

	return domain.CODE_REJECTED_GENERIC, err
}

func (p *Payment) rejectedCustomErr(cErr *domain.CustomError) (string, error) {
	p.debugLog(cErr.Error())

	_ = p.memoryLockRepository.Unlock(context.Background(), p.transactionLocked.Key)

	return cErr.Code, fmt.Errorf("failed to approve transaction: %s", cErr.Message)
}

func (p *Payment) debugLog(msg string) {
	if p.logger != nil {
		p.logger.Debug(msg)
	}
}
