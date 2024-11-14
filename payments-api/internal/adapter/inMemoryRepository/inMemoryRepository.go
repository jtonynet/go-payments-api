package inMemoryRepository

import (
	"context"
	"fmt"

	"github.com/jtonynet/go-payments-api/internal/adapter/inMemoryRepository/redisRepos"
	"github.com/jtonynet/go-payments-api/internal/core/port"
)

func NewMerchant(cacheConn port.InMemoryDBConn, mRepository port.MerchantRepository) (port.MerchantRepository, error) {
	var mr port.MerchantRepository

	strategy, err := cacheConn.GetStrategy(context.Background())
	if err != nil {
		return mr, fmt.Errorf("error: dont retrieve cache strategy: %v", err)
	}

	switch strategy {
	case "redis":
		return redisRepos.NewMerchant(cacheConn, mRepository)
	default:
		return mr, fmt.Errorf("cached repository strategy not suported: %s", strategy)
	}
}

func NewMemoryLock(lockConn port.InMemoryDBConn) (port.MemoryLockRepository, error) {
	var mlr port.MemoryLockRepository

	strategy, err := lockConn.GetStrategy(context.Background())
	if err != nil {
		return mlr, fmt.Errorf("error: dont retrieve memory lock strategy: %v", err)
	}

	switch strategy {
	case "redis":
		return redisRepos.NewMemoryLock(lockConn)
	default:
		return mlr, fmt.Errorf("memory lock repository strategy not suported: %s", strategy)
	}
}