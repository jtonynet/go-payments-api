package port

import (
	"context"
)

type MemoryLockEntity struct {
	Key         string //accountUID
	Transcation string //transactionUID
	Timestamp   int64  //startTimestamp.UnixMilli
}

/*
- Prevent two or more transactions from the same `accountUID` from occurring concurrently
  - Retrieve or create, if it doesn’t exist, a representation of `MemoryLockEntity` in my lock source
  - Remove a representation of `MemoryLockEntity` from my lock source
*/
type MemoryLockRepository interface {
	Lock(ctx context.Context, mle MemoryLockEntity) (MemoryLockEntity, error)
	Unlock(ctx context.Context, key string) error
}
