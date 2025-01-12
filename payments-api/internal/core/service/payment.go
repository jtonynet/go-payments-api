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

	log               logger.Logger
	transactionLocked port.MemoryLockEntity
}

func NewPayment(
	timeoutSLA port.TimeoutSLA,

	aRepository port.AccountRepository,
	mRepository port.MerchantRepository,
	mlRepository port.MemoryLockRepository,

	log logger.Logger,
) *Payment {
	return &Payment{
		timeoutSLA:           timeoutSLA,
		accountRepository:    aRepository,
		merchantRepository:   mRepository,
		memoryLockRepository: mlRepository,

		log: log,
	}
}

func (p *Payment) Execute(tpr port.TransactionPaymentRequest) (string, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(p.timeoutSLA),
	)
	ctx = context.WithValue(ctx, logger.CtxTransactionUIDKey, tpr.TransactionUID.String())
	ctx = context.WithValue(ctx, logger.CtxAccountUIDKey, tpr.AccountUID.String())
	defer cancel()

	transactionLocked, err := p.memoryLockRepository.Lock(
		ctx,
		mapTransactionRequestToMemoryLockEntity(tpr),
	)
	if err != nil {
		return p.rejectedGenericErr(
			ctx,
			fmt.Errorf("failed concurrent transaction locked: %w", err),
		)
	}
	p.transactionLocked = transactionLocked

	accountEntity, err := p.accountRepository.FindByUID(ctx, tpr.AccountUID)
	if err != nil {
		return p.rejectedGenericErr(
			ctx,
			fmt.Errorf("failed to retrieve account entity: %w", err),
		)
	}

	account := mapAccountEntityToDomain(accountEntity, p.log)

	var merchant domain.Merchant
	merchantEntity, err := p.merchantRepository.FindByName(ctx, tpr.Merchant)
	if err != nil {
		return p.rejectedGenericErr(
			ctx,
			fmt.Errorf("failed to retrieve merchant entity with name %s", tpr.Merchant),
		)
	}

	if merchantEntity != nil {
		merchant = mapMerchantEntityToDomain(merchantEntity)
	}

	transaction := merchant.NewTransaction(
		tpr.TransactionUID,
		tpr.MCC,
		tpr.TotalAmount,
		tpr.Merchant,
		account,
	)

	approvedTransactions, cErr := account.ApproveTransaction(ctx, transaction)
	if cErr != nil {
		return p.rejectedCustomErr(ctx, cErr)
	}

	err = p.accountRepository.SaveTransactions(
		ctx,
		mapTransactionDomainsToEntities(approvedTransactions),
	)
	if err != nil {
		return p.rejectedGenericErr(
			ctx,
			fmt.Errorf("failed to save transaction entity: %w", err),
		)
	}

	_ = p.memoryLockRepository.Unlock(ctx, p.transactionLocked.Key)

	return domain.CODE_APPROVED, nil
}

func (p *Payment) rejectedGenericErr(ctx context.Context, err error) (string, error) {
	p.log.Error(ctx, err.Error())

	_ = p.memoryLockRepository.Unlock(ctx, p.transactionLocked.Key)

	return domain.CODE_REJECTED_GENERIC, err
}

func (p *Payment) rejectedCustomErr(ctx context.Context, cErr *domain.CustomError) (string, error) {
	if cErr.Code == domain.CODE_REJECTED_GENERIC {
		p.log.Error(ctx, cErr.Error())
	} else {
		p.log.Warn(ctx, cErr.Error())
	}

	_ = p.memoryLockRepository.Unlock(ctx, p.transactionLocked.Key)

	return cErr.Code, fmt.Errorf("failed to approve transaction: %s", cErr.Message)
}
