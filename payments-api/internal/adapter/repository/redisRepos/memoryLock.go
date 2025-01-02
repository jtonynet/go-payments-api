package redisRepos

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jtonynet/go-payments-api/internal/adapter/database"
	"github.com/jtonynet/go-payments-api/internal/adapter/pubSub"
	"github.com/jtonynet/go-payments-api/internal/core/port"
)

type MemoryLock struct {
	lockConn database.InMemory
	pubsub   pubSub.PubSub
}

func NewMemoryLock(lockConn database.InMemory, pubsub pubSub.PubSub) (port.MemoryLockRepository, error) {
	return &MemoryLock{
		lockConn,
		pubsub,
	}, nil
}

func (ml *MemoryLock) Lock(
	ctx context.Context,
	mle port.MemoryLockEntity,
) (port.MemoryLockEntity, error) {
	expiration, err := ml.lockConn.GetDefaultExpiration(ctx)
	if err != nil {
		return port.MemoryLockEntity{}, err
	}

	isUnlocked, _ := ml.isUnlocked(ctx, mle)
	if isUnlocked {
		err = ml.lockConn.Set(ctx, mle.Key, mle.Timestamp, expiration)
		if err != nil {
			return port.MemoryLockEntity{}, err
		}

		return mle, nil
	}

	accountTransactionKey := fmt.Sprintf(`%s:%s`, mle.Key, mle.Transcation)
	unlockSubscription, err := ml.pubsub.Subscribe(context.Background(), accountTransactionKey)
	if err != nil {
		return port.MemoryLockEntity{}, err
	}
	defer func() {
		ml.pubsub.UnSubscribe(context.Background(), accountTransactionKey)
	}()

	deadline, ok := ctx.Deadline()
	if !ok {
		log.Fatalf("cannot acquire deadline from context")
	}

	for {
		select {
		case <-unlockSubscription:
			err = ml.lockConn.Set(ctx, mle.Key, mle.Timestamp, expiration)
			if err != nil {
				return port.MemoryLockEntity{}, err
			}
			return mle, nil
		case <-time.After(time.Until(deadline)):
			return port.MemoryLockEntity{}, fmt.Errorf("timeout waiting for lock release on key: %s", mle.Key)
		case <-ctx.Done():
			return port.MemoryLockEntity{}, ctx.Err()
		}
	}

}

func (ml *MemoryLock) Unlock(ctx context.Context, key string) error {
	return ml.lockConn.Expire(ctx, key, 0)
}

func (ml *MemoryLock) isUnlocked(ctx context.Context, mle port.MemoryLockEntity) (bool, error) {
	locked, err := ml.get(ctx, mle.Key)
	if err == nil {
		elapsedTime := time.Now().UnixMilli() - locked.Timestamp
		return false, fmt.Errorf("key is already locked by another process, held for %v ms", elapsedTime)
	}

	return true, nil
}

func (ml *MemoryLock) get(_ context.Context, key string) (port.MemoryLockEntity, error) {
	timestampStr, err := ml.lockConn.Get(context.Background(), key)
	if err != nil {
		return port.MemoryLockEntity{}, err
	}

	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return port.MemoryLockEntity{}, err
	}

	return port.MemoryLockEntity{
		Key:       key,
		Timestamp: timestamp,
	}, nil
}
