package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jtonynet/go-payments-api/internal/adapter/database"
	"github.com/jtonynet/go-payments-api/internal/adapter/pubSub"
	"github.com/jtonynet/go-payments-api/internal/adapter/repository/gormRepos"
	"github.com/jtonynet/go-payments-api/internal/adapter/repository/redisRepos"
	"github.com/jtonynet/go-payments-api/internal/core/port"
)

type AllRepos struct {
	Account  port.AccountRepository
	Merchant port.MerchantRepository
}

func GetAll(conn database.Conn) (AllRepos, error) {
	repos := AllRepos{}

	strategy, err := conn.GetStrategy(context.Background())
	if err != nil {
		return AllRepos{}, fmt.Errorf("error when instantiating merchant repository: %v", err)
	}

	switch strategy {
	case "gorm":
		account, err := gormRepos.NewGormAccount(conn)
		if err != nil {
			return AllRepos{}, fmt.Errorf("error when instantiating account repository: %v", err)
		}
		repos.Account = account

		Merchant, err := gormRepos.NewMerchant(conn)
		if err != nil {
			return AllRepos{}, fmt.Errorf("error when instantiating merchant repository: %v", err)
		}
		repos.Merchant = Merchant

		return repos, nil
	default:
		return AllRepos{}, errors.New("repository strategy not suported: " + strategy)
	}
}

func NewCachedMerchant(cacheConn database.InMemory, mRepository port.MerchantRepository) (port.MerchantRepository, error) {
	var mr port.MerchantRepository

	strategy, err := cacheConn.GetStrategy(context.Background())
	if err != nil {
		return mr, fmt.Errorf("error: dont retrieve cache strategy: %v", err)
	}

	switch strategy {
	case "redis":
		return redisRepos.NewRedisMerchant(cacheConn, mRepository)
	default:
		return mr, fmt.Errorf("cached repository strategy not suported: %s", strategy)
	}
}

func NewMemoryLock(lockConn database.InMemory, pubsub pubSub.PubSub) (port.MemoryLockRepository, error) {
	var mlr port.MemoryLockRepository

	strategy, err := lockConn.GetStrategy(context.Background())
	if err != nil {
		return mlr, fmt.Errorf("error: dont retrieve memory lock strategy: %v", err)
	}

	switch strategy {
	case "redis":
		return redisRepos.NewMemoryLock(lockConn, pubsub)
	default:
		return mlr, fmt.Errorf("memory lock repository strategy not suported: %s", strategy)
	}
}
